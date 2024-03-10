package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jpillora/backoff"
	"github.com/rpcxio/libkv/store/consul"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/util"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

func init() {
	consul.Register()
}

type Discovery struct {
	pairsMu    sync.RWMutex
	pairs      []*client.KVPair
	chans      []chan []*client.KVPair
	initialize chan int8 // 初始化管道，主要用于阻塞到第一次拉取服务完成
	mu         sync.Mutex

	filter client.ServiceDiscoveryFilter

	client *api.Client

	service string
	tag     string
	dc      string
}

func NewDiscovery(discoveryService DiscoveryService, client *api.Client) (*Discovery, error) {
	d := &Discovery{
		service:    discoveryService.Service,
		tag:        discoveryService.Tag,
		dc:         discoveryService.Dc,
		client:     client,
		initialize: make(chan int8, 0),
	}
	go d.watch()
	<-d.initialize
	return d, nil
}

func (d *Discovery) Clone(servicePath string) (client.ServiceDiscovery, error) {
	discoveryService := DiscoveryService{
		Service: servicePath,
		Tag:     d.tag,
		Dc:      d.dc,
	}
	return NewDiscovery(discoveryService, d.client)
}

// SetFilter sets the filer.
func (d *Discovery) SetFilter(filter client.ServiceDiscoveryFilter) {
	d.filter = filter
}

// GetServices returns the servers
func (d *Discovery) GetServices() []*client.KVPair {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()
	return d.pairs
}

// WatchService returns a nil chan.
func (d *Discovery) WatchService() chan []*client.KVPair {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan []*client.KVPair, 10)
	d.chans = append(d.chans, ch)
	return ch
}

func (d *Discovery) RemoveWatcher(ch chan []*client.KVPair) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var chans []chan []*client.KVPair
	for _, c := range d.chans {
		if c == ch {
			continue
		}

		chans = append(chans, c)
	}
	d.chans = chans
}

func (d *Discovery) watch() {
	initFlag := true
	bck := &backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    10 * time.Millisecond,
		Max:    time.Minute,
	}
	var lastIndex uint64
	for {
		ss, meta, err := d.client.Health().Service(
			d.service,
			d.tag,
			true,
			&api.QueryOptions{
				WaitIndex:         lastIndex,
				WaitTime:          time.Second * 5,
				Datacenter:        d.dc,
				AllowStale:        true,
				RequireConsistent: false,
			},
		)
		if err != nil {
			logx.Errorf("[Consul resolver] Couldn't fetch endpoints. target={%s}; error={%v}", d.service, err)
			time.Sleep(bck.Duration())
			continue
		}

		bck.Reset()
		lastIndex = meta.LastIndex
		logx.Debugf("[Consul resolver] %d endpoints fetched in(+wait) %s for target={%s}",
			len(ss),
			meta.RequestTime,
			d.service,
		)
		pairs := make([]*client.KVPair, 0, len(ss))

		for _, s := range ss {
			address := s.Service.Address
			if s.Service.Address == "" {
				address = s.Node.Address
			}

			key := fmt.Sprintf("%s:%d", address, s.Service.Port)
			pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(s.Service.Meta)}
			if d.filter != nil && !d.filter(pair) {
				continue
			}
			pairs = append(pairs, pair)
		}
		d.pairsMu.Lock()
		d.pairs = pairs
		if initFlag {
			d.initialize <- 0
			initFlag = false
		}
		d.pairsMu.Unlock()
		d.mu.Lock()
		for _, ch := range d.chans {
			ch := ch
			go func() {
				defer func() {
					recover()
				}()
				select {
				case ch <- d.pairs:
				case <-time.After(time.Minute):
					logx.Info("chan is full and new change has been dropped")
				}
			}()
		}
		d.mu.Unlock()
	}
}

func (d *Discovery) Close() {
	if d.chans != nil {
		for _, c := range d.chans {
			close(c)
		}
	}
	close(d.initialize)
}

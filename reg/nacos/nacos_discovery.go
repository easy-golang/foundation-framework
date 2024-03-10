package nacos

import (
	"fmt"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/util"
)

type Discovery struct {
	servicePath string
	// nacos client config
	ClientConfig constant.ClientConfig
	// nacos server config
	ServerConfig []constant.ServerConfig
	Cluster      string
	Group        string

	namingClient naming_client.INamingClient

	pairsMu sync.RWMutex
	pairs   []*client.KVPair
	chans   []chan []*client.KVPair
	mu      sync.Mutex

	filter                  client.ServiceDiscoveryFilter
	RetriesAfterWatchFailed int

	stopCh chan struct{}
}

func NewDiscovery(discoveryService DiscoveryService, clientConfig constant.ClientConfig,
	serverConfig []constant.ServerConfig) (client.ServiceDiscovery, error) {
	d := &Discovery{
		servicePath:  discoveryService.Service,
		Cluster:      discoveryService.Cluster,
		Group:        discoveryService.Group,
		ClientConfig: clientConfig,
		ServerConfig: serverConfig,
		stopCh:       make(chan struct{}),
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"clientConfig":  clientConfig,
		"serverConfigs": serverConfig,
	})
	if err != nil {
		logx.Errorf("failed to create NacosDiscovery: %v", err)
		return nil, err
	}

	d.namingClient = namingClient

	d.fetch()
	go d.watch()
	return d, nil
}

func (d *Discovery) fetch() {
	service, err := d.namingClient.GetService(vo.GetServiceParam{
		ServiceName: d.servicePath,
		Clusters:    []string{d.Cluster},
		GroupName:   d.Group,
	})
	if err != nil {
		logx.Must(comm.Newf("failed to get service %s: %+v", d.servicePath, err))
	}
	pairs := make([]*client.KVPair, 0, len(service.Hosts))
	for _, inst := range service.Hosts {
		if !inst.Enable {
			continue
		}
		//network := inst.Metadata["network"]
		ip := inst.Ip
		port := inst.Port
		key := fmt.Sprintf("%s:%d", ip, port)
		pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(inst.Metadata)}
		if d.filter != nil && !d.filter(pair) {
			continue
		}
		pairs = append(pairs, pair)
	}

	d.pairsMu.Lock()
	d.pairs = pairs
	d.pairsMu.Unlock()
}

func (d *Discovery) Clone(servicePath string) (client.ServiceDiscovery, error) {
	discoveryService := DiscoveryService{
		Service: servicePath,
		Cluster: d.Cluster,
		Group:   d.Group,
	}
	return NewDiscovery(discoveryService, d.ClientConfig, d.ServerConfig)
}

func (d *Discovery) SetFilter(filter client.ServiceDiscoveryFilter) {
	d.filter = filter
}

func (d *Discovery) GetServices() []*client.KVPair {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()

	return d.pairs
}

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
	param := &vo.SubscribeParam{
		ServiceName: d.servicePath,
		Clusters:    []string{d.Cluster},
		GroupName:   d.Group,
		SubscribeCallback: func(services []model.Instance, err error) {
			pairs := make([]*client.KVPair, 0, len(services))
			for _, inst := range services {
				//network := inst.Metadata["network"]
				ip := inst.Ip
				port := inst.Port
				key := fmt.Sprintf("%s:%d", ip, port)
				pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(inst.Metadata)}
				if d.filter != nil && !d.filter(pair) {
					continue
				}
				pairs = append(pairs, pair)
			}
			d.pairsMu.Lock()
			d.pairs = pairs
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
		},
	}

	err := d.namingClient.Subscribe(param)
	// if failed to Subscribe, retry
	if err != nil {
		var tempDelay time.Duration
		retry := d.RetriesAfterWatchFailed
		for d.RetriesAfterWatchFailed < 0 || retry >= 0 {
			err := d.namingClient.Subscribe(param)
			if err != nil {
				if d.RetriesAfterWatchFailed > 0 {
					retry--
				}
				if tempDelay == 0 {
					tempDelay = 1 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 30 * time.Second; tempDelay > max {
					tempDelay = max
				}
				logx.Infof("can not subscribe (with retry %d, sleep %v): %s: %v", retry, tempDelay, d.servicePath, err)
				time.Sleep(tempDelay)
				continue
			}
			break
		}
	}
}

func (d *Discovery) Close() {
	close(d.stopCh)
}

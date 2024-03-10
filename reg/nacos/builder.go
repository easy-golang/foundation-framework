package nacos

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
	"net"
	"strconv"
)

func init() {
	resolver.Register(&builder{})
}

const schemeName = "nacos"

// builder implements resolver.Builder and use for constructing all consul resolvers
type builder struct{}

func (b *builder) Build(url resolver.Target, conn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	tgt, err := parseURL(url.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong nacos URL")
	}

	host, ports, err := net.SplitHostPort(tgt.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed parsing address error: %v", err)
	}
	port, _ := strconv.ParseUint(ports, 10, 16)

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port),
	}

	cc := &constant.ClientConfig{
		AppName:     tgt.AppName,
		NamespaceId: tgt.NamespaceID,
		Username:    tgt.User,
		Password:    tgt.Password,
		TimeoutMs:   uint64(tgt.Timeout),
	}

	if tgt.CacheDir != "" {
		cc.CacheDir = tgt.CacheDir
	}
	if tgt.LogDir != "" {
		cc.LogDir = tgt.LogDir
	}
	if tgt.LogLevel != "" {
		cc.LogLevel = tgt.LogLevel
	}

	cli, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't connect to the nacos API")
	}

	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []string)
	go populateEndpoints(ctx, conn, pipe)
	// 解决nacos缓存阻塞问题
	service, err := cli.GetService(vo.GetServiceParam{
		Clusters:    tgt.Clusters,
		ServiceName: tgt.Service,
		GroupName:   tgt.GroupName,
	})
	callBackHandle := newWatcher(ctx, cancel, pipe).CallBackHandle
	callBackHandle(service.Hosts, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't get service from nacos cache")
	}

	go cli.Subscribe(&vo.SubscribeParam{
		ServiceName:       tgt.Service,
		Clusters:          tgt.Clusters,
		GroupName:         tgt.GroupName,
		SubscribeCallback: callBackHandle, // required
	})

	return &resolvr{cancelFunc: cancel}, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *builder) Scheme() string {
	return schemeName
}

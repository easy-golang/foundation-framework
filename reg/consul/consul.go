package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/smallnest/rpcx/server"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"strconv"
)

type Conf struct {
	consul.Conf
	JsonRpcConf JsonRpcConf `json:",optional"`
}

type JsonRpcConf struct {
	DiscoveryServices []DiscoveryService `json:",optional"`
	ListenOn          string             `json:",optional"`
	Tag               []string           `json:",optional"`
	TTL               int                `json:"ttl,optional"`
}

type DiscoveryService struct {
	Service string
	Tag     string `json:",optional"`
	Dc      string `json:",optional"`
}

type registry struct {
	conf     Conf
	listenOn string
	client   *api.Client
}

func NewRestRegistry(conf Conf, restConf rest.RestConf) *registry {
	rpcServerConf := zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: restConf.Name,
		},
		ListenOn: restConf.Host + ":" + strconv.Itoa(restConf.Port),
	}
	return NewRegistry(conf, rpcServerConf)
}

func NewRegistry(conf Conf, rpcServerConf zrpc.RpcServerConf) *registry {
	client, err := api.NewClient(&api.Config{Scheme: "http", Address: conf.Host, Token: conf.Token})
	if err != nil {
		logx.Must(comm.WrapNew("create consul client error", err))
	}
	return &registry{
		conf:     conf,
		listenOn: rpcServerConf.ListenOn,
		client:   client,
	}
}

func (r *registry) NewJsonRpcClient() *jsonrpc.Client {
	discoveryMap := make(map[string]jsonrpc.Selector)
	for _, discoveryService := range r.conf.JsonRpcConf.DiscoveryServices {

		discovery, err := NewDiscovery(discoveryService, r.client)
		if err != nil {
			logx.Must(comm.Wrap(err))
		}
		discoveryMap[discoveryService.Service] = jsonrpc.NewRoundRobinSelector(discoveryService.Service, discovery)
	}
	return jsonrpc.NewRpcClient(discoveryMap)
}

func (r *registry) NewJsonRpcRegister(server *server.Server) *jsonrpc.Register {
	registerPlugin, err := NewRegisterPlugin(r.client, r.conf)
	if err != nil {
		logx.Must(comm.Wrap(err))
	}
	return jsonrpc.NewRegister(r.conf.Key, server, registerPlugin)
}

func (r *registry) RegisterService() {
	RegisterService(r.listenOn, r.conf.Conf)
}

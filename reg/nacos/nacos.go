package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/rpcxio/rpcx-nacos/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/wangliujing/foundation-framework/util/ip"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
	"strconv"
)

type Conf struct {
	Hosts               []Host
	Namespace           string      `json:",default=public"`
	NotLoadCacheAtStart bool        `json:",default=true"`
	LogLevel            string      `json:",default=info"` // the level of log, it's must be debug,info,warn,error, default value is info
	LogDir              string      `json:",default=/tmp/nacos/log"`
	CacheDir            string      `json:",default=/tmp/nacos/cache"`
	JsonRpcConf         JsonRpcConf `json:",optional"`
}
type JsonRpcConf struct {
	DiscoveryServices []DiscoveryService `json:",optional"`
	ListenOn          string
}

type DiscoveryService struct {
	Service string
	Cluster string `json:",optional"`
	Group   string `json:",optional"`
}

type Host struct {
	Ip   string
	Port uint64
}

type registry struct {
	conf        Conf
	sc          []constant.ServerConfig
	cc          *constant.ClientConfig
	serviceName string
	listenOn    string
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
	// 注册服务
	var sc []constant.ServerConfig
	for _, host := range conf.Hosts {
		sc = append(sc, *constant.NewServerConfig(host.Ip, host.Port))
	}

	cc := &constant.ClientConfig{
		NamespaceId:         conf.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: conf.NotLoadCacheAtStart,
		LogDir:              conf.LogDir,
		CacheDir:            conf.CacheDir,
		LogLevel:            conf.LogLevel,
	}

	return &registry{
		conf:        conf,
		sc:          sc,
		cc:          cc,
		serviceName: rpcServerConf.Name,
		listenOn:    rpcServerConf.ListenOn,
	}
}

func (r *registry) NewJsonRpcRegister(server *server.Server) *jsonrpc.Register {
	pubListenOn := ip.FigureOutListenOn(r.conf.JsonRpcConf.ListenOn)
	registerPlugin := &serverplugin.NacosRegisterPlugin{
		ServiceAddress: "tcp@" + pubListenOn,
		ClientConfig:   *r.cc,
		ServerConfig:   r.sc,
	}
	err := registerPlugin.Start()
	if err != nil {
		logx.Must(comm.Wrap(err))
	}
	return jsonrpc.NewRegister(r.serviceName, server, registerPlugin)
}

func (r *registry) NewJsonRpcClient() *jsonrpc.Client {
	discoveryMap := make(map[string]jsonrpc.Selector)
	for _, discoveryService := range r.conf.JsonRpcConf.DiscoveryServices {
		discovery, err := NewDiscovery(discoveryService, *r.cc, r.sc)
		if err != nil {
			logx.Must(comm.Wrap(err))
		}
		discoveryMap[discoveryService.Service] = jsonrpc.NewRoundRobinSelector(discoveryService.Service, discovery)
	}
	return jsonrpc.NewRpcClient(discoveryMap)
}

func (r *registry) RegisterService() {
	opts := nacos.NewNacosConfig(r.serviceName, r.listenOn, r.sc, r.cc)
	nacos.RegisterService(opts)
}

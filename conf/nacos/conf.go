package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"reflect"
)

type Conf struct {
	Hosts               []Host
	Namespace           string `json:",default=public"`
	NotLoadCacheAtStart bool   `json:",default=true"`
	LogLevel            string `json:",default=info"` // the level of log, it's must be debug,info,warn,error, default value is info
	LogDir              string `json:",default=/tmp/nacos/log"`
	CacheDir            string `json:",default=/tmp/nacos/cache"`
	ConfigParam         ConfigParam
}

type ConfigParam struct {
	DataId string
	Group  string        `json:",default=DEFAULT_GROUP"`
	Type   vo.ConfigType `json:",default=yaml"`
}

type Host struct {
	Ip   string
	Port uint64
}

var loaders = map[vo.ConfigType]func([]byte, any) error{

	vo.JSON: conf.LoadFromJsonBytes,
	//".toml": conf.LoadFromTomlBytes,
	vo.YAML: conf.LoadFromYamlBytes,
	//".yml":  conf.LoadFromYamlBytes,
}

func MustLoad(configFile string, v any, opts ...Option) {
	var c Conf
	conf.MustLoad(configFile, &c)
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		logx.Must(comm.New("param v must be pointer type"))
	}

	configClient := initConfigClient(c)
	data, err := configClient.GetConfig(vo.ConfigParam{
		DataId: c.ConfigParam.DataId,
		Group:  c.ConfigParam.Group,
	})
	if err != nil {
		logx.Must(err)
	}
	loader, ok := loaders[c.ConfigParam.Type]
	if !ok {
		logx.Must(comm.Newf("unrecognized file type: %s", c.ConfigParam.Type))
	}
	var opt options
	for _, o := range opts {
		o(&opt)
	}
	parseContent(data, v, loader, opt)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: c.ConfigParam.DataId,
		Group:  c.ConfigParam.Group,
		OnChange: func(namespace, group, dataId, data string) {
			logx.Infof("nacos config changed:\n%s", data)
			parseContent(data, v, loader, opt)
		},
	})
}

func parseContent(data string, v any, loader func([]byte, any) error, opt options) {
	if opt.env {
		err := loader([]byte(os.ExpandEnv(data)), v)
		if err != nil {
			logx.Must(err)
		}
		return
	}
	loader([]byte(data), v)
}

func initConfigClient(conf Conf) config_client.IConfigClient {
	// nacos服务
	var sc []constant.ServerConfig
	for _, host := range conf.Hosts {
		sc = append(sc, *constant.NewServerConfig(host.Ip, host.Port))
	}
	cc := constant.ClientConfig{
		NamespaceId:         conf.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: conf.NotLoadCacheAtStart,
		LogDir:              conf.LogDir,
		CacheDir:            conf.CacheDir,
		LogLevel:            conf.LogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		constant.KEY_SERVER_CONFIGS: sc,
		constant.KEY_CLIENT_CONFIG:  cc,
	})
	if err != nil {
		logx.Must(err)
	}
	return configClient
}

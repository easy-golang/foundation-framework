package xxl

import (
	"github.com/easy-golang/foundation-framework/util/ip"
	"time"
)

type Options struct {
	ServerAddr   string        `json:"server_addr"`   //调度中心地址
	AccessToken  string        `json:"access_token"`  //请求令牌
	Timeout      time.Duration `json:"timeout"`       //接口超时时间
	ExecutorIp   string        `json:"executor_ip"`   //本地(执行器)IP(可自行获取)
	ExecutorPort string        `json:"executor_port"` //本地(执行器)端口
	RegistryKey  string        `json:"registry_key"`  //执行器名称
}

func newOptions(opts ...Option) Options {
	opt := Options{
		ExecutorIp:   ip.GetInternalIp(),
		ExecutorPort: DefaultExecutorPort,
		RegistryKey:  DefaultRegistryKey,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

type Option func(o *Options)

var (
	DefaultExecutorPort = "9999"
	DefaultRegistryKey  = "golang-jobs"
)

// AdminAddress 设置调度中心地址
func AdminAddress(addr string) Option {
	return func(o *Options) {
		o.ServerAddr = addr
	}
}

// Token 请求令牌
func Token(token string) Option {
	return func(o *Options) {
		o.AccessToken = token
	}
}

// ClientIp 设置执行器IP
func ClientIp(ip string) Option {
	return func(o *Options) {
		o.ExecutorIp = ip
	}
}

// ClientPort 设置执行器端口
func ClientPort(port string) Option {
	return func(o *Options) {
		o.ExecutorPort = port
	}
}

// RegistryKey 设置执行器标识
func RegistryKey(registryKey string) Option {
	return func(o *Options) {
		o.RegistryKey = registryKey
	}
}

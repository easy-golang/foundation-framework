package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/easy-golang/foundation-framework/common/dto"
	"github.com/easy-golang/foundation-framework/common/security"
	"github.com/easy-golang/foundation-framework/constant/symbol"
	"github.com/easy-golang/foundation-framework/err/comm"
	"github.com/easy-golang/foundation-framework/util/ip"
	"github.com/easy-golang/foundation-framework/util/ustring"
	"github.com/hashicorp/consul/api"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	rpcxUtil "github.com/smallnest/rpcx/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"net"
	"strconv"
	"time"
)

type RegisterPlugin struct {
	client            *api.Client
	host              string
	port              uint64
	serviceNamePrefix string
	jsonRpcConf       JsonRpcConf
}

func NewRegisterPlugin(client *api.Client, c Conf) (*RegisterPlugin, error) {
	if len(c.JsonRpcConf.ListenOn) == 0 {
		return nil, comm.New("jsonrpc port is not set")
	}
	pubListenOn := ip.FigureOutListenOn(c.JsonRpcConf.ListenOn)
	host, ports, err := net.SplitHostPort(pubListenOn)
	if err != nil {
		return nil, fmt.Errorf("failed parsing address error: %v", err)
	}
	port, _ := strconv.ParseUint(ports, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("create consul client error: %v", err)
	}

	return &RegisterPlugin{
		client:            client,
		host:              host,
		port:              port,
		serviceNamePrefix: c.Key + symbol.Colon,
		jsonRpcConf:       c.JsonRpcConf,
	}, nil
}

func (p *RegisterPlugin) Stop() error {
	return nil
}

func (p *RegisterPlugin) Register(name string, service any, metadata string) error {
	// 服务节点的名称
	serviceID := fmt.Sprintf("%s-%s-%d", name, p.host, p.port)
	if p.jsonRpcConf.TTL <= 0 {
		p.jsonRpcConf.TTL = 20
	}
	ttl := fmt.Sprintf("%ds", p.jsonRpcConf.TTL)
	expiredTTL := fmt.Sprintf("%ds", p.jsonRpcConf.TTL*3)
	meta := rpcxUtil.ConvertMeta2Map(metadata)
	reg := &api.AgentServiceRegistration{
		ID:      serviceID,         // 服务节点的名称
		Name:    name,              // 服务名称
		Tags:    p.jsonRpcConf.Tag, // tag，可以为空
		Meta:    meta,              // meta， 可以为空
		Port:    int(p.port),       // 服务端口
		Address: p.host,            // 服务 IP
		Checks: []*api.AgentServiceCheck{ // 健康检查
			{
				CheckID:                        serviceID, // 服务节点的名称
				TTL:                            ttl,       // 健康检查间隔
				Status:                         "passing",
				DeregisterCriticalServiceAfter: expiredTTL, // 注销时间，相当于过期时间
			},
		},
	}

	if err := p.client.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("initial register service '%s' host to consul error: %s", name, err.Error())
	}
	check := api.AgentServiceCheck{TTL: ttl, Status: "passing", DeregisterCriticalServiceAfter: expiredTTL}
	if err := p.client.Agent().CheckRegister(&api.AgentCheckRegistration{ID: serviceID, Name: name,
		ServiceID: serviceID, AgentServiceCheck: check}); err != nil {
		return fmt.Errorf("initial register service check to consul error: %s", err.Error())
	}
	ttlTicker := time.Duration(p.jsonRpcConf.TTL-1) * time.Second
	if ttlTicker < time.Second {
		ttlTicker = time.Second
	}
	// routine to update ttl
	go func() {
		ticker := time.NewTicker(ttlTicker)
		defer ticker.Stop()
		for {
			<-ticker.C
			err := p.client.Agent().UpdateTTL(serviceID, "", "passing")
			logx.Info("update ttl")
			if err != nil {
				logx.Infof("update ttl of service error: %v", err.Error())
			}
		}
	}()
	// consul deregister
	proc.AddShutdownListener(func() {
		err := p.client.Agent().ServiceDeregister(serviceID)
		if err != nil {
			logx.Info("deregister service error: ", err.Error())
		}
		logx.Info("deregistered service from consul server.")
	})
	return nil
}

func (p *RegisterPlugin) RegisterFunction(serviceName, fname string, fn any, metadata string) error {
	return p.Register(serviceName, fn, metadata)
}

func (p *RegisterPlugin) Unregister(name string) error {
	return nil
}

/*
PostReadRequest 预处理上下文
请求头：X-RPCX-Meta
请求头值：x-request-context=%7B%0A++++++++%22x-request-context%22%3A+%7B%0A++++++++++++%22userInfo%22%3A+%7B%0A++++++++++++++++%22userId%22%3A+%22123%22%2C%0A++++++++++++++++%22userName%22%3A+%22%E8%84%8F%E4%BB%B2%E8%A1%A8%22%2C%0A++++++++++++++++%22email%22%3A+%22123%40qq.com%22%2C%0A++++++++++++++++%22status%22%3A+%221%22%2C%0A++++++++++++++++%22preferredLanguage%22%3A+%22cn%22%0A++++++++++++%7D%2C%0A++++++++++++%22erpInfo%22%3A+null%0A++++++++%7D%2C%0A++++++++%22x-request-header%22%3A+%7B%0A++++++++++++%22User-Agent%22%3A+%22Mozilla%2F5.0+%28Windows+NT+10.0%3B+Win64%3B+x64%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Chrome%2F120.0.0.0+Safari%2F537.36%22%0A++++++++%7D%2C%0A++++++++%22x-request-id%22%3A+%22c21edeeb745a4d96b45b583b3f835b45.167.17026241022054781%22%0A++++%7D
*/
func (p *RegisterPlugin) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	shareContext, ok := ctx.(*share.Context)
	if !ok {
		return nil
	}
	requestContextStr := r.Metadata[security.RequestContextKey]
	if ustring.IsNotEmptyString(requestContextStr) {
		requestContext := new(dto.RequestContext)
		err := json.Unmarshal([]byte(requestContextStr), requestContext)
		if err != nil {
			return comm.WrapNew("jsonRpc上下文Json解析失败：%+v", err)
		}

		shareContext.Context = context.WithValue(shareContext.Context, security.RequestContextKey, requestContext)
		return nil
	}
	return nil
}

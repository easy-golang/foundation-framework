package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/easy-golang/foundation-framework/common/dto"
	"github.com/easy-golang/foundation-framework/common/security"
	"github.com/easy-golang/foundation-framework/err/comm"
	"github.com/easy-golang/foundation-framework/util/convertor"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"time"
)

func Start(listenOn string) *server.Server {
	server := server.NewServer()
	go func() {
		defer server.Shutdown(context.Background())
		err := server.Serve("tcp", listenOn)
		if err != nil {
			logx.Must(comm.Wrap(err))
		}
	}()
	return server
}

type Register struct {
	server          *server.Server
	applicationName string
}

func NewRegister(applicationName string, server *server.Server, registerPlugin server.RegisterPlugin) *Register {
	server.Plugins.Add(registerPlugin)
	return &Register{
		server:          server,
		applicationName: applicationName,
	}
}

func (r *Register) RegisterService(serviceName string, serviceObj any, metadata string) {
	r.server.RegisterName(r.applicationName+":"+serviceName, serviceObj, metadata)
}

type JsonRpcBody struct {
	Jsonrpc        string              `json:"jsonrpc"`
	Method         string              `json:"method"`
	Params         interface{}         `json:"params"`
	Id             string              `json:"id"`
	RequestContext *dto.RequestContext `json:"context"`
}

type JsonRpcResponse struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
}

type Error struct {
	Code    int32     `json:"code"`
	Message string    `json:"message"`
	Data    ErrorData `json:"data"`
}

type ErrorData struct {
	Class   string `json:"class"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type ResponseError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (responseError ResponseError) Error() string {
	return responseError.Message + " -> " + responseError.Details
}

type Client struct {
	discoveryMap map[string]Selector
}

func NewRpcClient(discoveryMap map[string]Selector) *Client {
	return &Client{
		discoveryMap: discoveryMap,
	}
}

func (c *Client) Call(ctx context.Context, serviceName string, serviceMethod string, args any, reply any) error {
	selector := c.discoveryMap[serviceName]
	if selector == nil {
		return comm.New("未订阅服务：" + serviceName)
	}
	address, err := selector.Selector()
	if err != nil {
		return err
	}
	requestBody := JsonRpcBody{
		Id:             convertor.Int64ToString(time.Now().UnixMilli()),
		Jsonrpc:        "2.0",
		Method:         serviceMethod,
		Params:         args,
		RequestContext: security.NewGiimallContext(ctx).GetRequestContext(),
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	requestUrl := "http://" + address
	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonData))
	logx.Debugf("JsonRpc 请求地址：%s  请求参数：%s", requestUrl, string(jsonData))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return comm.New("调用失败")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var result JsonRpcResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result.Error != nil {
		return ResponseError{Code: result.Error.Data.Code, Message: result.Error.Message,
			Details: result.Error.Data.Message}
	}

	data, err := json.Marshal(result.Result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, reply)
	if err != nil {
		return err
	}
	return nil
}

type Selector interface {
	Selector() (string, error)
}

type RoundRobinSelector struct {
	serviceName      string
	serviceDiscovery client.ServiceDiscovery
	index            int
}

func NewRoundRobinSelector(serviceName string, serviceDiscovery client.ServiceDiscovery) *RoundRobinSelector {
	return &RoundRobinSelector{
		serviceName:      serviceName,
		serviceDiscovery: serviceDiscovery,
	}
}

func (roundRobinSelector *RoundRobinSelector) Selector() (string, error) {
	index := roundRobinSelector.index
	services := roundRobinSelector.serviceDiscovery.GetServices()
	if len(services) == 0 {
		return "", comm.New("未找到服务节点：" + roundRobinSelector.serviceName)
	}
	index++
	if index > len(services)-1 {
		index = 0
	}
	roundRobinSelector.index = index
	pair := services[roundRobinSelector.index]
	return pair.Key, nil

}

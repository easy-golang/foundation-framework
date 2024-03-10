package reg

import (
	"github.com/smallnest/rpcx/server"
	"github.com/wangliujing/foundation-framework/jsonrpc"
)

type Registry interface {
	NewJsonRpcRegister(server *server.Server) *jsonrpc.Register // 获取jsonrpc服务注册器，可用于注册jsonrpc服务
	NewJsonRpcClient() *jsonrpc.Client                          // 创建jsonrpc客户端，可用于调用jsonrpc服务
	RegisterService()                                           // 应用级注册grpc服务
}

package system

import (
	"github.com/easy-golang/foundation-framework/err/biz"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"os/signal"
	"syscall"
)

func Start(service service.Service) {
	rpcServer, ok := service.(*zrpc.RpcServer)
	if ok {
		// 添加异常拦截器，兼容bizError
		rpcServer.AddUnaryInterceptors(biz.ErrorInterceptor)
	}
	defer service.Stop()
	go service.Start()
	// 捕获系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	// 关闭资源
	proc.Shutdown()
	logx.Info("服务器关闭...")
}

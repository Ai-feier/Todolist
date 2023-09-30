package shutdown

import (
	"context"
	"grpc-todolist/pkg/utils/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func GracefullyShutdown(server *http.Server) {
	// 创建系统信号接收器接收关闭信号
	done := make(chan os.Signal)
	// 监听 ctrl+c, kill 信号
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done  // 监听 goroutine
	logger.LogrusObj.Println("closing http server gracefully ...")
	
	if err := server.Shutdown(context.Background()); err != nil {
		logger.LogrusObj.Fatalln("closing http server gracefully failed:", err)
	}
}

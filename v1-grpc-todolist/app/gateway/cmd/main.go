package main

import (
	"fmt"
	"grpc-todolist/app/gateway/routers"
	"grpc-todolist/app/gateway/rpc"
	"grpc-todolist/config"
	"grpc-todolist/pkg/utils/shutdown"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	rpc.Init()
	
	// 启动服务
	go startListen()
	{
		osSignal := make(chan os.Signal)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		<-osSignal
		fmt.Println("exit gracefully...")
	}
	fmt.Println("gateway listen on :3000")
}

func startListen() {
	//opts := []grpc.DialOption{}
	//userConn, _ := grpc.Dial("127.0.0.1:10001", opts...)
	//userServer := user.NewUserServiceClient(userConn)
	// 加入熔断 TODO main太臃肿了
	// wrapper.NewServiceWrapper(userServiceName)
	// wrapper.NewServiceWrapper(taskServiceName)
	r := routers.NewRouter()
	server := &http.Server{
		Addr:                         config.Conf.Server.Port,
		Handler:                      r,
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  0,
		MaxHeaderBytes:               1 << 20,  // 限制数据交互量
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("gateway 启动失败")
		panic(err)
	}
	// 优雅退出
	go func() {
		shutdown.GracefullyShutdown(server)
	}()
}

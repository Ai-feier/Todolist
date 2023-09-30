package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc-todolist/app/task/internal/repository/dao"
	"grpc-todolist/app/task/internal/services"
	"grpc-todolist/config"
	"grpc-todolist/idl/pb/task"
	"grpc-todolist/pkg/discovery"
	"grpc-todolist/pkg/utils/logger"
	"net"
)

func main() {
	config.InitConfig()
	dao.InitDB()
	
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewRegister(etcdAddress, logger.LogrusObj)
	
	grpcAddress := config.Conf.Services["task"].Addr[0]  // grpc 连接地址
	defer etcdRegister.Stop()
	
	taskNode := discovery.Server{
		Name:    config.Conf.Domain["task"].Name,
		Addr:    grpcAddress,
	}
	
	// 启动一个 rpc 服务
	server := grpc.NewServer()
	defer server.Stop()
	
	// 绑定 service 
	task.RegisterTaskServiceServer(server, services.GetTaskSrv())
	
	// 获取一个监听
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(taskNode, 10); err != nil {  // 将当前服务注册到 etcd
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
















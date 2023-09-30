package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"micro-todolist/app/task/repository/db/dao"
	"micro-todolist/app/task/repository/mq"
	"micro-todolist/app/task/script"
	"micro-todolist/app/task/service"
	"micro-todolist/config"
	"micro-todolist/idl/pb"
	"micro-todolist/pkgs/logger"
)

func main() {
	config.Init()
	dao.InitDB()
	mq.InitRabbitMQ()
	logger.InitLog()
	loadingSrcipt()

	// etcd 注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.EtcdHost, config.EtcdPort)),
	)

	// 得到微服务实例
	microService := micro.NewService(
		micro.Name("rpcTaskService"),
		micro.Address(config.TaskServiceAddress),
		micro.Registry(etcdReg),
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterTaskServiceHandler(microService.Server(), service.GetTaskSrv())
	// 启动微服务
	_ = microService.Run()
}

func loadingSrcipt() {
	ctx := context.Background()
	go script.TaskCreateSync(ctx)
}

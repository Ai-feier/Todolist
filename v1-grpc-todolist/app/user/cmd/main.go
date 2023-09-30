package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc-todolist/app/user/internal/repository/db/dao"
	"grpc-todolist/app/user/internal/service"
	"grpc-todolist/config"
	"grpc-todolist/idl/pb/user"
	"grpc-todolist/pkg/discovery"
	"net"
)

func main() {
	config.InitConfig()
	dao.InitDB()

	// 将服务注册到 ETCD
	// 取出 ETCD 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务的注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["user"].Addr[0]
	userNode := discovery.Server{
		Name: config.Conf.Services["user"].Name,
		Addr: grpcAddress,
	}

	fmt.Println(userNode)
	// 开启一个 grpc 服务
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定服务
	user.RegisterUserServiceServer(server, service.GetUserSrv())

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err = etcdRegister.Register(userNode, 10); err != nil {
		panic(err)
	}
	// 对服务进行监听
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

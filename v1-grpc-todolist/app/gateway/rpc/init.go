package rpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc-todolist/config"
	"grpc-todolist/idl/pb/task"
	"grpc-todolist/idl/pb/user"
	"grpc-todolist/pkg/discovery"
	"log"
	"time"
)

var (
	Register *discovery.Resolver
	ctx context.Context
	CancelFunc context.CancelFunc
	
	UserClient user.UserServiceClient  // user 模块的客户端连接
	TaskClient task.TaskServiceClient  // task Module connection
)

func Init() {
	Register = discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
	resolver.Register(Register)
	
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	go Register.Close()
	initClient(config.Conf.Services["user"].Name, &UserClient)
	initClient(config.Conf.Services["task"].Name, &TaskClient)
}

func initClient(name string, client interface{}) {
	conn, err := connectServer(name)  // 获取 rpc 服务连接
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(conn.GetState())
	
	switch c := client.(type){
	case *user.UserServiceClient:
		*c = user.NewUserServiceClient(conn)  // 进行 user 模块的 rpc 服务链接
	case *task.TaskServiceClient:
		*c = task.NewTaskServiceClient(conn)
	default:
		panic("unsupported client type")
	}
}

func connectServer(name string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), name)
	

	// Load balance
	if config.Conf.Services[name].LoadBalance {
		log.Printf("load balance enabled for %s\n", name)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}
	
	conn, err = grpc.DialContext(ctx, addr, opts...)
	return  
}
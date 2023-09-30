package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type Register struct {
	EtcdAddr    []string // etcd 地址
	DialTimeout int      // 超时

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID                        // 租约
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 心跳检验, 确认服务是否存活

	srvInfo Server // 服务
	srvTTL  int64
	cli     *clientv3.Client
	logger  *logrus.Logger
}

// NewRegister 基于 ETCD 创建一个 Register
func NewRegister(etcdAddrs []string, logger *logrus.Logger) *Register {
	return &Register{
		EtcdAddr:    etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register 初始化自己的实例
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}

	// 初始化  新建连接
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddr,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl
	if err = r.register(); err != nil {  // 自定义的 ETCD 连接配置好后进行注册
		return nil, err
	}

	// 新建一个用于监听关闭的 chan
	r.closeCh = make(chan struct{})
	go r.keepAlive()  // 保证节点高可活
	return r.closeCh, nil
}

// register 新建ETCD自带的实例
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 定义一个新的租约
	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID

	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	// put 到 etcd 作为节点注册
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

func (r *Register) keepAlive() error {
	ticker := time.NewTicker(time.Duration(r.srvTTL))
	for {
		select {
		case <-r.closeCh:
			// 关闭则删除连接
			if err := r.unregister(); err != nil {
				fmt.Println("unregister failed error", err)
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				// 废除租约
				fmt.Println("revoke failed")
			}
		case res := <-r.keepAliveCh:
			// 租约过期, 重新注册
			if res == nil {
				// res 为 nil, 则重新注册一下
				if err := r.register(); err != nil {
					fmt.Println("register error")
				}
			}
		case <-ticker.C:
			// 定时器 避免以上管道阻塞  ->  保证高可活
			if r.keepAliveCh == nil {
				if err := r.register(); err!=nil {
					fmt.Println("register error")
				}
			}
		}
	}

}

// 在ETCD把这个服务器删除
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}

// Stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}
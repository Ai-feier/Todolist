package service

import (
	"context"
	"grpc-todolist/app/user/internal/repository/db/dao"
	"grpc-todolist/idl/pb/user"
	"grpc-todolist/pkg/e"
	"sync"
)

/**
	实现 grpc 暴露的服务器接口的实现
 */

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once  // 确保只进行一次初始化

type UserSrv struct {
	user.UnimplementedUserServiceServer
}

// GetUserSrv 返回服务
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

/**
	重写 grpc 方法
 */


func (u *UserSrv) UserLogin(ctx context.Context,req *user.UserRequest) (*user.UserDetailResponse, error) {
	resp := new(user.UserDetailResponse)
	resp.Code = e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	r, err := userDao.GetUserInfo(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	
	resp.UserDetail = &user.UserResponse{
		UserId:   r.UserId,
		NickName: r.NickName,
		UserName: r.UserName,
	}
	return resp, nil
}
func (u *UserSrv) UserRegister(ctx context.Context,req *user.UserRequest) (*user.UserCommonResponse, error) {
	resp := new(user.UserCommonResponse)
	resp.Code = e.SUCCESS
	err := dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = e.ERROR
		return resp, err
	}
	// 创建成功, 返回对应消息
	resp.Msg = e.GetMsg(int(resp.Code))
	return resp, nil
}
func (u *UserSrv) UserLogout(ctx context.Context,req *user.UserRequest) (resp *user.UserCommonResponse, err error) {
	return 
}
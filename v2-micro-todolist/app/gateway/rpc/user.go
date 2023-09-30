package rpc

import (
	"context"
	"errors"
	"micro-todolist/idl/pb"
	"micro-todolist/pkgs/e"
)

// UserLogin 用户登陆
func UserLogin(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	resp, err = UserService.UserLogin(ctx, req)  // 调用远程 rpc 方法
	if err != nil {
		return  
	}
	if resp.Code != e.SUCCESS {
		err = errors.New(e.GetMsg(int(resp.Code)))
		return
	}
	return
}

func UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	
	resp, err = UserService.UserRegister(ctx, req)
	if err != nil {
		return
	}
	return
}

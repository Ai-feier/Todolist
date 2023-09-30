package rpc

import (
	"context"
	"errors"
	"grpc-todolist/idl/pb/user"
	"grpc-todolist/pkg/e"
)

// UserRegister 需要获得 rpc 的客户端连接, 后将请求传入
func UserRegister(ctx context.Context,req *user.UserRequest) (resp *user.UserCommonResponse,err error) {
	resp, err = UserClient.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != e.SUCCESS {
		err = errors.New(resp.Msg)
		return 
	}
	return
}

func UserLogin(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	r, err := UserClient.UserLogin(ctx, req)
	if err != nil {
		return 
	}
	if r.Code != e.SUCCESS {
		err = errors.New("登陆失败")
		return 
	}
	return r.UserDetail, nil 
}
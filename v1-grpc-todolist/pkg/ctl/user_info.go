package ctl

import (
	"context"
	"errors"
)

type key int

var userKey key

type UserInfo struct {
	Id int64 `json:"id"`
}

func GetUserInfo(ctx context.Context)(*UserInfo, error) {
	user, ok := FormContext(ctx)
	if !ok {
		return nil, errors.New("获取用户信息失败")
	}
	return user, nil
}

func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func FormContext(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(userKey).(*UserInfo)
	return user, ok
}






func InitUserInfo(ctx context.Context) {
	// TODO 放缓存，之后的用户信息，走缓存
}




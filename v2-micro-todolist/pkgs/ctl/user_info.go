package ctl

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"micro-todolist/app/user/repository/cache"
	"time"
)

var userKey string = "user_key"

type UserInfo struct {
	Id uint `json:"id"`
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	// 从缓存中查询
	data, err := cache.RedisClient.Get(ctx, userKey).Bytes()
	if err == redis.Nil {
		// 缓存中不存在
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user UserInfo
	if err = json.Unmarshal(data, &user); err != nil {
		// 解码数据时出错
		return nil, err
	}
	return &user, nil
}

func InitUserInfo(ctx context.Context, user *UserInfo) error {
	cachedUser, err := GetUserInfo(ctx)
	if err != nil {
		return err
	}

	// 如果缓存中存在用户信息，直接返回
	if cachedUser != nil {
		return nil
	}

	// 如果缓存中不存在用户信息，将用户信息存储到缓存中
	err = SetUserInCache(ctx, user, cache.RedisClient)
	if err != nil {
		return err
	}
	return nil
}

func SetUserInCache(ctx context.Context, user *UserInfo, redisClient *redis.Client) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	if err = redisClient.Set(ctx, userKey, data, 30*time.Second).Err(); err != nil {
		return err
	}
	return nil
}

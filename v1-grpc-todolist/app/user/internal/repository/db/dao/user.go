package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"grpc-todolist/app/user/internal/repository/model"
	"grpc-todolist/idl/pb/user"
	"grpc-todolist/pkg/utils/logger"
)

type UserDao struct {
	*gorm.DB
}

// NewUserDao 返回一个对于 User 表的封装方法的 UserDao 类
func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

// GetUserInfo 获取用户信息
func (dao *UserDao) GetUserInfo(req *user.UserRequest) (r *model.User,err error) {
	err = dao.Model(&model.User{}).Where("user_name=?", req.UserName).First(&r).Error
	return
}

// CreateUser 用户创建
func (dao *UserDao) CreateUser(req *user.UserRequest) (err error) {
	var u model.User
	var count int64
	dao.Model(&model.User{}).Where("user_name = ?", req.GetUserName()).Count(&count)
	if count != 0 {
		return errors.New("UserName Exist")
	}
	u = model.User{
		UserName:       req.UserName,
		NickName:       req.NickName,
	}
	_ = u.SetPassword(req.Password)
	if err = dao.Model(&model.User{}).Create(&u).Error; err != nil {
		logger.LogrusObj.Error("Insert User Error:" + err.Error())
	}
	return
}

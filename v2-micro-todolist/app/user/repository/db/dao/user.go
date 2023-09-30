package dao

import (
	"context"
	"gorm.io/gorm"
	"micro-todolist/app/user/repository/db/model"
)

type UserDao struct {
	*gorm.DB
}

// NewUserDao 新建 user 业务持久层 
func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByUsername(username string) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("user_name=?", username).Find(&r).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *UserDao) CreateUser(in *model.User) error {
	err := dao.Model(&model.User{}).Create(&in).Error
	return err
}
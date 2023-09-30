package services

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"micro-todolist/app/user/repository/db/dao"
	"micro-todolist/app/user/repository/db/model"
	"micro-todolist/idl/pb"
	"micro-todolist/pkgs/e"
	"sync"
)

var (
	UserSrvInts *UserSrv
	UserSrvOnce sync.Once
)

type UserSrv struct {
}

func (user *UserSrv) UserLogin(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) error {
	resp.Code = e.SUCCESS
	// 查看是否有这个人
	res, err := dao.NewUserDao(ctx).FindUserByUsername(req.UserName)
	if err != nil {
		resp.Code = e.Error
		return err 
	}
	if res.ID == 0 {
		resp.Code = e.InvalidParams
		err = errors.New("用户不存在")
		return err
	}
	if !res.CheckPassword(req.Password) {
		resp.Code = e.InvalidParams
		return err
	}

	resp.UserDetail = BuildUser(res)
	return nil
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	if req.Password != req.PasswordConfirm {
		err = errors.New("两次密码输入不一致")
		return
	}
	resp.Code = e.SUCCESS
	_, err = dao.NewUserDao(ctx).FindUserByUsername(req.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 如果不存在就继续下去
			// ...continue
		} else {
			resp.Code = e.Error
			return
		}
	}
	user := &model.User{
		UserName: req.UserName,
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		resp.Code = e.Error
		return
	}
	if err = dao.NewUserDao(ctx).CreateUser(user); err != nil {
		resp.Code = e.Error
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

func BuildUser(item *model.User) *pb.UserModel {
	userModel := &pb.UserModel{
		Id: uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
	return userModel
}

// NewUserSrv 懒汉式加载
func NewUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvInts = &UserSrv{}
	})
	return UserSrvInts
}
//// NewUserSrv 饿汉式加载, 直接返回
//func NewUserSrv() *UserSrv {
//	if UserSrvInts == nil {
//		return new(UserSrv)
//	}
//	return UserSrvInts
//}
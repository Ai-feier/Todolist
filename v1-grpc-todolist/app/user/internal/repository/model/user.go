package model

import (
	"golang.org/x/crypto/bcrypt"
	"grpc-todolist/consts"
)

type User struct {
	UserId         int64   `gorm:"primarykey"`
	UserName       string `gorm:"unique"`
	NickName       string
	PasswordDigest string
}

func (*User) TableName() string {
	return "user"
}

// SetPassword 加密密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), consts.PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 检验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.PasswordDigest))
	return err == nil
}
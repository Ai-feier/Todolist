package dao

import (
	"grpc-todolist/app/user/internal/repository/model"
	"grpc-todolist/pkg/utils/logger"
	"os"
)

func migration() {
	// 实现自动迁移 user 服务下的表
	err := _db.Set("gorm:table_option", "charset=utf8mb4").
			AutoMigrate(
					&model.User{},
					)
	if err != nil {
		logger.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logger.LogrusObj.Infoln("register table success")
}

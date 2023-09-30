package dao

import (
	"grpc-todolist/app/task/internal/repository/model"
	"grpc-todolist/pkg/utils/logger"
	"os"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Task{},
		)
	if err != nil {
		logger.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logger.LogrusObj.Infoln("register table success")
}

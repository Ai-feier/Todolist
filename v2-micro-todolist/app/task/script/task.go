package script

import (
	"context"
	"micro-todolist/app/task/repository/mq/task"
	"micro-todolist/pkgs/logger"
)

func TaskCreateSync(ctx context.Context) {
	tSync := new(task.SyncTask)
	err := tSync.RunTaskCreate(ctx)
	if err != nil {
		logger.LogrusObj.Infof("RunTaskCreate:%s", err)
	}
}
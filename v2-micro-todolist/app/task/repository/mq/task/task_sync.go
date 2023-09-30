package task

import (
	"context"
	"encoding/json"
	"micro-todolist/app/task/repository/mq"
	"micro-todolist/app/task/service"
	"micro-todolist/consts"
	"micro-todolist/idl/pb"
	"micro-todolist/pkgs/logger"
)

type SyncTask struct {
}


func (s *SyncTask) RunTaskCreate(ctx context.Context) error {
	rabbitMqQueue := consts.RabbitMqTaskQueue
	msgs, err := mq.ConsumeMessage(ctx, rabbitMqQueue)
	if err != nil {
		return err
	}
	var forever chan struct{}
	
	go func() {
		for d := range msgs {
			logger.LogrusObj.Infof("Received run Task: %s", d.Body)

			// 落库
			reqRabbitMQ := new(pb.TaskRequest)
			err = json.Unmarshal(d.Body, reqRabbitMQ)
			if err != nil {
				logger.LogrusObj.Infof("Received run Task: %s", err)
			}
			err = service.TaskMQ2MySQL(ctx, reqRabbitMQ)
			if err != nil {
				logger.LogrusObj.Infof("Received run Task: %s", err)
			}
			d.Ack(false)
		}
	}()

	logger.LogrusObj.Infoln(err)
	<-forever
	return nil
}

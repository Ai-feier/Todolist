package dao

import (
	"context"
	"gorm.io/gorm"
	"grpc-todolist/app/task/internal/repository/model"
	"grpc-todolist/idl/pb/task"
	"grpc-todolist/pkg/utils/logger"
)

type TaskDao struct {
	*gorm.DB
}

func NewTaskDao(ctx context.Context) *TaskDao {
	return &TaskDao{NewDBClient(ctx)}
}

// CreateTask 创建任务
func (dao *TaskDao) CreateTask(req *task.TaskRequest) (err error) {
	t := &model.Task{
		TaskID:    req.TaskID,
		UserID:    req.UserID,
		Status: int(req.Status),
		Title:     req.Title,
		Context:   req.Content,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	if err = dao.Model(&model.Task{}).Create(&t).Error; err != nil {
		logger.LogrusObj.Error("Insert Task error ", err.Error())
		return  
	}
	return  
}

// UpdateTask 更新任务
func (dao *TaskDao) UpdateTask(req *task.TaskRequest) error {
	taskUpdateMap := make(map[string]interface{})
	taskUpdateMap["title"] = req.Title
	taskUpdateMap["context"] = req.Content
	taskUpdateMap["status"] = req.Status
	taskUpdateMap["start_time"] = req.StartTime
	taskUpdateMap["end_time"] = req.EndTime
	
	return dao.Model(&model.Task{}).Where("task_id=?", req.TaskID).Updates(taskUpdateMap).Error
}

// ListTaskByUserId 获取当前用户的所有任务
func (dao *TaskDao) ListTaskByUserId(req *task.TaskRequest) (r []*model.Task,err error) {
	err = dao.Model(&model.Task{}).Where("user_id=?", req.UserID).Find(&r).Error
	return 
}

// DeleteTaskById 删除任务
func (dao *TaskDao) DeleteTaskById(req *task.TaskRequest) (err error) {
	err = dao.Model(&model.Task{}).Where("user_id=? AND task_id=?", req.UserID, req.TaskID).
						Delete(&model.Task{}).Error
	return  
}











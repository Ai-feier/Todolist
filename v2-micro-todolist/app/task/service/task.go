package service

import (
	"context"
	"encoding/json"
	"micro-todolist/app/task/repository/db/dao"
	"micro-todolist/app/task/repository/db/model"
	"micro-todolist/app/task/repository/mq"
	"micro-todolist/idl/pb"
	"micro-todolist/pkgs/e"
	"micro-todolist/pkgs/logger"
	"sync"
)

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

type TaskSrv struct {
}

// CreateTask 创建备忘录，将备忘录信息生产，放到rabbitMQ消息队列中
func (t *TaskSrv) CreateTask(ctx context.Context, request *pb.TaskRequest, response *pb.TaskDetailResponse) error {
	body, _:= json.Marshal(request)
	response.Code = e.SUCCESS
	err := mq.SendMessage2MQ(body)
	if err != nil {
		response.Code = e.Error
		return err
	}
	return nil
}

func TaskMQ2MySQL(ctx context.Context, req *pb.TaskRequest) error {
	m := &model.Task{
		Uid: uint32(req.Uid),
		Title:     req.Title,
		Status:    int(req.Status),
		Content:   req.Content,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	return dao.NewTaskDao(ctx).CreateTask(m)
}

// GetTasksList 实现备忘录服务接口 获取备忘录列表
func (t *TaskSrv) GetTasksList(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskListResponse) error {
	resp.Code = e.SUCCESS
	if req.Limit == 0 {
		req.Limit = 10
	}
	// 查找备忘录
	r, count, err := dao.NewTaskDao(ctx).ListTaskByUserId(req.Uid, int(req.Start), int(req.Limit))
	if err != nil {
		resp.Code = e.Error
		logger.LogrusObj.Error("ListTaskByUserId err:%v", err)
		return err
	}
	// 返回proto里面定义的类型
	var taskRes []*pb.TaskModel
	for _, item := range r {
		taskRes = append(taskRes, BuildTask(item))
	}
	resp.TaskList = taskRes
	resp.Count = uint32(count)
	return nil
}

func (t *TaskSrv) GetTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	resp.Code = e.SUCCESS
	res, err := dao.NewTaskDao(ctx).GetTaskByTaskIdAndUserId(req.Id, req.Uid)
	if err != nil {
		resp.Code = e.Error
		logger.LogrusObj.Error("GetTask err:%v", err)
		return err
	}
	taskRes := BuildTask(res)
	resp.TaskDetail = taskRes
	return nil
}

func (t *TaskSrv) UpdateTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	resp.Code = e.SUCCESS
	taskData, err := dao.NewTaskDao(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = e.Error
		logger.LogrusObj.Error("UpdateTask err:%v", err)
		return err
	}
	resp.TaskDetail = BuildTask(taskData)
	return nil
}

func (t *TaskSrv) DeleteTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	resp.Code = e.SUCCESS
	err := dao.NewTaskDao(ctx).DeleteTaskByIdAndUserId(req.Id, req.Uid)
	if err != nil {
		resp.Code = e.Error
		logger.LogrusObj.Error("DeleteTask err:%v", err)
		return err
	}
	return nil
	
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns =  &TaskSrv{}
	})
	return TaskSrvIns
}


func BuildTask(item *model.Task) *pb.TaskModel {
	taskModel := pb.TaskModel{
		Id:         uint64(item.ID),
		Uid:        uint64(item.Uid),
		Title:      item.Title,
		Content:    item.Content,
		StartTime:  item.StartTime,
		EndTime:    item.EndTime,
		Status:     int64(item.Status),
		CreateTime: item.CreatedAt.Unix(),
		UpdateTime: item.UpdatedAt.Unix(),
	}
	return &taskModel
}
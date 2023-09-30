package services

import (
	"context"
	"grpc-todolist/app/task/internal/repository/dao"
	"grpc-todolist/idl/pb/task"
	"grpc-todolist/pkg/e"
	"sync"
)

var (
	TaskSrvInts *TaskSrv
	TaskSrvOnce sync.Once
)

type TaskSrv struct {
	task.UnimplementedTaskServiceServer
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvInts = &TaskSrv{}
	})
	return TaskSrvInts
}

func (t *TaskSrv) TaskCreate(ctx context.Context,taskReq *task.TaskRequest) (taskResp *task.TaskCommonResponse,err error) {
	taskResp = new(task.TaskCommonResponse)
	taskResp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).CreateTask(taskReq)
	if err != nil {
		taskResp.Code = e.ERROR
		taskResp.Msg = e.GetMsg(e.ERROR)
		taskResp.Data = err.Error()
		return  
	}
	taskResp.Msg = e.GetMsg(int(taskResp.Code))
	return  	
}

func (t *TaskSrv) TaskUpdate(ctx context.Context,req *task.TaskRequest) (resp *task.TaskCommonResponse,err error) {
	resp = new(task.TaskCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return  
	}
	resp.Msg = e.GetMsg(int(resp.Code))
	return  
}
func (t *TaskSrv) TaskShow(ctx context.Context,req *task.TaskRequest) (resp *task.TasksDetailResponse,err error) {
	resp = new(task.TasksDetailResponse)
	resp.Code = e.SUCCESS
	r, err := dao.NewTaskDao(ctx).ListTaskByUserId(req)
	if err != nil {
		resp.Code = e.ERROR
		return  
	}
	for _, m := range r {
		resp.TaskDetail = append(resp.TaskDetail, 
				&task.TaskModel{
					TaskID:    m.TaskID,
					UserID:    m.UserID,
					Status: int64(m.Status),
					Title:     m.Title,
					Content:   m.Context,
					StartTime: m.StartTime,
					EndTime:   m.EndTime,
				},
		)
	}
	return  
}
func (t *TaskSrv) TaskDelete(ctx context.Context,req *task.TaskRequest) (resp *task.TaskCommonResponse,err error) {
	resp = new(task.TaskCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).DeleteTaskById(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = e.GetMsg(e.ERROR)
		resp.Data = err.Error()
		return  
	}
	resp.Msg = e.GetMsg(int(resp.Code))
	return  
}























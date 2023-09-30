package rpc

import (
	"context"
	"errors"
	"grpc-todolist/idl/pb/task"
	"grpc-todolist/pkg/e"
)

func TaskList(ctx context.Context, req *task.TaskRequest) (resp *task.TasksDetailResponse, err error) {
	r, err := TaskClient.TaskShow(ctx, req)
	
	if err != nil {
		return  
	}
	
	if r.Code != e.SUCCESS {
		err = errors.New("获取失败")
		return  
	}
	return r, nil
}

func TaskDelete(ctx context.Context, req *task.TaskRequest) (resp *task.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskDelete(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != e.SUCCESS {
		err = errors.New(resp.Msg)
		return nil, err
	}
	return  
}

func TaskCreate(ctx context.Context, req *task.TaskRequest) (resp *task.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != e.SUCCESS {
		err = errors.New(resp.Msg)
		return nil, err
	}
	return  
}

func TaskUpdate(ctx context.Context, req *task.TaskRequest) (resp *task.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Code != e.SUCCESS{
		err = errors.New(resp.Msg)
		return nil, err
	}
	return  
}
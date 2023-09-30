package handler

import (
	"github.com/gin-gonic/gin"
	"grpc-todolist/app/gateway/rpc"
	"grpc-todolist/idl/pb/task"
	"grpc-todolist/pkg/ctl"
	"grpc-todolist/pkg/e"
	"net/http"
)

/**
	task router module
 */



// GetTaskList 获取任务列表
func GetTaskList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req task.TaskRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, ctl.RespError(c, err, "参数绑定失败"))
			return
		}
		// 通过请求获得 user 信息
		user, err := ctl.GetUserInfo(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, ctl.RespError(c, err, "获取用户信息失败"))
			return
		}
		req.UserID = user.Id
		// 通过 rpc 服务获得当前用户的所有任务列表
		r, err := rpc.TaskList(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctl.RespError(c, err, "RPC TaskList服务调用失败"))
			return 
		}
		c.JSON(e.SUCCESS, ctl.RespSuccess(c, r))
	}
} 

// CreateTask 创建任务
func CreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req task.TaskRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, ctl.RespError(c, err, "绑定参数错误"))
			return
		}
		user, err := ctl.GetUserInfo(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctl.RespError(c, err, "获取用户信息错误"))
			return
		}
		req.UserID = user.Id
		// 调用 rpc 远程服务
		resp, err := rpc.TaskCreate(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctl.RespError(c, err, "CreateTask RPC调用错误"))
			return
		}
		c.JSON(e.SUCCESS, ctl.RespSuccess(c, resp))
	}
}

// UpdateTask 更新任务
func UpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req task.TaskRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, ctl.RespError(c, err, "绑定参数错误"))
			return
		}
		user, err := ctl.GetUserInfo(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctl.RespError(c, err, "获取用户信息错误"))
			return
		}
		req.UserID = user.Id
		// 调用 rpc 服务
		resp, err := rpc.TaskUpdate(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctl.RespError(c, err, "UpdateTask RPC调用错误"))
			return
		}
		c.JSON(e.SUCCESS, ctl.RespSuccess(c, resp))
	}
}

// DeleteTask 删除任务
func DeleteTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req task.TaskRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, ctl.RespError(c, err, "绑定参数错误"))
			return
		}
		user, err := ctl.GetUserInfo(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctl.RespError(c, err, "获取用户信息错误"))
			return
		}
		req.UserID = user.Id
		// 调用 rpc 服务
		resp, err := rpc.TaskDelete(c, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, ctl.RespError(c, err, "DeleteTask RPC调用失败"))
		}
		c.JSON(e.SUCCESS, ctl.RespSuccess(c, resp))
	}
} 
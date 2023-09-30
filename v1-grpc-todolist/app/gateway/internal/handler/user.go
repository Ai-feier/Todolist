package handler

import (
	"github.com/gin-gonic/gin"
	"grpc-todolist/app/gateway/rpc"
	"grpc-todolist/idl/pb/user"
	"grpc-todolist/pkg/ctl"
	"grpc-todolist/pkg/res"
	"grpc-todolist/pkg/utils/jwt"
	"net/http"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq user.UserRequest
	if err := ctx.Bind(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "参数绑定失败"))
		return
	}
	
	r, err := rpc.UserRegister(ctx, &userReq)
	
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserRegister RPC服务错误"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	var req user.UserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "参数绑定失败"))
		return
	}
	
	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserLogin RPC服务错误"))
		return
	}
	
	// 生成 token 
	token, err := jwt.GenerateToken(userResp.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "加密错误"))
		return
	}
	
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, res.TokenData{
		User:  userResp,
		Token: token,
	}))
}

package http

import (
	"github.com/gin-gonic/gin"
	"micro-todolist/app/gateway/rpc"
	"micro-todolist/idl/pb"
	"micro-todolist/pkgs/ctl"
	"micro-todolist/pkgs/e"
	"micro-todolist/pkgs/logger"
	"micro-todolist/pkgs/utils"
	"micro-todolist/types"
	"net/http"
)

// UserRegisterHandler 用户注册
func UserRegisterHandler(ctx *gin.Context) {
	var req pb.UserRequest
	if err := ctx.ShouldBind(&req);err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "Register Bind 绑定参数失败"))
		return
	}
	// web -> rpc
	resp, err := rpc.UserRegister(ctx, &req)
	if err != nil {
		logger.LogrusObj.Errorf("UserRegister:%v", err)
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserRegister RPC 调用失败"))
		return 	
	}
	ctx.JSON(e.SUCCESS, ctl.RespSuccess(ctx, resp))
}

// UserLoginHandler 用户登录
func UserLoginHandler(ctx *gin.Context) {
	var req pb.UserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserLogin Bind 绑定参数失败"))
		return
	}
	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserLogin RPC 调用失败"))
		return
	}
	// 生成 token
	token, err := utils.GenerateToken(uint(userResp.UserDetail.Id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "GenerateToken 失败"))
		return
	}
	// 封装 token
	res := &types.TokenData{
		User:  userResp,
		Token: token,
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, res))
	
}

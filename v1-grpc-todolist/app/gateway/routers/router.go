package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"grpc-todolist/app/gateway/internal/handler"
	"grpc-todolist/app/gateway/middleware"
	"net/http"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.ErrorMiddleware())  // 加入跨域处理 handle...
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("session", store))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "success")
		})
		v1.POST("/user/register", handler.UserRegister)
		v1.POST("/user/login", handler.UserLogin)

		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.JWT())  // 引入 JWT 鉴权中间件
		{
			// task module 
			authed.GET("task", handler.GetTaskList())
			authed.POST("task", handler.CreateTask())
			authed.PUT("task", handler.UpdateTask())
			authed.DELETE("task", handler.DeleteTask())
		}
	}
	
	return ginRouter
}

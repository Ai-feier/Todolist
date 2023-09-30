package router

import (
	"github.com/gin-gonic/gin"
	"micro-todolist/app/gateway/http"
	"micro-todolist/app/gateway/middleware"
	"micro-todolist/pkgs/e"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	
	// 中间件
	ginRouter.Use(middleware.Cors())
	
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(e.SUCCESS, "success")
		})
		v1.POST("/user/register", http.UserRegisterHandler)
		v1.POST("/user/login", http.UserLoginHandler)
		
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("tasks", http.ListTaskHandler)
			authed.POST("task", http.CreateTaskHandler)
			authed.GET("task/:id", http.GetTaskHandler)
			authed.PUT("task/:id", http.UpdateTaskHandler)
			authed.DELETE("task/:id", http.DeleteTaskHandler)
			
		}
	}
	return ginRouter
}
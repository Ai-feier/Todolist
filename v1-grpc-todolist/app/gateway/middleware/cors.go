package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method  // 方法
		origin := ctx.Request.Header.Get("Origin")  // 请求头部
		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Origin", "*")   // 这是允许访问所有域
			
			// 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")

			//  header的类型
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")

			// 允许跨域设置               可以返回其他子段
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, " +
				"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires," +
				"Last-Modified,Pragma,FooBar") 				// 跨域关键设置 让浏览器可以解析

			ctx.Header("Access-Control-Max-Age", "172800")                         // 缓存请求信息 单位为秒
			ctx.Header("Access-Control-Allow-Credentials", "false")                 //  跨域请求是否需要带cookie信息 默认设置为true
			ctx.Set("content-type", "application/json")                              // 设置返回格式是json
		}
		
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Option Request!")
		}
		
		// 处理请求, 向后面的 handle 传递 gin.Context
		ctx.Next()
	}
}

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(200, gin.H {
					"code": 404,
					"msg": fmt.Sprintf("%s", r),
				})
				c.Abort()  // 用于终止HTTP请求的处理流程
			}
		}()
		c.Next()
	}
}

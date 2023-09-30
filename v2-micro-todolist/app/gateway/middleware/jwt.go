package middleware

import (
	"github.com/gin-gonic/gin"
	"micro-todolist/pkgs/ctl"
	"micro-todolist/pkgs/utils"
	"net/http"
	"time"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code uint32
		
		code = http.StatusOK
		token := ctx.GetHeader("Authorization")
		if token == "" {
			code = http.StatusNotFound
			ctx.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg": "鉴权失败",
			})
			return
		}
		
		claims, err := utils.ParseToken(token)
		if err != nil {
			code = http.StatusUnauthorized
			ctx.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg": "鉴权失败",
			})
			ctx.Abort()
			return
		}
		
		if time.Now().Unix() > claims.ExpiresAt {
			code = http.StatusUnauthorized
			ctx.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  "权限过期，请重新登陆",
			})
			ctx.Abort()
			return
		}
		ctl.InitUserInfo(ctx.Request.Context(), &ctl.UserInfo{Id: claims.Id})
		ctx.Next()
	}
}

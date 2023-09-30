package middleware

import (
	"github.com/gin-gonic/gin"
	"grpc-todolist/pkg/ctl"
	"grpc-todolist/pkg/e"
	"grpc-todolist/pkg/utils/jwt"
	"time"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		var data interface{}
		code = 200
		token := ctx.GetHeader("Authorization")
		if token == "" {
			code = 404
			ctx.JSON(200, gin.H{
				"status": code,
				"msg": e.GetMsg(code),
				"data": data,
			})
			ctx.Abort()
		}
		// 获得 JWT 的验证链, 进行加权检验
		claims, err := jwt.ParseToken(token)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		}
		
		// 加权验证错误
		if code != e.SUCCESS {
			ctx.JSON(200, gin.H{
				"status": code,
				"msg":e.GetMsg(code),
				"data":data,
			})
			ctx.Abort()
			return
		}
		ctx.Request = ctx.Request.WithContext(ctl.NewContext(ctx.Request.Context(), &ctl.UserInfo{Id: claims.UserID}))
		// 更新当前用户的最大上线时间
		ctl.InitUserInfo(ctx)
		ctx.Next()
	}
}
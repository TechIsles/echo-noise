package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("user_id")

		if userID == nil {
			// 如果只是分页获取首页留言，则不需要鉴权
			if strings.HasPrefix(ctx.Request.URL.Path, "/api/messages/page") {
				ctx.Set("user_id", uint(0))
				ctx.Next()
				return
			}

			// 查看留言详情也不需要鉴权
			if strings.HasPrefix(ctx.Request.URL.Path, "/api/messages/") {
				ctx.Set("user_id", uint(0))
				ctx.Next()
				return
			}

			ctx.JSON(http.StatusUnauthorized, dto.Fail[any]("未登录或登录已过期"))
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("user_id", userID.(uint))
		ctx.Set("username", session.Get("username"))
		ctx.Set("is_admin", session.Get("is_admin"))
		ctx.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		isAdmin := session.Get("is_admin")

		if isAdmin == nil || !isAdmin.(bool) {
			ctx.JSON(http.StatusForbidden, dto.Fail[any]("需要管理员权限"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
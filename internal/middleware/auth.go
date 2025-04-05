package middleware

import (
	"net/http"
    "strings"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/internal/dto"
    "github.com/lin-snow/ech0/internal/models"
    "github.com/lin-snow/ech0/internal/database"
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
// 添加 Token 认证中间件
func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusOK, dto.Fail[any]("未提供认证信息"))
            c.Abort()
            return
        }

        db, err := database.GetDB()
        if err != nil {
            c.JSON(http.StatusOK, dto.Fail[any]("系统错误"))
            c.Abort()
            return
        }

        // 查询用户
        var user models.User
        if err := db.Where("token = ?", token).First(&user).Error; err != nil {
            c.JSON(http.StatusOK, dto.Fail[any]("无效的token"))
            c.Abort()
            return
        }

        // 检查 token 是否为空
        if user.Token == "" {
            c.JSON(http.StatusOK, dto.Fail[any]("token已失效"))
            c.Abort()
            return
        }

        // 设置用户信息到上下文
        c.Set("user_id", user.ID)
        c.Set("username", user.Username)
        c.Set("is_admin", user.IsAdmin)
        c.Next()
    }
}
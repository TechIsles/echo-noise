package pkg

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/models"
)

const (
	UserKey = "user"
)

// 初始化 Session
func InitSession(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret_key"))
	r.Use(sessions.Sessions("ech0_session", store))
}

// 保存用户会话
func SaveUserSession(c *gin.Context, user models.User) error {
	session := sessions.Default(c)
	session.Set(UserKey, user)
	return session.Save()
}

// 获取用户会话
func GetUserSession(c *gin.Context) (models.User, bool) {
	session := sessions.Default(c)
	val := session.Get(UserKey)
	if val == nil {
		return models.User{}, false
	}
	user, ok := val.(models.User)
	return user, ok
}

// 清除用户会话
func ClearUserSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
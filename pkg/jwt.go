package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/lin-snow/ech0/internal/models"
)

// CreateSession 创建会话
func CreateSession(user models.User) *models.UserSession {
	return &models.UserSession{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}
}

// ValidateSession 验证会话
func ValidateSession(session *models.UserSession) bool {
	if session == nil {
		return false
	}
	// 可以添加其他验证逻辑，比如检查会话是否过期等
	return true
}

// MD5Encrypt MD5 加密
func MD5Encrypt(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
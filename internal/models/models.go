package models

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Username  string    `gorm:"type:varchar(100)" json:"username,omitempty"`
	ImageURL  string    `gorm:"type:text" json:"image_url,omitempty"`
	Private   bool      `gorm:"default:false" json:"private"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserStatus struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	UserName string `json:"username"` // 用户名
	IsAdmin  bool   `json:"is_admin"` // 是否是管理员
}

type Status struct {
	SysAdminID    uint         `json:"sys_admin_id"` // 系统管理员ID
	Username      string       `json:"username"`     // 系统管理员用户名
	Users         []UserStatus `json:"users"`        // 所有用户
	TotalMessages int          `json:"total_messages"`
}

type MyCliams struct {
	Userid   uint   `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type Setting struct {
	gorm.Model
	AllowRegistration bool `gorm:"default:true"`
}

type SiteConfig struct {
	gorm.Model
	SiteTitle        string `gorm:"type:varchar(100)"`
	SubtitleText     string `gorm:"type:varchar(200)"`
	AvatarURL        string `gorm:"type:varchar(500)"`
	Username         string `gorm:"type:varchar(50)"`
	Description      string `gorm:"type:varchar(200)"`
	Backgrounds      string `gorm:"type:text"`
	CardFooterTitle  string `gorm:"type:varchar(100)"`
	CardFooterLink   string `gorm:"type:varchar(100)"`
	PageFooterHTML   string `gorm:"type:text"`
	RSSTitle         string `gorm:"type:varchar(100)"`
	RSSDescription   string `gorm:"type:varchar(200)"`
	RSSAuthorName    string `gorm:"type:varchar(50)"`
	RSSFaviconURL    string `gorm:"type:varchar(500)"`
	WalineServerURL  string `gorm:"type:varchar(500)"`
}

// GetBackgroundsList 将逗号分隔的背景图片字符串转换为切片
func (s *SiteConfig) GetBackgroundsList() []string {
	if s.Backgrounds == "" {
		return []string{}
	}
	return strings.Split(s.Backgrounds, ",")
}
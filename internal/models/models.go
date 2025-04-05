package models

import (
	"crypto/rand"
    "encoding/hex"
    "errors"
    "strings"
    "time"
    "gorm.io/gorm"
)

var DB *gorm.DB
type UserStatus struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	Status   Status `json:"status"`
}

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
	Token    string `gorm:"type:varchar(255)" json:"token"`
}
// 生成 Token 的工具函数
func GenerateToken(length int) string {
    b := make([]byte, length/2)
    rand.Read(b)
    return hex.EncodeToString(b)
}
type Status struct {
	SysAdminID    uint         `json:"sys_admin_id"` 
	Username      string       `json:"username"`     
	Users         []UserStatus `json:"users"`        
	TotalMessages int          `json:"total_messages"`
}

type UserSession struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	LoginTime time.Time `json:"login_time"`
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
	Version          int    `json:"version"`
}

func (s *SiteConfig) GetBackgroundsList() []string {
	if s.Backgrounds == "" {
		return []string{}
	}
	return strings.Split(s.Backgrounds, ",")
}
func UpdateMessage(id string, content string) error {
    // 先查询消息是否存在
    var message Message
    result := DB.First(&message, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return errors.New("消息不存在")
        }
        return result.Error
    }

    // 更新消息内容
    result = DB.Model(&message).Updates(map[string]interface{}{
        "content": content,
    })
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("更新失败")
    }
    
    return nil
}
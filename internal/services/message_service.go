package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

// GetAllMessages 封装业务逻辑，获取所有笔记
func GetAllMessages(showPrivate bool) ([]models.Message, error) {
	return repository.GetAllMessages(showPrivate)
}

// GetMessageByID 根据 ID 获取笔记
func GetMessageByID(id uint, showPrivate bool) (*models.Message, error) {
	return repository.GetMessageByID(id, showPrivate)
}

// GetMessagesByPage 分页获取笔记
func GetMessagesByPage(page, pageSize int, showPrivate bool) (dto.PageQueryResult, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询数据库
	var messages []models.Message
	var total int64

	if showPrivate {
		// 如果是管理员，则不需要过滤私密笔记
		database.DB.Model(&models.Message{}).Count(&total)
		database.DB.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&messages)
	} else {
		// 如果不是管理员，则只查询公开的笔记
		database.DB.Model(&models.Message{}).Where("private = ?", false).Count(&total)
		database.DB.Limit(pageSize).Offset(offset).Where("private = ?", false).Order("created_at DESC").Find(&messages)
	}

	// 返回结果
	var PageQueryResult dto.PageQueryResult
	PageQueryResult.Total = total
	PageQueryResult.Items = messages

	return PageQueryResult, nil
}

// CreateMessage 发布一条笔记
// 允许所有注册登陆用户发布信息
func CreateMessage(message *models.Message) error {
    user, err := GetUserByID(message.UserID)
    if err != nil {
        return err
    }

    // 删除管理员权限检查，允许所有登录用户发布信息
    message.Username = user.Username // 获取用户名
    return repository.CreateMessage(message)
}

// DeleteMessage 根据 ID 删除笔记
func DeleteMessage(id uint, userID uint) error {
    // 获取笔记信息
    message, err := repository.GetMessageByID(id, true)
    if err != nil {
        return err
    }

    // 验证是否为笔记作者
    if message.UserID != userID {
        return fmt.Errorf("无权删除他人的笔记")
    }

    return repository.DeleteMessage(id)
}

// DeleteMessageByAdmin 管理员删除笔记（无需验证作者）
func DeleteMessageByAdmin(id uint) error {
    return repository.DeleteMessage(id)
}

func GenerateRSS(c *gin.Context) (string, error) {
	// 获取所有笔记
	showPrivate := false
	messages, err := GetAllMessages(showPrivate)
	if err != nil {
		return "", err
	}
// 获取站点配置
settings, err := GetFrontendConfig()
if err != nil {
	return "", err
}
// 从返回的 map 中正确获取配置
frontendSettings := settings["frontendSettings"].(map[string]interface{})

schema := "http"
if c.Request.TLS != nil {
	schema = "https"
}
host := c.Request.Host

feed := &feeds.Feed{
	Title: frontendSettings["rssTitle"].(string),
	Link: &feeds.Link{
		Href: fmt.Sprintf("%s://%s/", schema, host),
	},
	Image: &feeds.Image{
		Url: frontendSettings["rssFaviconURL"].(string),
	},
	Description: frontendSettings["rssDescription"].(string),
	Author: &feeds.Author{
		Name: frontendSettings["rssAuthorName"].(string),
	},
	Updated: time.Now(),
}

	for _, msg := range messages {
		renderedContent := pkg.MdToHTML([]byte(msg.Content))

		title := msg.Username + " - " + msg.CreatedAt.Format("2006-01-02")

		// 添加图片链接到正文前(scheme://host/api/ImageURL)
		if msg.ImageURL != "" {
			image := fmt.Sprintf("%s://%s/api%s", schema, host, msg.ImageURL)
			renderedContent = append([]byte(fmt.Sprintf("<img src=\"%s\" alt=\"Image\" style=\"max-width:100%%;height:auto;\" />", image)), renderedContent...)
		}

		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/api/messages/%d", schema, host, msg.ID)},
			Description: string(renderedContent),
			Author: &feeds.Author{
				Name: msg.Username,
			},
			Created: msg.CreatedAt,
		}
		feed.Items = append(feed.Items, item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return "", err
	}

	return atom, nil
}
// UpdateMessage 更新消息内容
func UpdateMessage(messageID uint, content string) error {
    // 获取消息
    message, err := repository.GetMessageByID(messageID, true)
    if err != nil {
        return fmt.Errorf("获取消息失败: %v", err)
    }
    if message == nil {
        return fmt.Errorf("消息不存在")
    }

    // 更新消息内容
    message.Content = content
    if err := database.DB.Save(message).Error; err != nil {
        return fmt.Errorf("更新消息失败: %v", err)
    }

    return nil
}
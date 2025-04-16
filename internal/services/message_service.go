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
    messages, err := GetAllMessages(false)
    if err != nil {
        return "", fmt.Errorf("获取消息失败: %v", err)
    }

    // 获取前端配置
    config, err := GetFrontendConfig()
    if err != nil {
        return "", fmt.Errorf("获取配置失败: %v", err)
    }

    // 从配置中安全获取值
    getConfigValue := func(key string, defaultValue string) string {
        if settings, ok := config["frontendSettings"].(map[string]interface{}); ok {
            if value, exists := settings[key].(string); exists && value != "" {
                return value
            }
        }
        return defaultValue
    }

    schema := "http"
    if c.Request.TLS != nil {
        schema = "https"
    }
    
    // 从配置中获取域名，如果没有配置则使用请求的 host
    host := getConfigValue("siteURL", c.Request.Host)
    
    feed := &feeds.Feed{
        Title: getConfigValue("rssTitle", "说说笔记"),
        Link: &feeds.Link{
            Href: fmt.Sprintf("%s://%s/", schema, host),
        },
        Image: &feeds.Image{
            Url: fmt.Sprintf("%s://%s%s", schema, host, getConfigValue("rssFaviconURL", "/favicon.ico")),
        },
        Description: getConfigValue("rssDescription", "一个说说笔记~"),
        Author: &feeds.Author{
            Name: getConfigValue("rssAuthorName", "Noise"),
        },
        Updated: time.Now(),
    }

    for _, msg := range messages {
        // 处理内容
        content := msg.Content
        if msg.ImageURL != "" {
            imageURL := fmt.Sprintf("%s://%s/api%s", schema, host, msg.ImageURL)
            content = fmt.Sprintf("![图片](%s)\n\n%s", imageURL, content)
        }

        // 渲染 Markdown
        htmlContent := pkg.MdToHTML([]byte(content))

        // 生成标题
        title := msg.Username
        if firstLine := pkg.GetFirstLine(msg.Content); firstLine != "" {
            title = firstLine
        }

        // 生成前端页面 URL
        // 使用配置的域名生成 URL
        pageURL := fmt.Sprintf("%s://%s/#/messages/%d", schema, host, msg.ID)
        
        item := &feeds.Item{
            Title:       title,
            Link:        &feeds.Link{Href: pageURL},
            Description: string(htmlContent),
            Author:      &feeds.Author{Name: msg.Username},
            Created:     msg.CreatedAt,
            Id:         pageURL,
        }

        feed.Items = append(feed.Items, item)
    }

    rss, err := feed.ToRss()
    if err != nil {
        return "", err
    }

    return rss, nil
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
func GetMessagesGroupByDate() ([]struct {
    Date  string `json:"date"`
    Count int    `json:"count"`
}, error) {
    var results []struct {
        Date  string `json:"date"`
        Count int    `json:"count"`
    }
    
    // 移除 deleted_at 条件，因为该列不存在
    err := database.DB.Table("messages").
        Select("DATE(created_at) as date, COUNT(*) as count").
        // 移除这一行: Where("deleted_at IS NULL").
        Group("DATE(created_at)").
        Order("date DESC").
        Scan(&results).Error
        
    if err != nil {
        fmt.Printf("获取消息日历数据失败: %v\n", err)
        return nil, err
    }
    
    // 如果结果为空，返回空数组而不是nil
    if len(results) == 0 {
        return []struct {
            Date  string `json:"date"`
            Count int    `json:"count"`
        }{}, nil
    }
    
    return results, nil
}
// GetMessagePage 获取消息详情页
func GetMessagePage(id uint) (*models.Message, error) {
    message, err := repository.GetMessageByID(id, false)
    if err != nil {
        return nil, err
    }
    
    if message == nil {
        return nil, fmt.Errorf("消息不存在")
    }
    
    // 如果是私密消息，需要进行额外处理
    if message.Private {
        return nil, fmt.Errorf("无权访问")
    }
    
    return message, nil
}

func SearchMessages(keyword string, page, pageSize int, showPrivate bool) (dto.PageQueryResult, error) {
    // 参数校验
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }

    // 直接使用服务层实现
    query := database.DB.Model(&models.Message{}).
        Select("id, content, created_at, username, image_url, private, user_id").
        Where("content LIKE ?", "%"+keyword+"%")
    
    if !showPrivate {
        query = query.Where("private = ?", false)
    }

    var total int64
    var messages []models.Message
    
    err := query.Count(&total).
        Limit(pageSize).
        Offset((page-1)*pageSize).
        Order("created_at DESC").
        Find(&messages).Error
        
    if err != nil {
        return dto.PageQueryResult{}, err
    }

    // 确保返回的数据结构符合前端期望
    return dto.PageQueryResult{
        Total: total,
        Items: messages,
    }, nil
}
package controllers

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/lin-snow/ech0/internal/models"
    "github.com/lin-snow/ech0/internal/database" 
    "regexp"
)
// GetMessagesByTag 获取指定标签的消息
// GetMessagesByTag 获取指定标签的消息
func GetMessagesByTag(c *gin.Context) {
    tag := c.Param("tag")
    if tag == "" {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "标签不能为空"})
        return
    }

    db, err := database.GetDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库连接失败"})
        return
    }

    var messages []models.Message
    // 使用 LIKE 进行初步筛选
    tagPattern := "%#" + tag + "%"
    if err := db.Where("content LIKE ?", tagPattern).Order("created_at DESC").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "获取消息失败"})
        return
    }

    // 使用正则表达式进行精确匹配
    var filteredMessages []models.Message
    tagRegex := regexp.MustCompile(`#` + regexp.QuoteMeta(tag) + `(?:[\s,.!?]|$)`)
    for _, msg := range messages {
        if tagRegex.MatchString(msg.Content) {
            filteredMessages = append(filteredMessages, msg)
        }
    }

    c.JSON(http.StatusOK, gin.H{"code": 1, "data": filteredMessages})
}

// GetAllTags 获取所有标签列表
func GetAllTags(c *gin.Context) {
    db, err := database.GetDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库连接失败"})
        return
    }

    var messages []models.Message
    if err := db.Select("content").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "获取消息失败"})
        return
    }

    // 提取并统计标签
    tagMap := make(map[string]int)
    for _, msg := range messages {
        // 使用正则表达式匹配 #标签
        tags := regexp.MustCompile(`#([^\s#]+)`).FindAllStringSubmatch(msg.Content, -1)
        for _, tag := range tags {
            if len(tag) > 1 {
                tagMap[tag[1]]++
            }
        }
    }

    // 转换为数组格式
    var tagList []map[string]interface{}
    for tag, count := range tagMap {
        tagList = append(tagList, map[string]interface{}{
            "name":  tag,
            "count": count,
        })
    }

    c.JSON(http.StatusOK, gin.H{"code": 1, "data": tagList})
}

// GetAllImages 获取所有图片列表.
func GetAllImages(c *gin.Context) {
    db, err := database.GetDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "数据库连接失败"})
        return
    }

    var messages []models.Message
    if err := db.Select("id", "content", "image_url", "created_at").Order("created_at DESC").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": "获取图片列表失败"})
        return
    }

    type ImageInfo struct {
        ID        uint      `json:"id"`
        ImageURL  string    `json:"image_url"`
        CreatedAt time.Time `json:"created_at"`
    }

    var allImages []ImageInfo
    
    // 正则表达式匹配 Markdown 格式的图片
    imageRegex := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

    for _, msg := range messages {
        // 添加 image_url 字段的图片
        if msg.ImageURL != "" {
            allImages = append(allImages, ImageInfo{
                ID:        msg.ID,
                ImageURL:  msg.ImageURL,
                CreatedAt: msg.CreatedAt,
            })
        }

        // 提取内容中的 Markdown 格式图片
        matches := imageRegex.FindAllStringSubmatch(msg.Content, -1)
        for _, match := range matches {
            if len(match) > 1 {
                allImages = append(allImages, ImageInfo{
                    ID:        msg.ID,
                    ImageURL:  match[1],
                    CreatedAt: msg.CreatedAt,
                })
            }
        }
    }

    c.JSON(http.StatusOK, gin.H{"code": 1, "data": allImages})
}


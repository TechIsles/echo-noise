package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
)

func UpdateMessage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	messageID := c.Param("id")
	id, err := strconv.ParseUint(messageID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的消息ID",
			"detail": err.Error(),
		})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
			"detail": err.Error(),
		})
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	// 获取原消息
	oldMessage, err := repository.GetMessageByID(uint(id), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取消息失败",
			"detail": err.Error(),
		})
		return
	}
	if oldMessage == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "消息不存在"})
		return
	}

	// 权限检查
	currentUser, err := repository.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败",
			"detail": err.Error(),
		})
		return
	}

	if oldMessage.UserID != userID.(uint) && !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限编辑此消息"})
		return
	}

	// 删除旧消息
	if err := repository.DeleteMessage(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除原消息失败",
			"detail": err.Error(),
		})
		return
	}

	// 创建新消息
	newMessage := &models.Message{
		UserID:    oldMessage.UserID,
		Username:  oldMessage.Username,
		Content:   content,
		ImageURL:  oldMessage.ImageURL,
		Private:   oldMessage.Private,
		CreatedAt: oldMessage.CreatedAt,
	}

	if err := repository.CreateMessage(newMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建新消息失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    newMessage,
	})
}
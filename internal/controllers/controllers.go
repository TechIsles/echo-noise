package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)
// controllers.go
func checkAdmin(c *gin.Context) (uint, error) {
    // 从上下文中获取 userid
    userID, exists := c.Get("userid")  // 确保你的JWT中间件设置了"userid"
    if !exists {
        return 0, errors.New("未授权访问")
    }

    // 检查用户是否为管理员
    user, err := services.GetUserByID(userID.(uint))
    if err != nil {
        return 0, err
    }

    if !user.IsAdmin {
        return 0, errors.New("需要管理员权限")
    }

    return userID.(uint), nil
}


// checkUser 获取当前用户信息，返回用户对象和错误信息
func checkUser(c *gin.Context) (*models.User, error) {
	user, err := services.GetUserByID(c.MustGet("userid").(uint))
	if err != nil {
		return nil, fmt.Errorf(models.UserNotFoundMessage)
	}
	return user, nil
}

// Login 处理 POST /login 请求，用户登录
func Login(c *gin.Context) {
	var user dto.LoginDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	token, err := services.Login(user)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(token, models.LoginSuccessMessage))
}

func Register(c *gin.Context) {
	var user dto.RegisterDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	if err := services.Register(user); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.RegisterSuccessMessage))
}

func GetMessages(c *gin.Context) {
	messages, err := services.GetAllMessages(false)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetAllMessagesFailMessage))
		return
	}
	c.JSON(http.StatusOK, dto.OK(messages, models.GetAllMessagesSuccess))
}

func GetMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	userID := c.MustGet("userid").(uint)
	showPrivate := false
	if userID != 0 {
		user, err := services.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
			return
		}
		if user.IsAdmin {
			showPrivate = true
		}
	}

	message, err := services.GetMessageByID(uint(id), showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetMessageByIDFailMessage))
		return
	}

	if message == nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.MessageNotFoundMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(message, models.GetMessageByIDSuccess))
}

func GetMessagesByPage(c *gin.Context) {
	var pageRequest dto.PageQueryDto
	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	userID := c.MustGet("userid").(uint)
	showPrivate := false

	if userID != 0 {
		user, err := services.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
			return
		}
		if user.IsAdmin {
			showPrivate = true
		}
	}

	pageQueryResult, err := services.GetMessagesByPage(pageRequest.Page, pageRequest.PageSize, showPrivate)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(pageQueryResult, models.GetMessagesByPageSuccess))
}

func PostMessage(c *gin.Context) {
	var message models.Message

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	message.UserID = c.MustGet("userid").(uint)
	if err := services.CreateMessage(&message); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(message, models.PostMessageSuccess))
}

func GetStatus(c *gin.Context) {
	status, err := services.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.GetStatusFailMessage))
		return
	}

	c.JSON(http.StatusOK, dto.OK(status, models.GetStatusSuccessMessage))
}

// DeleteMessage 删除留言
func DeleteMessage(c *gin.Context) {
    // 获取消息 ID
    id := c.Param("id")
    messageID, err := strconv.ParseUint(id, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的消息ID"})
        return
    }

    // 获取当前用户ID
    userID := c.GetUint("userid")
    
    // 检查是否为管理员
    user, err := services.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
        return
    }

    // 如果是管理员，直接删除；如果不是，则验证是否为留言作者
    if !user.IsAdmin {
        // 调用服务层删除消息，传入消息ID和用户ID进行验证
        if err := services.DeleteMessage(uint(messageID), userID); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
            return
        }
    } else {
        // 管理员可以直接删除任何留言
        if err := services.DeleteMessageByAdmin(uint(messageID)); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "删除成功"})
}

func GenerateRSS(c *gin.Context) {
	atom, err := services.GenerateRSS(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail[string](models.GenerateRSSFailMessage))
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(atom))
}

func UpdateUser(c *gin.Context) {
	user, err := checkUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	var userdto dto.UserInfoDto
	if err := c.ShouldBindJSON(&userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	if err := services.UpdateUser(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

func ChangePassword(c *gin.Context) {
	var userdto dto.UserInfoDto
	if err := c.ShouldBindJSON(&userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	user, err := checkUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	if err := services.ChangePassword(user, userdto); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.ChangePasswordSuccessMessage))
}

func UpdateUserAdmin(c *gin.Context) {
    _, err := checkAdmin(c) // 获取的是 uint 类型的 userID
    if err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
        return
    }
    idStr := c.Query("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil || id == 1 { // 移除了 user.ID 的比较
        c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
        return
    }
    if err := services.UpdateUserAdmin(uint(id)); err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
        return
    }
    c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateUserSuccessMessage))
}

func GetUserInfo(c *gin.Context) {
	user, err := checkUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, dto.OK(user, models.QuerySuccessMessage))
}

func UpdateSetting(c *gin.Context) {
	userID, err := checkAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	var setting dto.SettingDto
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 将 SettingDto 转换为 map[string]interface{}
	settingMap := map[string]interface{}{
		"allow_registration": setting.AllowRegistration,
		"frontend_settings": map[string]interface{}{
			"site_title":    setting.FrontendSettings.SiteTitle,
			"subtitle_text": setting.FrontendSettings.SubtitleText,
			"avatar_url":    setting.FrontendSettings.AvatarURL,
			"username":      setting.FrontendSettings.Username,
			"description":   setting.FrontendSettings.Description,
			"backgrounds":   setting.FrontendSettings.Backgrounds,
		},
	}

	if err := services.UpdateSetting(userID, settingMap); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateSettingSuccessMessage))
}
// ... existing code ...

func GetFrontendConfig(c *gin.Context) {
    config, err := services.GetFrontendConfig()
    if err != nil {
        c.JSON(http.StatusOK, dto.Fail[interface{}](models.QueryFailMessage))
        return
    }

    // 确保返回的结构与前端期望的匹配
    response := map[string]interface{}{
        "code": 1,
        "msg":  models.QuerySuccessMessage,
        "data": map[string]interface{}{
            "frontendSettings": config,
        },
    }

    c.JSON(http.StatusOK, response)
}


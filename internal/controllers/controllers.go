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

func checkAdmin(c *gin.Context) (uint, error) {
    userID, exists := c.Get("userid")
    if !exists {
        return 0, errors.New("未授权访问")
    }

    user, err := services.GetUserByID(userID.(uint))
    if err != nil {
        return 0, err
    }

    if !user.IsAdmin {
        return 0, errors.New("需要管理员权限")
    }

    return userID.(uint), nil
}

func checkUser(c *gin.Context) (*models.User, error) {
    userID, exists := c.Get("userid")
    if !exists {
        return nil, fmt.Errorf(models.UserNotFoundMessage)
    }

    user, err := services.GetUserByID(userID.(uint))
    if err != nil {
        return nil, fmt.Errorf(models.UserNotFoundMessage)
    }
    return user, nil
}

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

func GetMessage(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
        return
    }

    showPrivate := false
    userID, exists := c.Get("userid")
    if exists {
        user, err := services.GetUserByID(userID.(uint))
        if err == nil && user.IsAdmin {
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

    showPrivate := false
    userID, exists := c.Get("userid")
    if exists {
        user, err := services.GetUserByID(userID.(uint))
        if err == nil && user.IsAdmin {
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

    userID, exists := c.Get("userid")
    if !exists {
        c.JSON(http.StatusOK, dto.Fail[string]("未授权访问"))
        return
    }
    message.UserID = userID.(uint)

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

func DeleteMessage(c *gin.Context) {
    id := c.Param("id")
    messageID, err := strconv.ParseUint(id, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "无效的消息ID"})
        return
    }

    userID, exists := c.Get("userid")
    if !exists {
        c.JSON(http.StatusOK, dto.Fail[string]("未授权访问"))
        return
    }

    user, err := services.GetUserByID(userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
        return
    }

    if !user.IsAdmin {
        if err := services.DeleteMessage(uint(messageID), userID.(uint)); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "msg": err.Error()})
            return
        }
    } else {
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
    _, err := checkAdmin(c)
    if err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
        return
    }

    idStr := c.Query("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil || id == 1 {
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

    settingMap := map[string]interface{}{
        "allow_registration": setting.AllowRegistration,
        "frontend_settings": map[string]interface{}{
            "siteTitle":          setting.FrontendSettings.SiteTitle,
            "subtitleText":       setting.FrontendSettings.SubtitleText,
            "avatarURL":          setting.FrontendSettings.AvatarURL,
            "username":           setting.FrontendSettings.Username,
            "description":        setting.FrontendSettings.Description,
            "backgrounds":        setting.FrontendSettings.Backgrounds,
            "cardFooterTitle":    setting.FrontendSettings.CardFooterTitle,
            "cardFooterLink":     setting.FrontendSettings.CardFooterLink,
            "pageFooterHTML":     setting.FrontendSettings.PageFooterHTML,
            "rssTitle":          setting.FrontendSettings.RSSTitle,
            "rssDescription":    setting.FrontendSettings.RSSDescription,
            "rssAuthorName":     setting.FrontendSettings.RSSAuthorName,
            "rssFaviconURL":     setting.FrontendSettings.RSSFaviconURL,
            "walineServerURL":   setting.FrontendSettings.WalineServerURL,
        },
    }

    if err := services.UpdateSetting(userID, settingMap); err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
        return
    }

    c.JSON(http.StatusOK, dto.OK[any](nil, models.UpdateSettingSuccessMessage))
}

func GetFrontendConfig(c *gin.Context) {
    config, err := services.GetFrontendConfig()
    if err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](models.QueryFailMessage))
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": 1,
        "msg":  models.QuerySuccessMessage,
        "data": config,
    })
}

func GetMessages(c *gin.Context) {
    showPrivate := false
    userID, exists := c.Get("userid")
    if exists {
        user, err := services.GetUserByID(userID.(uint))
        if err == nil && user.IsAdmin {
            showPrivate = true
        }
    }

    messages, err := services.GetAllMessages(showPrivate)
    if err != nil {
        c.JSON(http.StatusOK, dto.Fail[string](models.GetAllMessagesFailMessage))
        return
    }

    c.JSON(http.StatusOK, dto.OK(messages, models.GetAllMessagesSuccess))
}
package models

import (
    "bytes"
    "log"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "html"      // 添加 html 包
    "io"
    "mime/multipart"  // 添加这一行
    "net/http"
    "strings"
    "time"
    "gorm.io/gorm"
    "regexp"
)

type NotifyConfig struct {
    gorm.Model
    WebhookEnabled  bool   `json:"webhookEnabled"`
    WebhookURL      string `gorm:"type:varchar(255)" json:"webhookURL"`
    TelegramEnabled bool   `json:"telegramEnabled"`
    TelegramToken   string `gorm:"type:varchar(255)" json:"telegramToken"`
    TelegramChatID  string `gorm:"type:varchar(100)" json:"telegramChatID"`
    WeworkEnabled   bool   `json:"weworkEnabled"`
    WeworkKey       string `gorm:"type:varchar(255)" json:"weworkKey"`
    FeishuEnabled   bool   `json:"feishuEnabled"`
    FeishuWebhook   string `gorm:"type:varchar(255)" json:"feishuWebhook"`
    FeishuSecret    string `gorm:"type:varchar(255)" json:"feishuSecret"`
}

// 获取配置
// 确保 GetNotifyConfig 返回最新配置
func GetNotifyConfig() *NotifyConfig {
    db := GetDB()
    var config NotifyConfig
    if err := db.First(&config).Error; err != nil {
        return &NotifyConfig{} // 确保返回空配置而不是nil
    }
    return &config
}


// 保存配置
func SaveNotifyConfig(config NotifyConfig) error {
    db := GetDB() // 使用 models 包中的 GetDB
    if db == nil {
        return errors.New("数据库连接未初始化")
    }

    // 验证配置有效性
    if err := validateNotifyConfig(config); err != nil {
        return err
    }

    var existingConfig NotifyConfig
    result := db.First(&existingConfig)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            // 创建新配置时，如果配置了对应渠道，则自动启用
            if config.WebhookURL != "" {
                config.WebhookEnabled = true
            }
            if config.TelegramToken != "" && config.TelegramChatID != "" {
                config.TelegramEnabled = true
            }
            if config.WeworkKey != "" {
                config.WeworkEnabled = true
            }
            if config.FeishuWebhook != "" {
                config.FeishuEnabled = true
            }
            return db.Create(&config).Error
        }
        return result.Error
    }
    
    // 更新现有配置时，如果配置了对应渠道，则自动启用
    if config.WebhookURL != "" {
        config.WebhookEnabled = true
    }
    if config.TelegramToken != "" && config.TelegramChatID != "" {
        config.TelegramEnabled = true
    }
    if config.WeworkKey != "" {
        config.WeworkEnabled = true
    }
    if config.FeishuWebhook != "" {
        config.FeishuEnabled = true
    }
    
    // 更新所有字段
    existingConfig.WebhookEnabled = config.WebhookEnabled
    existingConfig.WebhookURL = config.WebhookURL
    existingConfig.TelegramEnabled = config.TelegramEnabled
    existingConfig.TelegramToken = config.TelegramToken
    existingConfig.TelegramChatID = config.TelegramChatID
    existingConfig.WeworkEnabled = config.WeworkEnabled
    existingConfig.WeworkKey = config.WeworkKey
    existingConfig.FeishuEnabled = config.FeishuEnabled
    existingConfig.FeishuWebhook = config.FeishuWebhook
    existingConfig.FeishuSecret = config.FeishuSecret
    
    return db.Save(&existingConfig).Error
}
// 验证配置
func validateNotifyConfig(config NotifyConfig) error {
    if config.WebhookEnabled {
        if !strings.HasPrefix(config.WebhookURL, "http") {
            return errors.New("Webhook URL必须以http/https开头")
        }
    }
    if config.TelegramEnabled {
        if config.TelegramToken == "" {
            return errors.New("Telegram Token不能为空")
        }
        if config.TelegramChatID == "" {
            return errors.New("Telegram Chat ID不能为空")
        }
    }
    if config.WeworkEnabled && config.WeworkKey == "" {
        return errors.New("企业微信Key不能为空")
    }
    if config.FeishuEnabled {
        if config.FeishuWebhook == "" {
            return errors.New("飞书Webhook不能为空")
        }
    }
    return nil
}
// 发送Webhook消息
func SendWebhook(message string) error {
    config := GetNotifyConfig()
    if !config.WebhookEnabled || config.WebhookURL == "" {
        log.Printf("Webhook未启用或URL为空")
        return nil
    }
    
    payload := map[string]interface{}{
        "text": message,
        "markdown": true,  // 启用markdown支持
    }
    jsonData, _ := json.Marshal(payload)
    
    resp, err := http.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("webhook请求错误: %v", err)
    }
    defer resp.Body.Close()
    
    // 读取响应内容
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return fmt.Errorf("webhook响应解析失败: %v", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("webhook请求失败: %d, %v", resp.StatusCode, response)
    }
    
    return nil
}

// 发送Telegram消息
// SendTelegram 发送 Telegram 消息
func SendTelegram(content string, images []string) error {
    config := GetNotifyConfig()
    if config.TelegramToken == "" || config.TelegramChatID == "" {
        return nil
    }

    // 处理 Markdown 内容
    messageText := content
    
    // 如果有图片，添加到消息末尾
    if len(images) > 0 {
        messageText += "\n\n"
        for _, img := range images {
            // 使用 HTML 格式的图片标签
            messageText += fmt.Sprintf("<a href=\"%s\">&#8205;</a>", img)
        }
    }

    url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
    payload := map[string]interface{}{
        "chat_id": config.TelegramChatID,
        "text":    messageText,
        "parse_mode": "HTML",
        "disable_web_page_preview": false,  // 启用网页预览以显示图片
    }

    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("telegram请求错误: %v", err)
    }
    defer resp.Body.Close()

    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return fmt.Errorf("telegram响应解析失败: %v", err)
    }

    if !response["ok"].(bool) {
        return fmt.Errorf("telegram发送失败: %v", response["description"])
    }

    return nil
}

// SendWework 发送企业微信消息
func SendWework(content string, images []string) error {
    config := GetNotifyConfig()
    if config.WeworkKey == "" {
        return nil
    }

    if !strings.HasPrefix(config.WeworkKey, "https://") {
        config.WeworkKey = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + config.WeworkKey
    }

    // 企业微信的 markdown 格式
    messageText := content
    if len(images) > 0 {
        messageText += "\n"
        for _, img := range images {
            messageText += fmt.Sprintf("![image](%s)\n", img)
        }
    }

    payload := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]string{
            "content": messageText,
        },
    }

    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(config.WeworkKey, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("企业微信请求错误: %v", err)
    }
    defer resp.Body.Close()

    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return fmt.Errorf("企业微信响应解析失败: %v", err)
    }

    if response["errcode"].(float64) != 0 {
        return fmt.Errorf("企业微信发送失败: %v", response["errmsg"])
    }

    return nil
}

// 发送飞书消息
func SendFeishu(message string) error {
    config := GetNotifyConfig()
    if config.FeishuWebhook == "" {
        return nil
    }
    
    // 飞书使用富文本格式
    payload := map[string]interface{}{
        "msg_type": "interactive",
        "card": map[string]interface{}{
            "config": map[string]interface{}{
                "wide_screen_mode": true,
            },
            "elements": []map[string]interface{}{
                {
                    "tag": "markdown",
                    "content": message,
                },
            },
        },
    }
    
    jsonData, _ := json.Marshal(payload)
    client := &http.Client{}
    req, err := http.NewRequest("POST", config.FeishuWebhook, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("飞书请求创建失败: %v", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    if config.FeishuSecret != "" {
        timestamp := time.Now().Unix()
        sign := genFeishuSign(timestamp, config.FeishuSecret)
        req.Header.Set("X-Lark-Signature", sign)
        req.Header.Set("X-Lark-Request-Timestamp", fmt.Sprintf("%d", timestamp))
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("飞书请求错误: %v", err)
    }
    defer resp.Body.Close()
    
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return fmt.Errorf("飞书响应解析失败: %v", err)
    }
    
    if code, ok := response["code"].(float64); !ok || code != 0 {
        return fmt.Errorf("飞书请求失败: %v", response["msg"])
    }
    
    return nil
}
// 生成飞书签名
func genFeishuSign(timestamp int64, secret string) string {
    stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(stringToSign))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 测试推送函数 - 保持测试推送的特殊格式
// TestNotify 函数修改
func TestNotify(notifyType string) error {
    testMessage := fmt.Sprintf("Echo Noise 推送测试\n\n"+
        "类型: %s\n"+
        "时间: %s\n"+
        "这是一条测试消息，用于验证推送配置是否正确",
        notifyType,
        time.Now().Format("2006-01-02 15:04:05"))

    // 测试消息不包含图片
    var emptyImages []string

    var err error
    switch notifyType {
    case "webhook":
        err = SendWebhook(testMessage)
    case "telegram":
        err = SendTelegram(testMessage, emptyImages)
    case "wework":
        err = SendWework(testMessage, emptyImages)
    case "feishu":
        err = SendFeishu(testMessage)
    default:
        return fmt.Errorf("不支持的推送类型: %s", notifyType)
    }

    return err
}

// SendNotify 函数修改
func SendNotify(content string, images []string, config NotifyConfig) error {
    var errors []string

    // 保持原始 Markdown 格式
    if config.WebhookEnabled && config.WebhookURL != "" {
        if err := SendWebhook(content); err != nil {
            errors = append(errors, fmt.Sprintf("webhook: %v", err))
        }
    }
    
    if config.TelegramEnabled && config.TelegramToken != "" && config.TelegramChatID != "" {
        if err := SendTelegram(content, images); err != nil {
            errors = append(errors, fmt.Sprintf("telegram: %v", err))
        }
    }
    
    if config.WeworkEnabled && config.WeworkKey != "" {
        if err := SendWework(content, images); err != nil {
            errors = append(errors, fmt.Sprintf("wework: %v", err))
        }
    }
    
    if config.FeishuEnabled && config.FeishuWebhook != "" {
        // 飞书需要特殊处理，确保图片正确显示
        fullContent := content
        if len(images) > 0 {
            fullContent += "\n\n"
            for _, img := range images {
                fullContent += fmt.Sprintf("![image](%s)\n", img)
            }
        }
        if err := SendFeishu(fullContent); err != nil {
            errors = append(errors, fmt.Sprintf("feishu: %v", err))
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("推送失败: %s", strings.Join(errors, "; "))
    }

    return nil
}

// SendTelegramWithFormat 发送支持 HTML 格式的 Telegram 消息
func SendTelegramWithFormat(content string, images []string, parseHTML bool) error {
    config := GetNotifyConfig()  // 使用 config 替代 notifyConfig
    if !config.TelegramEnabled {
        return nil
    }

    // 处理 Markdown 链接格式 [text](url)
    content = markdownLinkRegex.ReplaceAllStringFunc(content, func(match string) string {
        parts := markdownLinkRegex.FindStringSubmatch(match)
        if len(parts) == 3 {
            return fmt.Sprintf("<a href=\"%s\">%s</a>", parts[2], html.EscapeString(parts[1]))
        }
        return match
    })

    // 处理 Markdown 图片格式 ![alt](url)
    content = markdownImageRegex.ReplaceAllStringFunc(content, func(match string) string {
        parts := markdownImageRegex.FindStringSubmatch(match)
        if len(parts) == 3 {
            return fmt.Sprintf("<a href=\"%s\">&#8205;</a>", parts[2])
        }
        return match
    })

    apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
    payload := map[string]interface{}{
        "chat_id":    config.TelegramChatID,
        "text":       content,
        "parse_mode": "HTML",
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 如果有额外的图片，单独发送
    for _, img := range images {
        if img == "" {
            continue
        }
        
        photoUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", config.TelegramToken)
        photoPayload := map[string]interface{}{
            "chat_id": config.TelegramChatID,
            "photo":   img,
        }

        jsonData, err = json.Marshal(photoPayload)
        if err != nil {
            return err
        }

        resp, err = http.Post(photoUrl, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            return err
        }
        resp.Body.Close()
    }

    return nil
}

var (
    markdownLinkRegex  = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
    markdownImageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

// SendTelegramPhotoWithCaption 发送图片和文本作为一条消息
func SendTelegramPhotoWithCaption(photoURL string, caption string) error {
    config := GetNotifyConfig()
    if config == nil || !config.TelegramEnabled {
        return fmt.Errorf("Telegram 推送未启用")
    }

    // 构建请求URL
    apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", config.TelegramToken)

    // 检查URL是否为本地URL
    isLocalURL := strings.HasPrefix(photoURL, "http://localhost") || 
                 strings.HasPrefix(photoURL, "https://localhost") ||
                 strings.Contains(photoURL, "127.0.0.1") ||
                 (!strings.HasPrefix(photoURL, "http://") && !strings.HasPrefix(photoURL, "https://"))

    var resp *http.Response
    var err error

    if isLocalURL {
        // 对于本地URL，先下载图片然后作为文件上传
        log.Printf("检测到本地图片URL: %s，尝试下载后上传", photoURL)
        
        // 创建一个multipart表单
        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        
        // 添加chat_id字段
        _ = writer.WriteField("chat_id", config.TelegramChatID)
        
        // 添加caption字段
        _ = writer.WriteField("caption", caption)
        
        // 添加parse_mode字段
        _ = writer.WriteField("parse_mode", "Markdown")
        
        // 尝试下载图片
        var imgResp *http.Response
        if strings.HasPrefix(photoURL, "http") {
            imgResp, err = http.Get(photoURL)
        } else {
            // 如果是相对路径，尝试从本地文件系统读取
            return fmt.Errorf("无法处理相对路径图片: %s", photoURL)
        }
        
        if err != nil {
            return fmt.Errorf("下载图片失败: %v", err)
        }
        defer imgResp.Body.Close()
        
        // 创建photo部分
        part, err := writer.CreateFormFile("photo", "image.jpg")
        if err != nil {
            return fmt.Errorf("创建表单文件失败: %v", err)
        }
        
        // 复制图片数据
        _, err = io.Copy(part, imgResp.Body)
        if err != nil {
            return fmt.Errorf("复制图片数据失败: %v", err)
        }
        
        // 完成multipart表单
        err = writer.Close()
        if err != nil {
            return fmt.Errorf("关闭表单失败: %v", err)
        }
        
        // 发送请求
        req, err := http.NewRequest("POST", apiURL, body)
        if err != nil {
            return fmt.Errorf("创建请求失败: %v", err)
        }
        req.Header.Set("Content-Type", writer.FormDataContentType())
        
        client := &http.Client{}
        resp, err = client.Do(req)
    } else {
        // 对于公网可访问的URL，直接使用URL
        // 构建请求体
        data := map[string]interface{}{
            "chat_id": config.TelegramChatID,
            "photo":   photoURL,
            "caption": caption,
            "parse_mode": "Markdown",
        }

        jsonData, err := json.Marshal(data)
        if err != nil {
            return err
        }

        // 发送请求
        resp, err = http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
    }

    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 检查响应
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("Telegram API 错误: %s", string(body))
    }

    return nil
}

// SendTelegramVideoWithCaption 发送视频和文本作为一条消息
func SendTelegramVideoWithCaption(videoURL string, caption string) error {
    config := GetNotifyConfig()
    if config == nil || !config.TelegramEnabled {
        return fmt.Errorf("Telegram 推送未启用")
    }

    // 构建请求URL
    apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendVideo", config.TelegramToken)

    // 检查URL是否为本地URL
    isLocalURL := strings.HasPrefix(videoURL, "http://localhost") || 
                 strings.HasPrefix(videoURL, "https://localhost") ||
                 strings.Contains(videoURL, "127.0.0.1") ||
                 (!strings.HasPrefix(videoURL, "http://") && !strings.HasPrefix(videoURL, "https://"))

    // 如果是本地URL，改为发送消息并附带链接
    if isLocalURL {
        log.Printf("检测到本地视频URL: %s，改为发送消息并附带链接", videoURL)
        
        // 构建消息内容
        messageText := caption + "\n\n视频链接: " + videoURL
        
        // 使用sendMessage API
        messageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
        messageData := map[string]interface{}{
            "chat_id": config.TelegramChatID,
            "text":    messageText,
            "parse_mode": "Markdown",
        }
        
        jsonData, err := json.Marshal(messageData)
        if err != nil {
            return err
        }
        
        resp, err := http.Post(messageURL, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        
        // 检查响应
        if resp.StatusCode != http.StatusOK {
            body, _ := io.ReadAll(resp.Body)
            return fmt.Errorf("Telegram API 错误: %s", string(body))
        }
        
        return nil
    }
    
    // 对于公网可访问的URL，直接使用URL
    // 构建请求体
    data := map[string]interface{}{
        "chat_id": config.TelegramChatID,
        "video":   videoURL,
        "caption": caption,
        "parse_mode": "Markdown",
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    // 发送请求
    resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 检查响应
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("Telegram API 错误: %s", string(body))
    }

    return nil
}

func SendTelegramMessage(message string) error {
    config := GetNotifyConfig()
    if config == nil || !config.TelegramEnabled {
        return fmt.Errorf("Telegram 推送未启用")
    }

    // 构建请求URL
    apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)

    // 构建请求体
    data := map[string]interface{}{
        "chat_id": config.TelegramChatID,
        "text":    message,
        "parse_mode": "Markdown",
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    // 发送请求
    resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 检查响应
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("Telegram API 错误: %s", string(body))
    }

    return nil
}


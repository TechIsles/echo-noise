package services

import (
    "fmt"
    "encoding/json"
    "github.com/lin-snow/ech0/internal/models"  
    "github.com/lin-snow/ech0/internal/database" 
)

// GetFrontendConfig 获取前端配置
func GetFrontendConfig() (map[string]interface{}, error) {
    db, err := database.GetDB()
    if err != nil {
        return getDefaultConfig(), nil
    }
    
    var config models.SiteConfig
    if err := db.Table("site_configs").First(&config).Error; err != nil {
        return getDefaultConfig(), nil
    }
    
    // 转换为前端所需的格式
    configMap := map[string]interface{}{
        "allowRegistration": true,
        "frontendSettings": map[string]interface{}{
            "siteTitle":          config.SiteTitle,
            "subtitleText":       config.SubtitleText,
            "avatarURL":          config.AvatarURL,
            "username":           config.Username,
            "description":        config.Description,
            "backgrounds":        config.GetBackgroundsList(),
            "cardFooterTitle":    config.CardFooterTitle,
            "cardFooterLink":     config.CardFooterLink,
            "pageFooterHTML":     config.PageFooterHTML,
            "rssTitle":          config.RSSTitle,
            "rssDescription":    config.RSSDescription,
            "rssAuthorName":     config.RSSAuthorName,
            "rssFaviconURL":     config.RSSFaviconURL,
            "walineServerURL":   config.WalineServerURL,
        },
    }
    
    return configMap, nil
}

// UpdateSetting 更新站点配置
func UpdateFrontendSetting(userID uint, settingMap map[string]interface{}) error {
    db, err := database.GetDB()
    if err != nil {
        return fmt.Errorf("数据库连接失败: %v", err)
    }
    
    frontendSettings, ok := settingMap["frontendSettings"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("无效的前端配置格式")
    }
    
    // 类型检查
    for key, value := range frontendSettings {
        if value == nil {
            return fmt.Errorf("配置项 %s 不能为空", key)
        }
        if _, ok := value.(string); !ok && key != "backgrounds" {
            return fmt.Errorf("配置项 %s 必须为字符串类型", key)
        }
    }

    config := models.SiteConfig{
        SiteTitle:       frontendSettings["siteTitle"].(string),
        SubtitleText:    frontendSettings["subtitleText"].(string),
        AvatarURL:       frontendSettings["avatarURL"].(string),
        Username:        frontendSettings["username"].(string),
        Description:     frontendSettings["description"].(string),
        CardFooterTitle: frontendSettings["cardFooterTitle"].(string),
        CardFooterLink:  frontendSettings["cardFooterLink"].(string),
        PageFooterHTML:  frontendSettings["pageFooterHTML"].(string),
        RSSTitle:        frontendSettings["rssTitle"].(string),
        RSSDescription:  frontendSettings["rssDescription"].(string),
        RSSAuthorName:   frontendSettings["rssAuthorName"].(string),
        RSSFaviconURL:   frontendSettings["rssFaviconURL"].(string),
        WalineServerURL: frontendSettings["walineServerURL"].(string),
    }

    // 处理背景图片列表
    if backgrounds, ok := frontendSettings["backgrounds"].([]interface{}); ok {
        backgroundsList := make([]string, 0, len(backgrounds))
        for _, bg := range backgrounds {
            if bgStr, ok := bg.(string); ok {
                backgroundsList = append(backgroundsList, bgStr)
            }
        }
        backgroundsJSON, err := json.Marshal(backgroundsList)
        if err != nil {
            return fmt.Errorf("背景图片列表序列化失败: %v", err)
        }
        config.Backgrounds = string(backgroundsJSON)
    }

    tx := db.Begin()
    if err := tx.Error; err != nil {
        return fmt.Errorf("开启事务失败: %v", err)
    }

    result := tx.Table("site_configs").Where("id = ?", 1).Updates(&config)
    if result.Error != nil {
        tx.Rollback()
        return fmt.Errorf("更新配置失败: %v", result.Error)
    }

    if result.RowsAffected == 0 {
        config.ID = 1
        if err := tx.Table("site_configs").Create(&config).Error; err != nil {
            tx.Rollback()
            return fmt.Errorf("创建配置失败: %v", err)
        }
    }

    return tx.Commit().Error
}

// 获取默认配置
func getDefaultConfig() map[string]interface{} {
    return map[string]interface{}{
        "allowRegistration": true,
        "frontendSettings": map[string]interface{}{
            "siteTitle":          "Noise的说说笔记",
            "subtitleText":       "欢迎访问，点击头像可更换封面背景！",
            "avatarURL":          "https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png",
            "username":           "Noise",
            "description":        "执迷不悟",
            "backgrounds":        []string{
                "https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg",
                "https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg",
                "https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg",
                "https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg",
                "https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg",
                "https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg",
                "https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg",
                "https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg",
                "https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg",
                "https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg",
                "https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg",
                "https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg",
                "https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg",
            },
            "cardFooterTitle":    "Noise·说说·笔记~",
            "cardFooterLink":     "note.noisework.cn",
            "pageFooterHTML":     `<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0</a>发布</div>`,
            "rssTitle":          "Noise的说说笔记",
            "rssDescription":    "一个说说笔记~",
            "rssAuthorName":     "Noise",
            "rssFaviconURL":     "/favicon.ico",
            "walineServerURL":   "https://app-production-80c1.up.railway.app",
        },
    }
}
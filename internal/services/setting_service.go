package services

import (
    "github.com/lin-snow/ech0/internal/models"  
    "github.com/lin-snow/ech0/internal/database" 
)
// GetFrontendConfig 获取前端配置
func GetFrontendConfig() (map[string]interface{}, error) {
    db := database.GetDB()
    var config models.SiteConfig
    
    // 从数据库获取配置
    if err := db.Table("site_configs").First(&config).Error; err != nil {
        // 如果没有配置，返回默认配置
        return getDefaultConfig(), nil
    }
    
    // 转换为前端所需的格式
    configMap := map[string]interface{}{
        "allowRegistration": true, // 这里可能需要从数据库获取实际值
        "frontendSettings": map[string]interface{}{
            "siteTitle":    config.SiteTitle,
            "subtitleText": config.SubtitleText,
            "avatarURL":    config.AvatarURL,
            "username":     config.Username,
            "description":  config.Description,
            "backgrounds":  config.GetBackgroundsList(),
        },
    }
    
    return configMap, nil
}

// 获取默认配置
func getDefaultConfig() map[string]interface{} {
    return map[string]interface{}{
        "allowRegistration": true,
        "frontendSettings": map[string]interface{}{
            "siteTitle":    "Noise的说说笔记",
            "subtitleText": "欢迎访问，点击头像可更换封面背景！",
            "avatarURL":    "https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png",
            "username":     "Noise",
            "description":  "执迷不悟",
            "backgrounds": []string{
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
        },
    }
}
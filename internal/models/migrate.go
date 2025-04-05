package models

import (
    "gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
    // 自动迁移所有模型
    err := db.AutoMigrate(&User{}, &Message{}, &Setting{}, &SiteConfig{})
    if err != nil {
        return err
    }
    // 为现有用户添加 Token 字段
    var users []User
    if err := db.Find(&users).Error; err != nil {
        return err
    }
    for _, user := range users {
        if user.Token == "" {
            newToken := GenerateToken(32)
            // 使用正确的更新方式
            if err := db.Model(&User{}).Where("id = ?", user.ID).Update("token", newToken).Error; err != nil {
                return err
            }
        }
    }
    // 初始化系统设置
    var setting Setting
    if err := db.First(&setting).Error; err != nil {
        defaultSetting := Setting{
            AllowRegistration: true,
        }
        if err := db.Create(&defaultSetting).Error; err != nil {
            return err
        }
    }

    // 初始化管理员账户
    var adminUser User
    if err := db.Where("is_admin = ?", true).First(&adminUser).Error; err != nil {
        defaultAdmin := User{
            Username: "admin",
            Password: HashPassword("admin"),
            IsAdmin:  true,
        }
        if err := db.Create(&defaultAdmin).Error; err != nil {
            return err
        }
    }

    // 初始化前端配置
    var config SiteConfig
    defaultBg := `https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg,
https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg,
https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg,
https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg,
https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg,
https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg,
https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg,
https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg,
https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg,
https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg,
https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg,
https://s2.loli.net/2025/03/27/7Zck3y6XTzhYPs5.jpg,
https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg`

    defaultConfig := SiteConfig{
        SiteTitle:        "Noise的说说笔记",
        SubtitleText:     "欢迎访问，点击头像可更换封面背景！",
        AvatarURL:        "https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png",
        Username:         "Noise",
        Description:      "执迷不悟",
        Backgrounds:      defaultBg,
        CardFooterTitle:  "Noise·说说·笔记~",
        CardFooterLink:   "note.noisework.cn",
        PageFooterHTML:   `<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0</a>发布</div>`,
        RSSTitle:         "Noise的说说笔记",
        RSSDescription:   "一个说说笔记~",
        RSSAuthorName:    "Noise",
        RSSFaviconURL:    "/favicon.ico",
        WalineServerURL:  "https://app-production-80c1.up.railway.app", 
    }

    // 检查是否存在配置
    if err := db.First(&config).Error; err != nil {
        // 如果不存在则创建新配置
        if err := db.Create(&defaultConfig).Error; err != nil {
            return err
        }
    } else {
        // 如果存在则更新配置
        updates := map[string]interface{}{
            "site_title":         defaultConfig.SiteTitle,
            "subtitle_text":      defaultConfig.SubtitleText,
            "avatar_url":         defaultConfig.AvatarURL,
            "username":           defaultConfig.Username,
            "description":        defaultConfig.Description,
            "backgrounds":        defaultConfig.Backgrounds,
            "card_footer_title":  defaultConfig.CardFooterTitle,
            "card_footer_link":   defaultConfig.CardFooterLink,
            "page_footer_html":   defaultConfig.PageFooterHTML,
            "rss_title":          defaultConfig.RSSTitle,
            "rss_description":    defaultConfig.RSSDescription,
            "rss_author_name":    defaultConfig.RSSAuthorName,
            "rss_favicon_url":    defaultConfig.RSSFaviconURL,
            "waline_server_url":  defaultConfig.WalineServerURL,
        }
        if err := db.Model(&config).Updates(updates).Error; err != nil {
            return err
        }
    }

    return nil
}
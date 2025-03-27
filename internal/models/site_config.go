package models

import (
    "strings"
)

type SiteConfig struct {
    ID            uint     `gorm:"primarykey" json:"id"`
    SiteTitle     string   `json:"siteTitle"`
    SubtitleText  string   `json:"subtitleText"`
    AvatarURL     string   `json:"avatarURL"`
    Username      string   `json:"username"`
    Description   string   `json:"description"`
    Backgrounds   string   `json:"backgrounds" gorm:"type:text"`
}

func (s *SiteConfig) GetBackgroundsList() []string {
    if s.Backgrounds == "" {
        return []string{}
    }
    return strings.Split(s.Backgrounds, ",")
}
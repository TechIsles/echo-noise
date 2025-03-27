// setting.go
package dto

type SettingDto struct {
    AllowRegistration bool     `json:"allowRegistration"`
    FrontendSettings  struct {
        SiteTitle    string   `json:"siteTitle"`
        SubtitleText string   `json:"subtitleText"`
        AvatarURL    string   `json:"avatarURL"`
        Username     string   `json:"username"`
        Description  string   `json:"description"`
        Backgrounds  []string `json:"backgrounds"`
    } `json:"frontendSettings"`
}

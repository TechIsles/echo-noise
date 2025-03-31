package dto

type UserInfoDto struct {
    Username     string `json:"username"`
    Password     string `json:"password"`
    SiteName     string `json:"siteName"`
    Theme        string `json:"theme"`
    AllowSignUp  bool   `json:"allowSignUp"`
    IsAdmin      bool   `json:"is_admin"`
}

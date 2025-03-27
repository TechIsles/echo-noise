package dto

type UserInfoDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	SiteName     string `json:"siteName"`
    Theme        string `json:"theme"`
    AllowSignUp  bool   `json:"allowSignUp"`
}

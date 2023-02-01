package models

type UserAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type UserLogoutRes struct {
	Success bool `json:"success"`
}

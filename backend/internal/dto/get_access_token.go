package dto

type GetAccessTokenRequest struct {
	Id   string `json:"id"`
	Url  string `json:"url"`
	Code string `json:"code"`
}

type GetAccessTokenResponse struct {
	Token string `json:"token"`
}

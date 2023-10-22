package models

type ApiRequest struct {
	Url string `json:"url"`
}

type ApiResponse struct {
	Result string `json:"result"`
}

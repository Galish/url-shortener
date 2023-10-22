package models

type APIRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

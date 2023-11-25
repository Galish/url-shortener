package models

type ShortLink struct {
	ID        string `json:"uuid"`
	Short     string `json:"short_url"`
	Original  string `json:"original_url"`
	User      string `json:"user_id,omitempty"`
	IsDeleted bool   `json:"is_deleted"`
}

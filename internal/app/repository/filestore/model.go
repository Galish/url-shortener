package filestore

type link struct {
	ID       string `json:"uuid"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

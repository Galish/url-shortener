package restapi

// APIRequest represents API request payload.
type APIRequest struct {
	URL string `json:"url"`
}

// APIResponse represents API response payload.
type APIResponse struct {
	Result string `json:"result"`
}

// APIBatchEntity represents an API batch request entity.
type APIBatchEntity struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

// APIStatsResponse represents an API statistics response payload.
type APIStatsResponse struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

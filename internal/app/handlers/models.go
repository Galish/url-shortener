package handlers

type apiRequest struct {
	URL string `json:"url"`
}

type apiResponse struct {
	Result string `json:"result"`
}

type apiBatchEntity struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

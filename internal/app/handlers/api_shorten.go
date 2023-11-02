package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
)

type apiRequest struct {
	URL string `json:"url"`
}

type apiResponse struct {
	Result string `json:"result"`
}

func (h *httpHandler) apiShorten(w http.ResponseWriter, r *http.Request) {
	var req apiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	if req.URL == "" {
		http.Error(w, "link not provided", http.StatusBadRequest)
		return
	}

	id := h.generateUniqueID(8)

	if err := h.repo.Set(id, req.URL); err != nil {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := apiResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}

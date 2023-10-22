package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/models"
)

func (h *httpHandler) apiShorten(w http.ResponseWriter, r *http.Request) {
	var req models.ApiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	if req.Url == "" {
		http.Error(w, "no URL provided", http.StatusBadRequest)
		return
	}

	id := h.generateUniqueID(8)
	h.repo.Set(id, req.Url)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.ApiResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
		return
	}
}
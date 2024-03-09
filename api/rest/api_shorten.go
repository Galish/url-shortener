package restapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/usecase"
	"github.com/Galish/url-shortener/pkg/logger"
)

// APIShorten is an API handler for creating a short URL.
//
//	POST /api/shorten
func (h *HTTPHandler) APIShorten(w http.ResponseWriter, r *http.Request) {
	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	url := &entity.URL{
		User:     r.Header.Get(middleware.AuthHeaderName),
		Original: req.URL,
	}

	err := h.usecase.Shorten(r.Context(), url)
	if errors.Is(err, usecase.ErrMissingURL) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil && !errors.Is(err, usecase.ErrConflict) {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, usecase.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := APIResponse{
		Result: h.usecase.ShortURL(url),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/middleware"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

// APIShorten is an API handler for creating a short link.
//
//	POST /api/shorten
func (h *HTTPHandler) APIShorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	if req.URL == "" {
		http.Error(w, "link not provided", http.StatusBadRequest)
		return
	}

	id := h.generateUniqueID(ctx, idLength)

	err := h.repo.Set(
		ctx,
		&model.ShortLink{
			Short:    id,
			Original: req.URL,
			User:     r.Header.Get(middleware.AuthHeaderName),
		},
	)
	errConflict := repoErr.AsErrConflict(err)

	if err != nil && errConflict == nil {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if errConflict != nil {
		id = errConflict.ShortURL
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := APIResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}

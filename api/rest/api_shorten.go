package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/middleware"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
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

	errConflict := repoErr.AsErrConflict(err)

	if err != nil && errConflict == nil {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	short := url.Short

	if errConflict != nil {
		short = errConflict.ShortURL
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	resp := APIResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseURL, short),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot encode request JSON body")
	}
}

package restapi

import (
	"errors"
	"io"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/usecase"
	"github.com/Galish/url-shortener/pkg/logger"
)

// Shorten generates and returns a short URL for the given one.
//
//	POST /
func (h *HTTPHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		logger.WithError(err).Debug("unable to read request body")
		return
	}

	url := &entity.URL{
		User:     r.Header.Get(middleware.AuthHeaderName),
		Original: string(rawBody),
	}

	err = h.usecase.Shorten(r.Context(), url)

	if errors.Is(err, usecase.ErrMissingURL) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil && !errors.Is(err, usecase.ErrConflict) {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if errors.Is(err, usecase.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	shortURL := h.usecase.ShortURL(url)

	w.Write([]byte(shortURL))
}

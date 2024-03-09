package restapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/middleware"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
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
	errConflict := repoErr.AsErrConflict(err)

	if errors.Is(err, usecase.ErrMissingURL) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil && errConflict == nil {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "text/html")

	short := url.Short

	if errConflict != nil {
		short = errConflict.ShortURL
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	shortURL := fmt.Sprintf("%s/%s", h.cfg.BaseURL, short)

	w.Write([]byte(shortURL))
}

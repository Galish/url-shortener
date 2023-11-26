package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/middleware"
	repoErr "github.com/Galish/url-shortener/internal/app/repository/errors"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (h *httpHandler) shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		logger.WithError(err).Debug("unable to read request body")
		return
	}

	url := string(rawBody)

	if url == "" {
		http.Error(w, "link not provided", http.StatusBadRequest)
		return
	}

	id := h.generateUniqueID(ctx, idLength)

	err = h.repo.Set(
		ctx,
		&model.ShortLink{
			Short:    id,
			Original: url,
			User:     r.Header.Get(middleware.AuthHeaderName),
		},
	)
	errConflict := repoErr.AsErrConflict(err)

	if err != nil && errConflict == nil {
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		logger.WithError(err).Debug("unable to write to repository")
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if errConflict != nil {
		id = errConflict.ShortURL
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	fullLink := fmt.Sprintf("%s/%s", h.cfg.BaseURL, id)

	w.Write([]byte(fullLink))
}

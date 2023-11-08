package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/repository"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (h *httpHandler) shorten(w http.ResponseWriter, r *http.Request) {
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

	id := h.generateUniqueID(8)
	err = h.repo.Set(id, url)
	errConflict := repository.AsErrConflict(err)

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

func (h *httpHandler) generateUniqueID(length int) string {
	for {
		id := generateID(length)

		if !h.repo.Has(id) {
			return id
		}
	}
}

func generateID(length int) string {
	id := make([]byte, length)

	for i := range id {
		id[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(id)
}

package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
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

	h.repo.Set(id, url)

	fullLink := fmt.Sprintf("%s/%s", h.cfg.BaseURL, id)

	w.WriteHeader(http.StatusCreated)
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

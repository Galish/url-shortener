package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (h *httpHandler) makeShortLink(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	link := string(rawBody)

	if link == "" {
		http.Error(w, "link not provided", http.StatusBadRequest)
		return
	}

	id := h.generateUniqueID(8)

	h.store.Set(id, link)

	fullLink := fmt.Sprintf("http://localhost:8080/%s", id)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullLink))
}

func (h *httpHandler) generateUniqueID(length int) string {
	for {
		id := generateID(length)

		if !h.store.Has(id) {
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

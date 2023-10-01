package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (s *shortenerService) makeShortLink(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	link := string(rawBody)
	id := generateUniqueId(8)

	s.store.Set(id, link)

	fullLink := fmt.Sprintf("http://localhost:8080/%s", id)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullLink))
}

func generateUniqueId(keyLength int) string {
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, keyLength)

	for i := range id {
		id[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(id)
}

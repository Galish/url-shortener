package handlers

import (
	"context"
	"math/rand"
)

const idLength = 8

func (h *HTTPHandler) generateUniqueID(ctx context.Context, length int) string {
	for {
		id := generateID(length)

		if !h.repo.Has(ctx, id) {
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

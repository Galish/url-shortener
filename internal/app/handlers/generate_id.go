package handlers

import "context"

const idLength = 8

func (h *httpHandler) generateUniqueID(ctx context.Context, length int) string {
	for {
		id := generateID(length)

		if !h.repo.Has(ctx, id) {
			return id
		}
	}
}

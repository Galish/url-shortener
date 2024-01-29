package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/Galish/url-shortener/internal/app/middleware"
	"github.com/Galish/url-shortener/internal/app/repository/model"
)

// APIDeleteUserLinks is an API handler for deleting user short links.
//
//	DELETE /api/user/urls
func (h *HTTPHandler) APIDeleteUserLinks(w http.ResponseWriter, r *http.Request) {
	var shortURLs []string
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	if len(shortURLs) != 0 {
		go func() {
			for _, shortURL := range shortURLs {
				h.messageCh <- &handlerMessage{
					action: "delete",
					shortLink: &model.ShortLink{
						Short: shortURL,
						User:  r.Header.Get(middleware.AuthHeaderName),
					},
				}
			}
		}()
	}

	w.WriteHeader(http.StatusAccepted)
}

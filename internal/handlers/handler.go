package handlers

import (
	"context"
	"time"

	"github.com/Galish/url-shortener/internal/config"
	"github.com/Galish/url-shortener/internal/logger"
	"github.com/Galish/url-shortener/internal/repository"
	"github.com/Galish/url-shortener/internal/repository/model"
)

// HTTPHandler represents API handler.
type HTTPHandler struct {
	cfg       *config.Config
	repo      repository.Repository
	messageCh chan *handlerMessage
	ticker    *time.Ticker
}

type handlerMessage struct {
	action    string
	shortLink *model.ShortLink
}

// NewHandler implements HTTP handlers.
func NewHandler(cfg *config.Config, repo repository.Repository) *HTTPHandler {
	handler := &HTTPHandler{
		cfg:       cfg,
		repo:      repo,
		messageCh: make(chan *handlerMessage, 100),
	}

	go handler.flushMessages()

	return handler
}

func (h *HTTPHandler) flushMessages() {
	h.ticker = time.NewTicker(2 * time.Second)

	var deleteLinks []*model.ShortLink

	for {
		select {
		case message := <-h.messageCh:

			switch message.action {
			case "delete":
				deleteLinks = append(deleteLinks, message.shortLink)

			}
		case <-h.ticker.C:
			if len(deleteLinks) == 0 {
				continue
			}

			if err := h.repo.Delete(context.TODO(), deleteLinks...); err != nil {
				logger.WithError(err).Debug("cannot delete messages")
				continue
			}

			deleteLinks = nil
		}
	}
}

// Close  is executed to release the memory
func (h *HTTPHandler) Close() {
	h.ticker.Stop()
	close(h.messageCh)
}

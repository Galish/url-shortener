package restapi

import (
	"context"
	"time"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/internal/app/repository/model"
	"github.com/Galish/url-shortener/pkg/logger"
)

// HTTPHandler represents API handler.
type HTTPHandler struct {
	cfg         *config.Config
	repo        repository.Repository
	messageCh   chan *handlerMessage
	deleteLinks []*model.ShortLink
	ticker      *time.Ticker
	close       chan struct{}
	done        chan struct{}
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
		close:     make(chan struct{}),
		done:      make(chan struct{}),
	}

	go handler.run()

	return handler
}

func (h *HTTPHandler) run() {
	h.ticker = time.NewTicker(2 * time.Second)

loop:
	for {
		select {
		case message := <-h.messageCh:
			if message == nil {
				continue
			}

			switch message.action {
			case "delete":
				h.deleteLink(message.shortLink)
			}

		case <-h.ticker.C:
			h.flush()

		case <-h.close:
			h.flush()
			break loop
		}
	}

	close(h.done)
}

func (h *HTTPHandler) flush() {
	if len(h.deleteLinks) == 0 {
		return
	}

	if err := h.repo.Delete(context.TODO(), h.deleteLinks...); err != nil {
		logger.WithError(err).Debug("cannot delete messages")
		return
	}

	h.deleteLinks = nil
}

func (h *HTTPHandler) deleteLink(sl *model.ShortLink) {
	h.deleteLinks = append(h.deleteLinks, sl)
}

// Close  is executed to release the memory
func (h *HTTPHandler) Close() error {
	logger.Info("shutting down API handler")

	close(h.messageCh)
	h.messageCh = nil

	close(h.close)

	<-h.done

	return nil
}

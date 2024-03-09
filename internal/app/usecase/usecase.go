package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Galish/url-shortener/internal/app/config"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/repository"
	"github.com/Galish/url-shortener/pkg/logger"
)

var (
	ErrMissingURL = errors.New("no URL provided")
	ErrConflict   = errors.New("URL already exists")
)

type Shortener interface {
	Shorten(context.Context, ...*entity.URL) error
	ShortURL(*entity.URL) string
	Get(context.Context, string) (*entity.URL, error)
	GetByUser(context.Context, string) ([]*entity.URL, error)
	GetStats(context.Context) (urls, users int, err error)
	Delete(context.Context, []*entity.URL)
}

// ShortenerUseCase implements shortener logic.
type ShortenerUseCase struct {
	cfg       *config.Config
	repo      repository.Repository
	messageCh chan *shortenerMessage
	deleteURL []*entity.URL
	ticker    *time.Ticker
	close     chan struct{}
	done      chan struct{}
}

type shortenerMessage struct {
	action string
	url    *entity.URL
}

// New configures and returns Shortener usecase.
func New(cfg *config.Config, repo repository.Repository) *ShortenerUseCase {
	uc := &ShortenerUseCase{
		cfg:       cfg,
		repo:      repo,
		messageCh: make(chan *shortenerMessage, 100),
		close:     make(chan struct{}),
		done:      make(chan struct{}),
	}

	go uc.run()

	return uc
}

func (uc *ShortenerUseCase) run() {
	uc.ticker = time.NewTicker(2 * time.Second)

loop:
	for {
		select {
		case message := <-uc.messageCh:
			if message == nil {
				continue
			}

			switch message.action {
			case "delete":
				uc.delete(message.url)
			}

		case <-uc.ticker.C:
			uc.flush()

		case <-uc.close:
			uc.flush()
			break loop
		}
	}

	close(uc.done)
}

func (uc *ShortenerUseCase) flush() {
	if len(uc.deleteURL) == 0 {
		return
	}

	if err := uc.repo.Delete(context.TODO(), uc.deleteURL...); err != nil {
		logger.WithError(err).Debug("cannot delete messages")
		return
	}

	uc.deleteURL = nil
}

func (uc *ShortenerUseCase) delete(sl *entity.URL) {
	uc.deleteURL = append(uc.deleteURL, sl)
}

// Close  is executed to release the memory
func (uc *ShortenerUseCase) Close() error {
	logger.Info("shutting down API handler")

	close(uc.messageCh)
	uc.messageCh = nil

	close(uc.close)

	<-uc.done

	return nil
}

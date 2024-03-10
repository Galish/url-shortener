package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
	"github.com/Galish/url-shortener/internal/app/entity"
)

func (s *ShortenerServer) Shorten(
	ctx context.Context,
	in *pb.ShortenRequest,
) (*pb.ShortenResponse, error) {
	var response pb.ShortenResponse

	url := &entity.URL{
		// User:     r.Header.Get(middleware.AuthHeaderName),
		Original: in.OriginalUrl,
	}

	if err := s.usecase.Shorten(ctx, url); err != nil {
		response.Error = err.Error()
	} else {
		response.ShortUrl = url.Short
	}

	return &response, nil
}

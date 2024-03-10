package grpc

import (
	"context"
	"errors"

	pb "github.com/Galish/url-shortener/api/proto"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/infrastructure/grpc/interceptors"
	"github.com/Galish/url-shortener/internal/app/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ShortenerServer) Shorten(
	ctx context.Context,
	in *pb.ShortenRequest,
) (*pb.ShortenResponse, error) {
	var response pb.ShortenResponse

	user := ctx.Value(interceptors.UserContextKey).(string)

	url := &entity.URL{
		User:     user,
		Original: in.OriginalUrl,
	}

	err := s.usecase.Shorten(ctx, url)

	if errors.Is(err, usecase.ErrMissingURL) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, usecase.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.ShortUrl = url.Short
	}

	return &response, nil
}

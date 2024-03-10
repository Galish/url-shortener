package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/infrastructure/grpc/interceptors"
)

// ShortenBatch generates short URLs in batches.
func (s *ShortenerServer) ShortenBatch(
	ctx context.Context,
	in *pb.ShortenBatchRequest,
) (*pb.ShortenBatchResponse, error) {
	var response pb.ShortenBatchResponse

	user := ctx.Value(interceptors.UserContextKey).(string)
	urls := make([]*entity.URL, len(in.Urls))

	for i, url := range in.Urls {
		urls[i] = &entity.URL{
			User:     user,
			Original: url.OriginalUrl,
		}
	}

	err := s.usecase.Shorten(ctx, urls...)
	if err != nil {
		response.Error = err.Error()

		return &response, nil
	}

	response.Urls = make([]*pb.ShortenBatchResponseEntity, len(in.Urls))

	for i, url := range urls {
		response.Urls[i] = &pb.ShortenBatchResponseEntity{
			CorrelationId: in.Urls[i].CorrelationId,
			ShortUrl:      url.Short,
		}
	}

	return &response, nil
}

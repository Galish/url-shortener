package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
	"github.com/Galish/url-shortener/internal/app/entity"
	"github.com/Galish/url-shortener/internal/app/infrastructure/grpc/interceptors"
)

// Delete deletes user URLs.
func (s *ShortenerServer) Delete(
	ctx context.Context,
	in *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {
	var response pb.DeleteResponse

	user := ctx.Value(interceptors.UserContextKey).(string)
	urls := make([]*entity.URL, len(in.ShortUrl))

	for i, id := range in.ShortUrl {
		urls[i] = &entity.URL{
			Short: id,
			User:  user,
		}
	}

	s.usecase.Delete(ctx, urls...)

	return &response, nil
}

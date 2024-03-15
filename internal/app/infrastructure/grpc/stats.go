package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Galish/url-shortener/api/proto"
)

// GetStats returns the number of users and shortened URLs.
func (s *ShortenerServer) GetStats(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.StatsResponse, error) {
	var response pb.StatsResponse

	urls, users, err := s.usecase.Stats(ctx)
	if err != nil {
		response.Error = err.Error()
	} else {
		response.Urls = int32(urls)
		response.Users = int32(users)
	}

	return &response, nil
}

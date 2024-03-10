package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ShortenerServer) GetStats(ctx context.Context, _ *emptypb.Empty) (*pb.APIStatsResponse, error) {
	var response pb.APIStatsResponse

	urls, users, err := s.usecase.GetStats(ctx)
	if err != nil {
		response.Error = err.Error()
	} else {
		response.Urls = int32(urls)
		response.Users = int32(users)
	}

	return &response, nil
}

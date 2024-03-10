package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
)

func (s *ShortenerServer) Get(
	ctx context.Context,
	in *pb.UrlRequest,
) (*pb.UrlResponse, error) {
	var response pb.UrlResponse

	url, err := s.usecase.Get(ctx, in.ShortUrl)
	if err != nil {
		response.Error = err.Error()
	} else if url.IsDeleted {
		response.Error = "URL has been removed"
	} else {
		response.OriginalUrl = url.Original
	}

	return &response, nil
}

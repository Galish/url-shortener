package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ShortenerServer) GetByUser(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.UserUrlResponse, error) {
	var response pb.UserUrlResponse

	urls, err := s.usecase.GetByUser(ctx, "userid")
	if err != nil {
		response.Error = err.Error()

		return &response, nil
	}

	// if len(urls) == 0 {
	// 	w.WriteHeader(http.StatusNoContent)
	// 	return
	// }

	response.Urls = make([]*pb.UserUrlResponseEntity, len(urls))

	for i, url := range urls {
		response.Urls[i] = &pb.UserUrlResponseEntity{
			ShortUrl:    url.Short,
			OriginalUrl: url.Original,
		}
	}

	return &response, nil
}

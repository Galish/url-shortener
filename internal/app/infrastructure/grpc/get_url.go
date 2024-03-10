package grpc

import (
	"context"

	pb "github.com/Galish/url-shortener/api/proto"
)

func (s *ShortenerServer) Get(ctx context.Context, in *pb.UrlRequest) (*pb.UrlResponse, error) {
	var response pb.UrlResponse

	url, err := s.usecase.Get(ctx, in.ShortUrl)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		// logger.WithError(err).Debug("unable to fetch url")
		// return
		response.Error = err.Error()
	} else {
		response.OriginalUrl = url.Original
	}

	// if url.IsDeleted {
	// 	w.WriteHeader(http.StatusGone)
	// 	return
	// }

	return &response, nil
}

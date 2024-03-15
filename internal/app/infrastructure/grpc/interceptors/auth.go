package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

var UserContextKey = contextKey("user")

var methodsUserRequired = map[string]bool{
	"/service.Shortener/Shorten":      true,
	"/service.Shortener/ShortenBatch": true,
	"/service.Shortener/GetByUser":    true,
	"/service.Shortener/Delete":       true,
}

// UserCheckInterceptor serves authentication error if user identifier not provided.
func UserCheckInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if isRequired := methodsUserRequired[info.FullMethod]; !isRequired {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user identifier")
	}

	values := md.Get("user")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing user identifier")
	}

	if values[0] == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid user identifier")
	}

	ctx = context.WithValue(ctx, UserContextKey, values[0])

	return handler(ctx, req)

}

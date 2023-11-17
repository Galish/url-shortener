package middleware

import "net/http"

func Apply(
	handlerFn http.HandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {
	wrappedFn := handlerFn

	for _, middleware := range middlewares {
		wrappedFn = middleware(wrappedFn)
	}

	return wrappedFn
}

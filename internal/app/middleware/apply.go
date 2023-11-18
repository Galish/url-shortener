package middleware

import "net/http"

func Apply(
	handlerFn http.HandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {
	wrappedFn := handlerFn

	for i := len(middlewares) - 1; i >= 0; i-- {
		wrappedFn = middlewares[i](wrappedFn)
	}

	return wrappedFn
}

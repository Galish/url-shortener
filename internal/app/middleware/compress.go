package middleware

import (
	"net/http"
	"strings"

	"github.com/Galish/url-shortener/internal/app/compress"
)

func WithCompression(h http.HandlerFunc) http.HandlerFunc {
	compressor := compress.NewCompressor(compress.NewGzipCompressor())

	return func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isCompressionSupported := strings.Contains(acceptEncoding, compressor.Algorithm)

		if isCompressionSupported {
			cw := compressor.NewWriter(w)
			w = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isCompressed := strings.Contains(contentEncoding, compressor.Algorithm)

		if isCompressed {
			cr, err := compressor.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h(w, r)
	}
}

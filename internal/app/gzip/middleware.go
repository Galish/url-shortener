package gzip

import (
	"net/http"
	"strings"
)

func WithCompression(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		contentType := r.Header.Get("Content-Type")
		isContentTypeSupported := contentType == "application/json" || contentType == "text/html"

		acceptEncoding := r.Header.Get("Accept-Encoding")
		isGzipSupported := strings.Contains(acceptEncoding, "gzip")

		if isContentTypeSupported && isGzipSupported {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isGzipped := strings.Contains(contentEncoding, "gzip")

		if isContentTypeSupported && isGzipped {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = cr

			defer cr.Close()
		}

		h(ow, r)
	}
}

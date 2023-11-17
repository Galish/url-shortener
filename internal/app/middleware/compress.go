package middleware

import (
	"io"
	"net/http"
	"strings"

	"github.com/Galish/url-shortener/internal/app/compress"
)

var supportedContentTypes = [2]string{
	"application/json",
	"text/html",
}

func WithCompressor(compressor compress.Compressor) func(http.HandlerFunc) http.HandlerFunc {
	return func(handlerFn http.HandlerFunc) http.HandlerFunc {
		return WithCompression(handlerFn, compressor)
	}
}

func WithCompression(h http.HandlerFunc, compressor compress.Compressor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isCompressionSupported := strings.Contains(acceptEncoding, compressor.String())

		if isCompressionSupported {
			cw := newCompressWriter(w, compressor)
			w = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isCompressed := strings.Contains(contentEncoding, compressor.String())

		if isCompressed {
			cr, err := newCompressReader(r.Body, compressor)
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

type compressReader struct {
	r  io.ReadCloser
	zr io.ReadCloser
}

func newCompressReader(r io.ReadCloser, compressor compress.Compressor) (*compressReader, error) {
	zr, err := compressor.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}

	return cr.zr.Close()
}

type compressWriter struct {
	w          http.ResponseWriter
	zw         io.WriteCloser
	compressor compress.Compressor
}

func newCompressWriter(w http.ResponseWriter, compressor compress.Compressor) *compressWriter {
	return &compressWriter{
		w:          w,
		compressor: compressor,
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	if !cw.isContentTypeSupported() {
		return cw.w.Write(p)
	}

	cw.zw = cw.compressor.NewWriter(cw.w)

	return cw.zw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if cw.isContentTypeSupported() {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}

	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	if cw.zw == nil {
		return nil
	}

	return cw.zw.Close()
}

func (cw *compressWriter) isContentTypeSupported() bool {
	contentType := cw.w.Header().Get("Content-Type")

	for _, v := range supportedContentTypes {
		if contentType == v {
			return true
		}
	}

	return false
}

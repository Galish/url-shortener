package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

type GzipCompressor struct {
	algorithm string
}

func NewGzipCompressor() *GzipCompressor {
	return &GzipCompressor{"gzip"}
}

func (gz *GzipCompressor) NewReader(r io.ReadCloser) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func (gz *GzipCompressor) NewWriter(w http.ResponseWriter) io.WriteCloser {
	return gzip.NewWriter(w)
}

func (gz *GzipCompressor) String() string {
	return gz.algorithm
}

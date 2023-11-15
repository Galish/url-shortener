package compress

import (
	"compress/gzip"
	"io"
)

type gzipCompressor struct {
	algorithm string
}

func NewGzipCompressor() *gzipCompressor {
	return &gzipCompressor{"gzip"}
}

func (gz *gzipCompressor) NewReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func (gz *gzipCompressor) NewWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}

func (gz *gzipCompressor) String() string {
	return gz.algorithm
}

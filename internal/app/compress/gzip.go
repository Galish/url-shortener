package compress

import (
	"compress/gzip"
	"io"
)

type gzipCompressor struct {
	algorithm string
}

// NewGzipCompressor creates a Gzip compressor.
func NewGzipCompressor() *gzipCompressor {
	return &gzipCompressor{"gzip"}
}

// NewReader returns a new Reader.
func (gz *gzipCompressor) NewReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

// NewWriter returns a new Writer.
func (gz *gzipCompressor) NewWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}

// String returns the name of the compression algorithm.
func (gz *gzipCompressor) String() string {
	return gz.algorithm
}

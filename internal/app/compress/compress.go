package compress

import (
	"io"
	"net/http"
)

type compressEngine interface {
	NewReader(io.ReadCloser) (io.ReadCloser, error)
	NewWriter(http.ResponseWriter) io.WriteCloser
	String() string
}

type compressor struct {
	engine    compressEngine
	Algorithm string
}

func NewCompressor(e compressEngine) *compressor {
	return &compressor{
		engine:    e,
		Algorithm: e.String(),
	}
}

func (c *compressor) NewReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := c.engine.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *compressor) NewWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:      w,
		engine: c.engine,
	}
}

package compress

import (
	"io"
)

type compressReader struct {
	r  io.ReadCloser
	zr io.ReadCloser
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

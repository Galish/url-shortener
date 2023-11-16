package compress

import (
	"fmt"
	"io"
)

type Compressor interface {
	NewReader(io.Reader) (io.ReadCloser, error)
	NewWriter(io.Writer) io.WriteCloser
	fmt.Stringer
}

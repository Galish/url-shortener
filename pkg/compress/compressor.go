// Package compress implements compression algorithms.
package compress

import (
	"fmt"
	"io"
)

// Compressor represents an instance of the compression algorithm.
type Compressor interface {
	NewReader(io.Reader) (io.ReadCloser, error)
	NewWriter(io.Writer) io.WriteCloser
	fmt.Stringer
}

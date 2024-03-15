package generator

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	ids := map[string]bool{}

	for i := 0; i < 10; i++ {
		id := NewID(8)

		assert.False(t, ids[id])

		ids[id] = true

		assert.Regexp(
			t,
			regexp.MustCompile("[0-9A-Za-z]{8}"),
			id,
		)
	}
}

func BenchmarkNewID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewID(10)
	}
}

func ExampleNewID() {
	NewID(10)
}

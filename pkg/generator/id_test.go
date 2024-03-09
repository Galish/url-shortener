package generator

import "testing"

func BenchmarkNewID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewID(10)
	}
}

func Example_NewID() {
	NewID(10)
}

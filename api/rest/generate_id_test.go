package restapi

import "testing"

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateID(10)
	}
}

func Example_generateID() {
	generateID(10)
}

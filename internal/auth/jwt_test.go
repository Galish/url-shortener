package auth

import (
	"testing"
)

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewToken(
			&JWTClaims{
				UserID: "395fd5f4-964d-4135-9a55-fbf91c4a1614",
			},
		)
	}
}

func BenchmarkParseToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIzOTVmZDVmNC05NjRkLTQxMzUtOWE1NS1mYmY5MWM0YTE2MTQifQ.ieGQmU4jMeUzxEFQ22wScIrQYe8ePTlSYSyBeSgHtTweyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIzOTVmZDVmNC05NjRkLTQxMzUtOWE1NS1mYmY5MWM0YTE2MTQifQ.ieGQmU4jMeUzxEFQ22wScIrQYe8ePTlSYSyBeSgHtTw")
	}
}

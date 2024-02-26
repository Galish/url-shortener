package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		name   string
		claims *JWTClaims
		hasErr bool
		want   string
	}{
		{
			"blank claims",
			&JWTClaims{},
			false,
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIifQ.nucTpEFKTGXpePK8UYy7GGVwmiYd29l86EMB-_8UQ28",
		},
		{
			"claims containing user identifier",
			&JWTClaims{
				UserID: "#12345",
			},
			false,
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIjMTIzNDUifQ.VsXLWQqA8jbGie95Oj58X0wObuiYcN-3qitptUUwNjQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.claims)

			assert.Equal(t, tt.hasErr, err != nil)
			assert.Equal(t, tt.want, token)
		})
	}
}

func TestParseToken(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		hasErr bool
		want   *JWTClaims
	}{
		{
			"blank claims",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIifQ.nucTpEFKTGXpePK8UYy7GGVwmiYd29l86EMB-_8UQ28",
			false,
			&JWTClaims{},
		},
		{
			"claims containing user identifier",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIjMTIzNDUifQ.VsXLWQqA8jbGie95Oj58X0wObuiYcN-3qitptUUwNjQ",
			false,
			&JWTClaims{
				UserID: "#12345",
			},
		},
		{
			"claims containing user identifier",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			true,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ParseToken(tt.token)

			assert.Equal(t, tt.hasErr, err != nil)
			assert.Equal(t, tt.want, token)
		})
	}
}

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

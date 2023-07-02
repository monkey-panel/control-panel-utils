package utils

import "testing"

func TestBcryptHash(t *testing.T) {
	for i := 0; i < 72; i++ {
		tc := RandomString(i)
		if hash := BcryptHash(tc); !BcryptCheck(tc, hash) {
			t.Errorf("BcryptHash(%s) error", tc)
		}
	}
}

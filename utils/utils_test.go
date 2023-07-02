package utils

import "testing"

func TestRandStringBytes(t *testing.T) {
	t.Parallel()
	t.Log(RandomString(64))
}

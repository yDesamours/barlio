package helper

import (
	"testing"
)

func TestHashPasswordTest(t *testing.T) {
	var s = "1234"
	var hash string
	var err error

	t.Run("hashtest", func(t *testing.T) {
		hash, err = HashPassword(s)
		t.Log(hash)
		if err != nil {
			t.Fail()
		}
	})

	t.Run("comparehash", func(t *testing.T) {
		if !CompareHash(s, hash) {
			t.Fail()
		}
	})

}

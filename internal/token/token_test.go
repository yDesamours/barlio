package token

import "testing"

func TestGenetateToken(t *testing.T) {
	_, _, err := GenerateToken()
	if err != nil {
		t.Fail()
	}
}

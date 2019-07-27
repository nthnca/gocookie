package tokenizer_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nthnca/gocookie/internal/tokenizer"
)

func SimpleHarness(tst *testing.T, in string, token int, expected string, ee error) {
	callback := func(t []byte) error {
		if expected != string(t) {
			tst.Errorf("Result = %q, expected %q", t, expected)
		}

		return nil
	}

	err :=
		tokenizer.CreateTokenizer(bytes.NewReader([]byte(in))).NextToken(
			[]tokenizer.TokenAction{{token, callback}})
	if err != nil {
		if ee == nil {
			tst.Errorf("Unexpected error: %q", err)
		}
	}
}

func TestIdentToken(tst *testing.T) {
	SimpleHarness(tst, "   input", tokenizer.IdentToken, "input", nil)
}

func TestIdentTokenWithColon(tst *testing.T) {
	SimpleHarness(tst, "   :in\nas", tokenizer.IdentToken, "", fmt.Errorf("hi"))
}

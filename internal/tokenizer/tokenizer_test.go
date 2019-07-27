package tokenizer

import (
	"bytes"
	"fmt"
	"testing"
)

func SimpleHarness(tst *testing.T, in string, token int, expected string, ee error) {
	callback := func(t []byte) error {
		if expected != string(t) {
			tst.Errorf("Result = %q, expected %q", t, expected)
		}

		return nil
	}

	err := CreateTokenizer(bytes.NewReader([]byte(in))).NextToken([]TokenAction{
		TokenAction{token, callback}})
	if err != nil {
		if ee == nil {
			tst.Errorf("Unexpected error: %q", err)
		}
	}
}

func TestIdentToken(tst *testing.T) {
	SimpleHarness(tst, "   input", IDENT_TOKEN, "input", nil)
}

func TestIdentTokenWithColon(tst *testing.T) {
	SimpleHarness(tst, "   :in\nas", IDENT_TOKEN, "", fmt.Errorf("hi"))
}

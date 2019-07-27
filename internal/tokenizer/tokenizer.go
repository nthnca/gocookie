// Package tokenizer provides functionality for splitting Cookie source code into
// *tokens*.
package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
)

const (
	// IdentToken is a cookie language identifier: [a-z_][a-z0-9_]*
	IdentToken = iota

	// IntToken is a cookie language integer: [-]?[1-9][0-9]*|0
	IntToken = iota

	// AssignToken is a cookie language assignment token: =
	AssignToken = iota

	// EolToken is a cookie language assignment end-of-line token: ;
	EolToken = iota

	// CurlyOpenToken is a cookie language method creation token: {
	CurlyOpenToken = iota

	// CurlyCloseToken is a cookie language method end token: {
	CurlyCloseToken = iota

	// FunctionToken is a cookie language function token: ()
	FunctionToken = iota
	maxToken      = iota
)

var (
	regex        [maxToken]*regexp.Regexp
	whiteRe      = regexp.MustCompile(`^\s*`)
	identRe      = regexp.MustCompile("^[a-z_][a-z0-9_]*")
	intRe        = regexp.MustCompile("^[-]?[1-9][0-9]*|0")
	assignRe     = regexp.MustCompile("^=")
	eolRe        = regexp.MustCompile("^;")
	curlyOpenRe  = regexp.MustCompile("^{")
	curlyCloseRe = regexp.MustCompile("^}")
	functionRe   = regexp.MustCompile("^[(][)]")
)

// Tokenizer is the basic handle that stores state for tokenizer a given IO stream.
type Tokenizer struct {
	reader   *bufio.Reader
	currLine []byte
	line     int
	column   int
}

// TokenAction is a simple mapping of a token type (IdentToken, IntToken, etc) to a
// method that will be invoked if that type of token is found. The []byte will be the
// token that was found.
type TokenAction struct {
	Token  int
	Action func([]byte) error
}

func init() {
	regex[IdentToken] = identRe
	regex[IntToken] = intRe
	regex[AssignToken] = assignRe
	regex[EolToken] = eolRe
	regex[CurlyOpenToken] = curlyOpenRe
	regex[CurlyCloseToken] = curlyCloseRe
	regex[FunctionToken] = functionRe
}

// CreateTokenizer initializes and returns a Tokenizer object that can be used to
// split the IO stream into a series of tokens.
func CreateTokenizer(rd io.Reader) *Tokenizer {
	return &Tokenizer{bufio.NewReader(rd), make([]byte, 0), 0, 0}
}

func (t *Tokenizer) nextChunk() ([]byte, error) {
	for {
		m := whiteRe.Find(t.currLine)
		if len(m) != len(t.currLine) {
			t.column += len(m)
			t.currLine = t.currLine[len(m):]
			return t.currLine, nil
		}

		line, _, err := t.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			log.Fatalf("Unexpected error while reading: %s", err)
		}
		t.currLine = line
		t.line++
		t.column = 0
	}
}

// NextToken takes a set of TokenActions and based on which Token is found next,
// invokes the given action associated with it. An error is returned if none of the
// given tokens is found next.
func (t *Tokenizer) NextToken(ra []TokenAction) error {
	chunk, err := t.nextChunk()
	if err == io.EOF {
		return io.EOF
	}

	for _, e := range ra {
		m := regex[e.Token].Find(chunk)
		if len(m) == 0 {
			continue
		}
		t.column += len(m)
		t.currLine = t.currLine[len(m):]
		return e.Action(m)
	}

	return fmt.Errorf("Unexpected token (line %d, col %d): '%s'",
		t.line, t.column, chunk)
}

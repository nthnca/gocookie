package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
)

const (
	IdentToken       = iota
	IntToken         = iota
	AssignToken      = iota
	EolToken         = iota
	CurlyOpenToken  = iota
	CurlyCloseToken = iota
	FunctionToken    = iota
	maxToken          = iota
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

type Tokenizer struct {
	reader   *bufio.Reader
	currLine []byte
	line     int
	column   int
}

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

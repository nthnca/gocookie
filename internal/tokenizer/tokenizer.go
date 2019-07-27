package tokenizer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
)

const (
	IDENT_TOKEN       = iota
	INT_TOKEN         = iota
	ASSIGN_TOKEN      = iota
	EOL_TOKEN         = iota
	CURLY_OPEN_TOKEN  = iota
	CURLY_CLOSE_TOKEN = iota
	FUNCTION_TOKEN    = iota
	max_token         = iota
)

var (
	regex          [max_token]*regexp.Regexp
	white_re       *regexp.Regexp = regexp.MustCompile(`^\s*`)
	ident_re       *regexp.Regexp = regexp.MustCompile("^[a-z_][a-z0-9_]*")
	int_re         *regexp.Regexp = regexp.MustCompile("^[-]?[1-9][0-9]*|0")
	assign_re      *regexp.Regexp = regexp.MustCompile("^=")
	eol_re         *regexp.Regexp = regexp.MustCompile("^;")
	curly_open_re  *regexp.Regexp = regexp.MustCompile("^{")
	curly_close_re *regexp.Regexp = regexp.MustCompile("^}")
	function_re    *regexp.Regexp = regexp.MustCompile("^[(][)]")
)

type Tokenizer struct {
	reader    *bufio.Reader
	curr_line []byte
	line      int
	column    int
}

type TokenAction struct {
	Token  int
	Action func([]byte) error
}

func init() {
	regex[IDENT_TOKEN] = ident_re
	regex[INT_TOKEN] = int_re
	regex[ASSIGN_TOKEN] = assign_re
	regex[EOL_TOKEN] = eol_re
	regex[CURLY_OPEN_TOKEN] = curly_open_re
	regex[CURLY_CLOSE_TOKEN] = curly_close_re
	regex[FUNCTION_TOKEN] = function_re
}

func CreateTokenizer(rd io.Reader) *Tokenizer {
	return &Tokenizer{bufio.NewReader(rd), make([]byte, 0), 0, 0}
}

func (t *Tokenizer) nextChunk() ([]byte, error) {
	for {
		m := white_re.Find(t.curr_line)
		if len(m) != len(t.curr_line) {
			t.column += len(m)
			t.curr_line = t.curr_line[len(m):]
			return t.curr_line, nil
		}

		line, _, err := t.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			log.Fatalf("Unexpected error while reading: %s", err)
		}
		t.curr_line = line
		t.line += 1
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
		t.curr_line = t.curr_line[len(m):]
		return e.Action(m)
	}

	return errors.New(fmt.Sprintf("Unexpected token (line %d, col %d): '%s'",
		t.line, t.column, chunk))
}

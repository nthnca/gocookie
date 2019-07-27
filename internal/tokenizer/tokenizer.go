package tokenizer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
)

var (
	WHITE_RE *regexp.Regexp = regexp.MustCompile(`^\s*`)
)

type Tokenizer struct {
	reader    *bufio.Reader
	curr_line []byte
	line      int
	column    int
}

type RegexpAction struct {
	Re     *regexp.Regexp
	Action func([]byte) error
}

func CreateTokenizer(rd io.Reader) *Tokenizer {
	return &Tokenizer{bufio.NewReader(rd), make([]byte, 0), 0, 0}
}

func (t *Tokenizer) nextChunk() ([]byte, error) {
	for {
		m := WHITE_RE.Find(t.curr_line)
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

func (t *Tokenizer) NextToken(ra []RegexpAction) error {
	chunk, err := t.nextChunk()
	if err == io.EOF {
		return io.EOF
	}

	for _, e := range ra {
		m := e.Re.Find(chunk)
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

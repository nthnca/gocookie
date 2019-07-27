package parser

import (
	"io"
	"log"
	"regexp"
	"strconv"

	"github.com/nthnca/gocookie/internal/tokenizer"
)

var (
	IDENT_RE       *regexp.Regexp = regexp.MustCompile("^[a-z_][a-z0-9_]*")
	INT_RE         *regexp.Regexp = regexp.MustCompile("^[-]?[1-9][0-9]*|0")
	ASSIGN_RE      *regexp.Regexp = regexp.MustCompile("^=")
	EOL_RE         *regexp.Regexp = regexp.MustCompile("^;")
	CURLY_OPEN_RE  *regexp.Regexp = regexp.MustCompile("^{")
	CURLY_CLOSE_RE *regexp.Regexp = regexp.MustCompile("^}")
	FUNCTION_RE    *regexp.Regexp = regexp.MustCompile("^[(][)]")

	OP_TYPE_ASSIGN int = 0
	OP_TYPE_INT    int = 1
	OP_TYPE_FUNC   int = 2
	OP_TYPE_METHOD int = 3
)

type Statement struct {
	Var       string
	OpType    int
	VarInt    int
	VarVar    string
	VarMethod []Statement
}

func GetNextStmt(tr *tokenizer.Tokenizer) (*Statement, error) {
	var stmt Statement

	done := func(_ []byte) error {
		return nil
	}

	curly_close := func(_ []byte) error {
		// TODO: This case is invalid if this isn't an embedded method.
		return io.EOF
	}

	integer := func(t []byte) error {
		var err error
		stmt.VarInt, err = strconv.Atoi(string(t))
		if err != nil {
			log.Fatalf("How..., %v", err)
		}
		stmt.OpType = OP_TYPE_INT
		// log.Printf("INTEGER: %s", t)
		err = tr.NextToken([]tokenizer.RegexpAction{
			tokenizer.RegexpAction{EOL_RE, done},
		})
		return err
	}

	function := func(_ []byte) error {
		stmt.OpType = OP_TYPE_FUNC
		err := tr.NextToken([]tokenizer.RegexpAction{
			tokenizer.RegexpAction{EOL_RE, done},
		})
		return err
	}

	function_no_assign := func(t []byte) error {
		stmt.VarVar = stmt.Var
		stmt.Var = "_"
		return function(t)
	}

	literal := func(t []byte) error {
		// log.Printf("LITERAL: %s", t)
		stmt.VarVar = string(t)
		err := tr.NextToken([]tokenizer.RegexpAction{
			tokenizer.RegexpAction{EOL_RE, done},
			tokenizer.RegexpAction{FUNCTION_RE, function},
		})
		return err
	}

	method := func(t []byte) error {
		stmt.OpType = OP_TYPE_METHOD
		stmt.VarMethod = GetMethod(tr)

		return nil
	}

	assign := func(t []byte) error {
		// log.Printf("ASSIGN: %s", t)
		err := tr.NextToken([]tokenizer.RegexpAction{
			tokenizer.RegexpAction{INT_RE, integer},
			tokenizer.RegexpAction{IDENT_RE, literal},
			tokenizer.RegexpAction{CURLY_OPEN_RE, method},
		})
		return err
	}

	ident := func(t []byte) error {
		// log.Printf("IDENT: %s", t)
		stmt.Var = string(t)
		err := tr.NextToken([]tokenizer.RegexpAction{
			tokenizer.RegexpAction{ASSIGN_RE, assign},
			tokenizer.RegexpAction{FUNCTION_RE, function_no_assign},
		})
		return err
	}

	err := tr.NextToken([]tokenizer.RegexpAction{
		tokenizer.RegexpAction{IDENT_RE, ident},
		tokenizer.RegexpAction{CURLY_CLOSE_RE, curly_close},
	})
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		log.Fatalf("End of FIle? %s", err)
	}
	return &stmt, nil
}

func GetMethod(t *tokenizer.Tokenizer) []Statement {
	method := []Statement{}
	for {
		stmt, err := GetNextStmt(t)
		if err != nil {
			if err == io.EOF {
				return method
			}
			log.Fatalf("Oops")
		}
		method = append(method, *stmt)
	}
	return method
}

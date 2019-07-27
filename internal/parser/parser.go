package parser

import (
	"io"
	"log"
	"strconv"

	"github.com/nthnca/gocookie/internal/tokenizer"
)

var (
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

func getNextStmt(tr *tokenizer.Tokenizer) (*Statement, error) {
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
		err = tr.NextToken([]tokenizer.TokenAction{
			tokenizer.TokenAction{tokenizer.EOL_TOKEN, done},
		})
		return err
	}

	function := func(_ []byte) error {
		stmt.OpType = OP_TYPE_FUNC
		err := tr.NextToken([]tokenizer.TokenAction{
			tokenizer.TokenAction{tokenizer.EOL_TOKEN, done},
		})
		return err
	}

	function_no_assign := func(t []byte) error {
		stmt.VarVar = stmt.Var
		stmt.Var = "_"
		return function(t)
	}

	literal := func(t []byte) error {
		stmt.VarVar = string(t)
		err := tr.NextToken([]tokenizer.TokenAction{
			tokenizer.TokenAction{tokenizer.EOL_TOKEN, done},
			tokenizer.TokenAction{tokenizer.FUNCTION_TOKEN, function},
		})
		return err
	}

	method := func(t []byte) error {
		stmt.OpType = OP_TYPE_METHOD
		stmt.VarMethod = GetMethod(tr)

		return nil
	}

	assign := func(t []byte) error {
		err := tr.NextToken([]tokenizer.TokenAction{
			tokenizer.TokenAction{tokenizer.INT_TOKEN, integer},
			tokenizer.TokenAction{tokenizer.IDENT_TOKEN, literal},
			tokenizer.TokenAction{tokenizer.CURLY_OPEN_TOKEN, method},
		})
		return err
	}

	ident := func(t []byte) error {
		stmt.Var = string(t)
		err := tr.NextToken([]tokenizer.TokenAction{
			tokenizer.TokenAction{tokenizer.ASSIGN_TOKEN, assign},
			tokenizer.TokenAction{tokenizer.FUNCTION_TOKEN, function_no_assign},
		})
		return err
	}

	err := tr.NextToken([]tokenizer.TokenAction{
		tokenizer.TokenAction{tokenizer.IDENT_TOKEN, ident},
		tokenizer.TokenAction{tokenizer.CURLY_CLOSE_TOKEN, curly_close},
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
		stmt, err := getNextStmt(t)
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

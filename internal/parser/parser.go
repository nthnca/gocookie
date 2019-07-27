package parser

import (
	"io"
	"log"
	"strconv"

	"github.com/nthnca/gocookie/internal/tokenizer"
)

const (
	// OpTypeAssign is for assigning the value of one variable to another variable.
	OpTypeAssign = iota

	// OpTypeInt is an integer assignment to a variable.
	OpTypeInt = iota

	// OpTypeFunc is for executing a function.
	OpTypeFunc = iota

	// OpTypeMethod is for the creation of a function.
	OpTypeMethod = iota
)

// Statement is the  internal representation of a line of cookie code.
type Statement struct {
	OpType    int
	Var       string
	VarInt    int
	VarVar    string
	VarMethod []Statement
}

func getNextStmt(tr *tokenizer.Tokenizer) (*Statement, error) {
	var stmt Statement

	done := func(_ []byte) error {
		return nil
	}

	curlyClose := func(_ []byte) error {
		// TODO: This case is invalid if this isn't an embedded method.
		return io.EOF
	}

	integer := func(t []byte) error {
		var err error
		stmt.VarInt, err = strconv.Atoi(string(t))
		if err != nil {
			log.Fatalf("How..., %v", err)
		}
		stmt.OpType = OpTypeInt
		err = tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.EolToken, done},
		})
		return err
	}

	function := func(_ []byte) error {
		stmt.OpType = OpTypeFunc
		err := tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.EolToken, done},
		})
		return err
	}

	functionNoAssign := func(t []byte) error {
		stmt.VarVar = stmt.Var
		stmt.Var = "_"
		return function(t)
	}

	literal := func(t []byte) error {
		stmt.VarVar = string(t)
		err := tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.EolToken, done},
			{tokenizer.FunctionToken, function},
		})
		return err
	}

	method := func(t []byte) error {
		stmt.OpType = OpTypeMethod
		stmt.VarMethod = GetMethod(tr)

		return nil
	}

	assign := func(t []byte) error {
		err := tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.IntToken, integer},
			{tokenizer.IdentToken, literal},
			{tokenizer.CurlyOpenToken, method},
		})
		return err
	}

	ident := func(t []byte) error {
		stmt.Var = string(t)
		err := tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.AssignToken, assign},
			{tokenizer.FunctionToken, functionNoAssign},
		})
		return err
	}

	err := tr.NextToken([]tokenizer.TokenAction{
		{tokenizer.IdentToken, ident},
		{tokenizer.CurlyCloseToken, curlyClose},
	})
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		log.Fatalf("End of FIle? %s", err)
	}
	return &stmt, nil
}

// GetMethod parses up to the end of the current method and returns the set of
// statements contained within. It determines the end of the method either by finding
// the closing '}' or the end-of-file.
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
}

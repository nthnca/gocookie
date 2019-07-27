package parser

import (
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/nthnca/gocookie/internal/tokenizer"
	"golang.org/x/xerrors"
)

const (
	// OpTypeAssign is for assigning the value of one variable to another variable.
	OpTypeAssign = iota

	// OpTypeInt is an integer assignment to a variable.
	OpTypeInt = iota

	// OpTypeFunc is for executing a function.
	OpTypeFunc = iota

	// OpTypeFuncCreate is for the creation of a function.
	OpTypeFuncCreate = iota
)

var (
	// ErrEndOfFunction is for symmetry to io.EOF and denotes the end of a
	// function was found.
	ErrEndOfFunction = errors.New("end of function was found")
)

// Statement is the  internal representation of a line of cookie code.
type Statement struct {
	OpType   int
	Var      string
	VarInt   int
	VarVar   string
	VarStmts []Statement
}

func getNextStmt(tr *tokenizer.Tokenizer) (*Statement, error) {
	var stmt Statement

	done := func(_ []byte) error {
		return nil
	}

	curlyClose := func(_ []byte) error {
		// TODO: This case is invalid if this isn't an embedded method.
		return ErrEndOfFunction
	}

	integer := func(t []byte) error {
		var err error
		stmt.VarInt, err = strconv.Atoi(string(t))
		if err != nil {
			return xerrors.Errorf(
				"unable to parse integer: %w", err)
		}
		stmt.OpType = OpTypeInt
		return tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.EolToken, done},
		})
	}

	function := func(_ []byte) error {
		stmt.OpType = OpTypeFunc
		return tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.EolToken, done},
		})
	}

	functionNoAssign := func(t []byte) error {
		stmt.VarVar = stmt.Var
		stmt.Var = "_"
		return function(t)
	}

	literal := func(t []byte) error {
		stmt.VarVar = string(t)
		return tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.EolToken, done},
			{tokenizer.FunctionToken, function},
		})
	}

	createFunc := func(t []byte) error {
		stmt.OpType = OpTypeFuncCreate
		stmt.VarStmts = getFunctionInternal(tr, true)

		return nil
	}

	assign := func(t []byte) error {
		return tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.IntToken, integer},
			{tokenizer.IdentToken, literal},
			{tokenizer.CurlyOpenToken, createFunc},
		})
	}

	ident := func(t []byte) error {
		stmt.Var = string(t)
		return tr.NextToken([]tokenizer.TokenAction{
			{tokenizer.AssignToken, assign},
			{tokenizer.FunctionToken, functionNoAssign},
		})
	}

	err := tr.NextToken([]tokenizer.TokenAction{
		{tokenizer.IdentToken, ident},
		{tokenizer.CurlyCloseToken, curlyClose},
	})
	if err != nil {
		return nil, xerrors.Errorf(
			"attempting to parse statement: %w", err)
	}
	return &stmt, nil
}

func getFunctionInternal(t *tokenizer.Tokenizer, embedded bool) []Statement {
	stmts := []Statement{}
	for {
		stmt, err := getNextStmt(t)
		if err != nil {
			if xerrors.Is(err, io.EOF) {
				if !embedded {
					return stmts
				}
				log.Fatalf("Unexpected EOF, %+v", err)
			}
			if xerrors.Is(err, ErrEndOfFunction) {
				if embedded {
					return stmts
				}
				log.Fatalf("Unexpected '}', %+v", err)
			}
			log.Fatalf("Unable to parse code %+v", err)
		}
		stmts = append(stmts, *stmt)
	}
}

// GetFunction parses up to the end of the current function and returns the set of
// statements contained within. It determines the end of the function either by finding
// the closing '}' or the end-of-file.
func GetFunction(t *tokenizer.Tokenizer) []Statement {
	return getFunctionInternal(t, false)
}

// Package main provides a command line interface for running "cookie" code.
// Example: "gocookie < source.cookie"
//
// For more documentation about this language see: https://github.com/nthnca/gocookie
package main

import (
	"os"

	"github.com/nthnca/gocookie/internal/cookie"
	"github.com/nthnca/gocookie/internal/parser"
	"github.com/nthnca/gocookie/internal/tokenizer"
)

func main() {
	t := tokenizer.CreateTokenizer(os.Stdin)
	stmts := parser.GetFunction(t)
	cookie.Run(stmts)
}

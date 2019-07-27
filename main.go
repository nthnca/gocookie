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

package cookie

import (
	"fmt"

	"github.com/nthnca/gocookie/internal/parser"
)

type op struct {
	vInt     int
	vBuiltin func()
	vFunc    []parser.Statement
}

var (
	variables = make(map[string]*op)
)

func printx() {
	x := variables["_1"].vInt
	fmt.Printf("%d\n", x)
}

func addx() {
	x := variables["_1"].vInt + variables["_2"].vInt
	variables["_r"] = &op{x, nil, nil}
}

func ifx() {
	if variables["_1"].vInt != 0 {
		exe(variables["_2"])
	}
}

func loopx() {
	x := variables["_1"]
	for {
		exe(x)
		if variables["_r"].vInt == 0 {
			break
		}
	}
}

func exe(op *op) {
	if op.vBuiltin != nil {
		op.vBuiltin()
	} else {
		run(op.vFunc)
	}
}

func run(m []parser.Statement) {
	for _, s := range m {
		switch s.OpType {
		case parser.OpTypeAssign:
			variables[s.Var] = variables[s.VarVar]
		case parser.OpTypeInt:
			variables[s.Var] = &op{s.VarInt, nil, nil}
		case parser.OpTypeFuncCreate:
			variables[s.Var] = &op{0, nil, s.VarStmts}
		case parser.OpTypeFunc:
			op := variables[s.VarVar]
			exe(op)
			variables[s.Var] = variables["_r"]
		}
	}
}

// Run is the entry point for executing a set of cookie Statements, or in other words
// it is the entry point for the created program.
func Run(m []parser.Statement) {
	variables["print"] = &op{0, printx, nil}
	variables["add"] = &op{0, addx, nil}
	variables["if"] = &op{0, ifx, nil}
	variables["loop"] = &op{0, loopx, nil}

	run(m)
}

package cookie

import (
	"fmt"

	"github.com/nthnca/gocookie/internal/parser"
)

type op struct {
	v_int     int
	v_builtin func()
	v_func    []parser.Statement
}

var (
	variables = make(map[string]*op)
)

func printx() {
	x := variables["_1"].v_int
	fmt.Printf("%d\n", x)
}

func addx() {
	x := variables["_1"].v_int + variables["_2"].v_int
	variables["_r"] = &op{x, nil, nil}
}

func ifx() {
	if variables["_1"].v_int != 0 {
		exe(variables["_2"])
	}
}

func loopx() {
	x := variables["_1"]
	for {
		exe(x)
		if variables["_r"].v_int == 0 {
			break
		}
	}
}

func exe(op *op) {
	if op.v_builtin != nil {
		op.v_builtin()
	} else {
		run(op.v_func)
	}
}

func run(m []parser.Statement) {
	for _, s := range m {
		switch s.OpType {
		case parser.OP_TYPE_ASSIGN:
			variables[s.Var] = variables[s.VarVar]
		case parser.OP_TYPE_INT:
			variables[s.Var] = &op{s.VarInt, nil, nil}
		case parser.OP_TYPE_METHOD:
			variables[s.Var] = &op{0, nil, s.VarMethod}
		case parser.OP_TYPE_FUNC:
			op := variables[s.VarVar]
			exe(op)
			variables[s.Var] = variables["_r"]
		}
	}
}

func Run(m []parser.Statement) {
	variables["print"] = &op{0, printx, nil}
	variables["add"] = &op{0, addx, nil}
	variables["if"] = &op{0, ifx, nil}
	variables["loop"] = &op{0, loopx, nil}

	run(m)
}

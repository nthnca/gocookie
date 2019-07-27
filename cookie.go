package main

import (
	"fmt"
	"os"
)

type Op struct {
	v_int     int
	v_builtin func()
	v_func    []Statement
}

var (
	Variables = make(map[string]*Op)
)

func printx() {
	x := Variables["_1"].v_int
	fmt.Printf("%d\n", x)
}

func addx() {
	x := Variables["_1"].v_int + Variables["_2"].v_int
	Variables["_r"] = &Op{x, nil, nil}
}

func ifx() {
	if Variables["_1"].v_int != 0 {
		exe(Variables["_2"])
	}
}

func loopx() {
	x := Variables["_1"]
	for {
		exe(x)
		if Variables["_r"].v_int == 0 {
			break
		}
	}
}

func exe(op *Op) {
	if op.v_builtin != nil {
		op.v_builtin()
	} else {
		run(op.v_func)
	}
}

func run(m []Statement) {
	for _, s := range m {
		switch s.OpType {
		case OP_TYPE_ASSIGN:
			Variables[s.Var] = Variables[s.VarVar]
		case OP_TYPE_INT:
			Variables[s.Var] = &Op{s.VarInt, nil, nil}
		case OP_TYPE_METHOD:
			Variables[s.Var] = &Op{0, nil, s.VarMethod}
		case OP_TYPE_FUNC:
			op := Variables[s.VarVar]
			exe(op)
			Variables[s.Var] = Variables["_r"]
		}
	}
}

func main() {
	t := CreateTokenizer(os.Stdin)
	m := GetMethod(t)

	Variables["print"] = &Op{0, printx, nil}
	Variables["add"] = &Op{0, addx, nil}
	Variables["if"] = &Op{0, ifx, nil}
	Variables["loop"] = &Op{0, loopx, nil}

	run(m)
}

# cookie

The goal of this project is to create a small little program that implements a very basic
programming language, but that could be used as a starting point for people interested in
learning to code, or learning to write a programming language.

I implemented the same language in python as well: http://github.com/nthnca/cookie


[![Go Report Card](https://goreportcard.com/badge/github.com/nthnca/gocookie?style=flat-square)](https://goreportcard.com/report/github.com/nthnca/gocookie)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/nthnca/gocookie)


## Basic syntax

Assign an integer, 1, to a variable, `v`.
```
v = l;
```

Define a function, `f`:
```
f = {
  [statements]
  ...
}
```

Call a function, and assign the result to `a`.
```
a = f();
```

Some further notes:
- All variables are globally scoped. So have fun!  :-)
- You can't pass parameters to a function, but variables are globally scoped so you can
  just use a normal variable.
- A function *returns* the value that is assigned to `_r` when the function exits.
- There are just four built in functions:
  - print(): prints the value in variable `_1`.
  - add(): returns the result of adding the values in variable `_1` and `_2`.
  - if(): calls the function assigned to variable `_2`, if `_1` is non-zero.
  - loop(): repeatedly calls the function assigned to variable `_1` until that function
    returns a non-zero value.


In this language, all variables are of a global scope, and nothing is ever freed,
although variables can be re-used.

By convention we use the variables `_1`, `_2`, etc to pass data into a function.


## Examples

See the tests directory for examples.


## Ideas

So here you go, have fun improving the language, some ideas:
- Add the ability to add comments // to a program.
- Allow passing parameters to functions.
- Make it so variables are scoped to the current function.
- Add a return statement.
- Add more builtin functions (less than, multiply, etc, etc)
  - Try not to add builtin types that you could implement in the language itself.
- Add more types: strings, arrays, maps, structs, etc.
- Allow importing code from other files.
- Whatever else you want...

For bonus points:
- Re-implement this little toy language in your favorite language.
- Re-implement this little toy language, in this toy language itself...  :-)

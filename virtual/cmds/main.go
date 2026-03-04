package main

import (
	"io"
	"os"
	"fmt"
	"errors"
)

//private types
type (
	wr_t struct{
		stdin io.Reader
		stdout io.Writer
		stderr io.Writer
	}
)

type (
	Fn func(exec *Executor, args []string) (int)
	Executor struct {
		fns map[string]Fn
		stdin io.Reader
		stdout io.Writer
		stderr io.Writer
		wr wr_t
	}
)

var (
	CmdNotFound = errors.New("unrecognized command") 
)

//private values
var (
	wr wr_t
)

func main() {
	fns := Default_fns()
	exec := Init(fns, os.Stdin, os.Stdout, os.Stderr)
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough args (usage: [cmd] [args...])")
		return
	}
	ret, e := exec.Exec(os.Args[1:])
	if e != nil { fmt.Fprintf(os.Stderr, "%v\n", e) }
	fmt.Printf("\nreturn code: %d\n", ret)
}

func Default_fns() map[string]Fn {
	return map[string]Fn {
		"date": date,
	}
}

func Init(
	fns map[string]Fn,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) Executor {
	return Executor {
		fns: fns,
		stdout: stdout,
		stderr: stderr,
		stdin: stdin,
		wr: wr_t {
			stdout: stdout,
			stderr: stderr,
			stdin: stdin,
		},
	}
}

func (exec Executor) Exec(args []string) (int, error) {
	if exec.fns[args[0]] != nil { 
		return exec.fns[args[0]](&exec, args), nil
	}
	return 127, CmdNotFound
}

func (w wr_t) out(str string, a ...any) {
	fmt.Fprintf(w.stdout, str, a...)
}
func (w wr_t) err(str string, a ...any) {
	fmt.Fprintf(w.stderr, str, a...)
}

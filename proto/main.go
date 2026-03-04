package main

import (
	"io"
	"os"
	"fmt"
	"time"
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

func date(exec *Executor, args []string) (int) {
	wr := exec.wr
	var res, format string
	var cur time.Time
	var e error
	switch len(args) {
	 case 3:
		cur, e = time.Parse(args[1], args[2])
		if e != nil {
			wr.err("%v\n", e)
			return 2
		}
		fallthrough
	 case 2:
		format = args[1]
		goto format
	 case 1:
		res = time.Now().Format(time.UnixDate)
		goto done
	 default:
		panic("NO ARGS")
	}
	
	format: {
		if cur.IsZero() { cur = time.Now() }
		res = cur.Format(format)
	}

	done: {
		wr.out("%s", res)
		return 0
	}
}

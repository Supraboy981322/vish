package main

import ("time")

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

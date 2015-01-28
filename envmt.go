//  envmt.go -- dynamic variables and procedure environment

package goaldi

import (
	"fmt"
	"io"
)

//  execution environment
type Env struct {
	Parent   *Env             // parent environment
	ThreadID int              // thread ID
	VarMap   map[string]Value // dynamic variable table
}

//  NewEnv(e) returns a new environment with parent e.
func NewEnv(e *Env) *Env {
	enew := &Env{}
	enew.Parent = e
	if e == nil {
		enew.ThreadID = <-TID
		enew.VarMap = StdEnv
	} else {
		enew.ThreadID = e.ThreadID
		enew.VarMap = make(map[string]Value)
	}
	return enew
}

//  Env.Lookup(s) -- look up dynamic variable s in environment tree
func (e *Env) Lookup(s string) Value {
	for ; e != nil; e = e.Parent {
		if v := e.VarMap[s]; v != nil {
			return v
		}
	}
	return nil
}

//  ThreadID production
var TID = make(chan int)

func init() {
	go func() {
		tid := 0
		for {
			tid++
			TID <- tid
		}
	}()
}

//  StdEnv is the initial environment
var StdEnv = make(map[string]Value)

//  EnvInit registers a standard environment value or variable at init time.
//  (Variables should be registered as trapped values).
func EnvInit(name string, v Value) {
	StdEnv[name] = v
}

//  Initial dynamic variables
func init() {

	// math constants
	EnvInit("e", E)
	EnvInit("phi", PHI)
	EnvInit("pi", PI)

	// standard files (mutable)
	EnvInit("stdin", Trapped(&STDIN))
	EnvInit("stdout", Trapped(&STDOUT))
	EnvInit("stderr", Trapped(&STDERR))
}

//	ShowEnvironment(f) -- list standard environment on file f
func ShowEnvironment(f io.Writer) {
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Standard Environment")
	fmt.Fprintln(f, "------------------------------")
	for k := range SortedKeys(StdEnv) {
		cv := "c"
		v := StdEnv[k]
		if t, ok := v.(*VTrapped); ok {
			cv = "v"
			v = t.Deref()
		}
		fmt.Fprintf(f, "%%%-8s %s  %#v\n", k, cv, v)
	}
}

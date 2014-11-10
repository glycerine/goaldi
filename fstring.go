//  fstring.go -- string functions

package goaldi

import (
	"fmt"
	"strings"
	"unicode"
)

//  Declare methods
var StringMethods = map[string]interface{}{
	"type":  (*VString).Type,
	"copy":  (*VString).Copy,
	"image": (*VString).GoString,
}

//  VString.Field implements methods
func (v *VString) Field(f string) Value {
	return GetMethod(StringMethods, v, f)
}

//  This init function adds a set of Go functions to the standard library
func init() {
	// Goaldi procedures
	LibProcedure("string", String)
	LibProcedure("char", Char)
	LibProcedure("ord", Ord)
	LibProcedure("reverse", Reverse)
	// Go library functions
	LibGoFunc("equalfold", strings.EqualFold)
	LibGoFunc("replace", strings.Replace)
	LibGoFunc("repl", strings.Repeat)
	LibGoFunc("toupper", strings.ToUpper)
	LibGoFunc("tolower", strings.ToLower)
	LibGoFunc("trim", strings.Trim)
}

//  String(x) -- return argument converted to string, or fail
func String(env *Env, a ...Value) (Value, *Closure) {
	// nonstandard entry; on panic, returns default nil values
	defer func() { recover() }()
	v := ProcArg(a, 0, NilValue)
	if s, ok := v.(Stringable); ok {
		return Return(s.ToString())
	} else if s, ok := v.(fmt.Stringer); ok {
		return Return(NewString(s.String()))
	} else {
		return Return(Import(v).(Stringable).ToString())
	}
}

//  Char(i) -- return one-character string with Unicode value i
func Char(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("char", a)
	var r [1]rune
	i := int(ProcArg(a, 0, NilValue).(Numerable).ToNumber().Val())
	if i < 0 || i > int(unicode.MaxRune) {
		panic(&RunErr{"character code out of range", a[0]})
	}
	r[0] = rune(i)
	return Return(RuneString(r[:]))
}

//  Ord(c) -- return Unicode value of one-character string
func Ord(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("ord", a)
	r := ProcArg(a, 0, NilValue).(Stringable).ToString().ToRunes()
	if len(r) != 1 {
		panic(&RunErr{"string length not 1", a[0]})
	}
	return Return(NewNumber(float64(r[0])))
}

//  Reverse(s) -- return mirror image of string
func Reverse(env *Env, a ...Value) (Value, *Closure) {
	defer Traceback("reverse", a)
	r := ProcArg(a, 0, NilValue).(Stringable).ToString().ToRunes()
	n := len(r)
	for i := 0; i < n/2; i++ {
		r[i], r[n-1-i] = r[n-1-i], r[i]
	}
	return Return(RuneString(r))
}

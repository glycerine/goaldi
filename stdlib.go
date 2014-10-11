//  stdlib.go -- standard library and miscellaneous functions

//  #%#% this initial set is for testing and illustration; it is NOT final!

package goaldi

import (
	"os"
	"strings"
)

//  StdLib is the set of procedures available at link time
var StdLib = make(map[string]*VProcedure)

//  LibProcedure registers a standard library procedure taking Goaldi arguments.
func LibProcedure(name string, p Procedure) {
	StdLib[name] = NewProcedure(name, p)
}

//  LibGoFunc registers a Go function as a standard library procedure.
//  This must be done before linking (e.g. via init func) to be effective.
func LibGoFunc(name string, f interface{}) {
	StdLib[name] = GoProcedure(name, f)
}

//  This init function adds a set of Go functions to the standard library
func init() {

	LibGoFunc("equalfold", strings.EqualFold)
	LibGoFunc("replace", strings.Replace)
	LibGoFunc("toupper", strings.ToUpper)
	LibGoFunc("tolower", strings.ToLower)
	LibGoFunc("trim", strings.Trim)

	LibGoFunc("exit", os.Exit)
	LibGoFunc("remove", os.Remove)
}

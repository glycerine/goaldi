//  vproc.go -- VProcedure, the Goaldi type "procedure"
//
//  A VProcedure is created by the linker for each Go or Goaldi
//  procedure or method, and as a constructor for each VRecord.
//  Additional procedure can be created at runtime by Goaldi
//  "procedure" and "lambda" expressions.

package goaldi

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var _ = fmt.Printf // enable debugging

//  Procedure value
type VProcedure struct {
	Name     string      // registered name
	Pnames   *[]string   // parameter names (nil if unknown)
	Variadic bool        // true if variadic
	RawCall  bool        // true if to use nonstandard raw argument lists
	GdProc   Procedure   // Goaldi-compatible function (possibly a shim)
	GoFunc   interface{} // underlying function
	Descr    string      // optional one-line description (used for stdlib)
}

//  DefProc constructs a procedure from a Goaldi function and a description.
func DefProc(entry Procedure, name string, pspec string, descr string) *VProcedure {
	pnames, isvar := ParmsFromSpec(pspec)
	return NewProcedure(name, pnames, isvar, entry, entry, descr)
}

//  ParmsFromSpec turns a parameter spec into a pnames list and variadic flag
func ParmsFromSpec(pspec string) (*[]string, bool) {
	isvariadic := strings.HasSuffix(pspec, "[]")
	if isvariadic {
		pspec = strings.TrimSuffix(pspec, "[]")
	}
	pnames := strings.Split(pspec, ",")
	return &pnames, isvariadic
}

//  NewProcedure -- construct a procedure value
//  The result is variadic only if allowvar is true *and* entry is variadic.
func NewProcedure(name string, pnames *[]string, allowvar bool,
	entry Procedure, ufunc interface{}, descr string) *VProcedure {
	isvar := allowvar && reflect.TypeOf(entry).IsVariadic()
	return &VProcedure{name, pnames, isvar, false, entry, ufunc, descr}
}

//  ProcedureType is the procedure instance of type type.
var ProcedureType = NewType(rProc, ProcCtor,
	"procedure", "p", "succeed if procedure")

//  VProcedure.String -- default conversion to Go string returns "P:procname"
func (v *VProcedure) String() string {
	return "P:" + v.Name
}

//  VProcedure.GoString -- convert to string for image() and printf("%#v")
func (v *VProcedure) GoString() string {
	s := "procedure " + v.Name + "("
	if v.Pnames == nil {
		return s + "?)" // params unknown
	}
	d := ""
	for _, t := range *v.Pnames {
		s = s + d + t
		d = ","
	}
	return s + ")"
}

//  VProcedure.ImplBy -- return name of implementing underlying function
func (v *VProcedure) ImplBy() string {
	if v.GoFunc == nil {
		return v.Name // no further information available
	} else {
		return runtime.FuncForPC(reflect.ValueOf(v.GoFunc).Pointer()).Name()
	}
}

//  VProcedure.Type -- return the procedure type
func (v *VProcedure) Type() IRanking {
	return ProcedureType
}

//  VProcedure.Copy returns itself
func (v *VProcedure) Copy() Value {
	return v
}

//  VProcedure.Import returns itself
func (v *VProcedure) Import() Value {
	return v
}

//  VProcedure.Export returns the underlying function
//  (#%#% at least for now. should we wrap it somehow?)
func (v *VProcedure) Export() interface{} {
	return v.GdProc
}

//  VProcedure.Call invokes a procedure
func (v *VProcedure) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	if v.RawCall {
		f := v.GoFunc.(func(*Env, []Value, []string) (Value, *Closure))
		return f(env, args, names)
	} else {
		args = ArgNames(v, args, names)
		return v.GdProc(env, args...)
	}
}

//  Declare methods
var ProcedureMethods = MethodTable([]*VProcedure{
	DefMeth((*VProcedure).Type, "type", "", "return procedure type"),
	DefMeth((*VProcedure).Copy, "copy", "", "return procedure value"),
	DefMeth((*VProcedure).String, "string", "", "return short string"),
	DefMeth((*VProcedure).GoString, "image", "", "return image string"),
})

//  VProcedure.Field implements methods
func (v *VProcedure) Field(f string) Value {
	return GetMethod(ProcedureMethods, v, f)
}

//  The "constructor" returns its argument if procedure and otherwise fails.
//  Note that this is tricky to access, but you can say:  type(main)(p).
func ProcCtor(env *Env, args ...Value) (Value, *Closure) {
	x := ProcArg(args, 0, NilValue)
	if p, ok := x.(*VProcedure); ok {
		return Return(p)
	} else {
		return Fail()
	}
}

//  Go methods already converted to Goaldi procedures
var KnownMethods = make(map[uintptr]*VProcedure)

//  ImportMethod(val, name, meth) -- construct a Goaldi method from a Go method.
//  meth is a method struct such as returned by reflect.Type.MethodByName(),
//  not a bound method value such as returned by reflect.Value.MethodByName().
func ImportMethod(val Value, name string, meth reflect.Method) Value {
	addr := meth.Func.Pointer()
	p := KnownMethods[addr]
	if p == nil {
		gofunc := meth.Func.Interface()
		proc := GoShim(name, gofunc)
		p = NewProcedure(name, nil, true, proc, gofunc, "")
		KnownMethods[addr] = p
	}
	return MethodVal(p, Deref(val))
}

//  GoProcedure(name, func) -- construct a procedure from a Go function
func GoProcedure(name string, f interface{}) *VProcedure {
	return NewProcedure(name, nil, true, GoShim(name, f), f, "")
}

//  GoShim(name, func) -- make a shim for converting args to a Go function
func GoShim(name string, f interface{} /*func*/) Procedure {

	//  get information about the Go function
	ftype := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)
	if fval.Kind() != reflect.Func {
		panic(NewExn("Not a func", f))
	}
	nargs := ftype.NumIn()
	nfixed := nargs
	if ftype.IsVariadic() {
		nfixed--
	}
	nrtn := ftype.NumOut()

	//  make an array of conversion functions, one per parameter
	passer := make([]func(Value) reflect.Value, nargs)
	for i := 0; i < nfixed; i++ {
		passer[i] = passfunc(ftype.In(i))
	}
	if nfixed < nargs { // if variadic
		passer[nfixed] = passfunc(ftype.In(nfixed).Elem())
	}

	// create a function that converts arguments and calls the underlying func
	return func(env *Env, args ...Value) (Value, *Closure) {
		//  set up traceback recovery
		defer Traceback(name, args)
		//  convert fixed arguments from Goaldi values to needed Go type
		in := make([]reflect.Value, 0, len(args))
		var v reflect.Value
		for i := 0; i < nfixed; i++ {
			a := NilValue
			if i < len(args) {
				a = args[i]
			}
			v = passer[i](a)
			if !v.IsValid() {
				panic(NewExn("Cannot convert argument", args[i]))
			}
			in = append(in, v)
		}
		//  convert additional variadic arguments to final type
		if nfixed < nargs {
			for i := nfixed; i < len(args); i++ {
				v = passer[nfixed](args[i])
				if !v.IsValid() {
					panic(NewExn("Cannot convert argument", args[i]))
				}
				in = append(in, v)
			}
		}
		//  call the Go function
		out := fval.Call(in)
		if nrtn == 0 {
			return Return(NilValue) // no return value: return %nil
		}
		r := Import(out[0].Interface()) // import the first return value
		if r == NilValue && nrtn == 2 { // if result is nil and there's one more
			if e, ok := out[1].Interface().(error); ok && e != nil { // if error
				return Fail() // then fail
			}
		}
		return Return(r) // return first value
	}
}

//  passfunc returns a function that converts a Goaldi value
//  into a Go reflect.Value of the specified type
func passfunc(t reflect.Type) func(Value) reflect.Value {
	k := t.Kind()
	switch k {
	case reflect.Bool:
		return func(v Value) reflect.Value {
			var b bool
			switch x := v.(type) {
			case bool:
				b = x
			case vnil:
				b = false
			case *VNumber:
				b = (x.Val() != 0)
			default:
				b = true
			}
			return reflect.ValueOf(b)
		}
	case reflect.Interface: // #%#% this assumes interface{}; should check
		// use default conversion
		break
	default:
		// check if convertible from numeric
		if reflect.TypeOf(1.0).ConvertibleTo(t) {
			return func(v Value) reflect.Value {
				if reflect.TypeOf(v).ConvertibleTo(t) {
					return reflect.ValueOf(v).Convert(t)
				} else {
					return reflect.ValueOf(
						v.(Numerable).ToNumber().Val()).Convert(t)
				}
			}
		}
		// otherwise, check if convertible from string
		if reflect.TypeOf("abc").ConvertibleTo(t) {
			return func(v Value) reflect.Value {
				if reflect.TypeOf(v).ConvertibleTo(t) {
					return reflect.ValueOf(v).Convert(t)
				} else {
					return reflect.ValueOf(
						v.(Stringable).ToString().ToUTF8()).Convert(t)
				}
			}
		}
		// otherwise, use default conversion
		break
	}
	// default conversion
	return func(v Value) reflect.Value {
		var inil interface{}
		x := Export(v) // default conversion
		if x == nil {
			return reflect.ValueOf(&inil).Elem() // nil is tricky
		} else {
			return reflect.ValueOf(x) // anything else
		}
	}
}

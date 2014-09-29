//  interp.go -- the interpreter main loop

package main

import (
	"fmt"
	g "goaldi"
)

//  procedure frame
type pr_frame struct {
	env    *g.Env             // dynamic execution enviromment
	info   *pr_Info           // static procedure information
	params []g.Value          // parameters
	locals []g.Value          // locals
	temps  map[string]g.Value // temporaries
	coord  *ir_coordinate     // last known source location
	offv   g.Value            // offending value for traceback
}

//  catchf -- annotate a panic value with procedure frame information
func catchf(p interface{}, f *pr_frame, args []g.Value) *g.CallFrame {
	return g.Catch(p, f.offv, f.coord.File, f.coord.Line, f.info.name, args)
}

//  interp -- interpret one procedure
func interp(env *g.Env, pr *pr_Info, args ...g.Value) (g.Value, *g.Closure) {

	if opt_trace {
		fmt.Printf("P: %s\n", pr.name)
	}

	// initialize procedure frame: params, locals, temps
	var f pr_frame
	f.env = env
	f.info = pr
	f.temps = make(map[string]g.Value)
	f.params = make([]g.Value, pr.nparams, pr.nparams)
	f.locals = make([]g.Value, pr.nlocals, pr.nlocals)
	for i := 0; i < len(f.params); i++ {
		if i < len(args) {
			f.params[i] = args[i]
		} else {
			f.params[i] = g.NewNil()
		}
	}
	for i := 0; i < len(f.locals); i++ {
		f.locals[i] = g.NewNil()
	}

	// set up tracback recovery
	defer func() {
		if p := recover(); p != nil {
			panic(catchf(p, &f, args))
		}
	}()

	// interpret the IR code
	label := pr.ir.CodeStart.Value
	for {
		if opt_trace {
			fmt.Printf("L: %s\n", label)
		}
		ilist := pr.insns[label]     // look up label
		for _, insn := range ilist { // execute insns in chunk
			if opt_trace {
				fmt.Printf("I: %T %v\n", insn, insn)
			}
			f.coord = nil //#%#% prudent, but should not be needed
			f.offv = nil  //#%#% prudent, but should not be needed
			switch i := insn.(type) {
			default:
				panic(i) // unrecognized or unimplemented
			case ir_Fail:
				return nil, nil
			case ir_IntLit:
				f.temps[i.Lhs.Name] =
					g.NewString(i.Val).ToNumber()
			case ir_RealLit:
				f.temps[i.Lhs.Name] =
					g.NewString(i.Val).ToNumber()
			case ir_StrLit:
				f.temps[i.Lhs.Name] =
					g.NewString(i.Val)
			case ir_Var:
				v := pr.dict[i.Name]
				switch t := v.(type) {
				case pr_local:
					v = g.Trapped(&f.locals[int(t)])
				case pr_param:
					v = g.Trapped(&f.params[int(t)])
				case nil:
					panic("nil in ir_Var; undeclared?")
				default:
					// global or static: already trapped
				}
				f.temps[i.Lhs.Name] = v
			case ir_OpFunction:
				f.coord = i.Coord
				argl := getArgs(&f, i.ArgList)
				f.offv = argl[0]
				v, c := opFunc(&f, i.Fn, argl)
				if v == nil && i.FailLabel.Value != "" {
					label = i.FailLabel.Value
					break
				}
				f.temps[i.Lhs.Name] = v
				f.temps[i.Lhsclosure.Name] = c
			case ir_Call:
				f.coord = i.Coord
				proc := f.temps[i.Fn.Name]
				argl := getArgs(&f, i.ArgList)
				f.offv = proc
				v, c := proc.(g.ICall).Call(env, argl...)
				if v == nil && i.FailLabel.Value != "" {
					label = i.FailLabel.Value
					break
				}
				f.temps[i.Lhs.Name] = v
				f.temps[i.Lhsclosure.Name] = c
			}
		}
	}
	return nil, nil
}

//  getArgs -- load values from heterogeneous ArgList slice field
func getArgs(f *pr_frame, arglist []interface{}) []g.Value {
	n := len(arglist)
	argl := make([]g.Value, n, n)
	for i, a := range arglist {
		switch t := a.(type) {
		case ir_Tmp:
			argl[i] = f.temps[t.Name]
		default:
			argl[i] = g.Deref(a)
		}
	}
	return argl
}

//  opFunc -- implement operator function
func opFunc(f *pr_frame, o *ir_operator, a []g.Value) (g.Value, *g.Closure) {
	op := o.Arity + o.Name
	switch op {
	default:
		panic("unimplemented operator: " + op)
	case "2+":
		return a[0].(g.IAdd).Add(a[1]), nil
	}
}
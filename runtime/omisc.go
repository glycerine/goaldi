//  omisc.go -- miscellaneous runtime operations

package runtime

import (
	"reflect"
)

//  Identical(a,b) implements the === operator.
//  NotIdentical(a,b) implements the ~=== operator.
//  Both call a.Identical(b) if implemented (interface IIdentical).
func Identical(a, b Value) Value {
	if _, ok := a.(IIdentical); ok {
		return a.(IIdentical).Identical(b)
	}
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if av.Type() != bv.Type() {
		return nil
	}
	same := false
	switch av.Kind() {
	default:
		same = (a == b)
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		same = (av.Pointer() == bv.Pointer())
	}
	if same {
		return b
	} else {
		return nil
	}
}

func NotIdentical(a, b Value) Value {
	if Identical(b, a) != nil {
		return nil
	} else {
		return b
	}
}

//  Size(x) calls x.Size() or falls back to calling len().
//  It panics on an inappropriate argument type.
func Size(x Value) Value {
	if t, ok := x.(ISize); ok {
		return t.Size()
	} else {
		return NewNumber(float64(reflect.ValueOf(x).Len()))
	}
}

//  Take(x) calls x.Take() or uses reflection for an arbitrary map or channel.
//  It panics on an inappropriate argument type.
func Take(x Value) Value {
	if t, ok := x.(ITake); ok {
		return t.Take()
	}
	k := reflect.ValueOf(x).Kind()
	if k == reflect.Chan {
		return TakeChan(x)
	} else {
		return x.(ITake).Take() // provoke panic
	}
}

//  VNumber.Call -- implement i(e1, e2, e3...)
func (v *VNumber) Call(env *Env, args []Value, names []string) (Value, *Closure) {
	if len(names) > 0 {
		panic(NewExn("Named arguments not allowed", v))
	}
	i := GoIndex(int(v.Val()), len(args))
	if i < len(args) {
		return Return(args[i])
	} else {
		return Fail()
	}
}

//  ToBy -- implement "e1 to e2 by e3"
func ToBy(e1 Value, e2 Value, e3 Value) (Value, *Closure) {
	n1 := e1.(Numerable).ToNumber()
	n2 := e2.(Numerable).ToNumber()
	n3 := e3.(Numerable).ToNumber()
	if *n3 == 0 {
		panic(NewExn("ToBy: bad increment", e3))
	}
	v1 := *n1
	v2 := *n2
	v3 := *n3
	v1 -= v3
	var f *Closure
	f = &Closure{func() (Value, *Closure) {
		v1 += v3
		if (v3 > 0 && v1 <= v2) || (v3 < 0 && v1 >= v2) {
			return NewNumber(float64(v1)), f
		} else {
			return Fail()
		}
	}}
	return f.Resume()
}
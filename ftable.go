//  ftable.go -- table functions and methods

package goaldi

import (
	"math/rand"
	"reflect"
)

//  Declare methods
var TableMethods = MethodTable([]*VProcedure{
	DefMeth("type", VTable.Type, []string{}, "return table type"),
	DefMeth("copy", VTable.Copy, []string{}, "duplicate table"),
	DefMeth("string", VTable.String, []string{}, "return short string"),
	DefMeth("image", VTable.GoString, []string{}, "return string image"),
	DefMeth("member", VTable.Member, []string{"x"}, "test membership"),
	DefMeth("delete", VTable.Delete, []string{"x"}, "remove entry"),
	DefMeth("sort", VTable.Sort, []string{}, "produce sorted list"),
})

//  VTable.Field implements method calls
func (m VTable) Field(f string) Value {
	return GetMethod(TableMethods, m, f)
}

//  init() declares the constructor function
func init() {
	// Goaldi procedures
	LibProcedure("table", Table)
}

//  Declare methods on Go Tables
var GoMapMethods = MethodTable([]*VProcedure{
	DefMeth("member", GoMapMember, []string{"x"}, "test membership"),
	DefMeth("delete", GoMapDelete, []string{"x"}, "remove entry"),
	DefMeth("sort", GoMapSort, []string{}, "produce sorted list"),
})

//  Table() returns a new table
func Table(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("table", args)
	return Return(NewTable())
}

//  VTable.Member(k) succeeds if k is an existing key
func (m VTable) Member(args ...Value) (Value, *Closure) {
	return GoMapMember(m, args...)
}

//  GoMapMember(m, k) succeeds if k is an existing key in m
func GoMapMember(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.member", args)
	key := ProcArg(args, 0, NilValue)
	if TrapMap(m, key).Exists() {
		return Return(key)
	} else {
		return Fail()
	}
}

//  VTable.Delete(k) deletes the entry, if any, with key k
func (m VTable) Delete(args ...Value) (Value, *Closure) {
	return GoMapDelete(m, args...)
}

//  GoMapDelete(m, k) deletes the entry in Go map m, if any, with key k
func GoMapDelete(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.delete", args)
	key := ProcArg(args, 0, NilValue)
	TrapMap(m, key).Delete()
	return Return(m)
}

//  VTable.Sort(i) produces [:!T:].sort(i)
func (m VTable) Sort(args ...Value) (Value, *Closure) {
	return GoMapSort(m, args...)
}

//  GoMapSort(m, i) produces [:!T:].sort(i)
func GoMapSort(m Value, args ...Value) (Value, *Closure) {
	defer Traceback("T.sort", args)
	i := ProcArg(args, 0, ONE).(Numerable).ToNumber()
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	vlist := make([]Value, mv.Len())
	for i, kv := range klist {
		k := Import(kv.Interface())
		v := Import(mv.MapIndex(kv).Interface())
		vlist[i] = kvRecord.New([]Value{k, v})
	}
	return InitList(vlist).Sort(i)
}

//  -------------------------- key/value pairs ---------------------

//  kvRecord defines the {key,value} struct returned by ?T and !T
var kvRecord = NewDefn("tableElem", []string{"key", "value"})

//  ChooseMap returns a key/value pair from any Goaldi table or Go map
func ChooseMap(m interface{} /*anymap*/) Value {
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	n := len(klist)
	if n == 0 { // if map empty
		return nil // fail
	}
	i := rand.Intn(n)
	k := Import(klist[i].Interface())
	v := Import(mv.MapIndex(klist[i]).Interface())
	return kvRecord.New([]Value{k, v})
}

//  DispenseMap generates key/value pairs for any Goaldi table or Go map
func DispenseMap(m interface{} /*anymap*/) (Value, *Closure) {
	mv := reflect.ValueOf(m)
	klist := mv.MapKeys()
	i := -1
	var c *Closure
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(klist) {
			k := Import(klist[i].Interface())
			v := Import(mv.MapIndex(klist[i]).Interface())
			return kvRecord.New([]Value{k, v}), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}
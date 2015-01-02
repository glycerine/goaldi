//  vdefn.go -- struct definition (constructor) information
//
//  Defines the interpretation of a vstruct object that points to it,
//  and constructs objects for it.

package goaldi

type VDefn struct {
	Name    string                 // type name
	Flist   []string               // ordered list of field names
	Methods map[string]*VProcedure // method list
	//#%#% could add a hash map for the fields ... but is it worth it?
}

//  NewDefn(name, fields) -- construct new definition
func NewDefn(name string, fields []string) *VDefn {
	return &VDefn{name, fields, make(map[string]*VProcedure)}
}

//  AddMethod(name, procedure) -- add a method for this struct type
//  Returns false if rejected as a duplicate.
func (v *VDefn) AddMethod(name string, p *VProcedure) bool {
	for _, s := range v.Flist {
		if s == name {
			return false
		}
	}
	if v.Methods[name] != nil {
		return false
	}
	v.Methods[name] = p
	return true
}

//  VDefn.New(values) -- create a new underlying struct object
func (v *VDefn) New(a []Value) *VStruct {
	r := &VStruct{v, make([]Value, len(v.Flist))}
	for i := range r.Data {
		if i < len(a) {
			r.Data[i] = a[i]
		} else {
			r.Data[i] = NilValue
		}
	}
	return r
}

//  VDefn.String -- conversion to Go string returns "C:name"
func (v *VDefn) String() string {
	return "C:" + v.Name
}

//  VDefn.GoString -- convert to Go string for image() and printf("%#v")
func (v *VDefn) GoString() string {
	s := "constructor " + v.Name + "("
	d := ""
	for _, t := range v.Flist {
		s = s + d + t
		d = ","
	}
	return s + ")"
}

//  VDefn.Rank returns rDefn
func (v *VDefn) Rank() int {
	return rDefn
}

//  VDefn.Type returns "constructor"
func (v *VDefn) Type() Value {
	return type_constructor
}

var type_constructor = NewString("constructor")

//  VDefn.Copy returns itself
func (v *VDefn) Copy() Value {
	return v
}

//  VDefn.Import returns itself
func (v *VDefn) Import() Value {
	return v
}

//  VDefn.Export returns itself
func (v *VDefn) Export() interface{} {
	return v
}

//  VDefn.Dispense() implements !D to generate the field names
func (v *VDefn) Dispense(unused IVariable) (Value, *Closure) {
	var c *Closure
	i := -1
	c = &Closure{func() (Value, *Closure) {
		i++
		if i < len(v.Flist) {
			return NewString(v.Flist[i]), c
		} else {
			return Fail()
		}
	}}
	return c.Resume()
}

//  VDefn.Call() implements a struct constructor
func (v *VDefn) Call(env *Env, args ...Value) (Value, *Closure) {
	return Return(v.New(args))
}

//  Declare required methods of the constructor (not the underlying type)
var DefnMethods = map[string]interface{}{
	"type":   (*VDefn).Type,
	"copy":   (*VDefn).Copy,
	"string": (*VDefn).String,
	"image":  (*VDefn).GoString,
}

//  VDefn.Field implements methods called *on the constructor*
func (v *VDefn) Field(f string) Value {
	return GetMethod(DefnMethods, v, f)
}

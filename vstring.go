//  vstring.go -- VString, the Goaldi type "string"
//
//  Strings contain sequences of Unicode characters (Code Points or Runes)
//

package goaldi

import (
	"fmt"
	"unicode/utf8"
)

// predefined constants
var (
	EMPTY = NewString("") // the empty string
)

//  A string is encoded by one (usually) or two parallel slices
type VString struct {
	low  []uint8  // required: low-order 8 bits of each rune
	high []uint16 // optional: high-order 13 bits of each rune
}

//  NewString -- construct a Goaldi string from a Go UTF8 string
func NewString(s string) *VString {
	n := len(s)
	low := make([]uint8, n, n)
	high := make([]uint16, n, n)
	h := '\000'
	i := 0
	for _, c := range s {
		h |= c
		low[i] = uint8(c)
		high[i] = uint16(c >> 8)
		i++
	}
	// #%#% could copy now to smaller underlying arrays if warranted
	if (h >> 8) == 0 {
		return &VString{low[:i], nil}
	} else {
		return &VString{low[:i], high[:i]}
	}
}

//  BinaryString -- construct a Goaldi string from Go Latin1 bytes
func BinaryString(s []byte) *VString {
	b := make([]uint8, len(s), len(s))
	copy(b, s)
	return &VString{b, nil}
}

//  VString.ToUTF8 -- convert Goaldi Unicode string to Go UTF8 string
func (v *VString) ToUTF8() string {
	b := make([]byte, 0, len(v.low))
	p := make([]byte, 8, 8)
	for i, c := range v.low {
		r := rune(c)
		if v.high != nil {
			r |= rune(v.high[i]) << 8
		}
		n := utf8.EncodeRune(p, r)
		b = append(b, p[:n]...)
	}
	return string(b)
}

//  VString.ToBinary -- convert Goaldi Unicode to 8-bit bytes by truncation
func (v *VString) ToBinary() []byte {
	return []byte(v.low)
}

//  VString.String -- return image of string, quoted, as a Go string
func (v *VString) String() string {
	return `"` + v.ToUTF8() + `"`
}

//  VString.ToString -- for a Goaldi string, this just returns self
func (v *VString) ToString() *VString {
	return v
}

//  VString.ToNumber -- return conversion to VNumber, or issue RunErr
func (v *VString) ToNumber() *VNumber {
	n := v.TryNumber()
	if n == nil {
		panic(&RunErr{"Cannot convert to number", v})
	}
	return n
}

//  VString.TryNumber -- return conversion to VNumber, or nil for failure
//#%#% this could probably be improved and should allow trailing spaces
func (v *VString) TryNumber() *VNumber {
	var f float64
	var b byte
	if v.high != nil { // if has exotic characters //#%#% bogus test?
		return nil // it can't be valid
	}
	n, _ := fmt.Sscanf(string(v.low), "%f%c", &f, &b)
	if n == 1 {
		return NewNumber(f)
	} else {
		return nil
	}
}

//  VString.Type -- return "string"
func (v *VString) Type() Value {
	return type_string
}

var type_string = NewString("string")

//  VString.Identical -- check equality for === operator
func (s *VString) Identical(x Value) Value {
	t, ok := x.(*VString)
	if !ok {
		return nil
	} else if s == t {
		return t
	} else {
		return s.StrEQ(t)
	}
}

//  VString.Export returns a Go string
func (v *VString) Export() interface{} {
	return v.ToUTF8()
}

//  -------------------------- trapped substrings ---------------------

type vSubStr struct {
	target IVariable // pointer to target
	i, j   int       // original subscripts
}

//  vSubStr.Deref() -- extract value of for use as an rvalue
func (ss *vSubStr) Deref() Value {
	return Deref(ss.target).(*VString).slice(nil, ss.i, ss.j)
}

//  vSubStr.String() -- show string representation: produces v[i:j]
func (ss *vSubStr) String() string {
	return fmt.Sprintf("(&%v[%d:%d])", ss.target, ss.i, ss.j)
}

//  vSubStr.Assign -- store value in target variable
func (ss *vSubStr) Assign(v Value) IVariable {
	src := Deref(ss.target).(*VString)
	ins := v.(Stringable).ToString()
	//#%#% check that i & j are still valid?
	snew := scat(src, 0, ss.i, ins, 0, ins.length(), src, ss.j, src.length())
	ss.target = ss.target.Assign(snew)
	ss.j = ss.i + ins.length()
	return ss
}

//  -------------------------- internal functions ---------------------

//  VString.length -- return string length as int
func (s *VString) length() int {
	return len(s.low)
}

//  VString.slice -- return substring given Go-style zero-based limits
//  If lval is non-null, generates a trapped slice reference.
func (s *VString) slice(lval IVariable, i int, j int) Value {
	if lval != nil {
		return &vSubStr{lval, i, j} // produce variable
	}
	// produce value
	r := &VString{s.low[i:j], nil}
	if s.high != nil && j > i {
		r.high = s.high[i:j]
		//#%#% remove if nothing there (all zeroes) ?
	}
	return r
}

//  VString.compare -- compare two strings, return <0, 0, or >0
func (s *VString) compare(t *VString) int {
	// check for easy case
	if s == t {
		return 0
	}
	// extract fields
	sl := s.low
	tl := t.low
	sh := s.high
	th := t.high
	sn := len(sl)
	tn := len(tl)
	// compare runes until one differs
	for i := 0; i < sn && i < tn; i++ {
		sr := rune(sl[i])
		tr := rune(tl[i])
		if sh != nil {
			sr |= rune(sh[i] << 8)
		}
		if th != nil {
			tr |= rune(th[i] << 8)
		}
		if sr != tr {
			return int(sr) - int(tr)
		}
	}
	// reached the end of one or both strings
	return sn - tn
}

//  scat -- general string concatenator.
//  produces x1[i1:j1] || s2[i2:j2] || s3[i3:j3]  (using Go indexing).
//  all arguments are assumed valid.
func scat(s1 *VString, i1, j1 int, s2 *VString, i2, j2 int,
	s3 *VString, i3, j3 int) *VString {
	n1 := j1 - i1
	n2 := j2 - i2
	n3 := j3 - i3
	nt := n1 + n2 + n3
	low := make([]uint8, nt, nt)
	copy(low[0:], s1.low[i1:j1])
	copy(low[n1:], s2.low[i2:j2])
	copy(low[n1+n2:], s3.low[i3:j3])
	if s1.high == nil && s2.high == nil && s3.high == nil {
		return &VString{low, nil}
	}
	high := make([]uint16, nt, nt)
	if s1.high != nil {
		copy(high[0:], s1.high[i1:j1])
	}
	if s2.high != nil {
		copy(high[n1:], s2.high[i2:j2])
	}
	if s3.high != nil {
		copy(high[n1+n2:], s3.high[i3:j3])
	}
	//#%#% could check here if "high" is really needed
	return &VString{low, high}
}

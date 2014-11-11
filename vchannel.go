//  vchannel.go -- VChannel, the Goaldi type "channel"

package goaldi

import (
	"fmt"
	"reflect"
)

//  VChannel implements a Goaldi channel, which just wraps a Go channel.
type VChannel chan Value

//  NewChannel -- construct a new Goaldi channel
func NewChannel(i int) VChannel {
	return VChannel(make(chan Value, i))
}

//  VChannel.String -- default conversion to Go string returns "M:size"
func (c VChannel) String() string {
	return fmt.Sprintf("CH:%d", cap(c))
}

//  VChannel.GoString -- convert to Go string for image() and printf("%#v")
func (c VChannel) GoString() string {
	return fmt.Sprintf("channel(%d)", cap(c))
}

//  VChannel.Rank -- return rChannel
func (v VChannel) Rank() int {
	return rChannel
}

//  VChannel.Type -- return "channel"
func (c VChannel) Type() Value {
	return type_channel
}

var type_channel = NewString("channel")

//  VChannel.Copy returns itself
func (c VChannel) Copy() Value {
	return c
}

//  VChannel.Identical checks equality for the === operator
func (c VChannel) Identical(x Value) Value {
	c2, ok := x.(VChannel)
	if ok && reflect.ValueOf(c).Pointer() == reflect.ValueOf(c2).Pointer() {
		return x
	} else {
		return nil
	}
}

//  VChannel.Import returns itself
func (v VChannel) Import() Value {
	return v
}

//  VChannel.Export returns itself.
func (v VChannel) Export() interface{} {
	return v
}

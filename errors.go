package rpsl

import (
	"fmt"
	"reflect"
)

// UnmarshalBinaryErr describes an invalid argument passed to rpsl.UnmarshalBinary.
// The argument to rpsl.UnmarshalBinary must be a non-nil pointer.
type UnmarshalBinaryErr struct {
	Type reflect.Type
}

// Error returns a string representation of the error.
func (e *UnmarshalBinaryErr) Error() string {
	if e.Type == nil {
		return "rpsl: UnmarshalBinary(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "rpsl: UnmarshalBinary(non-pointer " + e.Type.String() + ")"
	}
	return fmt.Sprintf("rpsl: UnmarshalBinary(%s %s)", e.Type.Kind(), e.Type.String())
}

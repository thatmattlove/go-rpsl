package rpsl

import (
	"go.mdl.wtf/rpsl/internal/serialize"
)

// UnmarshalBinaryErr describes an invalid argument passed to rpsl.UnmarshalBinary.
// The argument to rpsl.UnmarshalBinary must be a non-nil pointer.
type UnmarshalBinaryErr = serialize.UnmarshalBinaryErr

package rpsl

import (
	"go.mdl.wtf/rpsl/internal/serialize"
)

// MarshalBinary encodes an RPSL data structure as a byte string. The argument must be a pointer to
// a struct.
//
// Example:
//
//	b, err := rpsl.MarshalBinary(&route)
//	fmt.Println(string(b))
func MarshalBinary(o any) ([]byte, error) {
	s, err := serialize.RPSL(o)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

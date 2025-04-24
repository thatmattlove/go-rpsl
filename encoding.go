package rpsl

import "go.mdl.wtf/rpsl/internal/serialize"

// MarshalBinary encodes an RPSL data structure as a byte string. The argument must be a pointer to
// a struct.
//
// Example:
//
//	b, err := rpsl.MarshalBinary(&route)
//	fmt.Println(string(b))
func MarshalBinary(o any) ([]byte, error) {
	return serialize.Encode(o)
}

// UnmarshalBinary decodes a byte string of RPSL data to a Go RPSL object.
// The second argument must be a pointer to an RPSL struct.
//
// Example:
//
//	b := []byte(`aut-num: 65000`)
//	var autNum rpsl.AutNum
//	err := rpsl.UnmarshalBinary(b, &autNum)
func UnmarshalBinary(b []byte, o any) error {
	return serialize.Decode(b, o)
}

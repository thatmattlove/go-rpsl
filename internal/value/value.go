package value

import (
	"bytes"

	"go.mdl.wtf/rpsl/internal/serialize"
)

type V string

func (v V) String() string {
	return string(v)
}

type VCommaSpace []V

func (vc VCommaSpace) String() string {
	return serialize.JoinStringers(vc, ", ")
}

func (vc VCommaSpace) UnmarshalBinary(b []byte) (VCommaSpace, error) {
	parts := bytes.Split(b, []byte{0x2c, 0x20})
	out := make(VCommaSpace, len(parts))
	for i := range parts {
		out[i] = V(string(parts[i]))
	}
	vc = out
	return vc, nil
}

type VComma []V

func (vc VComma) String() string {
	return serialize.JoinStringers(vc, ",")
}

func (vc VComma) UnmarshalBinary(b []byte) (VComma, error) {
	parts := bytes.Split(b, []byte{0x2c})
	out := make(VComma, len(parts))
	for i := range parts {
		out[i] = V(string(parts[i]))
	}
	vc = out
	return vc, nil
}

type VNewline []V

func (vn VNewline) String() string {
	return serialize.JoinStringers(vn, "\n")
}

func (vn VNewline) UnmarshalBinary(b []byte) (VNewline, error) {
	parts := bytes.Split(b, []byte{0xa})
	partsV := make(VNewline, len(parts))
	for i := range parts {
		partsV[i] = V(string(parts[i]))
	}
	if len(vn) == 0 {
		vn = append(vn, partsV...)
		return vn, nil
	}
	vn = append(vn, partsV...)
	return vn, nil
}

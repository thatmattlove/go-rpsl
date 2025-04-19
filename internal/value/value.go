package value

import "go.mdl.wtf/rpsl/internal/serialize"

type V string

func (v V) String() string {
	return string(v)
}

type VCommaSpace []V

func (vc VCommaSpace) String() string {
	return serialize.JoinStringers(vc, ", ")
}

type VComma []V

func (vc VComma) String() string {
	return serialize.JoinStringers(vc, ",")
}

type VNewline []V

func (vn VNewline) String() string {
	return serialize.JoinStringers(vn, "\n")
}

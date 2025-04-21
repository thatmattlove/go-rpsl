package serialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mdl.wtf/rpsl/internal/serialize"
)

type S string

func (s S) String() string {
	return string(s)
}

func AsAny(a any) any {
	return a
}

func Test_StringsStringers(t *testing.T) {
	asStringers := []S{S("a"), S("b"), S("c")}
	asStrings := serialize.StringsStringers(asStringers)
	assert.True(t, len(asStringers) == len(asStrings))
	_, isStringSlice := AsAny(asStrings).([]string)
	assert.True(t, isStringSlice)
}

func Test_JoinStringers(t *testing.T) {
	asStringers := []S{S("a"), S("b"), S("c")}
	exp := "a--b--c"
	assert.Equal(t, exp, serialize.JoinStringers(asStringers, "--"))
}

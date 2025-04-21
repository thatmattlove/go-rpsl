package rpsl_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mdl.wtf/rpsl"
)

func Test_UnmarshalBinaryErr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		err := &rpsl.UnmarshalBinaryErr{reflect.TypeOf(nil)}
		assert.ErrorContains(t, err, "rpsl: UnmarshalBinary(nil)")
	})
	t.Run("non-pointer", func(t *testing.T) {
		t.Parallel()
		err := &rpsl.UnmarshalBinaryErr{reflect.TypeOf("")}
		assert.ErrorContains(t, err, "rpsl: UnmarshalBinary(non-pointer string")
	})
	t.Run("non-struct", func(t *testing.T) {
		t.Parallel()
		str := " "
		err := &rpsl.UnmarshalBinaryErr{reflect.TypeOf(&str)}
		assert.ErrorContains(t, err, "rpsl: UnmarshalBinary(ptr *string")
	})
	t.Run("default", func(t *testing.T) {
		t.Parallel()
		type TestStruct struct {
			Field int
		}
		s := TestStruct{}
		err := &rpsl.UnmarshalBinaryErr{reflect.TypeOf(&s)}
		assert.ErrorContains(t, err, "rpsl: UnmarshalBinary(ptr *rpsl_test.TestStruct)")
	})
}

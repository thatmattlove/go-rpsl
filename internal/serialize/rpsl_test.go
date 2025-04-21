package serialize_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl/internal/serialize"
)

type TypeWithRPSLMethod string

func (t TypeWithRPSLMethod) RPSL() string {
	parts := strings.Split(string(t), "\n")
	out := ""
	for i := range parts {
		out += "key: " + parts[i] + "\n"
	}
	return strings.TrimSuffix(out, "\n")
}

func Test_RPSL(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string `rpsl:"key1"`
			Key2 uint32 `rpsl:"key2"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: 65000,
		}
		result, err := serialize.RPSL(s)
		require.NoError(t, err)
		exp := `key1: value1
key2: 65000`
		assert.Equal(t, exp, result)
	})
	t.Run("omitempty", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string `rpsl:"key1"`
			Key2 uint32 `rpsl:"key2"`
			Key3 string `rpsl:"key3,omitempty"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: 65000,
		}
		result, err := serialize.RPSL(s)
		require.NoError(t, err)
		exp := `key1: value1
key2: 65000`
		assert.Equal(t, exp, result)
	})
	t.Run("non-struct", func(t *testing.T) {
		t.Parallel()
		_, err := serialize.RPSL("non-struct")
		assert.ErrorIs(t, err, serialize.ErrMustBeStruct)
	})
	t.Run("invalid tag", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string `rpsl:"key1"`
			Key2 uint32 `rpsl:"key2"`
			Key3 string `rpsl:",key3"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: 65000,
		}
		_, err := serialize.RPSL(s)
		assert.ErrorContains(t, err, "has an invalid")
	})
	t.Run("no tag", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string `rpsl:"key1"`
			Key2 uint32 `rpsl:"key2"`
			Key3 string
		}
		s := &Struct{
			Key1: "value1",
			Key2: 65000,
			Key3: "this should not show up",
		}
		result, err := serialize.RPSL(s)
		require.NoError(t, err)
		exp := `key1: value1
key2: 65000`
		assert.Equal(t, exp, result)
	})
	t.Run("with extra", func(t *testing.T) {
		type Struct struct {
			Key1  string            `rpsl:"key1"`
			Key2  uint32            `rpsl:"key2"`
			Extra map[string]string `rpsl:"-"`
		}
		s := &Struct{
			Key1:  "value1",
			Key2:  65000,
			Extra: map[string]string{"key3": "value3"},
		}
		exp := `key1: value1
key2: 65000
key3: value3`
		result, err := serialize.RPSL(s)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with rpsl method", func(t *testing.T) {
		type Struct struct {
			Key TypeWithRPSLMethod `rpsl:"key"`
		}
		s := &Struct{
			Key: `value1
value2
value3`,
		}
		exp := `key: value1
key: value2
key: value3`
		result, err := serialize.RPSL(s)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

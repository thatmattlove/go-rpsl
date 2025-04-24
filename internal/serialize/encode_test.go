package serialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl/internal/serialize"
)

func Test_Encode(t *testing.T) {
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
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: 65000`)
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
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: 65000`)
		assert.Equal(t, exp, result)
	})
	t.Run("non-struct", func(t *testing.T) {
		t.Parallel()
		_, err := serialize.Encode("non-struct")
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
		_, err := serialize.Encode(s)
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
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: 65000`)
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
		exp := []byte(`key1: value1
key2: 65000
key3: value3`)
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with as multiline string", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string `rpsl:"key1"`
			Key2 string `rpsl:"key2" as:"multiline"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: `value2-1
value2-2`,
		}
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: value2-1
key2: value2-2`)
		assert.Equal(t, exp, result)
	})
	t.Run("with as multiline string slice", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string   `rpsl:"key1"`
			Key2 []string `rpsl:"key2" as:"multiline"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: []string{"value2-1", "value2-2"},
		}
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: value2-1
key2: value2-2`)
		assert.Equal(t, exp, result)
	})
	t.Run("with as comma-space", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string   `rpsl:"key1"`
			Key2 []string `rpsl:"key2" as:"comma-space"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: []string{"value2-1", "value2-2"},
		}
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: value2-1, value2-2`)
		assert.Equal(t, exp, result)
	})
	t.Run("with as comma", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string   `rpsl:"key1"`
			Key2 []string `rpsl:"key2" as:"comma"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: []string{"value2-1", "value2-2"},
		}
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: value2-1,value2-2`)
		assert.Equal(t, exp, result)
	})
	t.Run("with as comma non-slice", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key1 string `rpsl:"key1"`
			Key2 string `rpsl:"key2" as:"comma"`
		}
		s := &Struct{
			Key1: "value1",
			Key2: "value2-1,value2-2",
		}
		result, err := serialize.Encode(s)
		require.NoError(t, err)
		exp := []byte(`key1: value1
key2: value2-1,value2-2`)
		assert.Equal(t, exp, result)
	})
}

package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl/internal/value"
)

func Test_V(t *testing.T) {
	t.Parallel()
	v := value.V(t.Name())
	assert.Equal(t, t.Name(), v.String())
}

func Test_VCommaSpace(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		vc := value.VCommaSpace{value.V("one"), value.V("two"), value.V("three")}
		exp := "one, two, three"
		assert.Equal(t, exp, vc.String())
	})
	t.Run("unmarshal", func(t *testing.T) {
		t.Parallel()
		b := []byte(`one, two, three`)
		exp := value.VCommaSpace{"one", "two", "three"}
		var v value.VCommaSpace
		result, err := v.UnmarshalBinary(b)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

func Test_VComma(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		vc := value.VComma{value.V("one"), value.V("two"), value.V("three")}
		exp := "one,two,three"
		assert.Equal(t, exp, vc.String())
	})
	t.Run("unmarshal", func(t *testing.T) {
		t.Parallel()
		b := []byte(`one,two,three`)
		exp := value.VComma{"one", "two", "three"}
		var v value.VComma
		result, err := v.UnmarshalBinary(b)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

func Test_VNewline(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		vc := value.VNewline{value.V("one"), value.V("two"), value.V("three")}
		exp := `one
two
three`
		assert.Equal(t, exp, vc.String())
	})
	t.Run("unmarshal empty", func(t *testing.T) {
		t.Parallel()
		b := []byte(`one
two
three`)
		exp := value.VNewline{"one", "two", "three"}
		var v value.VNewline
		result, err := v.UnmarshalBinary(b)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("unmarshal filled", func(t *testing.T) {
		t.Parallel()
		b := []byte(`two
three`)
		exp := value.VNewline{"one", "two", "three"}
		v := value.VNewline{"one"}
		result, err := v.UnmarshalBinary(b)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

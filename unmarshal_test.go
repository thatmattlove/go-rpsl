package rpsl_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
	"go.mdl.wtf/rpsl/internal/value"
)

type WillUnmarshalBinaryErr string

func (w WillUnmarshalBinaryErr) UnmarshalBinary([]byte) (*WillUnmarshalBinaryErr, error) {
	return nil, errors.New("an error")
}

func Test_UnmarshalBinary(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
`)
		var autNum rpsl.AutNum
		err := rpsl.UnmarshalBinary(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
	})
	t.Run("with append single line", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
descr: A Description`)
		var autNum rpsl.AutNum
		err := rpsl.UnmarshalBinary(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "A Description", string(autNum.Description))
	})
	t.Run("with append multiline", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
descr: Line 1
descr: Line 2`)
		var autNum rpsl.AutNum
		err := rpsl.UnmarshalBinary(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, "Line 1\nLine 2", string(autNum.Description))
	})
	t.Run("with vcomma", func(t *testing.T) {
		t.Parallel()
		b := []byte(`route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP`)
		var rs rpsl.RouteSet
		err := rpsl.UnmarshalBinary(b, &rs)
		require.NoError(t, err)
		exp := value.VComma{"192.0.2.0/24", "RS-CORP"}
		assert.Equal(t, exp, rs.Members)
	})
	t.Run("with vcommaspace", func(t *testing.T) {
		t.Parallel()
		b := []byte(`aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME`)
		var autNum rpsl.AutNum
		err := rpsl.UnmarshalBinary(b, &autNum)
		require.NoError(t, err)
		exp := value.VCommaSpace{"AS65001", "AS65002", "AS-ACME"}
		assert.Equal(t, exp, autNum.MemberOf)
	})
	t.Run("with vnewline", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS-ACME
members: AS65000
members: AS-65001`)
		var asSet rpsl.ASSet
		err := rpsl.UnmarshalBinary(b, &asSet)
		require.NoError(t, err)
		exp := value.VNewline{"AS65000", "AS-65001"}
		assert.Equal(t, exp, asSet.Members)
	})
	t.Run("with extra", func(t *testing.T) {
		t.Parallel()
		b := []byte(`as-set: AS-ACME
members: AS65000
members: AS-65001
extra1: value1`)
		var asSet rpsl.ASSet
		err := rpsl.UnmarshalBinary(b, &asSet)
		require.NoError(t, err)
		assert.Equal(t, "value1", asSet.Extra["extra1"])

	})
	t.Run("err non ptr", func(t *testing.T) {
		t.Parallel()
		err := rpsl.UnmarshalBinary([]byte(""), struct{}{})
		assert.ErrorContains(t, err, "non-pointer")
	})
	t.Run("err non struct", func(t *testing.T) {
		t.Parallel()
		o := []string{"one"}
		err := rpsl.UnmarshalBinary([]byte{0x1, 0x2}, &o)
		assert.ErrorContains(t, err, "ptr *[]string")
	})
	t.Run("skip malformed pair", func(t *testing.T) {
		t.Parallel()
		b := []byte(`
aut-num: 65000
as-name: AS-ACME-1
extra content`)
		var autNum rpsl.AutNum
		err := rpsl.UnmarshalBinary(b, &autNum)
		require.NoError(t, err)
		assert.Equal(t, rpsl.ASN(65000), autNum.AutNum)
		assert.Equal(t, "AS-ACME-1", autNum.ASName)
	})
	t.Run("skip empty tag", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			Key  string `rpsl:"key"`
			Skip string
		}
		b := []byte(`key: value
skip: skip`)
		var s Struct
		err := rpsl.UnmarshalBinary(b, &s)
		require.NoError(t, err)
		assert.Equal(t, "value", s.Key)
		assert.Zero(t, s.Skip)
	})
	t.Run("err nested unmarshal", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String string                 `rpsl:"string"`
			Error  WillUnmarshalBinaryErr `rpsl:"error"`
		}
		b := []byte(`string: value
error: irrelevant`)
		var s Struct
		err := rpsl.UnmarshalBinary(b, &s)
		assert.ErrorContains(t, err, "an error")
	})
	t.Run("uint32 field", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String string `rpsl:"string"`
			Uint32 uint32 `rpsl:"uint32"`
		}
		b := []byte(`string: value
uint32: 12345`)
		var s Struct
		err := rpsl.UnmarshalBinary(b, &s)
		require.NoError(t, err)
	})
	t.Run("uint32 invalid", func(t *testing.T) {
		t.Parallel()
		type Struct struct {
			String string `rpsl:"string"`
			Uint32 uint32 `rpsl:"uint32"`
		}
		b := []byte(`string: value
uint32: wrong`)
		var s Struct
		err := rpsl.UnmarshalBinary(b, &s)
		assert.ErrorContains(t, err, "value could not be parsed as uint32")
	})
}

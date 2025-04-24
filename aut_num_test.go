package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
)

func Test_ASN(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		a := rpsl.ASN(65000)
		assert.EqualValues(t, "AS65000", a.String())
	})
	t.Run("unmarshal no prefix", func(t *testing.T) {
		t.Parallel()
		b := []byte(`65000`)
		var a rpsl.ASN
		r, err := a.UnmarshalBinary(b)
		require.NoError(t, err)
		assert.EqualValues(t, rpsl.ASN(65000), r)
	})
	t.Run("unmarshal prefix", func(t *testing.T) {
		t.Parallel()
		b := []byte(`AS65000`)
		var a rpsl.ASN
		r, err := a.UnmarshalBinary(b)
		require.NoError(t, err)
		assert.EqualValues(t, uint32(65000), r)
	})
	t.Run("err unmarshal", func(t *testing.T) {
		t.Parallel()
		b := []byte(`abcd`)
		var a rpsl.ASN
		_, err := a.UnmarshalBinary(b)
		require.ErrorContains(t, err, "could not be parsed")
	})
}

func Test_ASNName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "AS65000", rpsl.ASNName(65000))
}

func Test_AutNumMembersOf(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"AS65000", "AS65001", "AS-SET"}, rpsl.AutNumMembersOf([]any{65000, "AS65001", "AS-SET"}))
}

func Test_AutNum(t *testing.T) {
	t.Parallel()
	autNum := rpsl.AutNum{
		AutNum:   rpsl.ASN(65000),
		ASName:   "AS-65000",
		MemberOf: []string{"AS65001", "AS65002", "AS-ACME"},
	}
	t.Run("base", func(t *testing.T) {
		result, err := rpsl.MarshalBinary(&autNum)
		require.NoError(t, err)
		exp := []byte(`aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME`)
		assert.Equal(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "AS65000", autNum.String())
	})
	t.Run("with extra", func(t *testing.T) {
		autNum.Source = "ARIN"
		autNum.AddExtra("extra", "value")
		require.NotNil(t, autNum.Extra)
		assert.Equal(t, "value", autNum.Extra["extra"])
		exp := []byte(`aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME
extra: value
source: ARIN`)
		result, err := rpsl.MarshalBinary(&autNum)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with descr", func(t *testing.T) {
		t.Parallel()
		autNum := rpsl.AutNum{
			AutNum: rpsl.ASN(65000),
			ASName: "AS-65000",
			Source: "ARIN",
			Description: `this is
a description`,
			MemberOf: []string{"AS65001"},
		}
		exp := []byte(`aut-num: AS65000
as-name: AS-65000
descr: this is
descr: a description
member-of: AS65001
source: ARIN`)
		result, err := rpsl.MarshalBinary(&autNum)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

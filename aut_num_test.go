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

func Test_AutNum(t *testing.T) {
	t.Parallel()
	autNum := &rpsl.AutNum{
		AutNum: rpsl.ASN(65000),
		ASName: "AS-65000",
		MemberOf: rpsl.AutNumMembers(
			rpsl.ASNName(65001),
			rpsl.AutNumMember("AS65002"),
			rpsl.ASSetName("AS-ACME"),
		)}
	t.Run("base", func(t *testing.T) {
		result, err := autNum.RPSL()
		require.NoError(t, err)
		exp := `aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME`
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
		exp := `aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME
extra: value
source: ARIN`
		result, err := autNum.RPSL()
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with descr", func(t *testing.T) {
		t.Parallel()
		autNum := &rpsl.AutNum{
			AutNum: rpsl.ASN(65000),
			ASName: "AS-65000",
			Source: "ARIN",
			Description: `this is
a description`,
			MemberOf: rpsl.AutNumMembers(rpsl.ASNName(65001)),
		}
		exp := `aut-num: AS65000
as-name: AS-65000
descr: this is
descr: a description
member-of: AS65001
source: ARIN`
		result, err := autNum.RPSL()
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

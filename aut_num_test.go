package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
)

func Test_ASN(t *testing.T) {
	t.Parallel()
	a := rpsl.ASN(65000)
	testza.AssertEqual(t, "AS65000", a.String())
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
		testza.AssertNoError(t, err)
		exp := `aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME`
		testza.AssertEqual(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		testza.AssertEqual(t, "AS65000", autNum.String())
	})
	t.Run("with extra", func(t *testing.T) {
		autNum.Source = "ARIN"
		autNum.AddExtra("extra", "value")
		testza.AssertNotNil(t, autNum.Extra)
		testza.AssertEqual(t, "value", autNum.Extra["extra"])
		exp := `aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME
extra: value
source: ARIN`
		result, err := autNum.RPSL()
		testza.AssertNoError(t, err)
		testza.AssertEqual(t, exp, result)
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
		testza.AssertNoError(t, err)
		testza.AssertEqual(t, exp, result)
	})
}

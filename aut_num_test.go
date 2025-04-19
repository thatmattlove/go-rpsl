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
	result, err := autNum.RPSL()
	testza.AssertNoError(t, err)
	exp := `aut-num: AS65000
as-name: AS-65000
member-of: AS65001, AS65002, AS-ACME`
	testza.AssertEqual(t, exp, result)
}

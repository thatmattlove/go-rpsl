package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
	"go.mdl.wtf/rpsl/internal/value"
)

func Test_RouteSet(t *testing.T) {
	t.Parallel()
	rs := &rpsl.RouteSet{
		RouteSet: "RS-ACME",
		Members:  rpsl.RSMembers(rpsl.RSMember("192.0.2.0/24"), rpsl.RSMember("RS-CORP")),
	}
	t.Run("base", func(t *testing.T) {
		exp := `route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP`
		result, err := rs.RPSL()
		testza.AssertNoError(t, err)
		testza.AssertEqual(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		testza.AssertEqual(t, "RS-ACME", rs.String())
	})
	t.Run("with extra", func(t *testing.T) {
		rs.AddExtra("extra", "value")
		testza.AssertNotNil(t, rs.Extra)
		testza.AssertEqual(t, "value", rs.Extra["extra"])
		exp := `route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP
extra: value`
		result, err := rs.RPSL()
		testza.AssertNoError(t, err)
		testza.AssertEqual(t, exp, result)
	})
}

func Test_RSSetName(t *testing.T) {
	t.Run("dash", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, value.V("RS-ACME"), rpsl.RSName("RS-ACME"))
	})
	t.Run("no dash", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, value.V("RS-ACME"), rpsl.RSName("RSACME"))
	})
	t.Run("no prefix", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, value.V("RS-ACME"), rpsl.RSName("ACME"))
	})
}

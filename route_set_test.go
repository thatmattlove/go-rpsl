package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
)

func Test_RouteSet(t *testing.T) {
	t.Parallel()
	rs := &rpsl.RouteSet{
		RouteSet: "RS-ACME",
		Members:  rpsl.RSMembers(rpsl.RSMember("192.0.2.0/24"), rpsl.RSMember("RS-CORP")),
	}
	exp := `route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP`
	result, err := rs.RPSL()
	testza.AssertNoError(t, err)
	testza.AssertEqual(t, exp, result)
}

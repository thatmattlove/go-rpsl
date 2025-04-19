package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
)

func TestRoute6_RPSL(t *testing.T) {
	r := &rpsl.Route6{
		Route6:      "2001:db8::/32",
		Origin:      65000,
		Description: "test",
		AdminPOC:    "TEST-ADMIN",
		TechPOC:     "TEST-TECH",
		MntBy:       "MNT-TEST",
	}
	result, err := r.RPSL()
	testza.AssertNoError(t, err)
	exp := `route: 2001:db8::/32
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST`
	testza.AssertEqual(t, exp, result)
}

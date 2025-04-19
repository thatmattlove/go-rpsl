package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
)

func TestRoute6_RPSL(t *testing.T) {
	t.Parallel()
	r := &rpsl.Route6{
		Route6:      "2001:db8::/32",
		Origin:      65000,
		Description: "test",
		AdminPOC:    "TEST-ADMIN",
		TechPOC:     "TEST-TECH",
		MntBy:       "MNT-TEST",
	}
	t.Run("base", func(t *testing.T) {
		result, err := r.RPSL()
		testza.AssertNoError(t, err)
		exp := `route: 2001:db8::/32
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST`
		testza.AssertEqual(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		testza.AssertEqual(t, "2001:db8::/32", r.String())
	})
	t.Run("with extra", func(t *testing.T) {
		r.Source = "ARIN"
		r.AddExtra("extra", "value")
		testza.AssertNotNil(t, r.Extra)
		testza.AssertEqual(t, "value", r.Extra["extra"])
		exp := `route: 2001:db8::/32
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST
extra: value
source: ARIN`
		result, err := r.RPSL()
		testza.AssertNoError(t, err)
		testza.AssertEqual(t, exp, result)
	})
}

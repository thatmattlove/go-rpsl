package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
)

func TestRoute_RPSL(t *testing.T) {
	t.Parallel()
	r := &rpsl.Route{
		Route:       "192.0.2.0/24",
		Origin:      65000,
		Description: "test",
		AdminPOC:    "TEST-ADMIN",
		TechPOC:     "TEST-TECH",
		MntBy:       "MNT-TEST",
	}
	t.Run("base", func(t *testing.T) {
		result, err := r.RPSL()
		testza.AssertNoError(t, err)
		exp := `route: 192.0.2.0/24
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST`
		testza.AssertEqual(t, exp, result)
	})
	t.Run("with extra", func(t *testing.T) {
		r.Source = "ARIN"
		r.AddExtra("extra", "value")
		testza.AssertNotNil(t, r.Extra)
		testza.AssertEqual(t, "value", r.Extra["extra"])
		exp := `route: 192.0.2.0/24
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

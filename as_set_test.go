package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
	"go.mdl.wtf/rpsl/internal/value"
)

func Test_ASSet(t *testing.T) {
	asSet := &rpsl.ASSet{
		ASSet:   "AS-ACME",
		Members: rpsl.ASSetMembers(rpsl.ASNName(65000), rpsl.ASSetName("AS-65001")),
	}
	t.Run("rpsl", func(t *testing.T) {
		t.Parallel()

		exp := `as-set: AS-ACME
members: AS65000
members: AS-65001`
		result, err := asSet.RPSL()
		testza.AssertNoError(t, err)
		testza.AssertEqual(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, "AS-ACME", asSet.String())
	})
}

func Test_ASSetName(t *testing.T) {
	t.Run("dash", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, value.V("AS-65000"), rpsl.ASSetName("AS-65000"))
	})
	t.Run("no dash", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, value.V("AS-65000"), rpsl.ASSetName("AS65000"))
	})
	t.Run("no prefix", func(t *testing.T) {
		t.Parallel()
		testza.AssertEqual(t, value.V("AS-65000"), rpsl.ASSetName("65000"))
	})
}

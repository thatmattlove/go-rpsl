package rpsl_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl"
)

func Test_Description(t *testing.T) {
	t.Parallel()
	s := rpsl.Description(`line one
line two
line three`)
	exp := `descr: line one
descr: line two
descr: line three`
	testza.AssertEqual(t, exp, s.RPSL())
}

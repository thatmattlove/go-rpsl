package value_test

import (
	"testing"

	"github.com/MarvinJWendt/testza"
	"go.mdl.wtf/rpsl/internal/value"
)

func Test_V(t *testing.T) {
	t.Parallel()
	v := value.V(t.Name())
	testza.AssertEqual(t, t.Name(), v.String())
}

func Test_VCommaSpace(t *testing.T) {
	t.Parallel()
	vc := value.VCommaSpace{value.V("one"), value.V("two"), value.V("three")}
	exp := "one, two, three"
	testza.AssertEqual(t, exp, vc.String())
}

func Test_VComma(t *testing.T) {
	t.Parallel()
	vc := value.VComma{value.V("one"), value.V("two"), value.V("three")}
	exp := "one,two,three"
	testza.AssertEqual(t, exp, vc.String())
}

func Test_VNewline(t *testing.T) {
	t.Parallel()
	vc := value.VNewline{value.V("one"), value.V("two"), value.V("three")}
	exp := `one
two
three`
	testza.AssertEqual(t, exp, vc.String())
}

package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mdl.wtf/rpsl"
)

func Test_Description(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		s := rpsl.Description(`line one
line two
line three`)
		exp := `descr: line one
descr: line two
descr: line three`
		assert.Equal(t, exp, s.RPSL())
	})
}

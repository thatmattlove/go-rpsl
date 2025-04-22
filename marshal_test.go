package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
)

func Test_MarshalBinary(t *testing.T) {
	t.Run("route", func(t *testing.T) {
		t.Parallel()
		route := rpsl.Route{
			Route:  "192.0.2.0/24",
			Origin: rpsl.ASN(65000),
			Description: `line one
line two`,
			Extra: map[string]string{"extra1": "value1"},
		}
		exp := []byte(`route: 192.0.2.0/24
origin: AS65000
descr: line one
descr: line two
extra1: value1`)
		result, err := rpsl.MarshalBinary(&route)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("as-set", func(t *testing.T) {
		t.Parallel()
		asSet := rpsl.ASSet{
			ASSet: "AS-ACME",
			Description: `line one
line two`,
			Members: rpsl.ASSetMembers(rpsl.ASSetName("AS-ACME"), rpsl.ASNName(65000)),
		}
		exp := []byte(`as-set: AS-ACME
descr: line one
descr: line two
members: AS-ACME
members: AS65000`)
		result, err := rpsl.MarshalBinary(&asSet)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("err non struct", func(t *testing.T) {
		t.Parallel()
		_, err := rpsl.MarshalBinary("")
		require.Error(t, err)
	})
}

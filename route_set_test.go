package rpsl_test

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
)

func Test_RouteSet(t *testing.T) {
	t.Parallel()
	rs := rpsl.RouteSet{
		RouteSet: "RS-ACME",
		Members:  []string{"192.0.2.0/24", "RS-CORP"},
	}
	t.Run("base", func(t *testing.T) {
		exp := []byte(`route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP`)
		result, err := rpsl.MarshalBinary(&rs)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "RS-ACME", rs.String())
	})
	t.Run("with extra", func(t *testing.T) {
		rs.AddExtra("extra", "value")
		assert.NotNil(t, rs.Extra)
		assert.Equal(t, "value", rs.Extra["extra"])
		exp := []byte(`route-set: RS-ACME
members: 192.0.2.0/24,RS-CORP
extra: value`)
		result, err := rpsl.MarshalBinary(&rs)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with descr", func(t *testing.T) {
		t.Parallel()
		rs := rpsl.RouteSet{
			RouteSet: "RS-ACME",
			Description: `123 Name Street
City, ST
12345
US`,
			Members: []string{"192.0.2.0/24"},
		}
		exp := []byte(`route-set: RS-ACME
descr: 123 Name Street
descr: City, ST
descr: 12345
descr: US
members: 192.0.2.0/24`)
		result, err := rpsl.MarshalBinary(&rs)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

func Test_RSSetName(t *testing.T) {
	t.Run("dash", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "RS-ACME", rpsl.RSName("RS-ACME"))
	})
	t.Run("no dash", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "RS-ACME", rpsl.RSName("RSACME"))
	})
	t.Run("no prefix", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "RS-ACME", rpsl.RSName("ACME"))
	})
}

func Test_RSMembers(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"192.0.2.0/24", "RS-SET"}, rpsl.RSMembers(netip.MustParsePrefix("192.0.2.0/24"), "RS-SET"))
}

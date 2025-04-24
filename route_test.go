package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
)

func TestRoute_RPSL(t *testing.T) {
	t.Parallel()
	r := rpsl.Route{
		Route:       "192.0.2.0/24",
		Origin:      65000,
		Description: "test",
		AdminPOC:    "TEST-ADMIN",
		TechPOC:     "TEST-TECH",
		MntBy:       "MNT-TEST",
	}
	t.Run("base", func(t *testing.T) {
		result, err := rpsl.MarshalBinary(&r)
		require.NoError(t, err)
		exp := []byte(`route: 192.0.2.0/24
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST`)
		assert.Equal(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "192.0.2.0/24", r.String())
	})
	t.Run("with extra", func(t *testing.T) {
		r.Source = "ARIN"
		r.AddExtra("extra", "value")
		assert.NotNil(t, r.Extra)
		assert.Equal(t, "value", r.Extra["extra"])
		exp := []byte(`route: 192.0.2.0/24
origin: AS65000
descr: test
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST
extra: value
source: ARIN`)
		result, err := rpsl.MarshalBinary(&r)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with long descr", func(t *testing.T) {
		r.Description = `123 Name Street
City, ST
12345
US`
		exp := []byte(`route: 192.0.2.0/24
origin: AS65000
descr: 123 Name Street
descr: City, ST
descr: 12345
descr: US
admin-c: TEST-ADMIN
tech-c: TEST-TECH
mnt-by: MNT-TEST
extra: value
source: ARIN`)
		result, err := rpsl.MarshalBinary(&r)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
)

func Test_ASSet(t *testing.T) {
	t.Parallel()
	asSet := rpsl.ASSet{
		ASSet:   "AS-ACME",
		Members: []string{"AS65000", "AS-65001"},
	}
	t.Run("rpsl", func(t *testing.T) {
		exp := []byte(`as-set: AS-ACME
members: AS65000
members: AS-65001`)
		result, err := rpsl.MarshalBinary(&asSet)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "AS-ACME", asSet.String())
	})
	t.Run("with extra", func(t *testing.T) {
		asSet.AddExtra("extra", "value")
		assert.NotNil(t, asSet.Extra)
		assert.Equal(t, "value", asSet.Extra["extra"])
		exp := []byte(`as-set: AS-ACME
members: AS65000
members: AS-65001
extra: value`)
		result, err := rpsl.MarshalBinary(&asSet)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
	t.Run("with descr", func(t *testing.T) {
		t.Parallel()
		asSet := rpsl.ASSet{
			ASSet:   "AS-ACME",
			Members: []string{"AS65000"},
			Description: `Some
Address`,
		}
		exp := []byte(`as-set: AS-ACME
descr: Some
descr: Address
members: AS65000`)
		result, err := rpsl.MarshalBinary(&asSet)
		require.NoError(t, err)
		assert.Equal(t, exp, result)
	})
}

func Test_ASSetName(t *testing.T) {
	t.Run("dash", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "AS-65000", rpsl.ASSetName("AS-65000"))
	})
	t.Run("no dash", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "AS-65000", rpsl.ASSetName("AS65000"))
	})
	t.Run("no prefix", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "AS-65000", rpsl.ASSetName("65000"))
	})
}

func Test_ASSetMembers(t *testing.T) {
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME"},
		rpsl.ASSetMembers("AS65000", "AS-ACME"),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME"},
		rpsl.ASSetMembers("65000", "AS-ACME"),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME"},
		rpsl.ASSetMembers(65000, "AS-ACME"),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME"},
		rpsl.ASSetMembers(uint32(65000), "AS-ACME"),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME", "AS65001", "AS-ACME-2"},
		rpsl.ASSetMembers(65000, "AS-ACME", []string{"65001", "AS-ACME-2"}),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME", "AS65001", "AS65002"},
		rpsl.ASSetMembers(65000, "AS-ACME", []uint32{65001, 65002}),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME", "AS65001", "AS65002"},
		rpsl.ASSetMembers(65000, "AS-ACME", []rpsl.ASN{65001, 65002}),
	)
	assert.Equal(t,
		[]string{"AS65000", "AS-ACME", "AS65001", "AS65002"},
		rpsl.ASSetMembers(65000, "AS-ACME", []any{65001, 65002}),
	)
}

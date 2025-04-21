package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mdl.wtf/rpsl"
)

func Test_RSMembers(t *testing.T) {
	t.Parallel()
	exp := "192.0.2.0/24,198.51.100.0/24"
	result := rpsl.RSMembers(rpsl.RSMember("192.0.2.0/24"), rpsl.RSMember("198.51.100.0/24"))
	assert.Equal(t, exp, result.String())
}

func Test_ASSetMembers(t *testing.T) {
	t.Parallel()
	exp := `AS65000
members: AS65001`
	result := rpsl.ASSetMembers(rpsl.ASNName(65000), rpsl.ASNName(65001))
	assert.Equal(t, exp, result.String())
}

func Test_AutNumMembers(t *testing.T) {
	t.Parallel()
	exp := "AS65001, AS65002, AS-ACME"
	result := rpsl.AutNumMembers(rpsl.ASNName(65001), rpsl.AutNumMember("AS65002"), rpsl.ASSetName("AS-ACME"))
	assert.Equal(t, exp, result.String())
}

package rpsl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mdl.wtf/rpsl"
)

var b = []byte(`route: 192.0.2.0/24
origin: AS65000`)

func Test_MarshalBinary(t *testing.T) {
	t.Parallel()
	obj := rpsl.Route{
		Route:  "192.0.2.0/24",
		Origin: 65000,
	}

	result, err := rpsl.MarshalBinary(&obj)
	require.NoError(t, err)
	assert.Equal(t, b, result)
}

func Test_UnmarshalBinary(t *testing.T) {
	var obj rpsl.Route
	err := rpsl.UnmarshalBinary(b, &obj)
	require.NoError(t, err)
	assert.Equal(t, "192.0.2.0/24", obj.Route)
	assert.Equal(t, rpsl.ASN(65000), obj.Origin)
}

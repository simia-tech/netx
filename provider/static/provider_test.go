package static_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/value"
)

func TestNewProviderFromURLs(t *testing.T) {
	provider := static.NewProvider()
	provider.Add("test", value.NewDial("tcp", "localhost:8080"))
	provider.Add("test", value.NewDial("quic", "localhost:8081"))

	dials, err := provider.Dials("test")
	require.NoError(t, err)
	require.Len(t, dials, 2)

	assert.Equal(t, "tcp", dials[0].Network())
	assert.Equal(t, "localhost:8080", dials[0].Address())
	assert.Equal(t, "quic", dials[1].Network())
	assert.Equal(t, "localhost:8081", dials[1].Address())
}

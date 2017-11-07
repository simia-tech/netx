package static_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx/provider/static"
)

func TestNewProviderFromURLs(t *testing.T) {
	provider, err := static.NewProviderFromURLs("tcp://localhost:8080", "quic://localhost:8081")
	require.NoError(t, err)

	addrs, err := provider.Addrs()
	require.NoError(t, err)
	require.Len(t, addrs, 2)

	assert.Equal(t, "tcp", addrs[0].Network())
	assert.Equal(t, "localhost:8080", addrs[0].String())
	assert.Equal(t, "quic", addrs[1].Network())
	assert.Equal(t, "localhost:8081", addrs[1].String())
}

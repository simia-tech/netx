package dnssrv_test

import (
	"testing"

	"github.com/simia-tech/netx/provider/consul/client"
	"github.com/simia-tech/netx/provider/dnssrv"
	"github.com/simia-tech/netx/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderEndpoints(t *testing.T) {
	c, err := client.NewConsul([]string{"tcp://127.0.0.1:8500"})
	require.NoError(t, err)
	id, err := c.Register("test", value.NewEndpoint("quic", "localhost:1000"))
	require.NoError(t, err)
	defer c.Deregister(id)

	p := dnssrv.NewProvider("127.0.0.1:8600", "quic")

	endpoints, err := p.Endpoints("test")
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	assert.Equal(t, "quic", endpoints[0].Network())
	assert.Equal(t, "localhost:1000", endpoints[0].Address())
}

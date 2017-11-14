package consul_test

import (
	"testing"

	"github.com/simia-tech/netx/provider/consul"
	"github.com/simia-tech/netx/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderEndpoints(t *testing.T) {
	p, err := consul.NewProvider([]string{"tcp://127.0.0.1:8500"})
	require.NoError(t, err)

	p.EndpointDialOptions = []value.DialOption{}

	serviceID, err := p.Client.Register("test", value.NewEndpoint("tcp", "localhost:1000"))
	require.NoError(t, err)
	defer p.Client.Deregister(serviceID)

	endpoints, err := p.Endpoints("test")
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	assert.Equal(t, "tcp", endpoints[0].Network())
	assert.Equal(t, "localhost:1000", endpoints[0].Address())
	assert.Equal(t, p.EndpointDialOptions, endpoints[0].Options())
}

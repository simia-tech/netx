package blacklist_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx/filter/blacklist"
	"github.com/simia-tech/netx/value"
)

func TestFilterNoEndpoints(t *testing.T) {
	f := blacklist.NewFilter(blacklist.ConstantBackoff(50 * time.Millisecond))
	endpoints := value.Endpoints{value.NewEndpoint("tcp", "localhost:1000")}

	filteredEndpoints, err := f.Filter(endpoints)
	require.NoError(t, err)

	assert.Equal(t, endpoints, filteredEndpoints)
}

func TestFilterFailedEndpoints(t *testing.T) {
	f := blacklist.NewFilter(blacklist.ConstantBackoff(50 * time.Millisecond))
	endpoint := value.NewEndpoint("tcp", "localhost:1000")
	endpoints := value.Endpoints{endpoint}

	require.NoError(t, f.Failure(endpoint))

	filteredEndpoints, err := f.Filter(endpoints)
	require.NoError(t, err)

	assert.Empty(t, filteredEndpoints)
}

func TestFilterRecoveredEndpoints(t *testing.T) {
	f := blacklist.NewFilter(blacklist.ConstantBackoff(50 * time.Millisecond))
	endpoint := value.NewEndpoint("tcp", "localhost:1000")
	endpoints := value.Endpoints{endpoint}

	require.NoError(t, f.Failure(endpoint))
	require.NoError(t, f.Success(endpoint))

	filteredEndpoints, err := f.Filter(endpoints)
	require.NoError(t, err)

	assert.Equal(t, endpoints, filteredEndpoints)
}

func TestFilterFailedEndpointsAfterBackoffDuration(t *testing.T) {
	f := blacklist.NewFilter(blacklist.ConstantBackoff(50 * time.Millisecond))
	endpoint := value.NewEndpoint("tcp", "localhost:1000")
	endpoints := value.Endpoints{endpoint}

	require.NoError(t, f.Failure(endpoint))
	time.Sleep(100 * time.Millisecond)

	filteredEndpoints, err := f.Filter(endpoints)
	require.NoError(t, err)

	assert.Equal(t, endpoints, filteredEndpoints)
}

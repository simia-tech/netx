package static_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/value"
)

func TestProviderEndpoints(t *testing.T) {
	provider := static.NewProvider()
	provider.Add("test", value.NewEndpoint("tcp", "localhost:8080"))
	provider.Add("test", value.NewEndpoint("quic", "localhost:8081"))

	endpoints, err := provider.Endpoints("test")
	require.NoError(t, err)
	require.Len(t, endpoints, 2)

	assert.Equal(t, "tcp", endpoints[0].Network())
	assert.Equal(t, "localhost:8080", endpoints[0].Address())
	assert.Equal(t, "quic", endpoints[1].Network())
	assert.Equal(t, "localhost:8081", endpoints[1].Address())
}

func TestProviderConcurrentEndpoints(t *testing.T) {
	provider := static.NewProvider()

	wg := sync.WaitGroup{}
	for index := 0; index < 5; index++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 20; i++ {
				provider.Add("test", value.NewEndpoint("tcp", fmt.Sprintf("localhost:%d", 8000+i)))

				_, err := provider.Endpoints("test")
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

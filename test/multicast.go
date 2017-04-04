package test

import (
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

// MulticastTest runs a multicast test suite with the provided options.
func MulticastTest(t *testing.T, options *Options) {
	t.Run("OneProducerOneConsumer", func(t *testing.T) {
		address := netx.RandomAddress("multicast-")
		if options.MulticastRequestAddress != "" {
			address = options.MulticastRequestAddress
		}

		conn, err := netx.ListenMulticast(options.MulticastNetwork, address, options.MulticastOptions...)
		require.NoError(t, err)

		go func() {
			conn, e := netx.DialMulticast(options.MulticastNetwork, address, options.MulticastOptions...)
			require.NoError(t, e)

			WriteBlock(conn, []byte("test"))
		}()

		data, err := ReadBlock(conn)
		require.NoError(t, err)

		assert.Equal(t, "test", string(data))
	})

	t.Run("ManyProducersOneConsumer", func(t *testing.T) {
		requestAddress := netx.RandomAddress("multicast-request-")
		if options.MulticastRequestAddress != "" {
			requestAddress = options.MulticastRequestAddress
		}
		responseAddress := netx.RandomAddress("multicast-response-")
		if options.MulticastResponseAddress != "" {
			responseAddress = options.MulticastResponseAddress
		}

		n := 5
		for index := 0; index < n; index++ {
			conn, err := netx.ListenAndDialMulticast(options.MulticastNetwork, requestAddress, responseAddress, options.MulticastOptions...)
			require.NoError(t, err)

			go func(index int, conn io.ReadWriteCloser) {
				data := RequireReadBlock(t, conn)
				RequireWriteBlock(t, conn, []byte(fmt.Sprintf("node %d: %s", index, data)))

				require.NoError(t, conn.Close())
			}(index, conn)
		}

		conn, err := netx.ListenAndDialMulticast(options.MulticastNetwork, responseAddress, requestAddress, options.MulticastOptions...)
		require.NoError(t, err)

		RequireWriteBlock(t, conn, []byte("test"))
		results := []string{}
		for index := 0; index < n; index++ {
			data := RequireReadBlock(t, conn)
			results = append(results, string(data))
		}

		sort.Strings(results)
		assert.Equal(t, []string{
			"node 0: test",
			"node 1: test",
			"node 2: test",
			"node 3: test",
			"node 4: test",
		}, results)
	})
}

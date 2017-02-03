package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func MulticastTest(t *testing.T, options *Options) {
	t.Run("OneProducerOneConsumer", func(t *testing.T) {
		address := netx.RandomAddress("multicast-")

		conn, err := netx.ListenMulticast(options.ListenMulticastNetwork, address, options.ListenMulticastOptions...)
		require.NoError(t, err)

		go func() {
			conn, err := netx.DialMulticast(options.DialMulticastNetwork, address, options.DialMulticastOptions...)
			require.NoError(t, err)

			WriteBlock(conn, []byte("test"))
		}()

		data, err := ReadBlock(conn)
		require.NoError(t, err)

		assert.Equal(t, "test", string(data))
	})
}

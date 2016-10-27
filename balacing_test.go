package netx_test

import (
	"testing"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/require"
)

func TestBalancing(t *testing.T) {
	address := netx.RandomAddress("echo-")

	listenerOne, _ := setUpTestEchoListener(t, address)
	defer listenerOne.Close()
	listenerTwo, _ := setUpTestEchoListener(t, address)
	defer listenerTwo.Close()

	for index := 0; index < 4; index++ {
		conn := setUpTestEchoClient(t, address)

		requireWrite(t, conn, []byte("test"))
		buffer := requireRead(t, conn)
		require.Equal(t, "test", string(buffer))

		require.NoError(t, conn.Close())
	}
}

func BenchmarkBalancing(b *testing.B) {
	address := netx.RandomAddress("echo-")

	listenerOne, _ := setUpTestEchoListener(b, address)
	defer listenerOne.Close()
	listenerTwo, _ := setUpTestEchoListener(b, address)
	defer listenerTwo.Close()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		conn := setUpTestEchoClient(b, address)

		requireWrite(b, conn, []byte("test"))
		buffer := requireRead(b, conn)
		require.Equal(b, "test", string(buffer))

		require.NoError(b, conn.Close())
	}
}

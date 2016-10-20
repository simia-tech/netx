package netx_test

import (
	"testing"

	n "github.com/nats-io/nats"
	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/require"
)

func TestBalancing(t *testing.T) {
	address := n.NewInbox()

	listenerOne, _ := setUpTestEchoListener(t, address)
	defer listenerOne.Close()
	listenerTwo, _ := setUpTestEchoListener(t, address)
	defer listenerTwo.Close()

	conn, err := netx.Dial(defaultNatsURL, address)
	require.NoError(t, err)

	requireWrite(t, conn, []byte("test"))
	buffer := requireRead(t, conn, 4)
	require.Equal(t, "test", string(buffer))
}

func BenchmarkBalancing(b *testing.B) {
	address := n.NewInbox()

	listenerOne, _ := setUpTestEchoListener(b, address)
	defer listenerOne.Close()
	listenerTwo, _ := setUpTestEchoListener(b, address)
	defer listenerTwo.Close()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		conn, err := netx.Dial(defaultNatsURL, address)
		require.NoError(b, err)

		requireWrite(b, conn, []byte("test"))
		buffer := requireRead(b, conn, 4)
		require.Equal(b, "test", string(buffer))

		require.NoError(b, conn.Close())
	}
}

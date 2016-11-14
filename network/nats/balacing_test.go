package nats_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/test"
)

func TestBalancing(t *testing.T) {
	address := netx.RandomAddress("echo-")

	listenerOne, _ := setUpTestEchoListener(t, address)
	defer listenerOne.Close()
	listenerTwo, _ := setUpTestEchoListener(t, address)
	defer listenerTwo.Close()

	for index := 0; index < 4; index++ {
		conn := setUpTestEchoClient(t, address)

		test.RequireWriteBlock(t, conn, []byte("test"))
		assert.Equal(t, "test", string(test.RequireReadBlock(t, conn)))

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

		test.RequireWriteBlock(b, conn, []byte("test"))
		assert.Equal(b, "test", string(test.RequireReadBlock(b, conn)))

		require.NoError(b, conn.Close())
	}
}

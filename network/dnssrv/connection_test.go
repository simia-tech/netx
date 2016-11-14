package dnssrv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx/test"
)

func TestConnection(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	for index := 0; index < 4; index++ {
		conn := setUpTestEchoClient(t, listener.Addr().String())

		test.RequireWriteBlock(t, conn, []byte("test"))
		assert.Equal(t, "test", string(test.RequireReadBlock(t, conn)))

		require.NoError(t, conn.Close())
	}
}

func BenchmarkConnection(b *testing.B) {
	listener, _ := setUpTestEchoListener(b)
	defer listener.Close()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		conn := setUpTestEchoClient(b, listener.Addr().String())

		test.RequireWriteBlock(b, conn, []byte("test"))
		assert.Equal(b, "test", string(test.RequireReadBlock(b, conn)))

		require.NoError(b, conn.Close())
	}
}

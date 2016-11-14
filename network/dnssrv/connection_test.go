package dnssrv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	for index := 0; index < 4; index++ {
		conn := setUpTestEchoClient(t, listener.Addr().String())

		requireWrite(t, conn, []byte("test"))
		assert.Equal(t, "test", string(requireRead(t, conn)))

		require.NoError(t, conn.Close())
	}
}

func BenchmarkConnection(b *testing.B) {
	listener, _ := setUpTestEchoListener(b)
	defer listener.Close()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		conn := setUpTestEchoClient(b, listener.Addr().String())

		requireWrite(b, conn, []byte("test"))
		buffer := requireRead(b, conn)
		require.Equal(b, "test", string(buffer))

		require.NoError(b, conn.Close())
	}
}

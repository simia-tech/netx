package nats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

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

func TestConnectionClientClose(t *testing.T) {
	listener, errChan := setUpTestEchoListener(t)
	defer listener.Close()

	conn := setUpTestEchoClient(t, listener.Addr().String())

	require.NoError(t, conn.Close())

	assert.Error(t, <-errChan)
}

func TestConnectionListenerClose(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	conn := setUpTestEchoClient(t, listener.Addr().String())
	defer conn.Close()

	test.RequireWriteBlock(t, conn, []byte("test"))
	buffer := test.RequireReadBlock(t, conn)
	require.Equal(t, "test", string(buffer))

	_, err := conn.Read(buffer[:])
	assert.Error(t, err)
}

func TestConnectionReadAfterClose(t *testing.T) {
	listener, errChan := setUpTestEchoListener(t)
	defer listener.Close()

	conn := setUpTestEchoClient(t, listener.Addr().String())
	defer conn.Close()

	require.NoError(t, conn.Close())
	require.Error(t, <-errChan)

	buffer := [4]byte{}
	_, err := conn.Read(buffer[:])
	require.Error(t, err)
	assert.Error(t, err)
}

func TestConnectionReadTimeout(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	conn := setUpTestEchoClient(t, listener.Addr().String())
	defer conn.Close()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(10*time.Millisecond)))

	buffer := [4]byte{}
	_, err := conn.Read(buffer[:])
	require.Error(t, err)
	assert.True(t, strings.HasSuffix(err.Error(), "timeout"))
}

func TestConnectionLargeTransfer(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	conn := setUpTestEchoClient(t, listener.Addr().String())
	defer conn.Close()

	data := make([]byte, 10000)
	test.RequireWriteBlock(t, conn, data)
	reply := test.RequireReadBlock(t, conn)

	assert.True(t, bytes.Equal(data, reply))
}

func BenchmarkConnection(b *testing.B) {
	listener, _ := setUpTestEchoListener(b)
	defer listener.Close()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		conn := setUpTestEchoClient(b, listener.Addr().String())

		test.RequireWriteBlock(b, conn, []byte("test"))
		buffer := test.RequireReadBlock(b, conn)
		require.Equal(b, "test", string(buffer))

		require.NoError(b, conn.Close())
	}
}

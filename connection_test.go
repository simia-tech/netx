package netx_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func TestConnection(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	for index := 0; index < 4; index++ {
		conn := setUpTestEchoClient(t, listener.Addr().String())

		requireWrite(t, conn, []byte("test"))
		buffer := requireRead(t, conn, 4)
		require.Equal(t, "test", string(buffer))

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

	requireWrite(t, conn, []byte("test"))
	buffer := requireRead(t, conn, 4)
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
	assert.Equal(t, err, io.ErrClosedPipe)
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
	assert.Equal(t, err.Error(), "nats: timeout")
}

func BenchmarkConnection(b *testing.B) {
	listener, _ := setUpTestEchoListener(b)
	defer listener.Close()

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		conn := setUpTestEchoClient(b, listener.Addr().String())

		requireWrite(b, conn, []byte("test"))
		buffer := requireRead(b, conn, 4)
		require.Equal(b, "test", string(buffer))

		require.NoError(b, conn.Close())
	}
}

func ExampleListen() {
	listener, _ := netx.Listen("nats://localhost:4222", "echo")
	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()

		buffer := make([]byte, 5)
		conn.Read(buffer)
		conn.Write(buffer)
	}()

	client, _ := netx.Dial("nats://localhost:4222", "echo")
	defer client.Close()

	fmt.Fprintf(client, "hello")

	buffer := make([]byte, 5)
	client.Read(buffer)

	fmt.Println(string(buffer))
	// Output: hello
}

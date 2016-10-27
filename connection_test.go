package netx_test

import (
	"strings"
	"testing"
	"time"

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
	buffer := requireRead(t, conn)
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

// func ExampleListen() {
// 	listener, _ := netx.Listen("nats://localhost:4222", "echo")
// 	go func() {
// 		conn, _ := listener.Accept()
// 		defer conn.Close()
//
// 		buffer := make([]byte, 5)
// 		conn.Read(buffer)
// 		conn.Write(buffer)
// 	}()
//
// 	client, _ := netx.Dial("nats://localhost:4222", "echo")
// 	defer client.Close()
//
// 	fmt.Fprintf(client, "hello")
//
// 	buffer := make([]byte, 5)
// 	client.Read(buffer)
//
// 	fmt.Println(string(buffer))
// 	// Output: hello
// }

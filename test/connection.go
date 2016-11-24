package test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ConnectionTest runs a connection test.
func ConnectionTest(t *testing.T, options *Options) {
	t.Run("Echo", func(t *testing.T) {
		listener, _ := makeEchoListener(t, "", options)
		defer listener.Close()
		makeEchoCalls(t, 4, listener.Addr().String(), options)
	})

	t.Run("ClientClose", func(t *testing.T) {
		listener, errChan := makeEchoListener(t, "", options)
		defer listener.Close()

		conn := makeEchoConn(t, listener.Addr().String(), options)

		require.NoError(t, conn.Close())
		assert.Error(t, <-errChan)
	})

	t.Run("ListenerClose", func(t *testing.T) {
		listener, _ := makeEchoListener(t, "", options)
		defer listener.Close()

		conn := makeEchoConn(t, listener.Addr().String(), options)
		defer conn.Close()

		RequireWriteBlock(t, conn, []byte("test"))
		buffer := RequireReadBlock(t, conn)
		require.Equal(t, "test", string(buffer))

		_, err := conn.Read(buffer[:])
		assert.Error(t, err)
	})

	t.Run("ReadAfterClose", func(t *testing.T) {
		listener, errChan := makeEchoListener(t, "", options)
		defer listener.Close()

		conn := makeEchoConn(t, listener.Addr().String(), options)
		defer conn.Close()

		require.NoError(t, conn.Close())
		require.Error(t, <-errChan)

		buffer := [4]byte{}
		_, err := conn.Read(buffer[:])
		require.Error(t, err)
		assert.Error(t, err)
	})

	t.Run("ReadTimeout", func(t *testing.T) {
		listener, _ := makeEchoListener(t, "", options)
		defer listener.Close()

		conn := makeEchoConn(t, listener.Addr().String(), options)
		defer conn.Close()

		require.NoError(t, conn.SetReadDeadline(time.Now().Add(10*time.Millisecond)))

		buffer := [4]byte{}
		_, err := conn.Read(buffer[:])
		require.Error(t, err)
		assert.True(t, strings.HasSuffix(err.Error(), "timeout"))
	})

	t.Run("LargeTransfer", func(t *testing.T) {
		listener, _ := makeEchoListener(t, "", options)
		defer listener.Close()

		conn := makeEchoConn(t, listener.Addr().String(), options)
		defer conn.Close()

		data := make([]byte, 10000)
		RequireWriteBlock(t, conn, data)
		reply := RequireReadBlock(t, conn)

		assert.True(t, bytes.Equal(data, reply))
	})
}

// ConnectionBenchmark runs a connection benchmark.
func ConnectionBenchmark(b *testing.B, options *Options) {
	address, close := makeEchoListeners(b, 1, options)
	defer close()
	b.ResetTimer()
	makeEchoCalls(b, b.N, address, options)
}
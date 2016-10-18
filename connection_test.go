package netx_test

import (
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

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	requireWrite(t, conn, []byte("test"))
	buffer := requireRead(t, conn, 4)

	assert.Equal(t, "test", string(buffer))
}

func TestConnectionClientClose(t *testing.T) {
	listener, errChan := setUpTestEchoListener(t)
	defer listener.Close()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)

	require.NoError(t, conn.Close())

	assert.Error(t, <-errChan)
}

func TestConnectionListenerClose(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)

	requireWrite(t, conn, []byte("test"))
	buffer := requireRead(t, conn, 4)
	require.Equal(t, "test", string(buffer))

	_, err = conn.Read(buffer[:])
	assert.Error(t, err)
}

func TestConnectionReadAfterClose(t *testing.T) {
	listener, errChan := setUpTestEchoListener(t)
	defer listener.Close()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)

	require.NoError(t, conn.Close())
	require.Error(t, <-errChan)

	buffer := [4]byte{}
	_, err = conn.Read(buffer[:])
	require.Error(t, err)
	assert.Equal(t, err, io.ErrClosedPipe)
}

func TestConnectionReadTimeout(t *testing.T) {
	listener, _ := setUpTestEchoListener(t)
	defer listener.Close()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(10*time.Millisecond)))

	buffer := [4]byte{}
	_, err = conn.Read(buffer[:])
	require.Error(t, err)
	assert.Equal(t, err.Error(), "timeout")
}

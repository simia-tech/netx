package netx_test

import (
	"testing"

	n "github.com/nats-io/nats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func TestConnection(t *testing.T) {
	listener := setUpTestListener(t)
	defer listener.Close()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	requireWrite(t, conn, []byte("test"))
	buffer := requireRead(t, conn, 4)

	assert.Equal(t, "test", string(buffer))
}

func TestConnectionClientClose(t *testing.T) {
	listener, err := netx.Listen(defaultNatsURL, n.NewInbox())
	require.NoError(t, err)

	signal := make(chan error)
	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)

		signal <- nil

		buffer := [4]byte{}
		_, err = conn.Read(buffer[:])
		signal <- err
	}()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)

	<-signal
	require.NoError(t, conn.Close())

	assert.Error(t, <-signal)
}

func TestConnectionListenerClose(t *testing.T) {
	listener, err := netx.Listen(defaultNatsURL, n.NewInbox())
	require.NoError(t, err)

	signal := make(chan error)
	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)

		<-signal

		require.NoError(t, conn.Close())
	}()

	conn, err := netx.Dial(defaultNatsURL, listener.Addr().String())
	require.NoError(t, err)

	signal <- nil

	buffer := [4]byte{}
	_, err = conn.Read(buffer[:])
	assert.Error(t, err)
}

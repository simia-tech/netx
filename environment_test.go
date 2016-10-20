package netx_test

import (
	"net"
	"testing"

	n "github.com/nats-io/nats"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

const defaultNatsURL = "nats://localhost:4222"

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	address := n.NewInbox()
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen(defaultNatsURL, address)
	require.NoError(tb, err)

	errChan := make(chan error, 1)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				errChan <- err
				return
			}

			buffer := [4]byte{}
			n, err := conn.Read(buffer[:])
			if err != nil {
				errChan <- err
				return
			}
			require.Equal(tb, len(buffer), n)

			n, err = conn.Write(buffer[:])
			if err != nil {
				errChan <- err
				return
			}
			require.Equal(tb, len(buffer), n)

			if err := conn.Close(); err != nil {
				errChan <- err
				return
			}
		}
	}()

	return listener, errChan
}

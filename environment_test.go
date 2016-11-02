package netx_test

import (
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	network := os.Getenv("LISTEN_NETWORK")
	if network == "" {
		tb.Skip("LISTEN_NETWORK is unset")
	}

	address := netx.RandomAddress("echo-")
	if len(addresses) > 0 {
		address = addresses[0]
	}

	nodes := strings.Split(os.Getenv("LISTEN_NETWORK_NODES"), ",")
	localAddress := os.Getenv("LISTEN_LOCAL_ADDRESS")

	listener, err := netx.Listen(network, address, netx.Nodes(nodes), netx.LocalAddress(localAddress))
	require.NoError(tb, err)

	errChan := make(chan error, 1)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				errChan <- err
				return
			}

			data, err := readBlock(conn)
			if err != nil {
				log.Printf("test echo listener read error: %v", err)
				errChan <- err
				return
			}
			if err := writeBlock(conn, data); err != nil {
				log.Printf("test echo listener write error: %v", err)
				errChan <- err
				return
			}

			if err := conn.Close(); err != nil {
				log.Printf("test echo listener close error: %v", err)
				errChan <- err
				return
			}
		}
	}()

	return listener, errChan
}

func setUpTestEchoClient(tb testing.TB, address string) net.Conn {
	network := os.Getenv("DIAL_NETWORK")
	if network == "" {
		tb.Skip("DIAL_NETWORK is unset")
	}

	nodes := strings.Split(os.Getenv("DIAL_NETWORK_NODES"), ",")

	conn, err := netx.Dial(network, address, netx.Nodes(nodes))
	require.NoError(tb, err)

	return conn
}

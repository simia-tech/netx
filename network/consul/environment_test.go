package consul_test

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/consul"
)

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	address := netx.RandomAddress("echo-")
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen("consul", address, netx.Nodes("http://localhost:8500"), netx.LocalAddress("localhost:0"))
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
	conn, err := netx.Dial("consul", address, netx.Nodes("http://localhost:8500"))
	require.NoError(tb, err)
	return conn
}

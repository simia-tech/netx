package test

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func makeEchoListeners(tb testing.TB, n int, options *Options) (string, func() []int, func()) {
	address := netx.RandomAddress("echo-")

	counters := []func() int{}
	listeners := []net.Listener{}
	for index := 0; index < n; index++ {
		listener, counter, _ := makeEchoListener(tb, address, options)
		listeners = append(listeners, listener)
		counters = append(counters, counter)
	}

	return address, func() []int {
			result := []int{}
			for _, counter := range counters {
				result = append(result, counter())
			}
			return result
		}, func() {
			for _, listener := range listeners {
				listener.Close()
			}
		}
}

func makeEchoListener(tb testing.TB, address string, options *Options) (net.Listener, func() int, chan error) {
	if address == "" {
		if options.ListenAddress == "" {
			address = netx.RandomAddress("echo-")
		} else {
			address = options.ListenAddress
		}
	}

	listener, err := netx.Listen(options.ListenNetwork, address, options.ListenOptions...)
	require.NoError(tb, err)

	counter := 0
	errChan := make(chan error, 1)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				errChan <- err
				return
			}

			data, err := ReadBlock(conn)
			if err != nil {
				log.Printf("test echo listener read error: %v", err)
				errChan <- err
				return
			}
			if err := WriteBlock(conn, data); err != nil {
				log.Printf("test echo listener write error: %v", err)
				errChan <- err
				return
			}

			if err := conn.Close(); err != nil {
				log.Printf("test echo listener close error: %v", err)
				errChan <- err
				return
			}

			counter++
		}
	}()

	return listener, func() int {
		return counter
	}, errChan
}

func makeEchoCalls(tb testing.TB, n int, address string, options *Options) {
	for index := 0; index < n; index++ {
		conn := makeEchoConn(tb, address, options)

		RequireWriteBlock(tb, conn, []byte("test"))
		assert.Equal(tb, "test", string(RequireReadBlock(tb, conn)))

		require.NoError(tb, conn.Close())
	}
}

func makeEchoConn(tb testing.TB, address string, options *Options) net.Conn {
	conn, err := netx.Dial(options.DialNetwork, address, options.DialOptions...)
	require.NoError(tb, err)
	return conn
}

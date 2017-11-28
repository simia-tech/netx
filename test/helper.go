package test

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/require"
)

type action func(net.Conn) error

func echoServer(conn net.Conn) error {
	data, err := ReadBlock(conn)
	if err != nil {
		return err
	}
	if err := WriteBlock(conn, data); err != nil {
		return err
	}
	return nil
}

func echoClient(conn net.Conn) error {
	if err := WriteBlock(conn, []byte("test")); err != nil {
		return err
	}
	bytes, err := ReadBlock(conn)
	if err != nil {
		return err
	}
	if string(bytes) != "test" {
		return fmt.Errorf("expected \"test\", got \"%s\"", bytes)
	}
	return nil
}

func makeListeners(tb testing.TB, n int, a action, options *Options) (string, func() []int, func()) {
	address := netx.RandomAddress("echo-")

	counters := []func() int{}
	listeners := []net.Listener{}
	for index := 0; index < n; index++ {
		listener, counter, _ := makeListener(tb, address, a, options)
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

func makeListener(tb testing.TB, address string, a action, options *Options) (net.Listener, func() int, chan error) {
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
				log.Printf("test echo listener accept error: %v", err)
				errChan <- err
				return
			}

			if err = a(conn); err != nil {
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

func makeCalls(n int, address string, a action, options *Options) error {
	for index := 0; index < n; index++ {
		conn, err := makeConn(address, options)
		if err != nil {
			return err
		}

		if err := a(conn); err != nil {
			return err
		}

		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func makeConn(address string, options *Options) (net.Conn, error) {
	conn, err := netx.Dial(options.DialNetwork, address, options.DialOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func sum(items []int) int {
	result := 0
	for _, item := range items {
		result += item
	}
	return result
}

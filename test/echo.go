package test

import (
	"fmt"
	"log"
	"net"
	"testing"

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

func makeFailingEchoListeners(tb testing.TB, n int, options *Options) (string, func() []int, func(int) error, func()) {
	address := netx.RandomAddress("echo-")

	counters := []func() int{}
	listeners := []net.Listener{}
	publicListeners := []net.Listener{}
	listenOptions := options.ListenOptions
	for index := 0; index < n; index++ {
		publicListener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(tb, err)
		publicListeners = append(publicListeners, publicListener)

		options.ListenOptions = append(listenOptions, netx.PublicListener(publicListener))
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
		}, func(index int) error {
			return publicListeners[index].Close()
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

func makeEchoCalls(n int, address string, options *Options) error {
	for index := 0; index < n; index++ {
		conn, err := makeEchoConn(address, options)
		if err != nil {
			return fmt.Errorf("make echo conn: %v", err)
		}

		if err := WriteBlock(conn, []byte("test")); err != nil {
			return err
		}
		bytes, err := ReadBlock(conn)
		if err != nil {
			return fmt.Errorf("read block: %v", err)
		}
		if string(bytes) != "test" {
			return fmt.Errorf("expected \"test\", got \"%s\"", bytes)
		}

		if err := conn.Close(); err != nil {
			return fmt.Errorf("conn close: %v", err)
		}
	}
	return nil
}

func makeEchoConn(address string, options *Options) (net.Conn, error) {
	conn, err := netx.Dial(options.DialNetwork, address, options.DialOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

package netx_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/value"
)

func setUpEndpoint(tb testing.TB, ports ...int) (value.Endpoint, net.Listener) {
	port := 0
	if len(ports) > 0 {
		port = ports[0]
	}
	l, err := netx.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	require.NoError(tb, err)
	return value.NewEndpoint(l.Addr().Network(), l.Addr().String()), l
}

func acceptConnections(ctx context.Context, l net.Listener, wg *sync.WaitGroup) *uint64 {
	wg.Add(1)
	counter := uint64(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				conn, err := l.Accept()
				if err != nil {
					log.Printf("accept: %v", err)
					continue
				}
				conn.Close()
				atomic.AddUint64(&counter, 1)
			}
		}
	}()
	return &counter
}

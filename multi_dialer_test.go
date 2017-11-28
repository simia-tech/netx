package netx_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/selector/roundrobin"
	"github.com/simia-tech/netx/value"
)

func TestMultiDialer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	ep, l := setUpEndpoint(t)
	c := acceptConnections(ctx, l, &wg)

	p := static.NewProvider()
	p.Add("test", ep)

	md, err := netx.NewMultiDialer(p, nil, roundrobin.NewSelector())
	require.NoError(t, err)

	_, err = md.Dial("test")
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)

	l.Close()
	cancel()
	wg.Wait()

	assert.Equal(t, uint64(1), *c)
}

func TestMultiDialerEndpointFailover(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	ep, l := setUpEndpoint(t)
	c := acceptConnections(ctx, l, &wg)

	p := static.NewProvider()
	p.Add("test", ep)
	p.Add("test", value.NewEndpoint("tcp", "127.0.0.1:5020")) // not existing

	md, err := netx.NewMultiDialer(p, nil, roundrobin.NewSelector())
	require.NoError(t, err)

	_, err = md.Dial("test") // should hit the listener
	require.NoError(t, err)
	_, err = md.Dial("test") // should fail first and hit the listener again
	require.NoError(t, err)

	cancel()
	wg.Wait()
	assert.Equal(t, uint64(2), *c)
}

func TestMultiDialerEndpointRecovering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	ep1, l1 := setUpEndpoint(t)
	acceptConnections(ctx, l1, &wg)

	ep2Ctx, ep2Cancel := context.WithCancel(context.Background())
	ep2Wg := sync.WaitGroup{}
	ep2, l2 := setUpEndpoint(t)
	acceptConnections(ep2Ctx, l2, &ep2Wg)

	p := static.NewProvider()
	p.Add("test", ep1)
	p.Add("test", ep2)

	md, err := netx.NewMultiDialer(p, nil, roundrobin.NewSelector())
	require.NoError(t, err)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := md.Dial("test")
				require.NoError(t, err)
				require.NoError(t, conn.Close())
			}
		}
	}()
	time.Sleep(100 * time.Millisecond)

	l2.Close()
	ep2Cancel()
	ep2Wg.Wait()
	time.Sleep(100 * time.Millisecond)

	_, l2 = setUpEndpoint(t, value.EndpointPort(ep2))
	c2 := acceptConnections(ctx, l2, &wg)
	time.Sleep(100 * time.Millisecond)

	l1.Close()
	l2.Close()
	cancel()
	wg.Wait()
	assert.True(t, *c2 > 0)
}

func TestMultiDialerConcurrentDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	ep1, l1 := setUpEndpoint(t)
	c1 := acceptConnections(ctx, l1, &wg)

	ep2, l2 := setUpEndpoint(t)
	c2 := acceptConnections(ctx, l2, &wg)

	p := static.NewProvider()
	p.Add("test", ep1)
	p.Add("test", ep2)

	md, err := netx.NewMultiDialer(p, nil, roundrobin.NewSelector())
	require.NoError(t, err)

	clientWg := sync.WaitGroup{}
	for index := 0; index < 5; index++ {
		clientWg.Add(1)
		go func() {
			for i := 0; i < 20; i++ {
				_, err := md.Dial("test")
				require.NoError(t, err)
			}
			clientWg.Done()
		}()
	}
	clientWg.Wait()

	time.Sleep(200 * time.Millisecond)
	l1.Close()
	l2.Close()
	cancel()
	wg.Wait()

	assert.Equal(t, uint64(100), (*c1)+(*c2))
}

package netx_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/selector/roundrobin"
	"github.com/simia-tech/netx/value"
)

func TestMultiDialer(t *testing.T) {
	l, err := netx.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	ch := make(chan bool)
	go func() {
		l.Accept()
		ch <- true
	}()

	p := static.NewProvider()
	p.Add("test", value.NewEndpoint(l.Addr().Network(), l.Addr().String()))

	md, err := netx.NewMultiDialer(p, roundrobin.NewSelector())
	require.NoError(t, err)

	_, err = md.Dial("test")
	require.NoError(t, err)

	assert.True(t, <-ch)
	close(ch)
}

func TestMultiDialerFailoverWithUnreachableEndpoint(t *testing.T) {
	l, err := netx.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	ch := make(chan bool)
	go func() {
		l.Accept()
		l.Accept()
		ch <- true
	}()

	p := static.NewProvider()
	p.Add("test", value.NewEndpoint(l.Addr().Network(), l.Addr().String()))
	p.Add("test", value.NewEndpoint("tcp", "127.0.0.1:5020")) // not existing

	md, err := netx.NewMultiDialer(p, roundrobin.NewSelector())
	require.NoError(t, err)

	_, err = md.Dial("test") // should hit the listener
	require.NoError(t, err)
	_, err = md.Dial("test") // should fail first and hit the listener again
	require.NoError(t, err)

	assert.True(t, <-ch)
	close(ch)
}

func TestMultiDialerConcurrentDial(t *testing.T) {
	l1, err := netx.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	l2, err := netx.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	counter := uint32(0)
	go func() {
		for {
			if _, e := l1.Accept(); e != nil {
				return
			}
			atomic.AddUint32(&counter, 1)
		}
	}()
	go func() {
		for {
			if _, e := l2.Accept(); e != nil {
				return
			}
			atomic.AddUint32(&counter, 1)
		}
	}()

	p := static.NewProvider()
	p.Add("test", value.NewEndpoint(l1.Addr().Network(), l1.Addr().String()))
	p.Add("test", value.NewEndpoint(l2.Addr().Network(), l2.Addr().String()))

	md, err := netx.NewMultiDialer(p, roundrobin.NewSelector())
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	for index := 0; index < 5; index++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 20; i++ {
				_, err := md.Dial("test")
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	l1.Close()
	l2.Close()

	assert.Equal(t, uint32(100), counter)
}

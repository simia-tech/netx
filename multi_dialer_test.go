package netx_test

import (
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
		_, err := l.Accept()
		require.NoError(t, err)
		ch <- true
	}()

	p := static.NewProvider()
	p.Add("test", value.NewDial(l.Addr().Network(), l.Addr().String()))

	md, err := netx.NewMultiDialer(p, roundrobin.NewSelector())
	require.NoError(t, err)

	_, err = md.Dial("test")
	require.NoError(t, err)

	assert.True(t, <-ch)
	close(ch)
}

package netx_test

import (
	"net"
	"sort"
	"testing"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/assert"
)

func TestAddrsSortIPv4OverIPv6(t *testing.T) {
	one := mustParseAddr("127.0.0.1:1000")
	two := mustParseAddr("[::1]:1000")

	a := netx.Addrs{one, two}
	sort.Sort(a)

	assert.Equal(t, netx.Addrs{two, one}, a)
}

func TestAddrsSortIPv4ByNumber(t *testing.T) {
	one := mustParseAddr("127.0.0.1:1000")
	two := mustParseAddr("127.0.0.3:1000")
	three := mustParseAddr("127.0.0.2:1000")

	a := netx.Addrs{one, two, three}
	sort.Sort(a)

	assert.Equal(t, netx.Addrs{one, three, two}, a)
}

func TestAddrsSortIPv4ByPort(t *testing.T) {
	one := mustParseAddr("127.0.0.1:1000")
	two := mustParseAddr("127.0.0.1:3000")
	three := mustParseAddr("127.0.0.1:2000")

	a := netx.Addrs{one, two, three}
	sort.Sort(a)

	assert.Equal(t, netx.Addrs{one, three, two}, a)
}

func mustParseAddr(input string) net.Addr {
	addr, err := net.ResolveTCPAddr("tcp", input)
	if err != nil {
		panic(err)
	}
	return addr
}
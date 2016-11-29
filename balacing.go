package netx

import (
	"math/rand"
	"net"
)

type BalancingFn func([]net.Addr) (net.Addr, error)

func RandomBalancing(addrs []net.Addr) (net.Addr, error) {
	switch l := len(addrs); l {
	case 0:
		return nil, nil
	case 1:
		return addrs[0], nil
	default:
		return addrs[rand.Intn(l)], nil
	}
}

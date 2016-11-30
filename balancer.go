package netx

import (
	"math/rand"
	"net"
)

type BalancerFn func([]net.Addr) (net.Addr, error)

func RandomBalancer() BalancerFn {
	return func(addrs []net.Addr) (net.Addr, error) {
		switch l := len(addrs); l {
		case 0:
			return nil, nil
		case 1:
			return addrs[0], nil
		default:
			return addrs[rand.Intn(l)], nil
		}
	}
}

func RoundRobinBalancer() BalancerFn {
	index := 0
	return func(addrs []net.Addr) (net.Addr, error) {
		length := len(addrs)
		if length == 0 {
			return nil, nil
		}
		if index >= length {
			index = 0
		}
		addr := addrs[index]
		index++
		return addr, nil
	}
}

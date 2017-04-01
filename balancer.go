package netx

import (
	"math/rand"
	"net"
)

// BalancerFn defines the address select function for a balancing strategy.
type BalancerFn func([]net.Addr) (net.Addr, error)

// RandomBalancer returns a BalancerFn that randomly picks one of the provided
// addresses and returns it. If no address is provided, nil is returned.
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

// RoundRobinBalancer returns a BalancerFn that returns the provided addresses
// one after the other. If no address is provided, nil is returned.
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

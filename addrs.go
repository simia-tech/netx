package netx

import (
	"math/big"
	"net"
)

type Addrs []net.Addr

func (a Addrs) Len() int {
	return len(a)
}

func (a Addrs) Less(i, j int) bool {
	one, two := a[i], a[j]

	tcpAddrOne, ok := one.(*net.TCPAddr)
	if !ok {
		return false
	}
	tcpAddrTwo, ok := two.(*net.TCPAddr)
	if !ok {
		return false
	}

	intOne := big.NewInt(0)
	intOne.SetBytes(tcpAddrOne.IP)

	intTwo := big.NewInt(0)
	intTwo.SetBytes(tcpAddrTwo.IP)

	result := big.NewInt(0)
	result.Sub(intOne, intTwo)

	switch result.Sign() {
	case -1:
		return true
	case 0:
		return tcpAddrOne.Port < tcpAddrTwo.Port
	}
	return false
}

func (a Addrs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

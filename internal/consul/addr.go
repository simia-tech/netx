package consul

type addr struct {
	network string
	address string
}

func (a *addr) Network() string {
	return a.network
}

func (a *addr) String() string {
	return a.address
}

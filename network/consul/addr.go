package consul

type addr struct {
	address string
}

func (a *addr) Network() string {
	return "consul"
}

func (a *addr) String() string {
	return a.address
}

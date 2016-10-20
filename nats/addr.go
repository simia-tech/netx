package nats

type addr struct {
	net     string
	address string
}

func (a *addr) Network() string {
	return a.net
}

func (a addr) String() string {
	return a.address
}

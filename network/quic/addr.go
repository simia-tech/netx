package quic

type addr struct {
	address string
}

func (a *addr) Network() string {
	return "quic"
}

func (a addr) String() string {
	return a.address
}

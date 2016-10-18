package nats

type addr struct {
	network *network
	address string
}

func (a *addr) Network() string {
	return a.network.conn.Opts.Name
}

func (a addr) String() string {
	return a.address
}

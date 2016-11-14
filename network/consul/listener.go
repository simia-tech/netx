package consul

import (
	"net"

	"github.com/pkg/errors"

	"github.com/simia-tech/netx"
)

type listener struct {
	listener net.Listener
	consul   *consul
	id       string
	address  string
}

func init() {
	netx.RegisterListen("consul", Listen)
}

// Listen starts a local tcp listener and registers its address and port under
// the provided address to the consul instance that is specified in the provided options.
func Listen(address string, options *netx.Options) (net.Listener, error) {
	consul, err := newConsulFrom(options.Nodes)
	if err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", options.LocalAddress)
	if err != nil {
		return nil, err
	}

	id, err := consul.register(address, l.Addr())
	if err != nil {
		return nil, errors.Wrapf(err, "register local listener address [%s] at consul failed", l.Addr())
	}

	return &listener{
		listener: l,
		consul:   consul,
		id:       id,
		address:  address,
	}, nil
}

func (l *listener) Accept() (net.Conn, error) {
	return l.listener.Accept()
}

func (l *listener) Close() error {
	if err := l.consul.deregister(l.id); err != nil {
		return err
	}
	return l.listener.Close()
}

func (l *listener) Addr() net.Addr {
	return &addr{address: l.address}
}

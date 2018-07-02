package quic

import (
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/value"
)

type listener struct {
	listener quic.Listener
}

func init() {
	netx.RegisterListen("quic", Listen)
}

// Listen starts a listener at the provided address.
func Listen(address string, options *value.Options) (net.Listener, error) {
	l, err := quic.ListenAddr(address, options.TLSConfig, nil)
	if err != nil {
		return nil, err
	}

	return &listener{
		listener: l,
	}, nil
}

func (l *listener) Accept() (net.Conn, error) {
	session, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	stream, err := session.AcceptStream()
	if err != nil {
		return nil, err
	}

	return &conn{
		session: session,
		stream:  stream,
	}, nil
}

func (l *listener) Addr() net.Addr {
	return &addr{address: l.listener.Addr().String()}
}

func (l *listener) Close() error {
	return l.listener.Close()
}

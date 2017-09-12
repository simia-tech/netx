package quic

import (
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/simia-tech/netx"
)

type listener struct {
	listener quic.Listener
}

func init() {
	netx.RegisterListen("quic", Listen)
}

func Listen(address string, options *netx.Options) (net.Listener, error) {
	l, err := quic.ListenAddr(address, options.TLSConfig, &quic.Config{})
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
		stream: stream,
	}, nil
}

func (l *listener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *listener) Close() error {
	return l.listener.Close()
}

package quic

import (
	"log"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/simia-tech/netx"
)

type listener struct {
	listener    quic.Listener
	sessionChan chan quic.Session
}

func init() {
	netx.RegisterListen("quic", Listen)
}

func Listen(address string, options *netx.Options) (net.Listener, error) {
	sessionChan := make(chan quic.Session)

	l, err := quic.ListenAddr(address, &quic.Config{
		TLSConfig: options.TLSConfig,
		ConnState: func(session quic.Session, connState quic.ConnState) {
			log.Printf("conn state %v", connState)
			if connState == quic.ConnStateSecure {
				sessionChan <- session
			}
		},
	})
	if err != nil {
		return nil, err
	}

	go func() {
		if err := l.Serve(); err != nil {
			log.Println(err)
		}
	}()

	return &listener{
		listener:    l,
		sessionChan: sessionChan,
	}, nil
}

func (l *listener) Accept() (net.Conn, error) {
	session := <-l.sessionChan

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
	return l.listener.Addr()
}

func (l *listener) Close() error {
	return l.listener.Close()
}

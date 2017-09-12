package quic

import (
	"fmt"
	"io"
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/qerr"
	"github.com/simia-tech/netx"
)

type conn struct {
	session quic.Session
	stream  quic.Stream
}

func init() {
	netx.RegisterDial("quic", Dial)
}

func Dial(address string, options *netx.Options) (net.Conn, error) {
	session, err := quic.DialAddr(address, options.TLSConfig, &quic.Config{})
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return nil, err
	}

	return &conn{
		session: session,
		stream:  stream,
	}, nil
}

func (c *conn) Read(data []byte) (int, error) {
	n, err := c.stream.Read(data)
	if qErr, ok := err.(*qerr.QuicError); ok {
		switch qErr.ErrorCode {
		case qerr.PeerGoingAway:
			err = io.EOF
		case qerr.NetworkIdleTimeout:
			err = fmt.Errorf("timeout")
		}
	}
	return n, err
}

func (c *conn) Write(data []byte) (int, error) {
	n, err := c.stream.Write(data)
	return n, err
}

func (c *conn) Close() error {
	if err := c.stream.Close(); err != nil {
		return err
	}
	if err := c.session.Close(nil); err != nil {
		return err
	}
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return nil
}

func (c *conn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}

func (c *conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *conn) String() string {
	return fmt.Sprintf("(%s -> %s)", c.LocalAddr(), c.RemoteAddr())
}

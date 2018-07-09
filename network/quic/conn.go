package quic

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/qerr"
	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/value"
)

type conn struct {
	session quic.Session
	stream  quic.Stream
}

func init() {
	netx.RegisterDial("quic", Dial)
}

// Dial opens a connection to the provided address.
func Dial(ctx context.Context, address string, options *value.Options) (net.Conn, error) {
	session, err := quic.DialAddrContext(ctx, address, options.TLSConfig, nil)
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
			err = fmt.Errorf("read timeout")
		}
	}
	if err != nil && err.Error() == "deadline exceeded" {
		err = fmt.Errorf("read timeout")
	}
	return n, err
}

func (c *conn) Write(data []byte) (int, error) {
	return c.stream.Write(data)
}

func (c *conn) Close() error {
	if c.stream == nil {
		return nil
	}
	if err := c.stream.Close(); err != nil {
		return err
	}
	if c.session == nil {
		return nil
	}
	if err := c.session.Close(nil); err != nil {
		return err
	}
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	if c.session == nil {
		return nil
	}
	return c.session.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	if c.session == nil {
		return nil
	}
	return c.session.RemoteAddr()
}

func (c *conn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}

func (c *conn) String() string {
	return fmt.Sprintf("(%s -> %s)", c.LocalAddr(), c.RemoteAddr())
}

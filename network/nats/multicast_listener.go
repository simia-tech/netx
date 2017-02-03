package nats

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	n "github.com/nats-io/nats"
	"github.com/pkg/errors"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/model"
)

type multicastListener struct {
	inbox        string
	subscription *n.Subscription

	readDeadline time.Time

	readBuffer []byte
}

func init() {
	netx.RegisterListenMulticast("nats", ListenMulticast)
}

func ListenMulticast(address string, options *netx.Options) (net.Conn, error) {
	o := []n.Option{}
	if options.TLSConfig != nil {
		o = append(o, n.Secure(options.TLSConfig))
	}
	conn, err := n.Connect(strings.Join(options.Nodes, ","), o...)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	return newMulticastListener(conn, host)
}

func newMulticastListener(nc *n.Conn, inbox string) (*multicastListener, error) {
	subscription, err := nc.SubscribeSync(inbox)
	if err != nil {
		return nil, err
	}

	return &multicastListener{
		inbox:        inbox,
		subscription: subscription,
	}, nil
}

func (ml *multicastListener) Read(readBuffer []byte) (int, error) {
	if len(ml.readBuffer) > 0 {
		n := copy(readBuffer, ml.readBuffer)
		if n < len(ml.readBuffer) {
			ml.readBuffer = ml.readBuffer[n:]
		} else {
			ml.readBuffer = nil
		}
		return n, nil
	}

	if ml.subscription == nil {
		return 0, io.EOF
	}

	packet, err := ml.receivePacket()
	if err != nil {
		return 0, err
	}
	switch packet.Type {
	case model.Packet_DATA:
		n := copy(readBuffer, packet.Payload)
		if n < len(packet.Payload) {
			ml.readBuffer = packet.Payload[n:]
		}
		return n, nil
	case model.Packet_CLOSE:
		return 0, io.EOF
	default:
		return 0, errors.Errorf("expected DATA packet, got %s", packet.Type)
	}
}

func (ml *multicastListener) Write(data []byte) (int, error) {
	return 0, netx.ErrNotSupported
}

func (ml *multicastListener) Close() error {
	return nil
}

func (ml *multicastListener) LocalAddr() net.Addr {
	return &addr{net: "nats", address: ml.inbox}
}

func (ml *multicastListener) RemoteAddr() net.Addr {
	return &addr{net: "nats", address: "multi"}
}

func (ml *multicastListener) SetDeadline(t time.Time) error {
	ml.readDeadline = t
	return nil
}

func (ml *multicastListener) SetReadDeadline(t time.Time) error {
	ml.readDeadline = t
	return nil
}

func (ml *multicastListener) SetWriteDeadline(t time.Time) error {
	return netx.ErrNotSupported
}

func (ml *multicastListener) String() string {
	return fmt.Sprintf("(%s -> %s)", ml.LocalAddr(), ml.RemoteAddr())
}

func (ml *multicastListener) receivePacket() (*model.Packet, error) {
	if ml.readDeadline.IsZero() {
		return receivePacket(ml.subscription, endlessTimeout)
	}
	return receivePacket(ml.subscription, ml.readDeadline.Sub(time.Now()))
}

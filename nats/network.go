package nats

import (
	"net"
	"time"

	n "github.com/nats-io/nats"
	"github.com/pkg/errors"
	"github.com/simia-tech/netx/model"
)

type network struct {
	conn *n.Conn
}

func JoinNetwork(url string) (*network, error) {
	conn, err := n.Connect(url)
	if err != nil {
		return nil, err
	}
	return &network{
		conn: conn,
	}, nil
}

func (n *network) Listen(address string) (net.Listener, error) {
	return &listener{
		network: n,
		address: address,
	}, nil
}

func (n *network) Dial(address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	message, err := n.conn.Request(host, []byte{}, 2*time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "requesting address from [%s] failed", host)
	}

	packet, err := receivePacket(message.Data)
	if err != nil {
		return nil, err
	}
	if packet.Type != model.Packet_ACCEPT {
		return nil, errors.Errorf("unexpected packet type %s", packet.Type)
	}

	remoteInbox := string(packet.Payload)

	return newConn(n, message.Subject, remoteInbox)
}

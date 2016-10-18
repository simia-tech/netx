package nats

import (
	"net"
	"time"

	natsio "github.com/nats-io/nats"
	"github.com/pkg/errors"
	"github.com/simia-tech/netx/model"
)

type network struct {
	conn *natsio.Conn
}

func JoinNetwork(url string) (*network, error) {
	conn, err := natsio.Connect(url)
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
	message, err := n.conn.Request(address, []byte{}, 100*time.Millisecond)
	if err != nil {
		return nil, err
	}

	packet := &model.Packet{}
	if err := packet.UnmarshalBinary(message.Data); err != nil {
		return nil, err
	}

	if packet.Type != model.Packet_ACCEPT {
		return nil, errors.Errorf("unexpected packet type %s", packet.Type)
	}

	remoteInbox := string(packet.Payload)

	return newConn(n, message.Subject, remoteInbox)
}

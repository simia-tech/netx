package netx

import (
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/simia-tech/netx/nats"
)

func Dial(net, address string) (net.Conn, error) {
	switch {
	case strings.HasPrefix(net, "nats:"):
		return nats.Dial(net, address)
	default:
		return nil, errors.Errorf("unknown network [%s]", net)
	}
}

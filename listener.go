package netx

import (
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/simia-tech/netx/nats"
)

func Listen(net, address string) (net.Listener, error) {
	switch {
	case strings.HasPrefix(net, "nats:"):
		return nats.Listen(net, address)
	default:
		return nil, errors.Errorf("unknown network [%s]", net)
	}
}

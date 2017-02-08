package udp_test

import (
	"testing"

	_ "github.com/simia-tech/netx/network/nats"
	"github.com/simia-tech/netx/test"
)

var options = &test.Options{
	MulticastNetwork:         "udp",
	MulticastRequestAddress:  "224.0.0.1:2000",
	MulticastResponseAddress: "224.0.0.1:3000",
}

func TestMulticast(t *testing.T) {
	test.MulticastTest(t, options)
}

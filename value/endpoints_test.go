package value_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simia-tech/netx/value"
)

func TestEndpointsSortIPv6OverIPv4(t *testing.T) {
	one := value.MustParseEndpointURL("tcp://127.0.0.1:1000")
	two := value.MustParseEndpointURL("tcp://[::1]:1000")

	a := value.Endpoints{one, two}
	sort.Sort(a)

	assert.Equal(t, value.Endpoints{one, two}, a)
}

func TestEndpointsSortIPv4ByNumber(t *testing.T) {
	one := value.MustParseEndpointURL("tcp://127.0.0.1:1000")
	two := value.MustParseEndpointURL("tcp://127.0.0.3:1000")
	three := value.MustParseEndpointURL("tcp://127.0.0.2:1000")

	a := value.Endpoints{one, two, three}
	sort.Sort(a)

	assert.Equal(t, value.Endpoints{one, three, two}, a)
}

func TestEndpointsSortIPv4ByPort(t *testing.T) {
	one := value.MustParseEndpointURL("tcp://127.0.0.1:1000")
	two := value.MustParseEndpointURL("tcp://127.0.0.1:3000")
	three := value.MustParseEndpointURL("tcp://127.0.0.1:2000")

	a := value.Endpoints{one, two, three}
	sort.Sort(a)

	assert.Equal(t, value.Endpoints{one, three, two}, a)
}

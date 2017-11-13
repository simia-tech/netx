package value_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simia-tech/netx/value"
)

func TestDialsSortIPv6OverIPv4(t *testing.T) {
	one := value.MustParseDialURL("tcp://127.0.0.1:1000")
	two := value.MustParseDialURL("tcp://[::1]:1000")

	a := value.Dials{one, two}
	sort.Sort(a)

	assert.Equal(t, value.Dials{one, two}, a)
}

func TestDialsSortIPv4ByNumber(t *testing.T) {
	one := value.MustParseDialURL("tcp://127.0.0.1:1000")
	two := value.MustParseDialURL("tcp://127.0.0.3:1000")
	three := value.MustParseDialURL("tcp://127.0.0.2:1000")

	a := value.Dials{one, two, three}
	sort.Sort(a)

	assert.Equal(t, value.Dials{one, three, two}, a)
}

func TestDialsSortIPv4ByPort(t *testing.T) {
	one := value.MustParseDialURL("tcp://127.0.0.1:1000")
	two := value.MustParseDialURL("tcp://127.0.0.1:3000")
	three := value.MustParseDialURL("tcp://127.0.0.1:2000")

	a := value.Dials{one, two, three}
	sort.Sort(a)

	assert.Equal(t, value.Dials{one, three, two}, a)
}

package roundrobin

import (
	"sync"

	"github.com/simia-tech/netx/selector"
	"github.com/simia-tech/netx/value"
)

type roundrobin struct {
	index int
	mutex sync.Mutex
}

// NewSelector returns a new round robin selector.
func NewSelector() selector.Interface {
	return &roundrobin{}
}

func (rr *roundrobin) Select(endpoints value.Endpoints) (value.Endpoint, error) {
	if len(endpoints) < 1 {
		return nil, selector.ErrNoEndpoint
	}
	rr.mutex.Lock()
	if rr.index >= len(endpoints) {
		rr.index = 0
	}
	endpoint := endpoints[rr.index]
	rr.index++
	rr.mutex.Unlock()
	return endpoint, nil
}

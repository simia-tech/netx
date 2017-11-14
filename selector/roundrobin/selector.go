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

func (rr *roundrobin) Select(dials value.Endpoints) (value.Endpoint, error) {
	if len(dials) < 1 {
		return nil, selector.ErrNoEndpoint
	}
	rr.mutex.Lock()
	if rr.index >= len(dials) {
		rr.index = 0
	}
	dial := dials[rr.index]
	rr.index++
	rr.mutex.Unlock()
	return dial, nil
}

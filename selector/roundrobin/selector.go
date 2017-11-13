package roundrobin

import (
	"github.com/simia-tech/netx/selector"
	"github.com/simia-tech/netx/value"
)

type roundrobin struct {
	index int
}

// NewSelector returns a new round robin selector.
func NewSelector() selector.Interface {
	return &roundrobin{}
}

func (rr *roundrobin) Select(dials value.Dials) (value.Dial, error) {
	if len(dials) < 1 {
		return nil, selector.ErrNoDial
	}
	if rr.index >= len(dials) {
		rr.index = 0
	}
	dial := dials[rr.index]
	rr.index++
	return dial, nil
}

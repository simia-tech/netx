package static

import (
	"fmt"

	"github.com/simia-tech/netx/value"
)

// Static implements a static provider.
type Static struct {
	dials map[string]value.Dials
}

// NewProvider returns a new static provider.
func NewProvider() *Static {
	return &Static{dials: make(map[string]value.Dials)}
}

// Add adds the provided dial to the provider.
func (p *Static) Add(service string, dial value.Dial) {
	if d, ok := p.dials[service]; ok {
		p.dials[service] = append(d, dial)
	} else {
		p.dials[service] = value.Dials{dial}
	}
}

// Dials returns the Dials for the provided service.
func (p *Static) Dials(service string) (value.Dials, error) {
	dials, ok := p.dials[service]
	if !ok {
		return nil, fmt.Errorf("unknown service [%s]", service)
	}
	return dials, nil
}

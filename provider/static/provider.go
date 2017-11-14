package static

import (
	"fmt"
	"sync"

	"github.com/simia-tech/netx/value"
)

// Static implements a static provider.
type Static struct {
	dials map[string]value.Endpoints
	mutex sync.RWMutex
}

// NewProvider returns a new static provider.
func NewProvider() *Static {
	return &Static{dials: make(map[string]value.Endpoints)}
}

// Add adds the provided dial to the provider.
func (p *Static) Add(service string, dial value.Endpoint) {
	p.mutex.Lock()
	if d, ok := p.dials[service]; ok {
		p.dials[service] = append(d, dial)
	} else {
		p.dials[service] = value.Endpoints{dial}
	}
	p.mutex.Unlock()
}

// Endpoints returns the Endpoints for the provided service.
func (p *Static) Endpoints(service string) (value.Endpoints, error) {
	p.mutex.RLock()
	dials, ok := p.dials[service]
	p.mutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown service [%s]", service)
	}
	return dials, nil
}

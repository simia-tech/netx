package selector

import "github.com/simia-tech/netx/value"

// Interface defines the selector interface.
type Interface interface {
	Select(value.Endpoints) (value.Endpoint, error)
}

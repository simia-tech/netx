package filter

import "github.com/simia-tech/netx/value"

// Interface defines the filter interface.
type Interface interface {
	Filter(value.Endpoints) (value.Endpoints, error)
	Success(value.Endpoint) error
	Failure(value.Endpoint) error
}

package blacklist

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/simia-tech/netx/value"
)

// Filter implements a blacklist filter.
type Filter struct {
	backoffFn BackoffFn
	blacklist sync.Map
}

// NewFilter returns a new blacklist filter.
func NewFilter(boFn BackoffFn) *Filter {
	return &Filter{
		backoffFn: boFn,
	}
}

// Filter filters the provided endpoints and returns the result.
func (f *Filter) Filter(endpoints value.Endpoints) (value.Endpoints, error) {
	now := time.Now()
	result := value.Endpoints{}
	for _, endpoint := range endpoints {
		if v, ok := f.blacklist.Load(value.EndpointURL(endpoint)); ok && v.(*entry).isValid(f.backoffFn, now) {
			continue
		}
		result = append(result, endpoint)
	}
	return result, nil
}

// Success signals the filter that the provided endpoint succeeded.
func (f *Filter) Success(endpoint value.Endpoint) error {
	f.blacklist.Delete(value.EndpointURL(endpoint))
	return nil
}

// Failure signals the filter that the provided endpoint failed.
func (f *Filter) Failure(endpoint value.Endpoint) error {
	v, _ := f.blacklist.LoadOrStore(value.EndpointURL(endpoint), newEntry())
	v.(*entry).increaseFailure()
	return nil
}

type entry struct {
	createdAt time.Time
	failures  uint64
}

func newEntry() *entry {
	return &entry{
		createdAt: time.Now(),
		failures:  0,
	}
}

func (e *entry) increaseFailure() {
	atomic.AddUint64(&(e.failures), 1)
}

func (e *entry) isValid(boFn BackoffFn, now time.Time) bool {
	f := atomic.LoadUint64(&(e.failures))
	if t := e.createdAt.Add(boFn(f)); t.Before(now) {
		return false
	}
	return true
}

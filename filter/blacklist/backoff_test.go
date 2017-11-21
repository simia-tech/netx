package blacklist_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/simia-tech/netx/filter/blacklist"
	"github.com/stretchr/testify/assert"
)

func TestConstantBackoff(t *testing.T) {
	boFn := blacklist.ConstantBackoff(10 * time.Millisecond)
	assert.Equal(t, 10*time.Millisecond, boFn(0))
}

func TestLinearBackoff(t *testing.T) {
	boFn := blacklist.LinearBackoff(10, 10*time.Millisecond, 10, 5*time.Second)

	tcs := []struct {
		failures       uint64
		expectDuration time.Duration
	}{
		{0, 10 * time.Millisecond},
		{1, 10 * time.Millisecond},
		{5, 10 * time.Millisecond},
		{10, 10 * time.Millisecond},
		{12, 1008 * time.Millisecond},
		{14, 2006 * time.Millisecond},
		{16, 3004 * time.Millisecond},
		{18, 4002 * time.Millisecond},
		{20, 5 * time.Second},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Failure-%d", tc.failures), func(t *testing.T) {
			assert.Equal(t, tc.expectDuration, boFn(tc.failures))
		})
	}
}

package blacklist

import "time"

// BackoffFn defines a function that calculates the duration after which
// an endpoint is removed from the blacklist, based on the provided number
// of failures.
type BackoffFn func(uint64) time.Duration

// ConstantBackoff creates a BackoffFn that always returns the provided duration.
func ConstantBackoff(d time.Duration) BackoffFn {
	return func(_ uint64) time.Duration {
		return d
	}
}

// LinearBackoff creates a BackoffFn that raises the returned duration after some
// failures from the provided lower to upper bound.
func LinearBackoff(index uint64, lower time.Duration, count uint64, upper time.Duration) BackoffFn {
	return func(f uint64) time.Duration {
		if f <= index {
			return lower
		} else if f >= index+count {
			return upper
		}
		d := upper - lower
		return lower + time.Duration((uint64(d)/count)*(f-index))
	}
}

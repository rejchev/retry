package retry

import "time"

// LinearDelay is linear
// @param x - time.Second / time.Milisecond ...
func LinearDelay(x time.Duration) DelayFn {
	return func(i, n int) time.Duration {
		return time.Duration(i * n) * x
	}
}

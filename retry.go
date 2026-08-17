package retry

import (
	"context"
	"time"
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

// DelayFn is delay provider
type DelayFn func(i, n int) time.Duration

// ExecFn is execution func
type ExecFn func(context.Context, unsafe.Pointer) errnov1.Code

type RetryFn func(context.Context, unsafe.Pointer, ExecFn) errnov1.Code

// Retry is main try func
func Try(delayFn DelayFn, max int) RetryFn {

	if delayFn == nil {
		delayFn = LinearDelay(time.Microsecond)
	}

	if max < 1 {
		max = 1
	}

	return func(ctx context.Context, buff unsafe.Pointer, execFn ExecFn) errnov1.Code {
		if execFn == nil {
			return errnov1.EINVAL
		}

		var errno errnov1.Code
		var i int
		var now, last time.Time = time.Now(), time.Now()

		for {
			select {
			case <-ctx.Done():
				return errno

			default:
				if now = time.Now(); now.Before(last) {
					continue
				}

				if errno = execFn(ctx, buff); errnov1.SUCCESS(errno) {
					return errno
				}

				i++

				if i >= max {
					return errno
				}

				if delayFn != nil {
					last = last.Add(delayFn(i, max))
				}
			}
		}
	}
}

// TryOnce is try 1 time
func TryOnce(delayFn DelayFn) RetryFn {
	return Try(delayFn, 1)
}

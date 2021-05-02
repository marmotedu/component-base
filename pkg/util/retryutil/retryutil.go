package retryutil

import (
	"context"
	"fmt"
	"math"
	"time"
)

var RetryAbleErr = fmt.Errorf("retry")
var TimeoutErr = fmt.Errorf("timeout")

func RetryUntilTimeout(ctx *context.Context, interval time.Duration, timeout time.Duration, do func() error) error {
	err := do()
	if err == nil {
		return nil
	}

	if err != RetryAbleErr {
		return err
	}

	if timeout == 0 {
		timeout = time.Duration(math.MaxInt64)
	}

	t := time.NewTimer(timeout)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			return TimeoutErr
		case <-time.After(interval):
			err := do()
			if err == nil {
				return nil
			}

			if err != RetryAbleErr {
				return err
			}
		}
	}
}

package retry

import (
	"context"
	"time"
)

// Do executes an operation with up to three retries after 1, 3, and 5 seconds.
func Do(ctx context.Context, operation func(context.Context) error, isRetriable func(error) bool) error {
	delays := []time.Duration{
		time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := operation(ctx)
		if err == nil {
			return nil
		}

		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		if !isRetriable(err) || attempt >= len(delays) {
			return err
		}

		timer := time.NewTimer(delays[attempt])
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

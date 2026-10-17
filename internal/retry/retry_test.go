package retry

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	t.Run("success after connection error", func(t *testing.T) {
		calls := 0
		err := Do(context.Background(), func(context.Context) error {
			calls++
			if calls == 1 {
				return errors.New("connection failed")
			}
			return nil
		}, func(error) bool { return true })
		require.NoError(t, err)
		require.Equal(t, 2, calls)
	})
	t.Run("non retriable error", func(t *testing.T) {
		calls := 0
		want := errors.New("invalid metric")
		err := Do(context.Background(), func(context.Context) error {
			calls++
			return want
		}, func(error) bool { return false })
		require.ErrorIs(t, err, want)
		require.Equal(t, 1, calls)
	})
	t.Run("canceled before operation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		calls := 0
		err := Do(ctx, func(context.Context) error {
			calls++
			return nil
		}, func(error) bool { return true })
		require.ErrorIs(t, err, context.Canceled)
		require.Zero(t, calls)
	})
}

func TestDoExhausted(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		var attempts []time.Duration
		want := errors.New("connection failed")
		err := Do(context.Background(), func(context.Context) error {
			attempts = append(attempts, time.Since(start))
			return want
		}, func(error) bool { return true })
		require.ErrorIs(t, err, want)
		require.Equal(t, []time.Duration{0, time.Second, 4 * time.Second, 9 * time.Second}, attempts)
	})
}

func TestDoCanceledDuringWait(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		start := time.Now()
		calls := 0
		err := Do(ctx, func(context.Context) error {
			calls++
			return errors.New("connection failed")
		}, func(error) bool { return true })
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, 1, calls)
		require.Equal(t, 100*time.Millisecond, time.Since(start))
	})
}

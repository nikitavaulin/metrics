package pingservice

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type pingFunc func(context.Context) error

func (f pingFunc) PingDB(ctx context.Context) error { return f(ctx) }

func TestPingDB(t *testing.T) {
	t.Run("not configured", func(t *testing.T) {
		require.EqualError(t, New(nil).PingDB(context.Background()), "postgres is not configured")
	})
	for _, want := range []error{nil, errors.New("database unavailable")} {
		name := "success"
		if want != nil {
			name = "storage error"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			called := false
			s := New(pingFunc(func(got context.Context) error {
				called = true
				require.Equal(t, ctx, got)
				return want
			}))
			err := s.PingDB(ctx)
			require.True(t, called)
			if want == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, want)
			}
		})
	}
}

package postgresstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestIsRetriableError(t *testing.T) {
	connectionErr := &pgconn.PgError{Code: "08006"}
	for _, tt := range []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"connection", connectionErr, true},
		{"wrapped connection", fmt.Errorf("query: %w", connectionErr), true},
		{"serialization failure", &pgconn.PgError{Code: pgerrcode.SerializationFailure}, true},
		{"deadlock", &pgconn.PgError{Code: pgerrcode.DeadlockDetected}, true},
		{"transaction rollback", &pgconn.PgError{Code: pgerrcode.TransactionRollback}, true},
		{"wrapped serialization failure", fmt.Errorf("query: %w", &pgconn.PgError{Code: pgerrcode.SerializationFailure}), true},
		{"commit serialization failure", &commitError{err: &pgconn.PgError{Code: pgerrcode.SerializationFailure}}, true},
		{"commit deadlock", &commitError{err: &pgconn.PgError{Code: pgerrcode.DeadlockDetected}}, true},
		{"commit unknown result", &commitError{err: &pgconn.PgError{Code: pgerrcode.StatementCompletionUnknown}}, false},
		{"unique violation", &pgconn.PgError{Code: "23505"}, false},
		{"canceled", context.Canceled, false},
		{"deadline", context.DeadlineExceeded, false},
		{"commit connection", fmt.Errorf("update: %w", &commitError{err: connectionErr}), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsRetriableError(tt.err))
		})
	}
}

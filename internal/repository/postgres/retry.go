package postgresstorage

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type commitError struct {
	err error
}

func (e *commitError) Error() string { return e.err.Error() }
func (e *commitError) Unwrap() error { return e.err }

func IsRetriableError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	var commitErr *commitError
	if errors.As(err, &commitErr) {
		// Only a confirmed serialization failure or deadlock is safe to
		// retry after Commit; connection loss leaves the result unknown.
		return pgErr.Code == pgerrcode.SerializationFailure ||
			pgErr.Code == pgerrcode.DeadlockDetected
	}

	return pgerrcode.IsConnectionException(pgErr.Code) ||
		pgerrcode.IsTransactionRollback(pgErr.Code)
}

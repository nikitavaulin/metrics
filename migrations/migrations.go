package migrations

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed *.sql
var files embed.FS

// Up applies the embedded migrations before the server starts accepting requests.
func Up(ctx context.Context, dsn string) (err error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse migration database configuration: %w", err)
	}
	db := stdlib.OpenDB(*config)
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to migration database: %w", err)
	}

	source, err := iofs.New(files, ".")
	if err != nil {
		return fmt.Errorf("open embedded migrations: %w", err)
	}
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		_ = source.Close()
		return fmt.Errorf("create migration database driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		_ = source.Close()
		_ = driver.Close()
		return fmt.Errorf("initialize migrations: %w", err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if closeErr := errors.Join(sourceErr, databaseErr); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close migrations: %w", closeErr))
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

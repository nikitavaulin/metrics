package migrations

import (
	"io"
	"testing"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedMigrations(t *testing.T) {
	source, err := iofs.New(files, ".")
	require.NoError(t, err)
	defer source.Close()

	version, err := source.First()
	require.NoError(t, err)
	require.Equal(t, uint(1), version)

	up, _, err := source.ReadUp(version)
	require.NoError(t, err)
	defer up.Close()
	sql, err := io.ReadAll(up)
	require.NoError(t, err)
	require.Contains(t, string(sql), "CREATE TABLE metrics")

	down, _, err := source.ReadDown(version)
	require.NoError(t, err)
	defer down.Close()
	sql, err = io.ReadAll(down)
	require.NoError(t, err)
	require.Contains(t, string(sql), "DROP TABLE IF EXISTS metrics")
}

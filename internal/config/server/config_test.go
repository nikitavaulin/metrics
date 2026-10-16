package serverconfig

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseDSNFlags(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		envDSN *string
		want   *string
	}{
		{name: "not provided"},
		{name: "explicitly empty", args: []string{"-d", ""}},
		{name: "empty equals", args: []string{"-d="}},
		{name: "flag", args: []string{"-d", "postgres://localhost/metrics"}, want: stringPointer("postgres://localhost/metrics")},
		{name: "old sentinel is a value", args: []string{"-d", "empty"}, want: stringPointer("empty")},
		{name: "environment priority", args: []string{"-d", "flag-dsn"}, envDSN: stringPointer("env-dsn"), want: stringPointer("env-dsn")},
		{name: "empty flag preserves environment", args: []string{"-d", ""}, envDSN: stringPointer("env-dsn"), want: stringPointer("env-dsn")},
		{name: "empty environment", envDSN: stringPointer("")},
		{name: "empty environment uses flag", args: []string{"-d", "flag-dsn"}, envDSN: stringPointer(""), want: stringPointer("flag-dsn")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			originalFlags, originalArgs := flag.CommandLine, os.Args
			t.Cleanup(func() {
				flag.CommandLine, os.Args = originalFlags, originalArgs
			})
			flag.CommandLine = flag.NewFlagSet("server", flag.ContinueOnError)
			os.Args = append([]string{"server"}, tt.args...)
			cfg := Config{DatabaseDsn: tt.envDSN}
			cfg.parseFlags()
			require.Equal(t, tt.want, cfg.DatabaseDsn)
		})
	}
}

func stringPointer(value string) *string { return &value }

package agentconfig

import (
	"flag"
	"os"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/require"
)

func TestKeyConfiguration(t *testing.T) {
	ptr := func(s string) *string { return &s }
	for _, tt := range []struct {
		name   string
		args   []string
		envKey *string
		want   *string
	}{
		{name: "optional"},
		{name: "flag", args: []string{"-k=flag-key"}, want: ptr("flag-key")},
		{name: "environment", envKey: ptr("env-key"), want: ptr("env-key")},
		{name: "environment priority", args: []string{"-k=flag-key"}, envKey: ptr("env-key"), want: ptr("env-key")},
		{name: "empty flag", args: []string{"-k="}, want: ptr("")},
		{name: "empty environment", envKey: ptr(""), want: ptr("")},
		{name: "empty environment priority", args: []string{"-k=flag-key"}, envKey: ptr(""), want: ptr("")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			originalFlags, originalArgs := flag.CommandLine, os.Args
			t.Cleanup(func() { flag.CommandLine, os.Args = originalFlags, originalArgs })
			flag.CommandLine = flag.NewFlagSet("config", flag.ContinueOnError)
			os.Args = append([]string{"config"}, tt.args...)
			t.Setenv("KEY", "")
			if tt.envKey == nil {
				require.NoError(t, os.Unsetenv("KEY"))
			} else {
				t.Setenv("KEY", *tt.envKey)
			}
			var cfg Config
			require.NoError(t, env.Parse(&cfg))
			cfg.parseFlags()
			require.Equal(t, tt.want, cfg.Key)
		})
	}
}

func TestRateLimitConfiguration(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []string
		env     string
		want    int
		wantErr bool
	}{
		{name: "default", want: 1},
		{name: "flag", args: []string{"-l=3"}, want: 3},
		{name: "environment", env: "4", want: 4},
		{name: "environment priority", args: []string{"-l=3"}, env: "4", want: 4},
		{name: "zero flag", args: []string{"-l=0"}, wantErr: true},
		{name: "negative flag", args: []string{"-l=-1"}, wantErr: true},
		{name: "zero environment", env: "0", wantErr: true},
		{name: "negative environment", env: "-1", wantErr: true},
		{name: "invalid environment", env: "invalid", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			originalFlags, originalArgs := flag.CommandLine, os.Args
			t.Cleanup(func() { flag.CommandLine, os.Args = originalFlags, originalArgs })
			flag.CommandLine = flag.NewFlagSet("config", flag.ContinueOnError)
			os.Args = append([]string{"config"}, tt.args...)
			t.Setenv("ADDRESS", "localhost:8080")
			t.Setenv("RATE_LIMIT", tt.env)
			if tt.env == "" {
				require.NoError(t, os.Unsetenv("RATE_LIMIT"))
			}

			cfg, err := New()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, *cfg.RateLimit)
		})
	}
}

package serverconfig

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPServerConfig struct {
	Address string
}

func New(addr string) (*HTTPServerConfig, error) {
	if err := ValidateServerAddress(addr); err != nil {
		return nil, fmt.Errorf("invalid server address: %w", err)
	}
	return &HTTPServerConfig{
		Address: addr,
	}, nil
}

func ValidateServerAddress(addr string) error {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address struct got: %s, want: %s", addr, "<addr>:<port>")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to convert port: %w\n", err)
	}

	if !(port >= 1 && port <= 65535) {
		return fmt.Errorf("invalid port value. got: %d, want: %d..%d", port, 1, 65535)
	}

	return nil
}

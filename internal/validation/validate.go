package validation

import (
	"fmt"
	"strconv"
	"strings"
)

func ValidateServerAddress(addr string) error {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address struct got: %s, want: %s", addr, "<addr>:<port>")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to convert port: %w\n", err)
	}

	if err := ValidateIntInBounds(port, 1, 65535); err != nil {
		return fmt.Errorf("invalid port value: %w", err)
	}

	return nil
}

func ValidateIntInBounds(value, minVal, maxVal int) error {
	if !(minVal <= value && value <= maxVal) {
		return fmt.Errorf("value '%d' not in bounds %d..%d", value, minVal, maxVal)
	}
	return nil
}

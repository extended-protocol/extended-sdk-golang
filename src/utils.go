package sdk

import (
	"errors"
	"fmt"
	"strings"
)

func isHexString(s string) error {
	if s == "" {
		return errors.New("empty hex string")
	}
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if len(s) == 0 {
		return errors.New("empty hex after 0x")
	}
	// Validate hex characters
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return fmt.Errorf("invalid hex char %q", c)
		}
	}
	return nil
}

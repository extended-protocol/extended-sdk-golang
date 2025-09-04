package sdk

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func parseHexBigInt(s string) (*big.Int, error) {
	if s == "" {
		return nil, errors.New("empty hex string")
	}
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if len(s) == 0 {
		return nil, errors.New("empty hex after 0x")
	}
	// Validate hex characters
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return nil, fmt.Errorf("invalid hex char %q", c)
		}
	}
	n := new(big.Int)
	_, ok := n.SetString(s, 16)
	if !ok {
		return nil, errors.New("failed to parse hex")
	}
	return n, nil
}
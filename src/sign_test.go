package sdk

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoGetOrderHash(t *testing.T) {
	// Test with the same parameters as in main.go
	hash, err := GetOrderHash(
		"100", "0x2", "100",
		"0x1", "-156",
		"0x1", "74",
		"100", "123",
		"0x5d05989e9302dcebc74e241001e3e3ac3f4402ccf2f8e6f74b034b07ad6a904", "Perpetuals", "v0", "SN_SEPOLIA", "1",
	)

	if err != nil {
		t.Fatalf("GetOrderHash failed: %v", err)
	}

	expected := "0x4de4c009e0d0c5a70a7da0e2039fb2b99f376d53496f89d9f437e736add6b48"
	if hash != expected {
		t.Errorf("GetOrderHash returned incorrect hash.\nExpected: %s\nGot:      %s", expected, hash)
	}

	sig, err := SignMessage(hash, "0x1234def56789012345678901234567890123456789012345678901234567890")
	if err != nil {
		t.Fatalf("SignMessage failed: %v", err)
	}

	log.Printf("Signature: %s", sig)
}

func TestStarkPerpetualAccountSign(t *testing.T) {
	// Use a well-known private key for deterministic testing
	privateKeyHex := "0x1234def56789012345678901234567890123456789012345678901234567890"
	publicKeyHex := "0x5d05989e9302dcebc74e241001e3e3ac3f4402ccf2f8e6f74b034b07ad6a904"

	// Create StarkPerpetualAccount
	account, err := NewStarkPerpetualAccount(100, privateKeyHex, publicKeyHex, "test-api-key")
	if err != nil {
		t.Fatalf("Failed to create StarkPerpetualAccount: %v", err)
	}

	// Use a known message hash
	msgHash := "0x4de4c009e0d0c5a70a7da0e2039fb2b99f376d53496f89d9f437e736add6b48"

	// Sign the message
	r, s, err := account.Sign(msgHash)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	// Verify signature has expected format (should be hex string with r, s, v components)
	if r == nil || s == nil {
		t.Errorf("Signature components are nil")
	}

	assert.Equal(t, r.String(), "2744225103614379349530169149569415648483556705538760809691766060588698917266", "R does not match")
	assert.Equal(t, s.String(), "575134845329043509424821214199431073576156064822439379079045654927136672163", "S does not match")
}

package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

type EndpointConfig struct {
	APIBaseURL string
}

var (
	ErrAPIKeyNotSet       = errors.New("api key is not set")
	ErrStarkAccountNotSet = errors.New("stark account is not set")
)

// BaseModule provides common functionality for API modules.
type BaseModule struct {
	endpointConfig EndpointConfig
	apiKey         string
	starkAccount   *StarkPerpetualAccount
	httpClient     *http.Client
	clientTimeout  time.Duration
}

// NewBaseModule constructs a BaseModule with all fields explicitly provided.
// Pass nil for httpClient to allow lazy creation. Pass nil for starkAccount if intentionally absent.
func NewBaseModule(
	cfg EndpointConfig,
	apiKey string,
	starkAccount *StarkPerpetualAccount,
	httpClient *http.Client,
	clientTimeout time.Duration,
) *BaseModule {
	return &BaseModule{
		endpointConfig: cfg,
		apiKey:         apiKey,
		starkAccount:   starkAccount,
		httpClient:     httpClient,
		clientTimeout:  clientTimeout,
	}
}

func (m *BaseModule) EndpointConfig() EndpointConfig {
	return m.endpointConfig
}

func (m *BaseModule) APIKey() (string, error) {
	if m.apiKey == "" {
		return "", ErrAPIKeyNotSet
	}
	return m.apiKey, nil
}

func (m *BaseModule) StarkAccount() (*StarkPerpetualAccount, error) {
	if m.starkAccount == nil {
		return nil, ErrStarkAccountNotSet
	}
	return m.starkAccount, nil
}

func (m *BaseModule) HTTPClient() *http.Client {
	if m.httpClient == nil {
		m.httpClient = &http.Client{
			Timeout: m.clientTimeout,
		}
	}
	return m.httpClient
}

// Close analogous to closing aiohttp session.
func (m *BaseModule) Close() {
	if m.httpClient != nil {
		m.httpClient.CloseIdleConnections()
		m.httpClient = nil
	}
}

// GetURL builds a full URL with optional query params.
func (m *BaseModule) GetURL(path string, query map[string]string) (string, error) {
	full := m.endpointConfig.APIBaseURL + path
	u, err := url.Parse(full)
	if err != nil {
		return "", err
	}
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

// DoRequest performs an HTTP request and unmarshals the JSON response into the provided object
// This function deduplicates common HTTP request logic across the SDK
func (m *BaseModule) DoRequest(ctx context.Context, method, url string, body io.Reader, result interface{}) error {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Only set Content-Type if we have a request body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add API key authentication if available
	if apiKey, err := m.APIKey(); err == nil {
		req.Header.Set("X-API-Key", apiKey)
	}

	// Execute request
	client := m.HTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse JSON response into the provided result object
	if err := json.Unmarshal(responseBody, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

type StarkPerpetualAccount struct {
	vault      uint64
	privateKey string
	publicKey  string
	apiKey     string
}

// NewStarkPerpetualAccount constructs the account, validating hex inputs.
func NewStarkPerpetualAccount(vault uint64, privateKeyHex, publicKeyHex, apiKey string) (*StarkPerpetualAccount, error) {
	if err := isHexString(privateKeyHex); err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	if err := isHexString(publicKeyHex); err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	// Ensure that private key and public key have 0x prefix
	if len(privateKeyHex) < 2 || privateKeyHex[:2] != "0x" {
		return nil, fmt.Errorf("private key must start with 0x")
	}
	if len(publicKeyHex) < 2 || publicKeyHex[:2] != "0x" {
		return nil, fmt.Errorf("public key must start with 0x")
	}

	// Check that API key does not start with 0x
	if len(apiKey) >= 2 && apiKey[:2] == "0x" {
		return nil, fmt.Errorf("api key should not start with 0x")
	}

	return &StarkPerpetualAccount{
		vault:      vault,
		privateKey: privateKeyHex,
		publicKey:  publicKeyHex,
		apiKey:     apiKey,
	}, nil
}

// Vault returns the vault id.
func (s *StarkPerpetualAccount) Vault() uint64 { return s.vault }

// PublicKey returns the public key as a string.
func (s *StarkPerpetualAccount) PublicKey() string { return s.publicKey }

// APIKey returns the API key string.
func (s *StarkPerpetualAccount) APIKey() string { return s.apiKey }

// Sign delegates to SignFunc, returning (r,s).
func (stark *StarkPerpetualAccount) Sign(msgHash string) (*big.Int, *big.Int, error) {
	if msgHash == "" {
		return big.NewInt(0), big.NewInt(0), errors.New("msgHash is empty")
	}

	sig, err := SignMessage(msgHash, stark.privateKey)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}

	// Extract r, s from the signature string.
	// Signature is in the format of {r}{s}{v}, where r, s and v are 64 chars each (192 hex chars).
	r, isGoodR := big.NewInt(0).SetString(sig[:64], 16)
	s, isGoodS := big.NewInt(0).SetString(sig[64:128], 16)

	if !isGoodR || !isGoodS {
		return big.NewInt(0), big.NewInt(0), errors.New("big int setting failed")
	}

	return r, s, nil
}

package sdk

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

// Replace/remove these placeholders if real types already exist elsewhere.
type EndpointConfig struct {
	APIBaseURL string
}

var (
	ErrAPIKeyNotSet       = errors.New("api key is not set")
	ErrStarkAccountNotSet = errors.New("stark account is not set")
)

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

type StarkPerpetualAccount struct {
	vault      uint64
	privateKey *big.Int
	publicKey  *big.Int
	apiKey     string
}

// NewStarkPerpetualAccount constructs the account, validating hex inputs.
func NewStarkPerpetualAccount(vault uint64, privateKeyHex, publicKeyHex, apiKey string) (*StarkPerpetualAccount, error) {
	priv, err := parseHexBigInt(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	pub, err := parseHexBigInt(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}
	return &StarkPerpetualAccount{
		vault:      vault,
		privateKey: priv,
		publicKey:  pub,
		apiKey:     apiKey,
	}, nil
}

// Vault returns the vault id.
func (s *StarkPerpetualAccount) Vault() uint64 { return s.vault }

// PublicKey returns the public key as a string.
func (s *StarkPerpetualAccount) PublicKey() string { return s.publicKey.String() }

// APIKey returns the API key string.
func (s *StarkPerpetualAccount) APIKey() string { return s.apiKey }

// Sign delegates to SignFunc, returning (r,s).
func (stark *StarkPerpetualAccount) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	if msgHash == nil {
		return big.NewInt(0), big.NewInt(0), errors.New("msgHash is nil")
	}
	sig, err := SignMessage(msgHash.String(), stark.privateKey.String())
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}

	// Extract r, s from the signature string.
	// Signature is in the format of {r}{s}{v}, where r, s and v are 64 chars each (192 hex chars).
	r, _ := big.NewInt(0).SetString(sig[:64], 16)
	s, _ := big.NewInt(0).SetString(sig[64:128], 16)
	return r, s, nil
}

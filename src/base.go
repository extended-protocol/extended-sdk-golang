package sdk

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BaseModel struct {
	PrivateKey string
}

// Replace/remove these placeholders if real types already exist elsewhere.
type EndpointConfig struct {
	APIBaseURL string
}

type StarkPerpetualAccount struct {
	// TODO: fill in fields
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

// GetURL builds a full URL with optional query params and simple {param} replacement.
func (m *BaseModule) GetURL(path string, query map[string]string, pathParams map[string]string) (string, error) {
	full := m.endpointConfig.APIBaseURL + path
	for k, v := range pathParams {
		full = strings.ReplaceAll(full, "{"+k+"}", v)
	}
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

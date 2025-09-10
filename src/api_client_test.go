package sdk

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { load() }

func load() {
	wd, _ := os.Getwd()
	for {
		p := filepath.Join(wd, ".env")
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return
		}
		wd = parent
	}
}
func createTestClient() *APIClient {
	cfg := EndpointConfig{
		APIBaseURL: "https://api.starknet.sepolia.extended.exchange/api/v1",
	}

	apiKey := os.Getenv("TEST_API_KEY")
	vaultStr := os.Getenv("TEST_VAULT")
	vault, _ := strconv.ParseUint(vaultStr, 10, 64)
	publicKey := os.Getenv("TEST_PUBLIC_KEY")
	privateKey := os.Getenv("TEST_PRIVATE_KEY")
	account, err := NewStarkPerpetualAccount(vault, apiKey, publicKey, privateKey)

	if err != nil {
		panic("Failed to create StarkPerpetualAccount: " + err.Error())
	}

	return NewAPIClient(cfg, apiKey, account, 30*time.Second)
}

func TestAPIClient_GetMarkets_SingleValidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.GetMarkets(ctx, []string{"BTC-USD"})

	require.NoError(t, err, "Should not error when requesting BTC-USD market")
	require.Equal(t, len(markets), 1, "Should return one market for valid request")
}

func TestAPIClient_GetMarkets_MultipleValidMarkets(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()
	requestedMarkets := []string{"BTC-USD", "ETH-USD"}

	markets, err := client.GetMarkets(ctx, requestedMarkets)

	require.NoError(t, err, "Should not error when requesting multiple valid markets")
	t.Logf("Requested %v, got %d markets", requestedMarkets, len(markets))

	require.Equal(t, len(markets), len(requestedMarkets), "Should return correct number of markets")
}

func TestAPIClient_GetMarkets_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.GetMarkets(ctx, []string{"INVALID-MARKET-NAME"})

	require.Error(t, err, "Should error when requesting invalid market")
	assert.Equal(t, len(markets), 0, "Should return zero markets for invalid request")
}

func TestAPIClient_GetMarkets_ContextTimeout(t *testing.T) {
	client := createTestClient()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.GetMarkets(ctx, []string{"BTC-USD"})

	require.Error(t, err, "Should error when context times out")
	t.Logf("Got expected timeout error: %v", err)
}

func TestAPIClient_GetMarkets_NetworkError(t *testing.T) {
	// Create client with invalid URL
	cfg := EndpointConfig{
		APIBaseURL: "http://invalid-url-that-does-not-exist.com",
	}
	client := NewAPIClient(cfg, "", nil, 5*time.Second)
	ctx := context.Background()

	_, err := client.GetMarkets(ctx, []string{"BTC-USD"})

	require.Error(t, err, "Should error when network request fails")
	t.Logf("Got expected network error: %v", err)
}

func TestAPIClient_GetMarketFee_ValidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.GetMarketFee(ctx, "BTC-USD")

	require.NoError(t, err, "Should not error when requesting fees for BTC-USD market")
	require.Equal(t, len(fees), 1, "Should return one fee entry for valid market")
	t.Logf("Got %d fees for BTC-USD", len(fees))

	for _, fee := range fees {
		t.Logf("Fee: %+v", fee)
	}
}

func TestAPIClient_GetMarketFee_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.GetMarketFee(ctx, "INVALID-MARKET-NAME")

	if err != nil {
		t.Logf("Got error for invalid market (this might be expected): %v", err)
		return
	}

	// If no error, should return empty list or no matching fees
	assert.Error(t, err, "Should error when requesting fees for invalid market")
	assert.Equal(t, len(fees), 0, "Should return zero fees for invalid market")
}

func TestAPIClient_GetMarketFee_ContextTimeout(t *testing.T) {
	client := createTestClient()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.GetMarketFee(ctx, "BTC-USD")

	require.Error(t, err, "Should error when context times out")
	t.Logf("Got expected timeout error: %v", err)
}

func TestAPIClient_GetMarketFee_NetworkError(t *testing.T) {
	// Create client with invalid URL
	cfg := EndpointConfig{
		APIBaseURL: "http://invalid-url-that-does-not-exist.com",
	}
	client := NewAPIClient(cfg, "", nil, 5*time.Second)
	ctx := context.Background()

	_, err := client.GetMarketFee(ctx, "BTC-USD")

	require.Error(t, err, "Should error when network request fails")
	t.Logf("Got expected network error: %v", err)
}

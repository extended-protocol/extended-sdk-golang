package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestClient() *APIClient {
	cfg := EndpointConfig{
		APIBaseURL: "https://api.starknet.sepolia.extended.exchange/api/v1",
	}
	return NewAPIClient(cfg, "", nil, 30*time.Second)
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

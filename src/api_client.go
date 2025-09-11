package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// APIClient provides REST API functionality for perpetual trading
// It embeds BaseModule to reuse common functionality like HTTP client, auth, etc.
type APIClient struct {
	*BaseModule
}

// NewAPIClient creates a new API client instance
func NewAPIClient(
	cfg EndpointConfig,
	apiKey string,
	starkAccount *StarkPerpetualAccount,
	clientTimeout time.Duration,
) *APIClient {
	baseModule := NewBaseModule(cfg, apiKey, starkAccount, nil, clientTimeout)
	return &APIClient{
		BaseModule: baseModule,
	}
}

// ===== Market Data Operations =====

// MarketResponse represents the API response for market data
type MarketResponse struct {
	Data   []MarketModel `json:"data"`
	Status string        `json:"status"`
}

// GetMarkets retrieves all available markets from the API
func (c *APIClient) GetMarkets(ctx context.Context, market []string) ([]MarketModel, error) {
	// Build the URL manually to handle multiple market parameters correctly
	baseURL := c.BaseModule.EndpointConfig().APIBaseURL + "/info/markets"

	if len(market) > 0 {
		baseURL += "?market=" + market[0]
		for i := 1; i < len(market); i++ {
			baseURL += "&market=" + market[i]
		}
	}

	// Use the new DoRequest method to handle the HTTP request and JSON parsing
	var marketResponse MarketResponse
	if err := c.BaseModule.DoRequest(ctx, "GET", baseURL, nil, &marketResponse); err != nil {
		return nil, err
	}

	// Check API status
	if marketResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %s", marketResponse.Status)
	}

	return marketResponse.Data, nil
}

// ===== Fee Data Operations =====

// FeeResponse represents the API response for trading fees
type FeeResponse struct {
	Data   []TradingFeeModel `json:"data"`
	Status string            `json:"status"`
}

// GetMarketFee retrieves current trading fees for a specific market
func (c *APIClient) GetMarketFee(ctx context.Context, market string) ([]TradingFeeModel, error) {
	baseUrl, err := c.GetURL("/user/fees", map[string]string{"market": market})
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Use the new DoRequest method to handle the HTTP request and JSON parsing
	var feeResponse FeeResponse
	if err := c.BaseModule.DoRequest(ctx, "GET", baseUrl, nil, &feeResponse); err != nil {
		return nil, err
	}

	if feeResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", feeResponse.Status)
	}

	return feeResponse.Data, nil
}

// ===== Order Operations =====

// OrderRequest represents the complete order submission request
type OrderRequest struct {
	Order     PerpetualOrderModel `json:"order"`
	Signature string              `json:"signature,omitempty"` // Additional API-level signature if needed
	Timestamp int64               `json:"timestamp"`           // Request timestamp for replay protection
}

// OrderResponse represents the API response after order submission
type OrderResponse struct {
	Status string `json:"status"`
	Data   struct {
		OrderID    uint   `json:"id"`
		ExternalID string `json:"externalId"`
	}
}

// SubmitOrder submits a perpetual order to the trading API
func (c *APIClient) SubmitOrder(ctx context.Context, order *PerpetualOrderModel) (*OrderResponse, error) {
	// Validate order object is complete and properly signed
	if order == nil {
		return nil, fmt.Errorf("order is nil")
	}

	baseUrl, err := c.GetURL("/user/order", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Marshal the order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order to JSON: %w", err)
	}

	// Create a buffer with the JSON data
	jsonData := bytes.NewBuffer(orderJSON)

	// Use the new DoRequest method to handle the HTTP request and JSON parsing
	var orderResponse OrderResponse
	if err := c.BaseModule.DoRequest(ctx, "POST", baseUrl, jsonData, &orderResponse); err != nil {
		return nil, err
	}

	if orderResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", orderResponse.Status)
	}

	if orderResponse.Data.ExternalID != order.ID {
		return nil, fmt.Errorf("mismatched order ID in response: got %s, expected %s", orderResponse.Data.ExternalID, order.ID)
	}

	return &orderResponse, nil
}

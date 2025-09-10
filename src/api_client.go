package sdk

import (
	"context"
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
	Data    []TradingFeeModel `json:"data"`
	Status string              `json:"status"`
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
	OrderID   string `json:"orderId"`
	Status    string `json:"status"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// SubmitOrder submits a perpetual order to the trading API
func (c *APIClient) SubmitOrder(ctx context.Context, order *PerpetualOrderModel) (*OrderResponse, error) {
	// TODO: Implement logic:
	// 1. Validate order object is complete and properly signed
	// 2. Create OrderRequest wrapper with timestamp
	// 3. Serialize to JSON
	// 4. Create POST request to orders endpoint (e.g., "/v1/orders")
	// 5. Add required headers:
	//    - Authorization: Bearer {apiKey} or X-API-Key: {apiKey}
	//    - Content-Type: application/json
	//    - X-Stark-Key: {starkAccount.PublicKey()} (if required)
	// 6. Execute request with context
	// 7. Parse response and handle various HTTP status codes:
	//    - 200/201: Success
	//    - 400: Bad request (validation errors)
	//    - 401: Unauthorized (API key issues)
	//    - 429: Rate limited
	//    - 500: Server errors
	// 8. Return structured response or appropriate error
	return nil, fmt.Errorf("not implemented")
}

// GetOrderStatus retrieves the current status of an order
func (c *APIClient) GetOrderStatus(ctx context.Context, orderID string) (*OrderResponse, error) {
	// TODO: Implement logic:
	// 1. Build URL with order ID (e.g., "/v1/orders/{orderID}")
	// 2. Create GET request with authentication
	// 3. Parse response to get current order state
	return nil, fmt.Errorf("not implemented")
}

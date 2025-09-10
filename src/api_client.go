package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient provides REST API functionality for perpetual trading
// It embeds BaseModule to reuse common functionality like HTTP client, auth, etc.
type APIClient struct {
	*BaseModule
}

// NewAPIClient creates a new API client instance
// Justification: Constructor pattern ensures proper initialization and validation
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
// Justification: Separate response struct allows for API metadata (pagination, status, etc.)
type MarketResponse struct {
	Data   []MarketModel `json:"data"`
	Status string        `json:"status"`
}

// GetMarkets retrieves all available markets from the API
func (c *APIClient) GetMarkets(ctx context.Context, market []string) ([]MarketModel, error) {
	// Build the URL manually to handle multiple market parameters correctly
	baseURL := c.BaseModule.EndpointConfig().APIBaseURL + "/api/v1/info/markets"

	if len(market) > 0 {
		baseURL += "?market=" + market[0]
		for i := 1; i < len(market); i++ {
			baseURL += "&" + market[i]
		}
	}

	// Create GET request
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if apiKey, err := c.BaseModule.APIKey(); err == nil {
		req.Header.Set("X-API-Key", apiKey)
	}

	// Execute request
	client := c.BaseModule.HTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var marketResponse MarketResponse
	if err := json.Unmarshal(body, &marketResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check API status
	if marketResponse.Status != "ok" {
		return nil, fmt.Errorf("API returned error status: %s", marketResponse.Status)
	}

	return marketResponse.Data, nil
}

// GetMarket retrieves a specific market by name
// Justification: Often need specific market data for order operations
func (c *APIClient) GetMarket(ctx context.Context, marketName string) (*MarketModel, error) {
	// TODO: Implement logic:
	// 1. Build URL with market name parameter (e.g., "/v1/markets/{marketName}")
	// 2. Similar to GetMarkets but for single market
	// 3. Return pointer to allow nil for not-found cases
	return nil, fmt.Errorf("not implemented")
}

// ===== Fee Data Operations =====

// FeeResponse represents the API response for trading fees
// Justification: Fees can change dynamically and must be fetched fresh
type FeeResponse struct {
	Fees    []TradingFeeModel `json:"fees"`
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
}

// GetTradingFees retrieves current trading fees for all markets
// Justification: Fee information is required for accurate order cost calculation
func (c *APIClient) GetTradingFees(ctx context.Context) ([]TradingFeeModel, error) {
	// TODO: Implement logic:
	// 1. Build URL for fees endpoint (e.g., "/v1/fees")
	// 2. Create GET request with authentication
	// 3. Parse response similar to markets
	// 4. Return fee models for use in order calculations
	return nil, fmt.Errorf("not implemented")
}

// GetMarketFees retrieves trading fees for a specific market
// Justification: More efficient when only need fees for one market
func (c *APIClient) GetMarketFees(ctx context.Context, marketName string) (*TradingFeeModel, error) {
	// TODO: Implement logic:
	// 1. Build URL with market parameter (e.g., "/v1/fees/{marketName}")
	// 2. Similar to GetTradingFees but for single market
	return nil, fmt.Errorf("not implemented")
}

// ===== Order Operations =====

// OrderRequest represents the complete order submission request
// Justification: Wraps the order model with any additional API-specific fields
type OrderRequest struct {
	Order     PerpetualOrderModel `json:"order"`
	Signature string              `json:"signature,omitempty"` // Additional API-level signature if needed
	Timestamp int64               `json:"timestamp"`           // Request timestamp for replay protection
}

// OrderResponse represents the API response after order submission
// Justification: API responses typically include status, order ID, and error details
type OrderResponse struct {
	OrderID   string `json:"orderId"`
	Status    string `json:"status"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// SubmitOrder submits a perpetual order to the trading API
// Justification: Core trading functionality - converts order object to API call
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

// CancelOrder cancels an existing order by ID
// Justification: Essential trading functionality for order management
func (c *APIClient) CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	// TODO: Implement logic:
	// 1. Build URL with order ID (e.g., "/v1/orders/{orderID}")
	// 2. Create DELETE request with authentication
	// 3. Handle response similar to SubmitOrder
	return nil, fmt.Errorf("not implemented")
}

// GetOrderStatus retrieves the current status of an order
// Justification: Needed to track order lifecycle and execution
func (c *APIClient) GetOrderStatus(ctx context.Context, orderID string) (*OrderResponse, error) {
	// TODO: Implement logic:
	// 1. Build URL with order ID (e.g., "/v1/orders/{orderID}")
	// 2. Create GET request with authentication
	// 3. Parse response to get current order state
	return nil, fmt.Errorf("not implemented")
}

// ===== Helper Methods =====

// buildAuthenticatedRequest creates an HTTP request with proper authentication
// Justification: Centralized auth logic to ensure consistency across all API calls
func (c *APIClient) buildAuthenticatedRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	// TODO: Implement logic:
	// 1. Create HTTP request with context
	// 2. Add standard headers:
	//    - User-Agent
	//    - Content-Type (if body present)
	//    - API key authentication
	// 3. Add Stark account signature headers if required by API
	// 4. Return configured request
	return nil, fmt.Errorf("not implemented")
}

// executeRequest executes an HTTP request and returns the response body
// Justification: Centralized request execution with error handling and response processing
func (c *APIClient) executeRequest(req *http.Request) ([]byte, error) {
	// TODO: Implement logic:
	// 1. Execute request using HTTP client
	// 2. Check response status code
	// 3. Read response body
	// 4. Handle common HTTP errors:
	//    - Network errors
	//    - Timeout errors
	//    - HTTP status errors
	// 5. Return body bytes or appropriate error
	return nil, fmt.Errorf("not implemented")
}

// parseAPIError attempts to parse API error responses into Go errors
// Justification: API errors often have structured format that should be preserved
func (c *APIClient) parseAPIError(statusCode int, body []byte) error {
	// TODO: Implement logic:
	// 1. Try to parse body as JSON error response
	// 2. Extract error message and code
	// 3. Return structured error with context
	// 4. Fall back to generic HTTP error if parsing fails
	return fmt.Errorf("API error: %d", statusCode)
}

// ===== Configuration and Validation =====

// ValidateConnection tests the API connection and authentication
// Justification: Useful for setup validation and health checks
func (c *APIClient) ValidateConnection(ctx context.Context) error {
	// TODO: Implement logic:
	// 1. Make simple authenticated request (e.g., GET /v1/health)
	// 2. Verify API key works
	// 3. Verify Stark account is properly configured
	// 4. Return error if any validation fails
	return fmt.Errorf("not implemented")
}

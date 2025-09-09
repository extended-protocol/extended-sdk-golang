package sdk

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// Test constants
const (
	TestPrivateKeyHex = "0x7a7ff6fd3cab02ccdcd4a572563f5976f8976899b03a39773795a3c486d4986"
	TestPublicKeyHex  = "0x61c5e7e8339b7d56f197f54ea91b776776690e3232313de0f2ecbd0ef76f466"
	TestVaultID       = 10002
	TestAPIKey        = "test-api-key"
	TestNonce         = 1473459052
)

// Test helper functions
func createTestAccount() (*StarkPerpetualAccount, error) {
	return NewStarkPerpetualAccount(TestVaultID, TestPrivateKeyHex, TestPublicKeyHex, TestAPIKey)
}

func createTestBTCUSDMarket() MarketModel {
	return MarketModel{
		Name:                     "BTC-USD",
		AssetName:                "BTC",
		AssetPrecision:           8,
		CollateralAssetName:      "USD",
		CollateralAssetPrecision: 6,
		Active:                   true,
		L2Config: L2ConfigModel{
			Type:                 "perpetual",
			CollateralID:         "0x31857064564ed0ff978e687456963cba09c2c6985d8f9300a1de4962fafa054",
			CollateralResolution: 1000000, // 6 decimals
			SyntheticID:          "0x4254432d3600000000000000000000",
			SyntheticResolution:  1000000, // 6 decimals
		},
	}
}

func createTestStarknetDomain() StarknetDomain {
	return StarknetDomain{
		Name:     "Perpetuals",
		Version:  "v0",
		ChainID:  "SN_SEPOLIA",
		Revision: "1",
	}
}

func createTestFrozenTime() time.Time {
	return time.Date(2024, 1, 5, 1, 8, 57, 0, time.UTC)
}

// OrdersTestSuite defines the test suite
type OrdersTestSuite struct {
	suite.Suite
	account        *StarkPerpetualAccount
	market         MarketModel
	starknetDomain StarknetDomain
	frozenTime     time.Time
	nonce          int
}

// SetupTest runs before each test
func (suite *OrdersTestSuite) SetupTest() {
	var err error
	suite.account, err = createTestAccount()
	suite.Require().NoError(err)

	suite.market = createTestBTCUSDMarket()
	suite.starknetDomain = createTestStarknetDomain()
	suite.frozenTime = createTestFrozenTime()
	suite.nonce = TestNonce
}

func (suite *OrdersTestSuite) TestCreateSellOrderWithDefaultExpiration() {
	// Create order parameters
	params := CreateOrderObjectParams{
		Market:                   suite.market,
		Account:                  *suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     OrderSideSell,
		Signer:                   suite.account.Sign,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               nil,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	// Create the order
	order, err := CreateOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	// Convert order to JSON for comparison
	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	// Parse JSON into a map for easier comparison
	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	// Expected JSON structure (matching Python test output)
	expectedOrder := map[string]interface{}{
		"id":                       "529621978301228831750156704671293558063128025271079340676658105549022202327",
		"market":                   "BTC-USD",
		"type":                     "limit",
		"side":                     "sell",
		"qty":                      "0.001",
		"price":                    "43445.1168",
		"reduceOnly":               false,
		"postOnly":                 false,
		"timeInForce":              "GTT",
		"fee":                      "0.0005",
		"nonce":                    "1473459052",
		"selfTradeProtectionLevel": "ACCOUNT",
		"cancelId":                 nil,
		"trigger":    nil,
		"tpSlType":   nil,
		"takeProfit": nil,
		"stopLoss":   nil,
		"builderFee": nil,
		"builderId":  nil,
	}

	// Assert JSON structure matches expected (excluding id since it's generated)
	suite.Equal(expectedOrder["market"], actualOrder["market"])
	suite.Equal(expectedOrder["type"], actualOrder["type"])
	suite.Equal(expectedOrder["side"], actualOrder["side"])
	suite.Equal(expectedOrder["qty"], actualOrder["qty"])
	suite.Equal(expectedOrder["price"], actualOrder["price"])
	suite.Equal(expectedOrder["reduceOnly"], actualOrder["reduceOnly"])
	suite.Equal(expectedOrder["postOnly"], actualOrder["postOnly"])
	suite.Equal(expectedOrder["timeInForce"], actualOrder["timeInForce"])
	suite.Equal(expectedOrder["fee"], actualOrder["fee"])
	suite.Equal(expectedOrder["nonce"], actualOrder["nonce"])
	suite.Equal(expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	suite.Equal(expectedOrder["cancelId"], actualOrder["cancelId"])
	suite.Equal(expectedOrder["trigger"], actualOrder["trigger"])
	suite.Equal(expectedOrder["tpSlType"], actualOrder["tpSlType"])
	suite.Equal(expectedOrder["takeProfit"], actualOrder["takeProfit"])
	suite.Equal(expectedOrder["stopLoss"], actualOrder["stopLoss"])
	suite.Equal(expectedOrder["builderFee"], actualOrder["builderFee"])
	suite.Equal(expectedOrder["builderId"], actualOrder["builderId"])

	// Verify ID is not empty (it's generated dynamically)
	suite.NotEmpty(actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCreateSellOrder() {
	// Set expiry time (1 hour from frozen time = 1704420537000 milliseconds)
	expiryTime := suite.frozenTime.Add(1 * time.Hour)

	// Create order parameters
	params := CreateOrderObjectParams{
		Market:                   suite.market,
		Account:                  *suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     OrderSideSell,
		Signer:                   suite.account.Sign,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               &expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	// Create the order
	order, err := CreateOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	// Convert order to JSON for comparison
	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	// Parse JSON into a map for easier comparison
	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	// Expected JSON structure (matching Python test output)
	expectedOrder := map[string]interface{}{
		"id":                       "529621978301228831750156704671293558063128025271079340676658105549022202327",
		"market":                   "BTC-USD",
		"type":                     "limit",
		"side":                     "sell",
		"qty":                      "0.001",
		"price":                    "43445.1168",
		"reduceOnly":               false,
		"postOnly":                 false,
		"timeInForce":              "GTT",
		"expiryEpochMillis":        float64(1704420537000), // JSON numbers become float64
		"fee":                      "0.0005",
		"nonce":                    "1473459052",
		"selfTradeProtectionLevel": "ACCOUNT",
		"cancelId":                 nil,
		"settlement": map[string]interface{}{
			"signature": map[string]interface{}{
				"r": "0x3d17d8b9652e5f60d40d079653cfa92b1065ea8cf159609a3c390070dcd44f7",
				"s": "0x76a6deccbc84ac324f695cfbde80e0ed62443e95f5dcd8722d12650ccc122e5",
			},
			"starkKey":           TestPublicKeyHex,
			"collateralPosition": "10002",
		},
		"trigger":    nil,
		"tpSlType":   nil,
		"takeProfit": nil,
		"stopLoss":   nil,
		"builderFee": nil,
		"builderId":  nil,
	}

	// Assert JSON structure matches expected (excluding id since it's generated)
	suite.Equal(expectedOrder["market"], actualOrder["market"])
	suite.Equal(expectedOrder["type"], actualOrder["type"])
	suite.Equal(expectedOrder["side"], actualOrder["side"])
	suite.Equal(expectedOrder["qty"], actualOrder["qty"])
	suite.Equal(expectedOrder["price"], actualOrder["price"])
	suite.Equal(expectedOrder["reduceOnly"], actualOrder["reduceOnly"])
	suite.Equal(expectedOrder["postOnly"], actualOrder["postOnly"])
	suite.Equal(expectedOrder["timeInForce"], actualOrder["timeInForce"])
	suite.Equal(expectedOrder["expiryEpochMillis"], actualOrder["expiryEpochMillis"])
	suite.Equal(expectedOrder["fee"], actualOrder["fee"])
	suite.Equal(expectedOrder["nonce"], actualOrder["nonce"])
	suite.Equal(expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	suite.Equal(expectedOrder["cancelId"], actualOrder["cancelId"])
	suite.Equal(expectedOrder["settlement"], actualOrder["settlement"])
	suite.Equal(expectedOrder["trigger"], actualOrder["trigger"])
	suite.Equal(expectedOrder["tpSlType"], actualOrder["tpSlType"])
	suite.Equal(expectedOrder["takeProfit"], actualOrder["takeProfit"])
	suite.Equal(expectedOrder["stopLoss"], actualOrder["stopLoss"])
	suite.Equal(expectedOrder["builderFee"], actualOrder["builderFee"])
	suite.Equal(expectedOrder["builderId"], actualOrder["builderId"])

	// Verify ID is not empty (it's generated dynamically)
	suite.NotEmpty(actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCreateBuyOrderWithClientProtection() {
	// Set expiry time (1 hour from frozen time)
	// @freeze_time("2024-01-05 01:08:56.860694")
	expiryTime := time.Date(2024, 1, 5, 1, 8, 56, 860694000, time.UTC).Add(14 * 24 * time.Hour)

	// Create order parameters for buy order
	params := CreateOrderObjectParams{
		Market:                   suite.market,
		Account:                  *suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     OrderSideBuy,
		Signer:                   suite.account.Sign,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               &expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	// Create the order
	order, err := CreateOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	// Convert order to JSON for comparison
	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	// Parse JSON into a map for easier comparison
	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	// Expected JSON structure for buy order
	expectedOrder := map[string]interface{}{
		"market":                   "BTC-USD",
		"type":                     "limit",
		"side":                     "buy",
		"qty":                      "0.001",
		"price":                    "43445.1168",
		"reduceOnly":               false,
		"postOnly":                 false,
		"timeInForce":              "GTT",
		"expiryEpochMillis":        float64(1705626536861),
		"fee":                      "0.0005",
		"nonce":                    "1473459052",
		"selfTradeProtectionLevel": "CLIENT",
		"cancelId":                 nil,
		"settlement": map[string]interface{}{
			"signature": map[string]interface{}{
				"r": "0xa55625c7d5f1b85bed22556fc805224b8363074979cf918091d9ddb1403e13",
				"s": "0x504caf634d859e643569743642ccf244434322859b2421d76f853af43ae7a46",
			},
			"starkKey":           TestPublicKeyHex,
			"collateralPosition": "10002",
		},
		"trigger":    nil,
		"tpSlType":   nil,
		"takeProfit": nil,
		"stopLoss":   nil,
		"builderFee": nil,
		"builderId":  nil,
	}

	// Assert key fields match expected
	suite.Equal(expectedOrder["market"], actualOrder["market"])
	suite.Equal(expectedOrder["type"], actualOrder["type"])
	suite.Equal(expectedOrder["side"], actualOrder["side"])
	suite.Equal(expectedOrder["qty"], actualOrder["qty"])
	suite.Equal(expectedOrder["price"], actualOrder["price"])
	suite.Equal(expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	suite.Equal(expectedOrder["settlement"], actualOrder["settlement"])
	suite.NotEmpty(actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCancelPreviousOrder() {
	// Set expiry time (1 hour from frozen time)
	expiryTime := suite.frozenTime.Add(1 * time.Hour)
	previousOrderID := "previous_custom_id"

	// Create order parameters with previous order ID
	params := CreateOrderObjectParams{
		Market:                   suite.market,
		Account:                  *suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     OrderSideBuy,
		Signer:                   suite.account.Sign,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               &expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  &previousOrderID,
		OrderExternalID:          nil,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	// Create the order
	order, err := CreateOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	// Convert order to JSON for comparison
	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	// Parse JSON into a map for easier comparison
	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	// Assert cancelId is set correctly
	suite.Equal(previousOrderID, actualOrder["cancelId"])
}

func (suite *OrdersTestSuite) TestExternalOrderID() {
	// Set expiry time (1 hour from frozen time)
	expiryTime := suite.frozenTime.Add(1 * time.Hour)
	customOrderID := "custom_id"

	// Create order parameters with custom order ID
	params := CreateOrderObjectParams{
		Market:                   suite.market,
		Account:                  *suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     OrderSideBuy,
		Signer:                   suite.account.Sign,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               &expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          &customOrderID,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	// Create the order
	order, err := CreateOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	// Convert order to JSON for comparison
	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	// Parse JSON into a map for easier comparison
	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	// Assert custom ID is set correctly
	suite.Equal(customOrderID, actualOrder["id"])
}

// TestOrdersTestSuite runs the test suite
func TestOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersTestSuite))
}

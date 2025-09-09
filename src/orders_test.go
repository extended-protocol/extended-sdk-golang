package sdk

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSellOrderWithDefaultExpiration(t *testing.T) {
	// Mock frozen time: "2024-01-05 01:08:57"
	frozenTime := time.Date(2024, 1, 5, 1, 8, 57, 0, time.UTC)
	
	// Mock frozen nonce
	frozenNonce := 1473459052
	
	// Create mock trading account (similar to create_trading_account fixture)
	privateKeyHex := "0x1234def56789012345678901234567890123456789012345678901234567890"
	publicKeyHex := "0x61c5e7e8339b7d56f197f54ea91b776776690e3232313de0f2ecbd0ef76f466"
	account, err := NewStarkPerpetualAccount(10002, privateKeyHex, publicKeyHex, "test-api-key")
	require.NoError(t, err)
	
	// Create mock BTC-USD market (similar to create_btc_usd_market fixture)
	btcUsdMarket := MarketModel{
		Name:                      "BTC-USD",
		AssetName:                 "BTC",
		AssetPrecision:           8,
		CollateralAssetName:       "USD",
		CollateralAssetPrecision: 6,
		Active:                   true,
		L2Config: L2ConfigModel{
			Type:                 "perpetual",
			CollateralID:         "0x1", // Mock collateral asset ID
			CollateralResolution: 1000000, // 6 decimals
			SyntheticID:          "0x2", // Mock synthetic asset ID
			SyntheticResolution:  100000000, // 8 decimals
		},
	}
	
	// Mock Starknet domain (similar to STARKNET_TESTNET_CONFIG.starknet_domain)
	starknetDomain := StarknetDomain{
		Name:     "Perpetuals",
		Version:  "v0",
		ChainID:  "SN_SEPOLIA",
		Revision: "1",
	}
	
	// Mock signer function that returns expected signature values
	signer := account.Sign
	
	// Set expiry time (1 hour from frozen time = 1704420537000 milliseconds)
	expiryTime := frozenTime.Add(1 * time.Hour)
	
	// Create order parameters
	params := CreateOrderObjectParams{
		Market:                   btcUsdMarket,
		Account:                  *account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     OrderSideSell,
		Signer:                   signer,
		PublicKey:                123456, // Mock public key
		StarknetDomain:           starknetDomain,
		ExpireTime:               &expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionAccount,
		Nonce:                    &frozenNonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}
	
	// Create the order
	order, err := CreateOrderObject(params)
	require.NoError(t, err)
	require.NotNil(t, order)
	
	// Convert order to JSON for comparison
	orderJSON, err := json.Marshal(order)
	require.NoError(t, err)
	
	// Parse JSON into a map for easier comparison
	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	require.NoError(t, err)
	
	// Expected JSON structure (matching Python test output)
	expectedOrder := map[string]interface{}{
		"id":                       "529621978301228831750156704671293558063128025271079340676658105549022202327",
		"market":                   "BTC-USD",
		"type":                     "limit",
		"side":                     "sell", 
		"qty":                      "0.00100000",
		"price":                    "43445.11680000",
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
			"starkKey":           publicKeyHex,
			"collateralPosition": "10002",
		},
		"trigger":     nil,
		"tpSlType":    nil,
		"takeProfit":  nil,
		"stopLoss":    nil,
		"builderFee":  nil,
		"builderId":   nil,
	}
	
	// Assert JSON structure matches expected (excluding id since it's generated)
	assert.Equal(t, expectedOrder["market"], actualOrder["market"])
	assert.Equal(t, expectedOrder["type"], actualOrder["type"])
	assert.Equal(t, expectedOrder["side"], actualOrder["side"])
	assert.Equal(t, expectedOrder["qty"], actualOrder["qty"])
	assert.Equal(t, expectedOrder["price"], actualOrder["price"])
	assert.Equal(t, expectedOrder["reduceOnly"], actualOrder["reduceOnly"])
	assert.Equal(t, expectedOrder["postOnly"], actualOrder["postOnly"])
	assert.Equal(t, expectedOrder["timeInForce"], actualOrder["timeInForce"])
	assert.Equal(t, expectedOrder["expiryEpochMillis"], actualOrder["expiryEpochMillis"])
	assert.Equal(t, expectedOrder["fee"], actualOrder["fee"])
	assert.Equal(t, expectedOrder["nonce"], actualOrder["nonce"])
	assert.Equal(t, expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	assert.Equal(t, expectedOrder["cancelId"], actualOrder["cancelId"])
	assert.Equal(t, expectedOrder["settlement"], actualOrder["settlement"])
	assert.Equal(t, expectedOrder["trigger"], actualOrder["trigger"])
	assert.Equal(t, expectedOrder["tpSlType"], actualOrder["tpSlType"])
	assert.Equal(t, expectedOrder["takeProfit"], actualOrder["takeProfit"])
	assert.Equal(t, expectedOrder["stopLoss"], actualOrder["stopLoss"])
	assert.Equal(t, expectedOrder["builderFee"], actualOrder["builderFee"])
	assert.Equal(t, expectedOrder["builderId"], actualOrder["builderId"])
	
	// Verify ID is not empty (it's generated dynamically)
	assert.NotEmpty(t, actualOrder["id"])
}
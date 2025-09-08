package sdk

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type OrderType string

const (
	OrderTypeLimit       OrderType = "limit"
	OrderTypeMarket      OrderType = "market"
	OrderTypeConditional OrderType = "conditional"
	OrderTypeTpsl        OrderType = "tpsl"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// TimeInForce represents the time-in-force setting
type TimeInForce string

const (
	TimeInForceGTT TimeInForce = "GTT" // Good till time
	TimeInForceFOK TimeInForce = "FOK" // Fill or kill
	TimeInForceIOC TimeInForce = "IOC" // Immediate or cancel
)

type SelfTradeProtectionLevel string

const (
	SelfTradeProtectionDisabled SelfTradeProtectionLevel = "DISABLED"
	SelfTradeProtectionAccount  SelfTradeProtectionLevel = "ACCOUNT"
	SelfTradeProtectionClient   SelfTradeProtectionLevel = "CLIENT"
)

type TriggerPriceType string

const (
	TriggerPriceTypeLast  TriggerPriceType = "last"
	TriggerPriceTypeMid   TriggerPriceType = "mid"
	TriggerPriceTypeMark  TriggerPriceType = "mark"
	TriggerPriceTypeIndex TriggerPriceType = "index"
)

type TriggerDirection string

const (
	TriggerDirectionUp   TriggerDirection = "up"
	TriggerDirectionDown TriggerDirection = "down"
)

// ExecutionPriceType represents the type of price used for order execution
type ExecutionPriceType string

const (
	ExecutionPriceTypeLimit  ExecutionPriceType = "limit"
	ExecutionPriceTypeMarket ExecutionPriceType = "market"
)

// TpSlType represents the TPSL type determining order size
type TpSlType string

const (
	TpSlTypeOrder    TpSlType = "order"
	TpSlTypePosition TpSlType = "position"
)

// Signature represents a cryptographic signature
type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
}

type Settlement struct {
	Signature          Signature `json:"signature"`
	StarkKey           string    `json:"starkKey"`
	CollateralPosition string    `json:"collateralPosition"`
}

type ConditionalTrigger struct {
	TriggerPrice       string             `json:"triggerPrice"`
	TriggerPriceType   TriggerPriceType   `json:"triggerPriceType"`
	Direction          TriggerDirection   `json:"direction"`
	ExecutionPriceType ExecutionPriceType `json:"executionPriceType"`
}

// TpSlTrigger represents take profit or stop loss trigger settings
type TpSlTrigger struct {
	TriggerPrice     string             `json:"triggerPrice"`
	TriggerPriceType TriggerPriceType   `json:"triggerPriceType"`
	Price            string             `json:"price"`
	PriceType        ExecutionPriceType `json:"priceType"`
	Settlement       Settlement         `json:"settlement"`
}

type PerpetualOrderModel struct {
	ID                       string                   `json:"id"`
	Market                   string                   `json:"market"`
	Type                     OrderType                `json:"type"`
	Side                     OrderSide                `json:"side"`
	Qty                      string                   `json:"qty"`
	Price                    string                   `json:"price"`
	TimeInForce              TimeInForce              `json:"timeInForce"`
	ExpiryEpochMillis        int64                    `json:"expiryEpochMillis"`
	Fee                      string                   `json:"fee"`
	Nonce                    string                   `json:"nonce"`
	Settlement               Settlement               `json:"settlement"`
	ReduceOnly               bool                     `json:"reduceOnly"`
	PostOnly                 bool                     `json:"postOnly"`
	SelfTradeProtectionLevel SelfTradeProtectionLevel `json:"selfTradeProtectionLevel"`
	Trigger                  *ConditionalTrigger      `json:"trigger,omitempty"`
	TpSlType                 *TpSlType                `json:"tpSlType,omitempty"`
	TakeProfit               *TpSlTrigger             `json:"takeProfit,omitempty"`
	StopLoss                 *TpSlTrigger             `json:"stopLoss,omitempty"`
	BuilderFee               *string                  `json:"builderFee,omitempty"`
	BuilderID                *int                     `json:"builderId,omitempty"`
	CancelID                 *string                  `json:"cancelId,omitempty"`
}

// CreateOrderObjectParams represents the parameters for creating an order object
type CreateOrderObjectParams struct {
	Market                   MarketModel
	Account                  StarkPerpetualAccount
	SyntheticAmount          decimal.Decimal
	Price                    decimal.Decimal
	Side                     OrderSide
	Signer                   func(string) (*big.Int, *big.Int) // Function that takes string and returns two values
	PublicKey                int
	StarknetDomain           StarknetDomain
	ExpireTime               *time.Time
	PostOnly                 bool
	PreviousOrderExternalID  *string
	OrderExternalID          *string
	TimeInForce              TimeInForce
	SelfTradeProtectionLevel SelfTradeProtectionLevel
	Nonce                    *int
	BuilderFee               *decimal.Decimal
	BuilderID                *int
}

// CreateOrderObject creates a PerpetualOrderModel with the given parameters
func CreateOrderObject(params CreateOrderObjectParams) (*PerpetualOrderModel, error) {
	market := params.Market

	if params.ExpireTime == nil {
		*params.ExpireTime = time.Now().Add(1 * time.Hour)
	}

	// Error if nonce is nil, we keep the input as a pointer so that
	// it is the same as the input to the function
	if params.Nonce == nil {
		return nil, fmt.Errorf("nonce must be provided")
	}

	// If we are buying, then we round up, otherwise we round down
	is_buying_synthetic := params.Side == OrderSideBuy
	collateral_amount := params.SyntheticAmount.Mul(params.Price)

	// For now we only use the default fee type
	// TODO: Allow users to add different fee types
	fees := DefaultFees

	total_fee := fees.TakerFeeRate
	if params.BuilderFee != nil {
		total_fee = total_fee.Add(*params.BuilderFee)
	}

	fee_amount := total_fee.Mul(collateral_amount)

	stark_collateral_amount_dec := collateral_amount.Mul(decimal.NewFromInt(market.L2Config.CollateralResolution))
	stark_synthetic_amount_dec := params.SyntheticAmount.Mul(decimal.NewFromInt(market.L2Config.SyntheticResolution))

	// Round accordingly
	if is_buying_synthetic {
		stark_collateral_amount_dec = stark_collateral_amount_dec.Ceil()
		stark_synthetic_amount_dec = stark_synthetic_amount_dec.Ceil()
	} else {
		stark_collateral_amount_dec = stark_collateral_amount_dec.Floor()
		stark_synthetic_amount_dec = stark_synthetic_amount_dec.Floor()
	}

	stark_collateral_amount := stark_collateral_amount_dec.IntPart()
	stark_synthetic_amount := stark_synthetic_amount_dec.IntPart()
	stark_fee_part := fee_amount.Mul(decimal.NewFromInt(market.L2Config.CollateralResolution)).Ceil().IntPart()

	if is_buying_synthetic {
		stark_collateral_amount = -stark_collateral_amount
	} else {
		stark_synthetic_amount = -stark_synthetic_amount
	}

	order_hash, err := HashOrder(HashOrderParams{
		AmountSynthetic:     stark_synthetic_amount,
		SyntheticAssetID:    market.L2Config.SyntheticID,
		AmountCollateral:    stark_collateral_amount,
		CollateralAssetID:   market.L2Config.CollateralID,
		MaxFee:              stark_fee_part,
		Nonce:               *params.Nonce,
		PositionID:          int(params.Account.vault),
		ExpirationTimestamp: *params.ExpireTime,
		PublicKey:           params.PublicKey,
		StarknetDomain:      params.StarknetDomain,
	})

	if err != nil {
		return nil, fmt.Errorf("hashing order failed: %w", err)
	}

	sig_r, sig_s := params.Signer(order_hash)
	if sig_r == nil || sig_s == nil {
		return nil, fmt.Errorf("signer function returned nil values")
	}

	settlement := Settlement{
		Signature: Signature{
			fmt.Sprintf("%064x", sig_r),
			fmt.Sprintf("%064x", sig_s),
		},
		StarkKey:           params.Account.PublicKey(),
		CollateralPosition: fmt.Sprintf("%d", params.Account.Vault()),
	}

	if params.OrderExternalID == nil {
		defaultID := order_hash
		params.OrderExternalID = &defaultID
	}

	var fee_builder_str *string
	if params.BuilderFee != nil {
		builderFeeStr := params.BuilderFee.String()
		fee_builder_str = &builderFeeStr
	}

	// Convert expire time to epoch milliseconds
	expiryEpochMillis := params.ExpireTime.UnixNano() / int64(time.Millisecond)

	order := &PerpetualOrderModel{
		ID:                       *params.OrderExternalID,
		Market:                   params.Market.Name,
		Type:                     OrderTypeLimit,
		Side:                     params.Side,
		Qty:                      params.SyntheticAmount.String(),
		Price:                    params.Price.String(),
		PostOnly:                 params.PostOnly,
		TimeInForce:              params.TimeInForce,
		ExpiryEpochMillis:        expiryEpochMillis,
		Fee:                      fees.TakerFeeRate.String(),
		SelfTradeProtectionLevel: params.SelfTradeProtectionLevel,
		Nonce:                    fmt.Sprintf("%d", *params.Nonce),
		CancelID:                 params.PreviousOrderExternalID,
		Settlement:               settlement,
		BuilderFee:               fee_builder_str,
		BuilderID:                params.BuilderID,
	}

	return order, nil
}

// HashOrderParams represents the parameters for hashing an order
type HashOrderParams struct {
	AmountSynthetic     int64
	SyntheticAssetID    string // hex string for asset ID
	AmountCollateral    int64
	CollateralAssetID   string // hex string for asset ID
	MaxFee              int64
	Nonce               int
	PositionID          int
	ExpirationTimestamp time.Time
	PublicKey           int
	StarknetDomain      StarknetDomain
}

// HashOrder computes the order hash using the provided parameters
// This mimics the Python hash_order function
func HashOrder(params HashOrderParams) (string, error) {
	// Add 14 days buffer to expiration timestamp
	expireTimeWithBuffer := params.ExpirationTimestamp.Add(14 * 24 * time.Hour)

	expireTimeAsSeconds := expireTimeWithBuffer.Unix()

	// Call the existing GetOrderHash function from sign.go
	hash, err := GetOrderHash(
		fmt.Sprintf("%d", params.PositionID),       // position_id
		params.SyntheticAssetID,                    // base_asset_id_hex
		fmt.Sprintf("%d", params.AmountSynthetic),  // base_amount
		params.CollateralAssetID,                   // quote_asset_id_hex
		fmt.Sprintf("%d", params.AmountCollateral), // quote_amount
		params.CollateralAssetID,                   // fee_asset_id_hex (same as collateral)
		fmt.Sprintf("%d", params.MaxFee),           // fee_amount
		fmt.Sprintf("%d", expireTimeAsSeconds),     // expiration
		fmt.Sprintf("%d", params.Nonce),            // salt (nonce)
		fmt.Sprintf("%d", params.PublicKey),        // user_public_key_hex
		params.StarknetDomain.Name,                 // domain_name
		params.StarknetDomain.Version,              // domain_version
		params.StarknetDomain.ChainID,              // domain_chain_id
		params.StarknetDomain.Revision,             // domain_revision
	)

	if err != nil {
		return "", fmt.Errorf("failed to compute order hash: %w", err)
	}

	return hash, nil
}

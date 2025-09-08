package sdk

import (
	"fmt"
	"math/big"
	"time"
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
	SyntheticAmount          float64
	Price                    float64
	Side                     OrderSide
	CollateralPositionID     int
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
	BuilderFee               *string // Using string for Decimal equivalent
	BuilderID                *int
}

// CreateOrderObject creates a PerpetualOrderModel with the given parameters
func CreateOrderObject(params CreateOrderObjectParams) (*PerpetualOrderModel, error) {
	if params.ExpireTime == nil {
		*params.ExpireTime = time.Now().Add(1 * time.Hour)
	}


	// Error if nonce is nil, we keep the input as a pointer so that
	// it is the same as the input to the function
	if params.Nonce == nil {
		return nil, fmt.Errorf("nonce must be provided")
	}

	// For now we only use the default fee type
	// TODO: Allow users to add different fee types
	// fees := DefaultFees


	order := &PerpetualOrderModel{
		Market:                   params.Market.Name,
		Type:                     OrderTypeLimit,
		Side:                     params.Side,
		Qty:                      fmt.Sprintf("%f", params.SyntheticAmount),
		Price:                    fmt.Sprintf("%f", params.Price),
		TimeInForce:              params.TimeInForce,
		PostOnly:                 params.PostOnly,
		SelfTradeProtectionLevel: params.SelfTradeProtectionLevel,
		BuilderFee:               params.BuilderFee,
		BuilderID:                params.BuilderID,
		CancelID:                 params.PreviousOrderExternalID,
	}

	return order, nil
}

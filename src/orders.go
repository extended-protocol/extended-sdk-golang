package sdk

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

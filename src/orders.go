package sdk

// Signature represents a cryptographic signature
type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
}

// Settlement represents settlement information for an order
type Settlement struct {
	Signature          Signature `json:"signature"`
	StarkKey           string    `json:"starkKey"`
	CollateralPosition string    `json:"collateralPosition"`
}

// ConditionalTrigger represents trigger conditions for conditional orders
type ConditionalTrigger struct {
	TriggerPrice       string `json:"triggerPrice"`
	TriggerPriceType   string `json:"triggerPriceType"`   // e.g., "LAST"
	Direction          string `json:"direction"`          // e.g., "UP", "DOWN"
	ExecutionPriceType string `json:"executionPriceType"` // e.g., "LIMIT", "MARKET"
}

// TpSlTrigger represents take profit or stop loss trigger settings
type TpSlTrigger struct {
	TriggerPrice     string     `json:"triggerPrice"`
	TriggerPriceType string     `json:"triggerPriceType"` // e.g., "LAST"
	Price            string     `json:"price"`
	PriceType        string     `json:"priceType"` // e.g., "LIMIT"
	Settlement       Settlement `json:"settlement"`
}

// PerpetualOrderModel represents a perpetual order
type PerpetualOrderModel struct {
	ID                       string              `json:"id"`
	Market                   string              `json:"market"`
	Type                     string              `json:"type"` // e.g., "CONDITIONAL", "LIMIT", "MARKET"
	Side                     string              `json:"side"` // e.g., "BUY", "SELL"
	Qty                      string              `json:"qty"`
	Price                    string              `json:"price"`
	TimeInForce              string              `json:"timeInForce"` // e.g., "GTT", "GTC", "IOC"
	ExpiryEpochMillis        int64               `json:"expiryEpochMillis"`
	Fee                      string              `json:"fee"`
	Nonce                    string              `json:"nonce"`
	Settlement               Settlement          `json:"settlement"`
	ReduceOnly               bool                `json:"reduceOnly"`
	PostOnly                 bool                `json:"postOnly"`
	SelfTradeProtectionLevel string              `json:"selfTradeProtectionLevel"`    // e.g., "ACCOUNT"
	Trigger                  *ConditionalTrigger `json:"trigger,omitempty"`           // Optional for conditional orders
	TpSlType                 *string             `json:"tpSlType,omitempty"`          // e.g., "ORDER"
	TakeProfit               *TpSlTrigger        `json:"takeProfit,omitempty"`        // Optional take profit settings
	StopLoss                 *TpSlTrigger        `json:"stopLoss,omitempty"`          // Optional stop loss settings
	BuilderFee               *string             `json:"builderFee,omitempty"`        // Optional builder fee
	BuilderID                *int                `json:"builderId,omitempty"`         // Optional builder ID
	CancelID                 *string             `json:"cancel_id,omitempty"`         // Optional field
	DebuggingAmounts         interface{}         `json:"debugging_amounts,omitempty"` // TODO: Define StarkDebuggingOrderAmountsModel
}

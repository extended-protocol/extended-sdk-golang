package sdk

import "math/big"

// PerpetualOrderModel represents a perpetual order
type PerpetualOrderModel struct {
	ID                       string      `json:"id"`
	Market                   string      `json:"market"`
	Type                     interface{} `json:"type"`  // TODO: Define OrderType enum
	Side                     interface{} `json:"side"`  // TODO: Define OrderSide enum
	Qty                      *big.Float  `json:"qty"`   // Using big.Float for Decimal
	Price                    *big.Float  `json:"price"` // Using big.Float for Decimal
	ReduceOnly               bool        `json:"reduce_only"`
	PostOnly                 bool        `json:"post_only"`
	TimeInForce              interface{} `json:"time_in_force"` // TODO: Define TimeInForce enum
	ExpiryEpochMillis        int64       `json:"expiry_epoch_millis"`
	Fee                      *big.Float  `json:"fee"`                         // Using big.Float for Decimal
	Nonce                    *big.Float  `json:"nonce"`                       // Using big.Float for Decimal
	SelfTradeProtectionLevel interface{} `json:"self_trade_protection_level"` // TODO: Define SelfTradeProtectionLevel enum
	CancelID                 *string     `json:"cancel_id,omitempty"`         // Optional field
	Settlement               interface{} `json:"settlement,omitempty"`        // TODO: Define StarkSettlementModel
	Trigger                  interface{} `json:"trigger,omitempty"`           // TODO: Define CreateOrderConditionalTriggerModel
	TpSlType                 interface{} `json:"tp_sl_type,omitempty"`        // TODO: Define OrderTpslType enum
	TakeProfit               interface{} `json:"take_profit,omitempty"`       // TODO: Define CreateOrderTpslTriggerModel
	StopLoss                 interface{} `json:"stop_loss,omitempty"`         // TODO: Define CreateOrderTpslTriggerModel
	DebuggingAmounts         interface{} `json:"debugging_amounts,omitempty"` // TODO: Define StarkDebuggingOrderAmountsModel
	BuilderFee               *big.Float  `json:"builderFee,omitempty"`        // Optional, using big.Float for Decimal
	BuilderID                *int        `json:"builderId,omitempty"`         // Optional field
}

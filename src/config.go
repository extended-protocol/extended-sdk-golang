package sdk

import "github.com/shopspring/decimal"

type StarknetDomain struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	ChainID  string `json:"chain_id"`
	Revision string `json:"revision"`
}

// TradingFeeModel represents trading fees for a market
type TradingFeeModel struct {
	Market         string  `json:"market"`
	MakerFeeRate   decimal.Decimal `json:"maker_fee_rate"`
	TakerFeeRate   decimal.Decimal `json:"taker_fee_rate"`
	BuilderFeeRate decimal.Decimal `json:"builder_fee_rate"`
}

var DefaultFees = TradingFeeModel{
	Market:         "BTC-USD",
	MakerFeeRate:   decimal.NewFromFloat(0.0002), // 2/10000 = 0.0002
	TakerFeeRate:   decimal.NewFromFloat(0.0005), // 5/10000 = 0.0005
	BuilderFeeRate: decimal.NewFromFloat(0),      // 0
}



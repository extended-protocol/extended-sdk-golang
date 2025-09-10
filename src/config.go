package sdk

import "github.com/shopspring/decimal"

type StarknetDomain struct {
	Name     string
	Version  string
	ChainID  string
	Revision string
}

// TradingFeeModel represents trading fees for a market
type TradingFeeModel struct {
	Market         string `json:"market"`
	MakerFeeRate   decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate   decimal.Decimal `json:"takerFeeRate"`
	BuilderFeeRate decimal.Decimal `json:"builderFeeRate"`
}

var DefaultFees = TradingFeeModel{
	Market:         "BTC-USD",
	MakerFeeRate:   decimal.NewFromFloat(0.0002), // 2/10000 = 0.0002
	TakerFeeRate:   decimal.NewFromFloat(0.0005), // 5/10000 = 0.0005
	BuilderFeeRate: decimal.NewFromFloat(0),      // 0
}

package sdk

import "github.com/shopspring/decimal"

type StarknetDomain struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	ChainID  string `json:"chainId"`
	Revision string `json:"revision"`
}

// TradingFeeModel represents trading fees for a market
type TradingFeeModel struct {
	Market         string          `json:"market"`
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

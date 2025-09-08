package sdk

type StarknetDomain struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	ChainID  string `json:"chain_id"`
	Revision string `json:"revision"`
}

// TradingFeeModel represents trading fees for a market
type TradingFeeModel struct {
	Market         string  `json:"market"`
	MakerFeeRate   float64 `json:"maker_fee_rate"`   // Using float64 for decimal rates
	TakerFeeRate   float64 `json:"taker_fee_rate"`   // Using float64 for decimal rates
	BuilderFeeRate float64 `json:"builder_fee_rate"` // Using float64 for decimal rates
}

var DefaultFees = TradingFeeModel{
	Market:         "BTC-USD",
	MakerFeeRate:   0.0002, // 2/10000 = 0.0002
	TakerFeeRate:   0.0005, // 5/10000 = 0.0005
	BuilderFeeRate: 0,      // 0
}

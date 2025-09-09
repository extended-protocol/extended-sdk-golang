package sdk

type L2ConfigModel struct {
	Type                 string `json:"type"`
	CollateralID         string `json:"collateral_id"`
	CollateralResolution int64  `json:"collateral_resolution"`
	SyntheticID          string `json:"synthetic_id"`
	SyntheticResolution  int64  `json:"synthetic_resolution"`
}

type MarketModel struct {
	Name                     string        `json:"name"`
	AssetName                string        `json:"asset_name"`
	AssetPrecision           int           `json:"asset_precision"`
	CollateralAssetName      string        `json:"collateral_asset_name"`
	CollateralAssetPrecision int           `json:"collateral_asset_precision"`
	Active                   bool          `json:"active"`
	L2Config                 L2ConfigModel `json:"l2_config"`
}

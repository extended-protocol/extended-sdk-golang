package sdk

type L2ConfigModel struct {
	Type                 string `json:"type"`
	CollateralID         string `json:"collateralId"`
	CollateralResolution int64  `json:"collateralResolution"`
	SyntheticID          string `json:"syntheticId"`
	SyntheticResolution  int64  `json:"syntheticResolution"`
}

type MarketModel struct {
	Name                     string        `json:"name"`
	AssetName                string        `json:"assetName"`
	AssetPrecision           int           `json:"assetPrecision"`
	CollateralAssetName      string        `json:"collateralAssetName"`
	CollateralAssetPrecision int           `json:"collateralAssetPrecision"`
	Active                   bool          `json:"active"`
	L2Config                 L2ConfigModel `json:"l2Config"`
}

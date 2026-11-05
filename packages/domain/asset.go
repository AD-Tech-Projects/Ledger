package domain

import "time"

type Asset struct {
	ID string

	Name   string
	Type   AssetType
	Status AssetStatus

	Code AssetCode
	// NOTE - number of decimal places for the asset, such as 2 for USD or 8 for BTC
	DecimalPlaces uint8

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AssetType string

const (
	ASSET_FIAT     AssetType = "FIAT"
	ASSET_CRYPTO   AssetType = "CRYPTO"
	ASSET_UNBACKED AssetType = "UNBACKED"
)

type AssetCode string

const (
	// * FIAT
	ASSET_USD AssetCode = "USD"
	ASSET_EUR AssetCode = "EUR"
	ASSET_BRL AssetCode = "BRL"

	// * CRYPTO
	ASSET_BTC  AssetCode = "BTC"
	ASSET_ETH  AssetCode = "ETH"
	ASSET_USDT AssetCode = "USDT"

	// * UNBACKED
	ASSET_POINTS AssetCode = "POINTS"
)

type AssetStatus string

const (
	ASSET_ACTIVE   AssetStatus = "ACTIVE"
	ASSET_INACTIVE AssetStatus = "INACTIVE"
)

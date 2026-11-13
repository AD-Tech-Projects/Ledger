package domain

type Asset struct {
	ID string

	Name   string // Name is the human readable name (e.g. "Brazilian Real", "Bitcoin").
	Type   AssetType
	Status AssetStatus

	Code AssetCode // Code is the unique ticker of the asset (e.g. BRL, BTC).
	// NOTE - number of decimal places for the asset, such as 2 for USD or 8 for BTC
	// Amounts are stored as integers in the smallest unit: BRL 100.00 = 10000; BTC 0.00000001 = 1.
	DecimalPlaces uint8
}

type AssetType string

const (
	ASSET_FIAT     AssetType = "FIAT"     // Government-issued currency (BRL, USD, EUR)
	ASSET_CRYPTO   AssetType = "CRYPTO"   // Crypto asset or stablecoin (BTC, ETH, USDT)
	ASSET_UNBACKED AssetType = "UNBACKED" // Asset without backing, issued by the platform (e.g. points)
)

type AssetCode string

const (
	// * FIAT

	ASSET_USD AssetCode = "USD" // US Dollar, FIAT, 2 decimal places
	ASSET_EUR AssetCode = "EUR" // Euro, FIAT, 2 decimal places
	ASSET_BRL AssetCode = "BRL" // Brazilian Real, FIAT, 2 decimal places

	// * CRYPTO

	ASSET_BTC  AssetCode = "BTC"  // Bitcoin, CRYPTO, 8 decimal places
	ASSET_ETH  AssetCode = "ETH"  // Ethereum, CRYPTO, 18 decimal places
	ASSET_USDT AssetCode = "USDT" // Tether Dollar, CRYPTO, 6 decimal places

	// * UNBACKED

	ASSET_POINTS AssetCode = "POINTS" // ASSET_POINTS is the platform internal currency (UNBACKED, 0 decimal places): 1 point = 0.50 BRL
)

type AssetStatus string

const (
	ASSET_ACTIVE   AssetStatus = "ACTIVE"   // Can be used by new wallets and transactions
	ASSET_INACTIVE AssetStatus = "INACTIVE" // No longer supported for new operations
)

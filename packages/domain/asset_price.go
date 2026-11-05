package domain

import "time"

type AssetPrice struct {
	ID string

	Price      int64
	RecordedAt time.Time

	// NOTE - the base asset for which this price is recorded, such as BTC or ETH
	AssetID string
	Asset   *Asset

	// NOTE - the reference asset for which this price is recorded, such as USD or BRL
	QuoteAssetID string
	QuoteAsset   *Asset
}

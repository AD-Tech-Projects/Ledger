package domain

import "time"

type AssetPrice struct {
	ID string

	Price      int64     // Price is the value of one unit of Asset in the smallest unit of QuoteAsset.
	RecordedAt time.Time // RecordedAt is when the quote was observed.

	// NOTE - the base asset for which this price is recorded, such as BTC or ETH
	AssetID string
	Asset   *Asset

	// NOTE - the reference asset for which this price is recorded, such as USD or BRL
	QuoteAssetID string
	QuoteAsset   *Asset
}

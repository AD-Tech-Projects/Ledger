package services

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

// # Main Flux
func TestValidateAssetPriceCases(t *testing.T) {
	tests := []struct {
		name   string
		change func(*domain.AssetPrice)
	}{
		{"positive price", nil},
		{"smallest price", func(p *domain.AssetPrice) { p.Price = 1 }},
		{"largest price", func(p *domain.AssetPrice) { p.Price = math.MaxInt64 }},
		{"asset quoted in itself", func(p *domain.AssetPrice) { p.QuoteAssetID = p.AssetID }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			price := newTestAssetPrice()
			if test.change != nil {
				test.change(price)
			}
			want := *price

			if err := validateAssetPrice(price); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertAssetPriceUnchanged(t, price, &want)
		})
	}
}

// # Failures
func TestValidateAssetPriceRejectsInvalidPrices(t *testing.T) {
	tests := []struct {
		name   string
		change func(*domain.AssetPrice)
	}{
		{"zero price", func(p *domain.AssetPrice) { p.Price = 0 }},
		{"negative price", func(p *domain.AssetPrice) { p.Price = -1 }},
		{"smallest int64 price", func(p *domain.AssetPrice) { p.Price = math.MinInt64 }},
		{"empty asset id", func(p *domain.AssetPrice) { p.AssetID = "" }},
		{"invalid asset id", func(p *domain.AssetPrice) { p.AssetID = "BTC" }},
		{"empty quote asset id", func(p *domain.AssetPrice) { p.QuoteAssetID = "" }},
		{"invalid quote asset id", func(p *domain.AssetPrice) { p.QuoteAssetID = "USD" }},
		{"without recorded at", func(p *domain.AssetPrice) { p.RecordedAt = time.Time{} }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			price := newTestAssetPrice()
			test.change(price)
			want := *price

			err := validateAssetPrice(price)
			var validationErr *assetPriceValidationError
			if !errors.As(err, &validationErr) || !errors.Is(err, ErrInvalidAssetPrice) {
				t.Fatalf("expected asset price validation error, got %v", err)
			}
			assertAssetPriceUnchanged(t, price, &want)
		})
	}
}

func TestValidateAssetPriceRejectsNilPrice(t *testing.T) {
	if err := validateAssetPrice(nil); !errors.Is(err, ErrInvalidAssetPrice) {
		t.Fatalf("expected invalid asset price error, got %v", err)
	}
}

// helpers
func newTestAssetPrice() *domain.AssetPrice {
	return &domain.AssetPrice{
		ID:           uuid.NewString(),
		Price:        9500000,
		RecordedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AssetID:      uuid.NewString(),
		QuoteAssetID: uuid.NewString(),
	}
}

func assertAssetPriceUnchanged(t *testing.T, price *domain.AssetPrice, want *domain.AssetPrice) {
	t.Helper()
	if price.Price != want.Price || price.AssetID != want.AssetID || price.QuoteAssetID != want.QuoteAssetID ||
		!price.RecordedAt.Equal(want.RecordedAt) {
		t.Fatalf("asset price changed: got %+v, want %+v", price, want)
	}
}

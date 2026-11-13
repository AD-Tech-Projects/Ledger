-- Seed of the supported assets and their reference prices (data only, the schema lives in migrations/).
-- Idempotent: assets are unique by code and prices are guarded by (asset, quote asset, recorded_at).

-- Supported assets (amounts are stored as value x 10^decimal_places in int64)
INSERT INTO "assets" ("id", "created_at", "updated_at", "name", "type", "status", "code", "decimal_places") VALUES
  (gen_random_uuid(), now(), now(), 'Brazilian Real', 'FIAT',     'ACTIVE', 'BRL',    2),
  (gen_random_uuid(), now(), now(), 'US Dollar',      'FIAT',     'ACTIVE', 'USD',    2),
  (gen_random_uuid(), now(), now(), 'Euro',           'FIAT',     'ACTIVE', 'EUR',    2),
  (gen_random_uuid(), now(), now(), 'Bitcoin',        'CRYPTO',   'ACTIVE', 'BTC',    8),
  (gen_random_uuid(), now(), now(), 'Ethereum',       'CRYPTO',   'ACTIVE', 'ETH',    18),
  (gen_random_uuid(), now(), now(), 'Tether',         'CRYPTO',   'ACTIVE', 'USDT',   6),
  (gen_random_uuid(), now(), now(), 'Points',         'UNBACKED', 'ACTIVE', 'POINTS', 0)
ON CONFLICT ("code") DO NOTHING;

-- Reference prices quoted in USD, recorded on 2026-11-12.
-- NOTE: approximate values for development only, NOT real-time market quotes.
-- "price" is the value of ONE unit of the asset in the smallest unit of USD (cents, 2 decimal places):
-- 1 BTC = US$ 95,000.00 -> 9500000. POINTS = R$ 0.50 x BRL price (18 cents) = 9 cents.
INSERT INTO "asset_prices" ("id", "created_at", "updated_at", "price", "recorded_at", "asset_id", "quote_asset_id")
SELECT gen_random_uuid(), now(), now(), p.price, TIMESTAMPTZ '2026-11-12 00:00:00+00', a."id", q."id"
FROM (VALUES
  ('USD',       100), -- US$ 1.00 (self-quote)
  ('EUR',       111), -- US$ 1.11
  ('BRL',        18), -- US$ 0.18
  ('BTC',   7102100), -- US$ 71,021.00
  ('ETH',    200972), -- US$ 2,009.72
  ('USDT',      100), -- US$ 1.00
  ('POINTS',      9)  -- R$ 0.50 x US$ 0.18 = US$ 0.09
) AS p ("code", "price")
JOIN "assets" a ON a."code" = p."code"
JOIN "assets" q ON q."code" = 'USD'
WHERE NOT EXISTS (
  SELECT 1 FROM "asset_prices" ap
  WHERE ap."asset_id" = a."id"
    AND ap."quote_asset_id" = q."id"
    AND ap."recorded_at" = TIMESTAMPTZ '2026-11-12 00:00:00+00'
);

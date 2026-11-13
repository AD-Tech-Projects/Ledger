-- Seed of the initial chart of accounts, see docs/chart-of-accounts.md (data only, the schema lives in migrations/).
-- Idempotent: ledger accounts are unique by code ("BASE-ASSET"). Assets are joined by code, so a missing asset produces no row.
-- Run after the asset seed.

INSERT INTO "ledger_accounts" ("id", "created_at", "updated_at", "type", "status", "normal_balance_direction", "name", "code", "asset_id")
SELECT gen_random_uuid(), now(), now(), b."type", 'ACTIVE', b."direction", b."name" || ' (' || a."code" || ')', b."code" || '-' || a."code", a."id"
FROM (VALUES
  ('1001', 'Reservas da plataforma',             'ASSET',     'DEBIT',  'BRL,USD,EUR,BTC,ETH,USDT'),
  ('2001', 'Saldo disponível dos usuários',      'LIABILITY', 'CREDIT', 'BRL,USD,EUR,BTC,ETH,USDT,POINTS'),
  ('2002', 'Saldo reservado dos usuários',       'LIABILITY', 'CREDIT', 'BRL,USD,EUR,BTC,ETH,USDT,POINTS'),
  ('2003', 'Saldo retido dos usuários',          'LIABILITY', 'CREDIT', 'BRL,USD,EUR,BTC,ETH,USDT,POINTS'),
  ('3101', 'Conversão de ativos',                'EQUITY',    'CREDIT', 'BRL,USD,EUR,BTC,ETH,USDT,POINTS'),
  ('4001', 'Receita com tarifas',                'REVENUE',   'CREDIT', 'BRL,USD,EUR,BTC,ETH,USDT'),
  ('5001', 'Despesas com rendimentos',           'EXPENSE',   'DEBIT',  'BRL,USD,EUR,BTC,ETH,USDT'),
  ('5002', 'Despesas com programa de pontos',    'EXPENSE',   'DEBIT',  'POINTS')
) AS b ("code", "name", "type", "direction", "assets")
JOIN "assets" a ON a."code" = ANY (string_to_array(b."assets", ','))
ON CONFLICT ("code") DO NOTHING;

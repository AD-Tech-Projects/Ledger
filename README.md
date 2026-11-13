# Ledger

A financial system in Go built around a **double-entry bookkeeping ledger**, with multi-asset wallets (fiat currencies, crypto assets and unbacked assets), financial operations (transfers, deposits, withdrawals, yields, fees, reserves and holds) and a focus on consistency, auditability and traceability.

> The project purpose, the domain concepts and the roadmap are described in [`docs/project-overview.md`](./docs/project-overview.md).

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.27.2 |
| HTTP | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL |
| ORM | [GORM](https://gorm.io) |
| Migrations | [Atlas](https://atlasgo.io) + `atlas-provider-gorm` |
| Logging | [`Go-Libs/logger`](https://github.com/AD-Tech-Projects/Go-Libs) |
| Cache | Redis (provisioned in `docker-compose.yml`, not used yet) |

## Prerequisites

- [Go](https://go.dev/dl/) **1.27.2** or later
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Atlas CLI](https://atlasgo.io/getting-started#installation): `curl -sSf https://atlasgo.sh | sh`
- `make`

## Getting started

### 1. Environment variables

Copy the example file and adjust it if needed:

```bash
cp .env.local .env
```

The `Makefile` and `docker-compose.yml` read `.env`. If the file does not exist, the application uses the system environment variables.

### 2. Start the infrastructure

```bash
docker compose up -d database   # PostgreSQL
# docker compose up -d          # PostgreSQL + Redis
```

### 3. Install the dependencies

```bash
make tidy   # go mod tidy
```

The shared library `github.com/AD-Tech-Projects/Go-Libs` is downloaded as a regular Go module. To develop the library locally, clone it into `modules/Go-Libs` and uncomment the `replace` directive in `go.mod`.

### 4. Create the database schema (migrations)

Atlas generates the schema from the GORM models listed in `loader/main.go`, and the versioned migrations in `migrations/` are in sync with those models.

```bash
make migrate-apply   # applies the migrations to the database configured in .env
```

Atlas requires a database dedicated to the Ledger. It refuses to apply migrations to a database that already has tables from another project.

When the models change:

```bash
make migrate-create name=<description>   # atlas migrate diff --env gorm (uses a temporary postgres:15 container; Docker must be running)
make migrate-hash                        # recomputes migrations/atlas.sum after a migration is edited manually
```

### 5. Register the assets

Wallet creation looks up the asset by its code, so the supported assets must exist. The seed in `seeds/` registers `BRL`, `USD`, `EUR`, `BTC`, `ETH`, `USDT` and `POINTS` (internal currency, 1 point = R$ 0.50), plus one reference price per asset quoted in USD:

```bash
make seed   # atlas migrate apply --dir file://seeds (revisions kept in the ledger_seeds schema)
```

The seed is idempotent and must run after `make migrate-apply`. The prices are approximate values for development, not real-time quotes.

### 6. Run the application

```bash
make run            # go run ./main.go
# or
make build && ./out/app
```

The API listens on `http://localhost:${APP_PORT}` (default `3000`):

```bash
curl http://localhost:3000/health   # OK
```

To debug in VS Code, use the **Go: Debug** configuration in `.vscode/launch.json`.

## Tests

```bash
make test            # go test -v ./...
make test-coverage   # generates coverage.out and opens the HTML report
```

The current tests are unit tests and do not use a database. They cover:

- the wallet balance operations, including their errors;
- the effect of each transaction type, including those involving a single wallet or wallets of different types;
- the reversal of each type, which must restore the initial balances;
- the preservation of the balances when the balance is insufficient;
- the validation of the input and of the wallets of each transaction type.

## API

The routes are documented in [`packages/api/doc.go`](./packages/api/doc.go), in godoc format. It includes request examples, the accepted fields, the wallets and balance effect of each transaction type, reversals and error codes.

```bash
go doc ./packages/api
```

## Project structure

```text
.
├── main.go                         # Loads .env, connects to the database, wires the handlers and starts the server
├── loader/main.go                  # Exports the GORM models schema to Atlas
├── migrations/                     # Versioned SQL migrations (Atlas)
├── docs/                           # Project documentation (purpose, domain, roadmap)
├── packages/
│   ├── api/                        # Gin routes, HTTP logging middleware and recovery (documented in api/doc.go)
│   ├── handlers/                   # HTTP input/output, payload validation and error mapping
│   ├── services/                   # Use cases, validations, balance operations and transaction control
│   ├── domain/                     # Entities and enums (documented in domain/doc.go)
│   └── infra/database/
│       ├── connection.go           # GORM/PostgreSQL connection
│       ├── models/                 # Persistence models and domain <-> database mapping
│       └── repositories/           # Generic repository (CRUD, filters, joins, transactions)
├── docker-compose.yml              # PostgreSQL and Redis
├── atlas.hcl                       # Atlas configuration
└── Makefile
```

## Known issues

- **`DB_*` variables exported in the shell take precedence over `.env`**: this applies to `docker compose` and to the application as well (`godotenv` does not override variables that are already set). Only `make migrate-*` reads `.env` directly. Check `echo $DB_NAME` before starting the database or running the app.

# Project Context

Ledger — a financial system in Go built around a double-entry bookkeeping ledger, with multi-asset wallets and financial operations (deposits, withdrawals, transfers, yields, fees, reserves, holds and reversals). This document centralizes the technical information an agent needs to work on the project.

Project-wide references:

- [`docs/project-overview.md`](../docs/project-overview.md) (Portuguese) — purpose, domain concepts, implemented and planned features, and the feature checklist.
- [`README.md`](../README.md) — how to run the project.
- [`packages/domain/doc.go`](../packages/domain/doc.go) — meaning of each domain entity (`go doc ./packages/domain`).
- [`packages/api/doc.go`](../packages/api/doc.go) — HTTP routes, request fields, transaction rules and error codes (`go doc ./packages/api`).

---

## System Overview

### Main Technologies

| Technology            | Role                                                        |
| --------------------- | ----------------------------------------------------------- |
| **Go**                | Language and runtime (version in `go.mod`)                  |
| **Gin**               | HTTP framework                                              |
| **GORM**              | ORM — PostgreSQL driver                                     |
| **PostgreSQL**        | Relational (SQL) database                                   |
| **Atlas**             | Versioned migrations generated from the GORM models         |
| **Go-Libs/logger**    | Shared logging library (`github.com/AD-Tech-Projects/Go-Libs`) |
| **Redis**             | Cache — provisioned in `docker-compose.yml`, not used yet   |

## Execution

Full reference: [`README.md`](../README.md)

### Quick Start

```shell
cp .env.local .env
docker compose up -d database
make tidy
make migrate-apply
make run
```

Register the supported assets and the initial chart of accounts with `make seed` (after `make migrate-apply`).

### Useful Commands

| Command                              | Purpose                                                        |
| ------------------------------------ | -------------------------------------------------------------- |
| `make run`                           | Run the application (`go run ./main.go`)                       |
| `make build`                         | Build the binary into `out/app`                                |
| `make test`                          | Run all tests (`go test -v ./...`)                             |
| `make test-coverage`                 | Generate `coverage.out` and open the HTML report               |
| `make migrate-create name=<name>`    | Generate a migration from the GORM models (`atlas migrate diff`) |
| `make migrate-hash`                  | Recompute `migrations/atlas.sum` after a manual migration edit |
| `make migrate-apply`                 | Apply the migrations to the database configured in `.env`      |
| `make seed`                          | Seed the supported assets, their USD reference prices and the initial chart of accounts (`seeds/`) |
| `gofmt -l .` / `go vet ./...`        | Formatting and static checks                                   |

`DB_*` variables exported in the shell take precedence over `.env` for `docker compose` and for the application (`godotenv` does not override them). Only `make` reads `.env` directly.

---

## Go Architecture

Architecture rules, layer definitions, naming conventions, and code quality hard rules are defined in [`architecture-flow-and-code-style.rule.md`](./rules/architecture-flow-and-code-style.rule.md).

The mandatory execution order is: `api (routes) → handlers → services → repositories`.

### Go Usage Patterns

- Build dependencies explicitly with constructors (`NewXService`, `NewXHandler`) and wire them in `main.go` — there is no DI container.
- Pass `context.Context` as the first parameter of every service and repository method that does IO.
- Return errors instead of panicking. Wrap them with `fmt.Errorf("...: %w", err)` and compare with `errors.Is` / `errors.As`.
- Services expose sentinel errors (`ErrInvalidTransaction`, `ErrWalletNotFound`, `ErrTransactionFailed`); handlers map them to HTTP status codes (400, 404, 422, 500).
- Database transactions are opened in the service with `repository.Transaction(ctx, func(tx *gorm.DB) error { ... })`, and every repository used inside it is rebound with `WithDB(tx)`.
- Rows that must not change concurrently are locked with the `repositories.ForUpdate[M]()` query option, always in the same (sorted) order.
- Monetary amounts are `int64` in the smallest unit of the asset — never floats, and amounts of different assets are never combined.
- Read environment variables only at startup (`main.go`, `database/connection.go`), never inside services.

---

## Infrastructure

| Path                  | Purpose                                                         |
| --------------------- | --------------------------------------------------------------- |
| `docker-compose.yml`  | Local services: `database` (PostgreSQL) and `cache` (Redis)     |
| `atlas.hcl`           | Atlas configuration (`gorm` env, dev database in a temporary container) |
| `loader/main.go`      | Exports the GORM models schema for Atlas — every new model must be registered here |
| `migrations/`         | Versioned SQL migrations and `atlas.sum`                        |
| `.vscode/launch.json` | **Go: Debug** configuration                                     |

### Docker

- Local development uses `docker-compose.yml` at the project root.
- Atlas needs Docker running to generate migrations (it diffs against a temporary Postgres container).
- The database is dedicated to the Ledger: Atlas refuses to apply migrations to a database with tables from other projects.

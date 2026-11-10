---
name: architecture-flow-and-code-style
description: Use when implementing or refactoring backend features in Go. Enforces the architecture flow api -> handlers -> services -> repositories, layer responsibilities, naming conventions, data integrity and clean code practices.
paths:
  - "**/*.go"
---

# Architecture Flow And Code Style

Use when implementing or refactoring backend features in Go. Enforces the architecture flow api -> handlers -> services -> repositories, layer responsibilities, naming conventions, data integrity and clean code practices.

## Mandatory Architecture Flow

- Respect this execution order: `api (routes) -> handlers -> services -> repositories`.
- Do not skip layers and do not call repositories directly from handlers.
- Services may delegate to auxiliary files of the same package (`*_validate.go`, `*_operations.go`) but must not bypass the layer contract.
- `packages/domain` holds only entities, enumerations and their meaning. It must not import any other project package, GORM, Gin or infrastructure code.

### Trigger Layer — `packages/api` and `packages/handlers`

Allowed types:

| Type          | Description                                                         |
| ------------- | ------------------------------------------------------------------- |
| `route`       | Gin route registration and middlewares (`packages/api`)             |
| `httpHandler` | Activated by HTTP requests (`packages/handlers/<entity>_handler.go`) |

Responsibilities:

- Bind and validate the request payload (request body structs with `json` tags live in the handler).
- Convert the payload into the service input (e.g. `services.CreateTransactionInput`) and delegate all business logic to the service.
- Map service errors to HTTP status codes with `errors.Is` (`writeBadRequest`, `writeNotFound`, `writeUnprocessableEntity`, `writeInternalServerError`).
- Avoid business rules and persistence logic in this layer.

### Logic Layer — `packages/services`

Allowed types:

| Type          | Description                                                                 |
| ------------- | --------------------------------------------------------------------------- |
| `service`     | Orchestrates business rules (`<entity>_service.go`)                         |
| `validation`  | Input and entity validation, validation error types (`<entity>_validate.go`) |
| `operations`  | Pure in-memory business operations, such as balance effects (`<entity>_operations.go`) |

Responsibilities:

- Orchestrate rules, open the database transaction and coordinate the repositories.
- Lock the rows that must not change concurrently (`repositories.ForUpdate`) in a deterministic order before changing them.
- Keep validations and operations as pure functions over domain entities, so they can be unit tested without a database.
- Must not build SQL or use GORM directly — persistence goes through the repositories.

### Data Layer — `packages/infra/database`

Allowed types:

| Type         | Description                                                                 |
| ------------ | --------------------------------------------------------------------------- |
| `repository` | Database access (`repositories.GormRepository[M, E]` and its query options) |
| `model`      | GORM persistence model with tags, constraints and `toDomain`/`fromDomain` mappers |
| `connection` | GORM/PostgreSQL connection setup                                            |

Responsibilities:

- Encapsulate all IO and infrastructure concerns and return domain entities to the upper layers.
- Declare schema rules in the model tags (`not null`, `index`, `uniqueIndex` with `priority`, `check`), so Atlas generates them.
- Register new models in `loader/main.go` and generate migrations with `make migrate-create name=<name>`. Never edit an applied migration; create a new one (and run `make migrate-hash` after manual edits).

## Code Quality (Hard Rules)

- Preserve data integrity during manipulation and persistence: balance changes, the transaction record and its status are persisted in the same database transaction.
- Monetary amounts are `int64` in the smallest unit of the asset. Never use `float32`/`float64` for money and never combine amounts of different assets.
- Never mutate inputs or values shared with other callers. When a function must normalize its input, take an explicit pointer and make it obvious from the function name or doc comment.
- Avoid global mutable state. Package-level maps are allowed only as read-only lookup tables (e.g. `transactionWalletRules`) and must never be written at runtime.
- Save complex or compound conditions in a named `bool` variable before branching.
- Return errors instead of panicking; wrap them with `%w` and never discard them silently (`_ =` only for deliberate best-effort rollbacks).
- Do not use `interface{}`/`any` to bypass the type system. If strictly necessary, document the reason and assert to the concrete type immediately.
- Pass `context.Context` as the first parameter of functions that do IO.

## Documentation And Comments (Hard Rules)

- Package and entity documentation lives in a `doc.go` file of the package (Go doc syntax: `# Heading`, tab-indented code blocks, lists and `[Links]`).
- In code, only add comments that explain what constant values and field values mean, plus short `// NOTE - ` comments for non-obvious decisions.
- Exported identifiers that need explanation get a Go doc comment starting with the identifier name.

## Naming Conventions (Hard Rules)

| Context                             | Convention                                | Example                                    |
| ----------------------------------- | ----------------------------------------- | ------------------------------------------ |
| Unexported variables and functions  | `camelCase`                               | `lockWallets`, `isActive`                  |
| Exported functions, types, fields   | `PascalCase`                              | `NewWalletService`, `CreateTransactionInput` |
| Enumeration constants               | `SNAKE_CASE` (in UPPER_CASE), type prefix | `WALLET_ACTIVE`, `TRANSACTION_DEPOSIT`     |
| Sentinel errors                     | `Err` + `PascalCase`                      | `ErrWalletNotFound`                        |
| Acronyms                            | Keep the same case                        | `ID`, `UUID`, `walletID`                   |
| Files                               | `snake_case`, suffixed by role            | `wallet_service.go`, `wallet_model.go`     |
| Packages                            | Short, lowercase, single word             | `services`, `repositories`                 |

## Formatting and Lint

- All generated and modified code must be formatted with **gofmt** (`gofmt -l .` must print nothing) and pass `go vet ./...`.
- Group imports as standard library, third-party modules and project packages, separated by blank lines.
- Do not hand-format code against gofmt, and do not suppress `go vet` findings without a documented reason.

## Implementation Checklist

- [ ] Layering follows `api -> handlers -> services -> repositories`.
- [ ] Each type belongs to its correct layer and `domain` imports no project or infrastructure package.
- [ ] No inputs or shared values are mutated across function boundaries without an explicit pointer contract.
- [ ] No global mutable variables introduced (package-level maps are read-only).
- [ ] Complex conditions are extracted to named booleans.
- [ ] Errors are wrapped with `%w`, mapped by the handler and never swallowed.
- [ ] Money is `int64` in the smallest unit; balance changes and their records are persisted atomically.
- [ ] New models are registered in `loader/main.go` and have a generated migration.
- [ ] Naming conventions (`camelCase`, `PascalCase`, `SNAKE_CASE`, `ErrX`) are respected.
- [ ] Code passes `gofmt -l .` and `go vet ./...`.

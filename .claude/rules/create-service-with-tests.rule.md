---
name: create-service-with-tests
description: Use when implementing a new service, validation or operation in Go and its test file. Enforces the project test structure, table-driven tests, scenario coverage (success, failure, errors) and the standard testing package.
paths:
  - "packages/**/*.go"
---

# Create Service With Tests

## Unit Test File Structure

Create test files next to the code they test, in the same package: `packages/<package>/<file>_test.go` (e.g. `packages/services/transaction_validate_test.go`). Use only the standard `testing` package.

### File Structure Pattern

```go
package services

import (
	"errors"
	"testing"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
)

// # Main Flux
func Test<Function>Cases(t *testing.T) {
	tests := []struct {
		name    string
		input   <InputType>
		want    <OutputType>
	}{
		{"<describe success scenario>", <input>, <want>},
		// ...
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange: build fresh fixtures for this case only
			// act: call the function under test
			// assert: compare the returned value AND the resulting state
		})
	}
}

// # Failures
func Test<Function>Rejects<Condition>(t *testing.T) {
	tests := []struct {
		name  string
		input <InputType>
	}{
		{"<describe invalid scenario>", <input>},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := <function>(test.input)
			if !errors.Is(err, <ErrSentinel>) {
				t.Fatalf("expected <ErrSentinel>, got %v", err)
			}
			// assert: state is unchanged after the failure
		})
	}
}

// helpers
func newTest<Entity>(...) *domain.<Entity> { ... }

func assert<Something>(t *testing.T, ...) {
	t.Helper()
	// ...
}
```

## Mandatory Rules

- Use **table-driven tests** with `t.Run(test.name, ...)`; one row per scenario with a descriptive name.
- Build **fresh fixtures per case** (helper constructors such as `newTestWallet`) — never share mutable state between cases or tests.
- Assertion helpers call `t.Helper()` and report the actual and expected values.
- Assert errors with `errors.Is` / `errors.As` against the sentinel errors (`ErrInvalidTransaction`, `ErrTransactionFailed`...), never by comparing error strings.
- Each test must assert **both** the returned value/error **and** the resulting state (e.g. wallet balances), including that the state is **unchanged** after a failure.
- Test pure functions (`*_validate.go`, `*_operations.go`). Code that depends on the concrete GORM repositories is not unit tested and must not be mocked: keep the business logic out of it.
- Do not add assertion or mock libraries; use the standard `testing` package.
- Group scenarios as **Main Flux** (success paths), **Failures** (rule violations returning sentinel errors) and **Errors** (dependency or IO failures being propagated).
- Test function names follow `Test<Subject><Behaviour>` (e.g. `TestApplyTransactionInsufficientBalanceKeepsWalletsUnchanged`).

## Scenarios to Cover

Derive from the business rules of the function:

- All success paths, including every enumeration value it handles (e.g. every transaction type).
- Each guard/validation that returns a validation or business error.
- Each dependency error that must be propagated (verify it with `errors.Is`).
- Edge cases: `nil` pointers, zero values, empty strings and slices, negative and zero amounts, `int64` limits, invalid UUIDs, inactive entities, mismatched assets.

## Completion Checklist

- [ ] Test file created at `packages/<package>/<file>_test.go`, in the same package as the code.
- [ ] Table-driven tests with `t.Run` and fresh fixtures per case.
- [ ] Every branch of the function covered by at least one case.
- [ ] Assertions cover both the returned value/error and the resulting state.
- [ ] `nil`, zero-value and boundary scenarios included.
- [ ] Errors asserted with `errors.Is` / `errors.As`.
- [ ] `go test ./...` passes and the code is formatted with `gofmt`.

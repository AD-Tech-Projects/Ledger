---
name: implement-feature
description: "Generic feature implementation session. Accepts a feature specification in English or Portuguese, routes work to the appropriate agent(s), and applies the relevant rules and skills."
argument-hint: "Paste the feature specification. Optionally specify: agent (programmer | reviewer | tester | all), rules, and skills to apply."
user-invocable: true
---

<!--
-------------------- How to Use This Skill --------------------
Provide a feature specification and optionally configure which agent(s) and additional rules/skills to activate. The session will route work to the correct agent based on the requested task type.
-->

# Feature Implementation Session

### Feature Specification

> Describe the feature in English or Portuguese. Include: package/entity, business rules, inputs, outputs, dependencies, and acceptance criteria.

### Agent

`<programmer | reviewer | tester | all>`  
**Default**: all (runs programmer → reviewer → tester in sequence)

### Rules (optional)

> List any additional rule files to apply, e.g. architecture-flow-and-code-style, create-service-with-tests

### Skills (optional)

> List any skills to activate, e.g. find-edge-cases, staged-security-review

---

## Agent Routing

### `programmer`

Activates the [Programmer Agent](../../agents/programmer.agent.md).

- Reads [`CONTEXT.md`](../../CONTEXT.md) for project context before implementing.
- Follows [`architecture-flow-and-code-style.rule.md`](../../rules/architecture-flow-and-code-style.rule.md).
- Creates files in order: domain entity/enum → model + mappers → service → handler → route, generating migrations with `make migrate-create`.
- Does NOT create tests.

### `reviewer`

Activates the [Reviewer Agent](../../agents/reviewer.agent.md).

- Applies [`find-edge-cases`](../find-edge-cases/SKILL.md) on the staged diff.
- Applies [`staged-security-review`](../staged-security-review/SKILL.md) on the staged diff.
- Proposes 3 commit message options using the [`commit-message`](../commit-message/SKILL.md) skill rules.
- Requires staged changes (`git add` before running).

### `tester`

Activates the [Tester Agent](../../agents/tester.agent.md).

- Creates unit tests next to the code (`<file>_test.go`), in the same package.
- Follows [`create-service-with-tests.rule.md`](../../rules/create-service-with-tests.rule.md).
- Covers all scenarios including `nil`, zero-value and boundary inputs.
- Does NOT modify source code.

### `all` (default)

Runs all three agents in sequence:

1. **Programmer** — implements the feature.
2. **Reviewer** — reviews the staged changes for edge cases and security.
3. **Tester** — creates tests covering all scenarios.

---

## Examples

### Example — Full Feature Session

```
### Feature Specification
Implement a service to deactivate an owner.
Package: services (owner_service.go), handlers, api
Input: ownerId (string, required, UUID)
Business rules:
  - The owner must exist; otherwise return ErrOwnerNotFound (HTTP 404).
  - If the owner still has active wallets, return a validation error (HTTP 400).
  - Add an OwnerStatus enum (ACTIVE, INACTIVE) to the owner, with its model column and migration.
  - Set the owner status to INACTIVE.
Output: updated domain.Owner

### Agent
all

### Rules
architecture-flow-and-code-style

### Skills
find-edge-cases
```

---

### Example — Programmer Only

```
### Feature Specification
Add a GET /wallet/:id/balance endpoint that returns the free, reserved and hold balances.
Package: services (wallet_service.go), handlers, api
Dependencies: wallet repository

### Agent
programmer
```

---

### Example — Reviewer Only (after staging changes)

```
### Feature Specification
Review the changes I just staged for the transaction service.

### Agent
reviewer

### Skills
staged-security-review
find-edge-cases
```

---

### Example — Tester Only

```
### Feature Specification
Create tests for validateCreateTransactionInput in packages/services.
Success: every transaction type accepts its wallet combination.
Failure: returns ErrInvalidTransaction for missing/forbidden wallets, invalid UUIDs and non-positive amounts.

### Agent
tester
```

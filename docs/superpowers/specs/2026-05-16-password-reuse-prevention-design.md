# Password Reuse Prevention — Design

**Date:** 2026-05-16
**Branch:** `password-age-new`
**Status:** Approved by user, ready for implementation planning
**Compliance driver:** NIST SP 800-63B alignment + common framework requirements (PCI DSS, ISO 27001) that bundle password history with complexity rules.

## Summary

Add an opt-in password-reuse prevention rule to zitadel. Operators configure `history_count: N` on the `PasswordComplexityPolicy` (instance or org scope). When a user changes their own password (self-service or reset-with-code), the new password must not match any of their last `N` previous passwords. Implemented by replaying the user's password-change event history into a `PreviousHashes []string` slice on the human-password write model, then verifying the new plaintext against each stored hash.

## Context: existing work on sibling branch

A `password-age` branch already implements this on `PasswordAgePolicy`. After review we decided to move the field to `PasswordComplexityPolicy` (better semantics: reuse prevention is a "what makes a valid password" rule, like min length). The `password-age` branch's policy-side code (~80%) is therefore reference material, not portable. The reusable concepts are:
- Write-model pattern that accumulates `PreviousHashes` from event replay
- The `checkPasswordHistory` function shape
- UI hint copy and Cypress test structure

## Non-goals

- No reuse check on admin `SetPassword` (override path for lockout recovery; common product convention).
- No reuse check on initial password set (register, invite, init, email-verify) — no prior history exists; check is a no-op.
- No new aggregate, no new projection table. Reuse data is derived from the existing event log.
- No support for legacy `crypto.CryptoValue`-encoded passwords in history comparisons (see Architecture).

## Architecture

### Policy state

`PasswordComplexityPolicy` gains a single field: `history_count uint64`. Lives on the same complexity aggregate at instance and org scope, following the exact pattern of `min_length` / `has_symbol` / etc. Default value `0` means feature disabled. Operators opt in by setting a non-zero value.

### History storage

No new table. The user's password history is reconstructed on demand by replaying the user-aggregate's password events into `HumanPasswordWriteModel.PreviousHashes []string`, ordered most-recent-first. Events consumed:

- `human.password.changed` (`HumanPasswordChangedEvent`)
- `user.human.password.changed` legacy variant (`UserV1PasswordChangedType`)

Reduce semantics: when the model encounters one of these events, if the model's current `EncodedHash` is non-empty, prepend it to `PreviousHashes` before overwriting `EncodedHash` with the new event's hash. Result: `PreviousHashes[0]` is the most-recently-replaced hash, `PreviousHashes[N-1]` is the oldest known.

**Legacy-payload handling.** Older events carry `Secret *crypto.CryptoValue` (encrypted via instance key, not in passwap format). Comparing a new passwap-format plaintext attempt against a decrypted-and-rehashed legacy secret is impractical (would require knowing the original algorithm and re-running it). Decision: **legacy `Secret`-only events are skipped** during reduce — they do not enter `PreviousHashes`. Trade-off: a user whose stored history is all legacy can technically reuse a legacy password on their first post-upgrade change. Acceptable because (a) after one change they have a fresh passwap hash, (b) the feature is opt-in so existing deployments aren't broken on upgrade, and (c) the alternative is a complex multi-algorithm verify path with poor ROI.

### Check site

A single private method on the command layer:

```
func (c *Commands) checkPasswordHistory(
    ctx context.Context,
    newPlaintext string,
    previousHashes []string,
    policy *domain.PasswordComplexityPolicy,
) error
```

Builds an in-memory check list `[currentEncodedHash] ++ previousHashes` (current first), then iterates the first `min(policy.HistoryCount, len(checkList))` entries, calling `c.userPasswordHasher.Verify(hash, newPlaintext)`. On any match → `zerrors.ThrowInvalidArgument(nil, "COMMAND-PwReuse", "Errors.User.Password.Reused")`. On `passwap.ErrPasswordMismatch` continue. On any other passwap error, return wrapped. The current hash is included in the check list because (a) NIST/SOC2 convention treats "last N passwords" as inclusive of the current, and (b) `SetPasswordWithVerifyCode` has no pre-existing same-password rejection (unlike `ChangePassword` which inherits `passwap.ErrPasswordNoChange` from `VerifyAndUpdate`). Empty-string entries (e.g. a user with no prior hash) are skipped.

Signature:

```
func (c *Commands) checkPasswordHistory(
    ctx context.Context,
    newPlaintext string,
    currentEncodedHash string,
    previousHashes []string,
    policy *domain.PasswordComplexityPolicy,
) error
```

Invoked from `setPasswordCommand` only when the caller passes a non-nil `previousHashes` slice and the resolved complexity policy has `HistoryCount > 0`. Both conditions must hold. (The caller also passes `wm.EncodedHash` as the current.)

### Caller wiring

| Caller | Path | Passes `previousHashes`? |
|---|---|---|
| `ChangePassword` | self-service, old password verified | **yes** (`wm.PreviousHashes`) |
| `SetPasswordWithVerifyCode` | reset via emailed code | **yes** (`wm.PreviousHashes`) |
| `SetPassword` | admin set, one-time | no (`nil`) |
| `setPasswordWithVerifyCode` callers in email/init/invite/v2 register flows | initial password set | no (`nil`) |

Non-enforcing callers retain the same write-model load they already do; we just don't pass the slice through.

## Schema changes

All changes are additive — no breaking changes.

### Proto

- `proto/zitadel/policy.proto` — `PasswordComplexityPolicy`: add `uint64 history_count` at next available field number.
- `proto/zitadel/admin.proto` — `AddCustomPasswordComplexityPolicyRequest`, `UpdatePasswordComplexityPolicyRequest`: add `history_count`.
- `proto/zitadel/management.proto` — same for org-scope variants.
- `proto/zitadel/settings/v2/password_settings.proto` + `v2beta` — `PasswordComplexitySettings`: add `history_count`.

### Domain

- `internal/domain/policy_password_complexity.go` — `PasswordComplexityPolicy` struct gains `HistoryCount uint64`.

### Events

- `internal/repository/policy/policy_password_complexity.go` — base events:
  - `PasswordComplexityPolicyAddedEvent` payload gains `HistoryCount uint64`.
  - `PasswordComplexityPolicyChangedEvent` payload gains `*uint64 HistoryCount` (partial-update pointer) plus a `ChangeHistoryCount(uint64)` option func.
- `internal/repository/instance/policy_password_complexity.go` — wrap and forward the new field.
- `internal/repository/org/policy_password_complexity.go` — same.

### Write models

- `internal/command/policy_password_complexity_model.go` — `PasswordComplexityPolicyWriteModel` gains `HistoryCount uint64`; `Reduce()` handles the new field; `NewChangedEvent` includes change-detection.
- `internal/command/instance_policy_password_complexity_model.go`, `org_policy_password_complexity_model.go` — pass through.
- `internal/command/user_human_password_model.go` — `HumanPasswordWriteModel` gains `PreviousHashes []string`; `Reduce()` handles `HumanPasswordChangedEvent` + `UserV1PasswordChangedType` per the rules in *History storage*. `Query()` already subscribes to these events; no change needed.

### Projection

- `internal/query/projection/password_complexity_policy.go` — new column constant `ComplexityPolicyHistoryCountCol = "history_count"`, included in `Init()` table definition with `handler.Default(0)`, and written in `reduceAdded` / `reduceChanged`.

### Query

- `internal/query/password_complexity_policy.go` — `PasswordComplexityPolicy` result struct gains `HistoryCount uint64`; new column var `PasswordComplexityColHistoryCount`; scanned in `preparePasswordComplexityPolicyQuery`.

### Migration

- `cmd/setup/<next-number>.go` + `.sql` — `ALTER TABLE IF EXISTS projections.password_complexity_policies2 ADD COLUMN IF NOT EXISTS history_count INT8 NOT NULL DEFAULT 0;`
- Register in `cmd/setup/config.go` and sequence in `cmd/setup/setup.go`. Implementing sub-agent picks the next available step number at the time of work (main may have advanced past 70).

### gRPC converters

- `internal/api/grpc/policy/password_complexity_policy.go` — `ModelPasswordComplexityPolicyToPb` includes `HistoryCount`.
- `internal/api/grpc/admin/policy_password_complexity_converter.go` + `internal/api/grpc/management/policy_password_complexity_converter.go` — request→domain converters include the field.
- `internal/api/grpc/settings/v2/settings_converter.go` and v2beta equivalent — `PasswordComplexitySettings` mapping includes the field.

## Error handling & UX

### i18n key (canonical)

```yaml
Errors:
  User:
    Password:
      Reused: "Password recently used, choose another"
```

Generic wording. No count interpolation in the error (avoids leaking the policy value through error responses). Added to every locale file in `internal/static/i18n/`. Non-English locales receive the English string as a placeholder for translators.

### Error path

`checkPasswordHistory` returns `zerrors.ThrowInvalidArgument(nil, "COMMAND-PwReuse", "Errors.User.Password.Reused")`. Propagates through gRPC as `INVALID_ARGUMENT` carrying the i18n key — same shape as existing complexity violations.

### Login-app UI hints

Every change-password / set-password component shows the hint when `historyCount > 0`. The hint **can** mention the count (informational, not an error response):

```
You can't reuse your last {count} {count, plural, one {password} other {passwords}}.
```

Components updated:
- `apps/login/src/components/change-password-form.tsx`
- `apps/login/src/components/set-password-form.tsx` (covers reset-with-code page)
- `apps/login/src/components/set-register-password-form.tsx` (covers register/invite)

i18n key `password.complexity.historyHint` added to `apps/login/locales/en.json` plus sibling locale files (English placeholder where untranslated).

`getPasswordComplexitySettings` in `apps/login/src/lib/zitadel.ts` already returns the full proto message; the new field flows through automatically once proto is regenerated.

### Console (Angular)

- `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.{ts,html}` — add `historyCount` form control + numeric input (min 0).
- `console/src/assets/i18n/en.json` — `POLICY.PWD_COMPLEXITY.HISTORYCOUNT: "Password history (generations)"` + helper description.
- Sibling locale JSON files: English placeholder.
- `console/src/app/modules/password-complexity-view/password-complexity-view.component.ts` — if it renders inline hints anywhere a user enters a new password in the admin UI, render the history hint too.

## Testing strategy

### Unit (Go)

- `internal/command/instance_policy_password_complexity_test.go` — extend `AddDefault…` / `ChangeDefault…` cases to include `historyCount` in expected events.
- `internal/command/org_policy_password_complexity_test.go` — same for org scope.
- `internal/command/user_human_password_model_test.go` — new test `TestHumanPasswordWriteModel_PreviousHashesAccumulate`: three sequential password-changed events → assert ordering, prepend semantics, legacy `Secret`-only events skipped.
- `internal/command/user_human_password_test.go` — new cases inside `TestCommandSide_ChangePassword` and `TestCommandSide_SetPasswordWithVerifyCode`:
  - `history_count=0` → reuse permitted (regression guard for default-disabled).
  - `history_count=3`, new password matches `previous[0]` → `INVALID_ARGUMENT / Errors.User.Password.Reused`.
  - `history_count=3`, new password matches `previous[1]` → rejected.
  - `history_count=3`, new password matches `previous[2]` (4th-from-now incl. current) → permitted.
  - `history_count=3`, `SetPasswordWithVerifyCode` with new password equal to **current** stored password → rejected (proves current-hash inclusion in check list).
  - admin `SetPassword` with `history_count=3` and matching new password → permitted (proves admin path bypasses).

### Projection

- `internal/query/projection/password_complexity_policy_test.go` — extend `reduceAdded` / `reduceChanged` expected SQL to include `history_count`.

### Query

- `internal/query/password_complexity_policy_test.go` — extend scan expectations.

### Integration (gRPC)

- `internal/api/grpc/user/v2/integration_test/password_test.go` + `v2beta` — new scenario: set instance complexity policy `history_count=2`, change user password three times, fourth change reusing the first-set password → `INVALID_ARGUMENT`; reusing the now-third-back password → permitted.
- `internal/api/grpc/admin/integration_test/...` (if a complexity-policy admin integration test exists) — new field round-trips.

### E2E (Cypress)

- `tests/functional-ui/cypress/e2e/settings/password-complexity.cy.ts` — extend with: history input visible, save `historyCount=3`, save `historyCount=0`, org-scope custom policy with `historyCount=5`.

### Verification gate

Every implementation session ends by running the relevant test suite for the files touched in that session and reporting pass/fail. Sub-agents do **not** mark work complete without green output (aligns with `verification-before-completion` harness skill).

## Rollout

- Schema migration is additive with `DEFAULT 0` — existing deployments upgrade silently.
- Feature is disabled by default (`history_count = 0`). Operators opt in.
- No data backfill needed. `PreviousHashes` is replayed lazily on the next password change for any user.
- No flag-gated rollout needed — the policy field itself is the toggle.

## Session decomposition

| # | Scope | Driver | Gate |
|---|---|---|---|
| 1 | This design + plan + handoff prompt | Architect (Opus) | User-approved spec committed |
| 2 | Schema, events, projection, migration | Sonnet sub-agent | `go test ./internal/command/... ./internal/query/...` for affected packages green |
| 3 | Proto + gRPC converters | Sonnet sub-agent | Proto regen clean, `go build ./...`, converter tests green |
| 4 | Password change wiring + history check | Sonnet sub-agent | New write-model + ChangePassword/SetPasswordWithVerifyCode tests green |
| 5 | Login-app UI | Sonnet sub-agent | Login app typecheck/build clean |
| 6 | Console UI | Sonnet sub-agent | Console build clean |
| 7 | Integration + E2E tests | Sonnet sub-agent | Integration green, Cypress green |
| 8 | Review pass + ship | Architect | Diff review, `/review` skill, PR opened |

Each session's `next_prompt.md` (at repo root) is overwritten by the previous session before exit. Architect halts cleanly on scope blow-up and writes a recovery prompt.

## Open questions deferred to implementation

- Exact next-available proto field numbers in each `.proto` file (sub-agent verifies at the moment of work).
- Exact next-available `cmd/setup/` step number (sub-agent picks at the moment of work).
- Whether `password-complexity-view.component.ts` renders the hint or not depends on its current usage — sub-agent inspects and decides during Session 6.
- Whether a Cypress reuse-rejection e2e is worth adding on top of the policy-save coverage — sub-agent decides in Session 7 based on existing fixture support; default is skip.

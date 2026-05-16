# Password Reuse Prevention Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add opt-in password-reuse prevention to zitadel by attaching `history_count` to `PasswordComplexityPolicy`, enforced on `ChangePassword` and `SetPasswordWithVerifyCode`, with UI hints everywhere users pick a password.

**Architecture:** New `history_count uint64` on the complexity policy (instance + org scope). History is reconstructed by replaying `HumanPasswordChangedEvent` / `UserV1PasswordChangedType` into a `PreviousHashes []string` slice on the human-password write model — no new table. A new `checkPasswordHistory` command-side method verifies the new plaintext against `[current] ++ PreviousHashes` truncated to `historyCount`, throwing `Errors.User.Password.Reused` on match. Sibling branch `password-age` is reference material (decision was to place this on complexity, not age, so ~80% of its code is not directly portable).

**Tech Stack:** Go 1.23 (event-sourced command/query layer), CockroachDB projections (`cmd/setup` migrations), buf-generated proto (admin/management v1, settings v2/v2beta), Next.js login app (`apps/login`), Angular console (`console/`), Cypress e2e (`tests/functional-ui`).

**Source documents:** Spec at `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` is authoritative. Sibling branch `password-age` contains reference code for the write-model pattern and UI copy — sub-agents may `git show password-age:<path>` to look at it but must adapt to ComplexityPolicy.

---

## File map

**Created:**
- `cmd/setup/<next-number>.go` + `<next-number>.sql` — migration step adding `history_count` column.

**Modified — schema & domain:**
- `internal/domain/policy_password_complexity.go` — `HistoryCount uint64` field.
- `internal/repository/policy/policy_password_complexity.go` — base events gain field.
- `internal/repository/instance/policy_password_complexity.go` — instance wrapper.
- `internal/repository/org/policy_password_complexity.go` — org wrapper.

**Modified — command layer:**
- `internal/command/policy_password_complexity_model.go` — write model + change event helper.
- `internal/command/instance_policy_password_complexity_model.go` — instance write model.
- `internal/command/org_policy_password_complexity_model.go` — org write model.
- `internal/command/instance_policy_password_complexity.go` — `AddDefault…` / `ChangeDefault…` accept new arg.
- `internal/command/org_policy_password_complexity.go` — org variants accept new arg.
- `internal/command/user_human_password_model.go` — `HumanPasswordWriteModel.PreviousHashes`.
- `internal/command/user_human_password.go` — `checkPasswordHistory`, wired into `setPasswordCommand` when `previousHashes != nil`; `ChangePassword` + `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`.
- `internal/command/instance_converter.go` — write model → domain mapping for new field.

**Modified — projection & query:**
- `internal/query/projection/password_complexity_policy.go` — new column.
- `internal/query/password_complexity_policy.go` — scan new column.

**Modified — proto & gRPC:**
- `proto/zitadel/policy.proto` — `PasswordComplexityPolicy` adds `history_count`.
- `proto/zitadel/admin.proto` — request messages add field.
- `proto/zitadel/management.proto` — request messages add field.
- `proto/zitadel/settings/v2/password_settings.proto` + `v2beta` — settings message adds field.
- `internal/api/grpc/policy/password_complexity_policy.go` — model→pb converter.
- `internal/api/grpc/admin/policy_password_complexity_converter.go` — request→domain.
- `internal/api/grpc/management/policy_password_complexity_converter.go` — request→domain.
- `internal/api/grpc/settings/v2/settings_converter.go` + v2beta — settings round-trip.

**Modified — login app UI:**
- `apps/login/src/components/change-password-form.tsx` — hint.
- `apps/login/src/components/set-password-form.tsx` — hint.
- `apps/login/src/components/set-register-password-form.tsx` — hint.
- `apps/login/locales/en.json` + sibling locales — new key `password.complexity.historyHint`.

**Modified — console UI:**
- `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.ts` — form control.
- `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.html` — input.
- `console/src/assets/i18n/en.json` + siblings — new key.
- `console/src/app/modules/password-complexity-view/password-complexity-view.component.{ts,html}` — only if inline-hint render exists.

**Modified — i18n:**
- `internal/static/i18n/en.yaml` + all siblings — new key `Errors.User.Password.Reused`.

**Modified — tests:**
- `internal/command/instance_policy_password_complexity_test.go`
- `internal/command/org_policy_password_complexity_test.go`
- `internal/command/user_human_password_model_test.go`
- `internal/command/user_human_password_test.go`
- `internal/query/projection/password_complexity_policy_test.go`
- `internal/query/password_complexity_policy_test.go`
- `internal/api/grpc/user/v2/integration_test/password_test.go`
- `internal/api/grpc/user/v2beta/integration_test/password_test.go`
- `tests/functional-ui/cypress/e2e/settings/password-complexity.cy.ts`

---

## Execution model

Each task is **one full session** dispatched to a Sonnet sub-agent by the architect (Opus). Sub-agents have full file access; they write the code following the spec. The architect dispatches one task, reviews the diff, writes the next session's `next_prompt.md` at repo root, exits the session.

Sub-agent dispatch template (used by architect):

> Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — it's authoritative for design decisions. Then execute Task N from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md`. Run the gate command at the end and report pass/fail with the actual output. Do not mark complete on red. Commit each logical chunk with a clear message.

**Verification rule:** every session ends green on its gate or the session is reported as failed. No "tests are flaky" handwaves — flake gets root-caused or escalated to architect.

---

## Task 1: Schema, events, projection, migration

**Why this is task 1:** Everything else builds on the policy data carrying the new field. Lands the data plane first; nothing user-visible yet.

**Files:**
- Modify: `internal/domain/policy_password_complexity.go`
- Modify: `internal/repository/policy/policy_password_complexity.go`
- Modify: `internal/repository/instance/policy_password_complexity.go`
- Modify: `internal/repository/org/policy_password_complexity.go`
- Modify: `internal/command/policy_password_complexity_model.go`
- Modify: `internal/command/instance_policy_password_complexity_model.go`
- Modify: `internal/command/org_policy_password_complexity_model.go`
- Modify: `internal/command/instance_policy_password_complexity.go`
- Modify: `internal/command/org_policy_password_complexity.go`
- Modify: `internal/command/instance_converter.go`
- Modify: `internal/query/projection/password_complexity_policy.go`
- Modify: `internal/query/password_complexity_policy.go`
- Create: `cmd/setup/<next-number>.go` + `.sql`
- Modify: `cmd/setup/config.go`, `cmd/setup/setup.go`
- Modify: `internal/command/instance_policy_password_complexity_test.go`
- Modify: `internal/command/org_policy_password_complexity_test.go`
- Modify: `internal/query/projection/password_complexity_policy_test.go`
- Modify: `internal/query/password_complexity_policy_test.go`

**Reference:** `git show password-age:internal/repository/policy/policy_password_age.go` etc. show the analogous shape on the age aggregate — pattern is identical, just on a different policy.

- [ ] **Step 1: Add `HistoryCount uint64` to domain struct.** Following the existing field ordering convention in `policy_password_complexity.go`.
- [ ] **Step 2: Update base event types in `internal/repository/policy/`.** `PasswordComplexityPolicyAddedEvent` payload gains `HistoryCount uint64`; `PasswordComplexityPolicyChangedEvent` payload gains `*uint64 HistoryCount` plus a `ChangeHistoryCount(uint64)` option function. Constructor signature for `NewPasswordComplexityPolicyAddedEvent` gains a trailing `historyCount uint64` arg (or follow the prevailing ordering — match how `password-age` did it on the age policy).
- [ ] **Step 3: Update instance + org wrappers.** They forward args through to base.
- [ ] **Step 4: Update write models.** `PasswordComplexityPolicyWriteModel` gains `HistoryCount uint64`; `Reduce()` handles the new field for both added and changed events; the `NewChangedEvent` helper detects changes and emits `ChangeHistoryCount`.
- [ ] **Step 5: Update command handlers.** `AddDefaultPasswordComplexityPolicy`, `ChangeDefaultPasswordComplexityPolicy`, `AddPasswordComplexityPolicy`, `ChangePasswordComplexityPolicy` accept `historyCount uint64`. `instance_converter.go` write-model→domain helper maps the field.
- [ ] **Step 6: Update projection.** Add column constant `ComplexityPolicyHistoryCountCol = "history_count"`. Include in `Init()` table definition with `handler.Default(0)`. Write column in `reduceAdded` and `reduceChanged`.
- [ ] **Step 7: Update query layer.** `PasswordComplexityPolicy` query result gains `HistoryCount uint64`; new column var `PasswordComplexityColHistoryCount`; scanned in `preparePasswordComplexityPolicyQuery`.
- [ ] **Step 8: Add migration step.** Find the next available step number (`ls cmd/setup/[0-9]*.go | sort -n | tail`). Create `<N>.go` mirroring the structure of an existing recent step (e.g. `cmd/setup/70.go` on `password-age` shows the exact pattern for an `ALTER TABLE` step). Create `<N>.sql` containing `ALTER TABLE IF EXISTS projections.password_complexity_policies2 ADD COLUMN IF NOT EXISTS history_count INT8 NOT NULL DEFAULT 0;`. Register in `cmd/setup/config.go` and sequence in `cmd/setup/setup.go`.
- [ ] **Step 9: Update unit tests.** `instance_policy_password_complexity_test.go` and `org_policy_password_complexity_test.go`: extend `AddDefault…` / `ChangeDefault…` / `Add…` / `Change…` table-driven cases so that expected events include `historyCount`. Add at least one case per command with non-zero `historyCount`.
- [ ] **Step 10: Update projection test.** `password_complexity_policy_test.go`: extend reduceAdded / reduceChanged expected SQL to include `history_count`.
- [ ] **Step 11: Update query test.** Extend scan expectations with `HistoryCount`.

**Commit boundaries:** Steps 1–4 as one commit (data shapes), 5–7 as another (handlers + projection + query), 8 alone (migration), 9–11 as the test commit.

**Gate:**

```bash
go test ./internal/domain/... ./internal/command/... ./internal/query/... ./internal/repository/... ./cmd/setup/... 2>&1 | tee /tmp/task1.log
```

Expected: all green. If `cmd/setup` has no Go tests for the new step (typical), that's fine — the build must compile (`go build ./cmd/setup/...`).

---

## Task 2: Proto + gRPC converters + regen

**Why this is task 2:** With domain state ready, expose it through the wire formats so UI tasks have something to read.

**Files:**
- Modify: `proto/zitadel/policy.proto`
- Modify: `proto/zitadel/admin.proto`
- Modify: `proto/zitadel/management.proto`
- Modify: `proto/zitadel/settings/v2/password_settings.proto`
- Modify: `proto/zitadel/settings/v2beta/password_settings.proto`
- Modify: `internal/api/grpc/policy/password_complexity_policy.go`
- Modify: `internal/api/grpc/admin/policy_password_complexity_converter.go`
- Modify: `internal/api/grpc/management/policy_password_complexity_converter.go`
- Modify: `internal/api/grpc/settings/v2/settings_converter.go`
- Modify: `internal/api/grpc/settings/v2beta/settings_converter.go` (if exists)

**Reference:** `git show password-age:proto/zitadel/policy.proto` and `password-age:proto/zitadel/settings/v2/password_settings.proto` show the field-number convention used on the age policy; mirror on complexity.

- [ ] **Step 1: Inspect current proto state.** `grep -n "PasswordComplexityPolicy" proto/zitadel/policy.proto` and adjacent admin/management/settings files to find current max field number per message.
- [ ] **Step 2: Add `history_count` field to `PasswordComplexityPolicy` in `proto/zitadel/policy.proto`.** Use the next free field number. Add a buf-style comment matching the existing style.
- [ ] **Step 3: Add field to admin + management request messages.** Both `AddCustomPasswordComplexityPolicyRequest` and `UpdatePasswordComplexityPolicyRequest` (or equivalents — verify names).
- [ ] **Step 4: Add field to v2 + v2beta `PasswordComplexitySettings`.**
- [ ] **Step 5: Regenerate proto.** Run the repo's proto-gen command. Typical: `buf generate` or `make generate` — discover the canonical command from `Makefile` or `buf.gen.yaml`.
- [ ] **Step 6: Update model→pb converter** in `internal/api/grpc/policy/password_complexity_policy.go` to include `HistoryCount: policy.HistoryCount`.
- [ ] **Step 7: Update request→domain converters** in admin + management packages.
- [ ] **Step 8: Update v2 + v2beta settings converters** to round-trip the field.
- [ ] **Step 9: Update any existing converter tests** that broke due to new field expectations.

**Commit boundaries:** Steps 1–5 as one commit (proto + regen), 6–9 as another (converters + tests).

**Gate:**

```bash
go build ./... 2>&1 | tee /tmp/task2-build.log
go test ./internal/api/grpc/... 2>&1 | tee /tmp/task2-test.log
```

Expected: clean build, all green tests.

---

## Task 3: Password-change wiring + history check + i18n

**Why this is task 3:** Now the enforcement bite. Smallest possible vertical slice that makes the feature behaviorally real.

**Files:**
- Modify: `internal/command/user_human_password_model.go`
- Modify: `internal/command/user_human_password.go`
- Modify: `internal/static/i18n/en.yaml` and all sibling locale files in `internal/static/i18n/`
- Modify: `internal/command/user_human_password_model_test.go`
- Modify: `internal/command/user_human_password_test.go`

**Reference:** `git show password-age:internal/command/user_human_password_model.go` and `password-age:internal/command/user_human_password.go` show the reduce logic and check-function shape — adapt to read `historyCount` from `*domain.PasswordComplexityPolicy` instead of `*domain.PasswordAgePolicy`.

- [ ] **Step 1: Add `PreviousHashes []string` to `HumanPasswordWriteModel`.**
- [ ] **Step 2: Update `Reduce()`** to handle `HumanPasswordChangedEvent` and `UserV1PasswordChangedType`. Semantics per spec: when these events fire, if `wm.EncodedHash != ""`, prepend it to `PreviousHashes`; then update `EncodedHash` to the event's hash. Skip events that carry only the legacy `Secret *crypto.CryptoValue` field with no `EncodedHash`.
- [ ] **Step 3: Confirm `Query()` already subscribes** to both event types. If not, add the missing subscription.
- [ ] **Step 4: Write the failing test** `TestHumanPasswordWriteModel_PreviousHashesAccumulate` in `user_human_password_model_test.go`: three successive `HumanPasswordChangedEvent`s with hashes `h0`, `h1`, `h2`; assert `EncodedHash == h2` and `PreviousHashes == ["h1", "h0"]`. Add a sub-case with a legacy `Secret`-only event in the middle — confirm it's skipped from `PreviousHashes`.
- [ ] **Step 5: Run the test to verify red, then green.**
- [ ] **Step 6: Add `checkPasswordHistory`** to `user_human_password.go` with the spec signature: `func (c *Commands) checkPasswordHistory(ctx context.Context, newPlaintext, currentEncodedHash string, previousHashes []string, policy *domain.PasswordComplexityPolicy) error`. Build `checkList := append([]string{currentEncodedHash}, previousHashes...)`; iterate first `min(int(policy.HistoryCount), len(checkList))` entries; skip empty strings; call `c.userPasswordHasher.Verify(hash, newPlaintext)`; on success → `zerrors.ThrowInvalidArgument(nil, "COMMAND-PwReuse", "Errors.User.Password.Reused")`; on `passwap.ErrPasswordMismatch` continue; on other errors return wrapped.
- [ ] **Step 7: Wire `checkPasswordHistory` into `setPasswordCommand`.** Add a parameter `previousHashes []string`. When non-nil and the resolved complexity policy has `HistoryCount > 0`, call `checkPasswordHistory` *after* complexity check passes but *before* hashing the new password. All existing callers of `setPasswordCommand` pass `nil` for the new parameter except the two below.
- [ ] **Step 8: Update `ChangePassword`** to pass `wm.PreviousHashes` (load is already happening for old-password verify).
- [ ] **Step 9: Update `SetPasswordWithVerifyCode`** to load `HumanPasswordWriteModel` and pass `wm.PreviousHashes`.
- [ ] **Step 10: Add i18n keys.** In `internal/static/i18n/en.yaml`, under `Errors.User.Password`: add `Reused: "Password recently used, choose another"`. Add the same key (English string as placeholder) to every other locale file in `internal/static/i18n/`. Likewise add `COMMAND-PwReuse` description if the repo has a separate error-codes file (search `grep -r "COMMAND-PwReuse" internal/static/i18n/ proto/`).
- [ ] **Step 11: Add command tests** in `user_human_password_test.go`. Five new sub-cases per the spec's testing section:
  - `history_count=0` → reuse permitted
  - `history_count=3`, new matches `previous[0]` → rejected
  - `history_count=3`, new matches `previous[1]` → rejected
  - `history_count=3`, new matches `previous[2]` → permitted
  - `SetPasswordWithVerifyCode` with `history_count=3` and new == current → rejected
  - Admin `SetPassword` with `history_count=3` and matching new password → permitted
- [ ] **Step 12: Run and verify green.**

**Commit boundaries:** Steps 1–5 (write-model + its test), 6–9 (check function + wiring), 10 (i18n), 11–12 (command tests).

**Gate:**

```bash
go test ./internal/command/... 2>&1 | tee /tmp/task3.log
```

Expected: all green. If existing tests broke from the new parameter on `setPasswordCommand`, fix them in-place (pass `nil`).

---

## Task 4: Login-app UI

**Why this is task 4:** With server-side complete, surface the hint to end users.

**Files:**
- Modify: `apps/login/src/components/change-password-form.tsx`
- Modify: `apps/login/src/components/set-password-form.tsx`
- Modify: `apps/login/src/components/set-register-password-form.tsx`
- Modify: `apps/login/locales/en.json` + all sibling locale files in `apps/login/locales/`

**Reference:** `git show password-age:apps/login/src/components/change-password-form.tsx` shows the `password-age` branch's hint impl on the change-password page only. Same pattern, broader application here.

- [ ] **Step 1: Confirm `getPasswordComplexitySettings`** in `apps/login/src/lib/zitadel.ts` returns `historyCount`. After Task 2 regenerated proto, the field should be auto-mapped. If it isn't (e.g. an explicit select-only mapping somewhere), add it.
- [ ] **Step 2: Update `change-password-form.tsx`.** Add `historyCount?: number` prop; render an info alert above the submit button when `historyCount && historyCount > 0` using the new translation key. Match the existing complexity-hint UI style (look at `<PasswordComplexity>` component).
- [ ] **Step 3: Update `set-password-form.tsx` and `set-register-password-form.tsx`** with the same hint logic. For these forms, decide where to source `historyCount`: if they don't already fetch complexity settings, add the fetch following the existing pattern in their parent server components (mirror `change/page.tsx`).
- [ ] **Step 4: Add translation key.** In `apps/login/locales/en.json`, under `password.complexity`: `"historyHint": "You can't reuse your last {count} {count, plural, one {password} other {passwords}}."`. Add same key (English placeholder) to every sibling locale.
- [ ] **Step 5: Wire the prop in parent server components.** `apps/login/src/app/(login)/password/change/page.tsx` already does this on `password-age`; mirror for the equivalent set/register pages.
- [ ] **Step 6: Typecheck + build.**

**Commit boundaries:** Step 1–3 as one commit (components), 4–5 as another (i18n + wiring).

**Gate:**

```bash
pnpm -F login typecheck 2>&1 | tee /tmp/task4-tc.log
pnpm -F login build 2>&1 | tee /tmp/task4-build.log
```

If the repo uses a different command for the login app, discover via `apps/login/package.json` scripts.

---

## Task 5: Console UI

**Why this is task 5:** Admin can now configure the value.

**Files:**
- Modify: `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.ts`
- Modify: `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.html`
- Modify: `console/src/assets/i18n/en.json` + sibling locales
- Optional: `console/src/app/modules/password-complexity-view/password-complexity-view.component.{ts,html}` if it currently renders inline hints in any admin UI password-input context

**Reference:** `git show password-age:console/src/app/modules/policies/password-age-policy/password-age-policy.component.ts` shows the form-control + numeric input pattern; same shape applies to complexity.

- [ ] **Step 1: Add `historyCount` form control** to the complexity-policy form group with `[null, []]` initialization. Getter mirroring siblings (e.g. `minLength`, `hasSymbol`).
- [ ] **Step 2: Update submit handlers** to include `this.historyCount?.value ?? 0` in the gRPC request body for both add and update paths.
- [ ] **Step 3: Add `<input>` element** to `.component.html` of type `number` with `min="0"`, `formControlName="historyCount"`, helper text from i18n. Match the existing layout (typically `<cnsl-form-field>` wrappers).
- [ ] **Step 4: Add i18n key** `POLICY.PWD_COMPLEXITY.HISTORYCOUNT: "Password history (generations)"` to `console/src/assets/i18n/en.json`. Add description if siblings have one (e.g. `HISTORYCOUNT_DESC`). Add English placeholder to every sibling locale JSON.
- [ ] **Step 5: Decide on `password-complexity-view.component`.** Inspect whether it renders the complexity hints in any admin password-change context. If yes, add the history hint; if no, leave alone.
- [ ] **Step 6: Typecheck + build.**

**Commit boundaries:** Steps 1–3 (component), 4 (i18n), 5 (optional view) as a separate commit if touched.

**Gate:**

```bash
pnpm -F console typecheck 2>&1 | tee /tmp/task5-tc.log
pnpm -F console build 2>&1 | tee /tmp/task5-build.log
```

Discover canonical commands via `console/package.json` if different.

---

## Task 6: Integration + E2E tests

**Why this is task 6:** Prove the feature works end-to-end before claiming completion.

**Files:**
- Modify: `internal/api/grpc/user/v2/integration_test/password_test.go`
- Modify: `internal/api/grpc/user/v2beta/integration_test/password_test.go`
- Modify: `tests/functional-ui/cypress/e2e/settings/password-complexity.cy.ts`
- Optionally: integration test for admin/management complexity-policy RPCs if one exists with similar shape — verify field round-trip.

- [ ] **Step 1: gRPC integration test (v2 + v2beta).** New scenario: instance complexity policy set to `history_count=2`. Change a fixture user's password three times via the v2 RPC, recording the plaintexts. Fourth attempt: change to the original plaintext (now 3-back, outside window) → expect success. Fifth attempt on a fresh fixture user: change twice, then attempt to change to the first password → expect `INVALID_ARGUMENT` with `Errors.User.Password.Reused`. Use existing test scaffolding (`Tester` / `Instance` helpers — discover via existing tests in the file).
- [ ] **Step 2: gRPC integration test — current-hash inclusion.** Specifically: set `history_count=1`, password reset via verify code, new password equal to current → expect `INVALID_ARGUMENT`.
- [ ] **Step 3: Cypress e2e.** Extend `password-complexity.cy.ts`: history-count input is visible on instance complexity policy page; save with `history_count=3`; save with `history_count=0`; org-scope custom complexity policy with `history_count=5`. Mirror the structure of existing complexity tests.
- [ ] **Step 4: Run both suites.**

**Commit boundaries:** One commit per test file.

**Gate:**

```bash
# Integration (requires test DB up — see existing test docs)
go test ./internal/api/grpc/user/v2/integration_test/... ./internal/api/grpc/user/v2beta/integration_test/... -run Password 2>&1 | tee /tmp/task6-int.log

# Cypress (in tests/functional-ui)
cd tests/functional-ui && pnpm cypress run --spec "cypress/e2e/settings/password-complexity.cy.ts" 2>&1 | tee /tmp/task6-cy.log
```

If integration tests require infra (Docker, test DB), discover the existing make target — typically `make test-integration` or similar.

---

## Task 7: Review pass + PR

**Why this is task 7:** Architect-driven final pass, not sub-agent territory. Listed so the plan is closed out cleanly.

- [ ] **Step 1:** Architect (Opus) reads the full diff against `main`.
- [ ] **Step 2:** Dispatch `/review` skill (or the project's `superpowers:review` if it exists) for a structural pass.
- [ ] **Step 3:** Fix anything flagged.
- [ ] **Step 4:** Verify all six gate commands above still pass against the merged-up branch state.
- [ ] **Step 5:** Open PR with summary referencing the spec doc.

**Gate:** PR opened, CI green.

---

## Self-review notes

**Spec coverage check:** Every section of the spec maps to a task:
- Architecture / data flow → Task 1 (schema), Task 3 (check site + wiring).
- Schema changes → Task 1 (Go side) + Task 2 (proto/gRPC).
- Error handling & UX → Task 3 (i18n + check), Task 4 (login UI), Task 5 (console UI).
- Testing strategy → unit/projection/query tests in Tasks 1 + 3; integration + e2e in Task 6.
- Rollout → Task 1 (additive migration, default 0). No flag-gated rollout needed per spec.

**Placeholder scan:** Two intentional `<next-number>` placeholders for the migration step number and proto field numbers — explicitly delegated to sub-agent inspection at the time of work (spec also acknowledges these as deferred). No other TBDs.

**Type/name consistency:** `historyCount` (Go field `HistoryCount`, proto `history_count`, TS prop `historyCount`, Angular form control `historyCount`, SQL column `history_count`, i18n key `historyHint`). `Errors.User.Password.Reused` i18n key matches in Go i18n yaml and stays consistent across tests. `COMMAND-PwReuse` error code is the single source.

**Bite-size check:** Each task is one session per the architect's `next_prompt.md` handoff model; internal steps are 2-5 minute units suitable for a Sonnet sub-agent following the TDD pattern.

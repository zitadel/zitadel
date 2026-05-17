# Session 7 handoff — password reuse prevention

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're mid-feature: password reuse prevention for NIST alignment, multi-session, architect-driven (you are the architect, you do not write code — you dispatch Sonnet sub-agents and review). Branch: `password-age-new`.

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan with file maps and gate commands. Each task is one session.

**Sessions 1–6 (already done):**
- Session 1: spec + plan + Session 2 handoff. Commits `8a1c4fc1d`, `cb62b9631`, `85114cc51`.
- Session 2: **Task 1 — Schema, events, projection, migration.** Commits `6283c6ecd`, `bdef6e80b`, `8d09be34c`, `82254b056`. `HistoryCount uint64` is live across domain, repository events, write models, projection, query, and migration (`cmd/setup/70`).
- Session 3: **Task 2 — Proto + gRPC converters + regen.** Commits `b6debd810`, `b12aeb7b5`. `history_count` is wire-visible end to end across `policy.proto`, `admin.proto`, `management.proto`, and settings v2 + v2beta. Converters in `internal/api/grpc/policy/`, `admin/`, `management/`, `settings/v2/`, `settings/v2beta/` all map the field. Generated `.pb.go` and `.ts` files are gitignored — a fresh clone needs `PATH=.artifacts/bin/linux/amd64:$PATH buf generate` at the repo root, plus `buf generate` inside `packages/zitadel-proto/` for the TS client.
- Session 4: **Task 3 — Password change wiring + history check + i18n.** Commits `09d5231a6`, `c9d192c01`, `730d35bee`, `0797dd49b`. `HumanPasswordWriteModel` accumulates `PreviousHashes []string` via `Reduce()`. `checkPasswordHistory` is wired into `setPasswordCommand`; only `ChangePassword` and `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`. `Errors.User.Password.Reused` is in all 22 Go locale yaml files. Gate green.
- Session 5: **Task 4 — Login-app UI.** Commits `6cf83ed53`, `326105b39`. `historyCount?: number` prop on three form components; hint renders `<Alert type={AlertType.INFO}>` only when `historyCount > 0` using i18n key `password.complexity.historyHint` (ICU plural). All 14 `apps/login/locales/*.json` files have the new key. Three parent server components pass `historyCount={Number(passwordComplexity.historyCount)}`. Gate green.
- Session 6: **Task 5 — Console UI.** Commits `a9c3a8fb8`, `17450b72c`. Console password-complexity policy form takes a `historyCount` numeric input; submit handlers for admin (default) + mgmt (custom add/update) all wire it. `POLICY.PWD_COMPLEXITY.HISTORYCOUNT: "Password history (generations)"` is in all 22 `console/src/assets/i18n/*.json` files. Both gates green (`pnpm exec tsc --noEmit` silent, `pnpm build` "Application bundle generation complete" in ~78s, only pre-existing CommonJS bailout warnings).

  Important Session 6 nuances for the next session:
  - **The plan said the complexity component uses Angular reactive forms; it does not.** It uses object-mutation against `complexityData?: PasswordComplexityPolicy.AsObject` with `[(ngModel)]` two-way binding (look at `incrementLength()`/`decrementLength()` for the prior pattern). The sub-agent adapted correctly — exposed `historyCount` as a TS getter/setter on the component that reads/writes `complexityData.historyCount`. The sibling `password-age` policy form uses reactive forms; complexity does NOT.
  - The numeric input uses a raw `<input type="number" min="0">` wrapped in `.length-wrapper`, not the `+/-` icon-button pattern that `minLength` uses. Defensible (history counts can be 24+ where button-clicking is tedious) and the plan was OK with "wrap in a `<cnsl-form-field>` or whatever wrapper the existing inputs use."
  - The sub-agent also touched `console/src/app/services/admin.service.ts` and `mgmt.service.ts` to add a `historyCount` parameter on three RPC wrapper methods. This wasn't in the plan's listed files but was structurally required to pass the form value to the gRPC call. Correct.
  - The view component (`console/src/app/modules/password-complexity-view/`) was deliberately left alone — it renders live client-side validation hints as a user types (length progress, symbol/number/case checks). History reuse can't be validated client-side (requires hash comparison server-side), so no hint is appropriate there.
  - No `HISTORYCOUNT_DESC` helper key was added; siblings in `POLICY.PWD_COMPLEXITY` don't have `_DESC` keys.
  - The sub-agent **manually patched gitignored proto stubs** (`policy_pb.js/.d.ts`, `admin_pb.js/.d.ts`, `management_pb.js/.d.ts`) because it tried `pnpm generate` without `.artifacts/bin/linux/amd64` on PATH and `protoc` was missing. Those manual edits are obsolete: a real `cd console && PATH=/home/zacharya/projects/zitadel/.artifacts/bin/linux/amd64:$PATH pnpm generate` works (just produces some harmless "path … is not equal to" warnings) and writes correct generated files containing `historyCount` from the committed `.proto` source. The architect verified this and re-ran both gates against a freshly-regenerated baseline — still green.

**Your job in Session 7: Execute Task 6 from the plan — "Integration + E2E tests."**

This is the proof-the-feature-actually-works pass. Two gRPC integration tests + one Cypress extension. After this lands, Task 7 (architect review + PR) is the final wrap.

**Dispatch instructions:**

1. Confirm branch is `password-age-new` and clean (`git status`). Confirm Session 6's two commits are on top of Session 5's two (`git log --oneline -10`).
2. Read the spec doc fully (re-read the "Testing strategy → Integration (gRPC)" and "Cypress" subsections), then re-read Task 6 in the plan.
3. **Decide infra strategy before dispatching.** Integration tests in this repo need a running test DB (Postgres) + the auth fixtures. Discover what's expected:
   - `grep -rn "TestMain\|integrationTest\|require.*Integration" internal/api/grpc/user/v2/integration_test/ | head -20` — see how existing tests bootstrap.
   - `cat Makefile | grep -A3 "test-integration\|integration"` — find the make target. Common targets: `make test-integration` or a `docker compose` invocation.
   - If Docker / a long-running Postgres is required and not already up, the sub-agent should NOT spin up infra blindly. Tell it to write the tests and run only what's runnable locally; the architect will run the full suite manually on the architect's machine.
4. Dispatch ONE Sonnet sub-agent (Agent tool, `subagent_type: general-purpose`, `model: sonnet`) with this prompt:

   > Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — authoritative. Then execute Task 6 from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` end-to-end: 4 steps, commit at the boundaries the task specifies (one commit per test file). Run the gate commands at the end and report pass/fail with the actual output.
   >
   > Tasks 1–5 already landed `HistoryCount` end-to-end through Go (domain, events, projection, query, command-layer enforcement, command-layer unit tests), proto/gRPC, the Next.js login app, and the Angular console. You're proving the feature works via integration + e2e tests.
   >
   > Specific guards:
   > - Files to create/extend:
   >   - `internal/api/grpc/user/v2/integration_test/password_test.go` — extend existing tests, do not create a new file. Inspect surrounding tests in the same file to see how complexity-policy + change-password are exercised. Use existing `Tester` / `Instance` / `IAMOwnerCtx` helpers — discover via `grep -n "func Test\|Instance.\|UserV2\|Tester" internal/api/grpc/user/v2/integration_test/password_test.go`.
   >   - `internal/api/grpc/user/v2beta/integration_test/password_test.go` — mirror the v2 changes. v2beta is typically a thinner copy.
   >   - `tests/functional-ui/cypress/e2e/settings/password-complexity.cy.ts` — extend existing test, do not create a new file. Inspect the file first to see the structure for the policy-edit flow.
   > - Scenarios required (from the plan + spec):
   >   - **gRPC scenario A**: instance policy `history_count=2`, fixture user changes password three times (track plaintexts), fourth change uses the original plaintext (now 3-back, outside window) → expect SUCCESS. Fresh fixture user with `history_count=2`, change twice, attempt to change to the first password → expect `INVALID_ARGUMENT` containing `Errors.User.Password.Reused` (or the code `COMMAND-PwReuse` — match how other invalid-argument assertions in this file check the error).
   >   - **gRPC scenario B (current-hash inclusion)**: `history_count=1`, reset via verify code with new password equal to **current** stored password → expect `INVALID_ARGUMENT`. This proves the current hash is in the comparison list, not just previous events.
   >   - **Admin (instance) gRPC round-trip** (optional per plan, do it if a similar test exists): set `history_count=N` via admin RPC, read it back via management/admin getter, assert the value round-trips. Look for an existing complexity-policy admin integration test before writing a new one.
   >   - **Cypress**: visit the instance complexity policy page, assert the history-count input is visible, save with `history_count=3`, reload, assert it persists. Save with `history_count=0`, reload, assert it persists. If there's an org-scope custom complexity policy page in existing tests, exercise `history_count=5` there too. Mirror the structure of existing complexity tests (don't invent a new pattern).
   > - The console field on the form is bound via `[(ngModel)]="historyCount"`; the input element has no special selector beyond `class="history-count-input"`. Cypress should target it via that class or via the parent row's i18n label.
   > - Do not exceed Task 6's scope: no production-code edits, no review commits (Task 7), no doc updates.
   >
   > **Commit boundaries (per plan):** one commit per test file. Three commits expected (v2 + v2beta + cypress); a fourth if you add an admin round-trip test.
   >
   > **Gate (must run, report output):**
   > ```
   > # Integration (Go) — discover the right invocation; common forms:
   > go test ./internal/api/grpc/user/v2/integration_test/... -run Password -v 2>&1 | tail -80
   > go test ./internal/api/grpc/user/v2beta/integration_test/... -run Password -v 2>&1 | tail -80
   > # If infra is required (Docker/Postgres) and not up, REPORT THAT — do not start infra. The architect will run it.
   >
   > # Cypress — discover canonical invocation in tests/functional-ui/package.json
   > cd tests/functional-ui && pnpm exec cypress run --spec "cypress/e2e/settings/password-complexity.cy.ts" 2>&1 | tail -60
   > # If Cypress requires a running app (the typical case), REPORT THAT — do not start the app. The architect will run e2e.
   > ```
   > If the suites can't run locally because they need infra, that is OK — the deliverable is correct, complete, compileable test code. Report:
   > - For Go: `go vet ./internal/api/grpc/user/v2/integration_test/... ./internal/api/grpc/user/v2beta/integration_test/...` must be clean. Report the output.
   > - For Cypress: `cd tests/functional-ui && pnpm exec tsc --noEmit` (or whatever the canonical typecheck for cypress files is) — clean for new code. Report the output.
   >
   > **Final report:**
   > - Commit SHAs and one-line subjects.
   > - For each new test scenario, one paragraph: which scenario, what asserts what, which helpers used.
   > - Output of `go vet` (or `go test` if runnable).
   > - Output of cypress typecheck (or `cypress run` if runnable).
   > - Any judgement calls (e.g. whether an admin round-trip test was added, whether existing fixture helpers needed extending).
   > - If anything is red — STOP. Report status. Don't paper over a real failure.

5. When the sub-agent returns, **verify the diff yourself** before trusting the "done" claim:
   - `git log --oneline -10` — should see 3–4 new commits from this task on top of Session 6's two.
   - `git diff HEAD~3 --stat` (or `~4`) — expect: 1 `internal/.../user/v2/integration_test/password_test.go`, 1 `internal/.../user/v2beta/integration_test/password_test.go`, 1 `tests/functional-ui/cypress/e2e/settings/password-complexity.cy.ts`, possibly an admin-integration test.
   - Spot-check each test file for: (a) `history_count=2` scenario with the 3-back permitted assertion, (b) `Errors.User.Password.Reused` or `COMMAND-PwReuse` assertion on reuse attempts, (c) current-hash-inclusion test, (d) cypress assertions that persist + reload.
   - `go vet ./internal/api/grpc/user/v2/integration_test/... ./internal/api/grpc/user/v2beta/integration_test/...` — must be clean.
   - If the architect has Docker + the test DB available, run `go test` for the new `-run` names. If not, document that as "ready for CI" in the Task 7 handoff.

6. **If green and diff looks sane:** overwrite `next_prompt.md` with a Session 8 handoff for Task 7 (architect review + PR — diff review, spec/plan reconciliation, PR creation). Commit the new handoff. Stop.

7. **If red or diff looks wrong:** do NOT dispatch a second sub-agent reflexively. Read the failure, decide if it's a sub-agent mistake (re-dispatch with a corrective prompt) or a plan/spec gap (update the spec or plan first, then re-dispatch). Common failure modes to watch for: (a) integration test uses a stale fixture helper that doesn't set instance policy correctly — verify `SetDefaultPasswordComplexityPolicy` or equivalent is being called with `historyCount` set; (b) test asserts on `INVALID_ARGUMENT` gRPC code but the actual returned error wraps a Zitadel-internal error code — check existing reuse-rejection tests in the command layer (`internal/command/user_human_password_test.go`) to see how they assert; (c) Cypress targets a `formControlName` selector — there is no form control, it's `[(ngModel)]` with `class="history-count-input"`; (d) v2beta tests fail because v2beta uses a different test scaffold (occasionally `IAMOwnerCtx` vs a different ctx).

**Hard rules:**

- You don't write code. Only design docs / plan files / next_prompt.md.
- All implementation goes through Sonnet sub-agents.
- One sub-agent at a time per session — they need a clean handoff, not parallel races on the same files.
- Verify before trusting. Sub-agents will claim success. Read the diff and re-run the gate yourself (or document why you can't, e.g. infra requirements).
- If scope grows past Task 6, stop and write a recovery `next_prompt.md`. Don't try to do Tasks 6 + 7 in one session.

**User preferences (from memory):**

- Sonnet for structured/mechanical work like extending existing test files with new scenarios. Opus only for hard reasoning. ✓
- Skill priority: invoke any matching skill before action. `superpowers:subagent-driven-development` may apply — check, but the user has consistently preferred the lighter "architect-verification" path (skip spec-reviewer/code-quality-reviewer sub-agents, do the verification yourself by reading the diff and re-running the gate). Default to that unless the user opts up.

**Useful context from Sessions 1–6 (so you don't re-derive):**

- Field naming is consistent: Go `HistoryCount uint64`, SQL column `history_count`, proto `uint64 history_count`. Login-app i18n key is `password.complexity.historyHint`; Console i18n key is `POLICY.PWD_COMPLEXITY.HISTORYCOUNT`. Backend error key is `Errors.User.Password.Reused`, error code `COMMAND-PwReuse`.
- Console has 22 locale files; login-app has 14 (different sets).
- Command-layer unit tests for reuse logic already exist (Session 4 commit `0797dd49b`) in `internal/command/user_human_password_test.go` — read them first; they're the closest existing model for how the gRPC integration tests should assert.
- The Console form is **not reactive forms** for complexity — `[(ngModel)]="historyCount"` against a getter/setter that proxies `complexityData.historyCount`. Cypress targets via `class="history-count-input"` or the row's i18n label.
- After Task 6 lands, Task 7 is the architect's review pass + PR creation. Task 7 does NOT delegate to a sub-agent — it's the closing pass that lives entirely in the architect's hands.

Good luck. Bite-sized scope. Be the architect.

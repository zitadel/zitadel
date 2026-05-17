# Session 8 handoff — password reuse prevention (final pass)

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're at the end of a multi-session feature: password reuse prevention for NIST alignment. Branch: `password-age-new`. Tasks 1–6 are done. Your job, Session 8, is **Task 7 — architect review pass + PR creation**. This is the closing pass; it lives entirely in your hands (no sub-agent).

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan. Re-read Task 7 specifically.

**Sessions 1–7 (already done):**
- Session 1: spec + plan + Session 2 handoff. Commits `8a1c4fc1d`, `cb62b9631`, `85114cc51`.
- Session 2: **Task 1 — Schema, events, projection, migration.** Commits `6283c6ecd`, `bdef6e80b`, `8d09be34c`, `82254b056`. `HistoryCount uint64` live across domain, repository events, write models, projection, query, and migration (`cmd/setup/70`).
- Session 3: **Task 2 — Proto + gRPC converters + regen.** Commits `b6debd810`, `b12aeb7b5`. `history_count` is wire-visible across `policy.proto`, `admin.proto`, `management.proto`, and settings v2 + v2beta. Note: generated `.pb.go` and `.ts` files are gitignored — a fresh clone needs `PATH=.artifacts/bin/linux/amd64:$PATH buf generate` at the repo root, plus `buf generate` inside `packages/zitadel-proto/` for the TS client. Console needs `cd console && PATH=/home/zacharya/projects/zitadel/.artifacts/bin/linux/amd64:$PATH pnpm generate`.
- Session 4: **Task 3 — Password change wiring + history check + i18n.** Commits `09d5231a6`, `c9d192c01`, `730d35bee`, `0797dd49b`. `HumanPasswordWriteModel.PreviousHashes` accumulates via `Reduce()`. `checkPasswordHistory` is wired into `setPasswordCommand`; `ChangePassword` and `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`; admin `SetPassword` does not (passes `nil`). `Errors.User.Password.Reused` in all 22 Go locale yaml files.
- Session 5: **Task 4 — Login-app UI.** Commits `6cf83ed53`, `326105b39`. `historyCount?: number` prop on three form components renders `<Alert type={AlertType.INFO}>` only when `historyCount > 0` using i18n key `password.complexity.historyHint` (ICU plural). All 14 `apps/login/locales/*.json` files have the new key.
- Session 6: **Task 5 — Console UI.** Commits `a9c3a8fb8`, `17450b72c`. Console password-complexity policy form takes a `historyCount` numeric input bound via `[(ngModel)]` against a getter/setter that proxies `complexityData.historyCount` (NOT reactive forms — sibling `password-age` form IS reactive forms, complexity is not). The input has `class="history-count-input"`. Submit handlers for admin (default) + mgmt (custom add/update) all wire it. `POLICY.PWD_COMPLEXITY.HISTORYCOUNT: "Password history (generations)"` is in all 22 `console/src/assets/i18n/*.json` files.
- Session 7: **Task 6 — Integration + E2E tests.** Commits `ffdfcb6de`, `36406e51b`, `98e2fc98f`, `6c5b9f0a9`. `TestServer_PasswordHistoryReuse` added to v2 + v2beta integration test files with two sub-scenarios (A: outside-window permitted + in-window rejected; B: current-hash inclusion via verify-code path). Cypress `password-complexity.cy.ts` replaces a stub with 3 instance-scope tests (history-count input visible; save `history_count=3` and persist after reload + restore to 0; save 0 and persist) plus 1 org-scope visibility test. `go vet -tags integration` clean; Cypress typecheck clean for the new file (3 pre-existing errors in `cypress/support/commands.ts` are untouched).

  **Important Session 7 nuances:**
  - The first sub-agent dispatch produced a buggy "in-window" assertion (chain `pw0→pw1→pw2` then `pw2→pw0` with `history_count=2`). With `checkPasswordHistory`'s "first N entries of `[current] ++ PreviousHashes`" semantic, `history_count=2` only checks `[pw2, pw1]` — pw0 at index 2 is OUTSIDE the window. The architect caught this against the Session 4 unit test `user_human_password_test.go:1596` ("history_count=3 permits new password just outside history window (previous[2])"), which establishes that `HistoryCount=N` covers `current + first (N-1) previous` = N entries. A second sub-agent dispatch (commit `6c5b9f0a9`) shortened the in-window chain to `pw0→pw1` then `pw1→pw0`, which correctly places pw0 at `PreviousHashes[0]` and triggers rejection.
  - **Implication for any future reuse-window scenarios:** `HistoryCount=N` includes the current hash + the most-recent `N-1` previous hashes. To force a hit, ensure the target plaintext sits at `PreviousHashes[k]` where `k <= N-2`.
  - Neither `go test -tags integration` nor `cypress run` was actually executed locally — both require a dedicated zitadel test stack (DB + server + login app) that this dev environment does not have spun up dedicated to this repo. The test code compiles and vets clean; it must be exercised in CI before claiming green.
  - Admin round-trip test was deliberately skipped per the plan ("if no parallel test exists, skip") — `internal/api/grpc/admin/integration_test/` has no existing complexity-policy test file to extend.

**Your job in Session 8: Task 7 — architect review + PR.**

This is **not** a sub-agent task. The architect (you) reads the full diff, runs the review skill, fixes anything flagged, and opens the PR. The plan defines 5 steps.

**Step-by-step:**

1. **Confirm branch state.** `git status` → clean. `git log --oneline main..HEAD` → should show ~18 commits across the 7 tasks. Sanity-check that nothing's missing.

2. **Read the full diff.** `git diff main...HEAD --stat` for the shape, then `git diff main...HEAD` for the substance. Spot-check:
   - **Domain & events:** `HistoryCount uint64` is on `PasswordComplexityPolicy`, the `PasswordComplexityPolicyAddedEvent` payload, and as `*uint64` on `PasswordComplexityPolicyChangedEvent` with a `ChangeHistoryCount` option func.
   - **Write model:** `HumanPasswordWriteModel.PreviousHashes []string` is appended in `Reduce()` for `HumanPasswordChangedEvent`, with the "prepend old EncodedHash if non-empty" semantic. Legacy `Secret`-only events are skipped from history.
   - **Check function:** `checkPasswordHistory` builds `[current] ++ previousHashes`, truncates to `HistoryCount`, iterates calling `c.userPasswordHasher.Verify`, returns `INVALID_ARGUMENT / Errors.User.Password.Reused / COMMAND-PwReuse` on match.
   - **Caller wiring:** `ChangePassword` and `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`. Admin `SetPassword` and initial-set paths (register, invite, init, email-verify) pass `nil`.
   - **Projection + migration:** `history_count INT8 NOT NULL DEFAULT 0` column, `cmd/setup/70` step registered in `config.go` and `setup.go`.
   - **Proto:** `history_count` on `PasswordComplexityPolicy`, both admin + management request messages, and v2 + v2beta `PasswordComplexitySettings`.
   - **UI:** Login app forms (three) render the hint via `password.complexity.historyHint`. Console complexity form has the numeric input + `[(ngModel)]="historyCount"` binding.
   - **Tests:** Six new command-layer test cases (commit `0797dd49b`), two integration scenarios + verify-code scenario across v2 + v2beta (commits `ffdfcb6de`, `36406e51b`, `6c5b9f0a9`), three Cypress instance-scope tests + 1 org-scope visibility test (commit `98e2fc98f`).
   - **i18n coverage:** 22 Go locales (`Errors.User.Password.Reused`), 14 login-app locales (`password.complexity.historyHint`), 22 console locales (`POLICY.PWD_COMPLEXITY.HISTORYCOUNT`).

3. **Run the review skill.** Try `/review` (if loaded) or read `~/.claude/skills/` for the canonical name in this environment. If a `superpowers:review` skill loads it as documented earlier, use that. Capture the output. The skill scans the diff for SQL safety, trust-boundary issues, conditional side effects, and structural problems. Fix anything that looks legitimately flagged; ignore false positives but document why in the PR description.

4. **Reconcile spec + plan.** Read both docs once more and confirm:
   - Every "Architecture" section maps to a commit.
   - Every "Schema changes" entry has a corresponding diff.
   - The two intentional placeholders called out in the spec (proto field numbers and migration step number) were resolved (sub-agents picked at time of work).
   - Open questions in the spec are answered: the verify decision on `password-complexity-view.component` (left alone, deliberate — Session 6 decided correctly that client-side hash comparison isn't possible); the Cypress decision (one-policy-save test added, reuse-rejection e2e skipped per the spec's default).
   - If anything in the spec is now contradicted by what shipped (e.g. a different file got touched), update the spec to reflect reality. Commit any doc edits as their own commit.

5. **Verify gates one more time before PR.** Run, with the architect's hands, each gate command listed in the plan (Tasks 1–6). Capture each command + output. If any are red:
   - For the integration test gate: it needs Docker + a zitadel test stack. If the architect's machine doesn't have one running, document this as "needs CI" in the PR description rather than blocking the PR.
   - For the Cypress gate: same — `cypress run` needs the login app + zitadel server up. Document.
   - For all Go unit-test gates, login-app build, console build: these should run cleanly on the architect's machine. They MUST be green before opening the PR.

6. **Open the PR.** `gh pr create` against `main`. PR body must reference the spec doc and summarize:
   - One-line problem statement (NIST 800-63B password reuse prevention).
   - Bullet list of what shipped per layer (domain, events, projection, proto, command, login UI, console UI, tests).
   - Test plan (which gates passed locally, which need CI).
   - "Generated with Claude Code" footer per the standard PR template.

**Hard rules:**

- This is the architect's own work — no sub-agent dispatch in this session.
- If anything in the review surfaces a real bug (not a style nit), **stop and decide**: fix it inline (small) or write a recovery prompt for Session 9 (large). Don't ship known bugs.
- If the spec or plan needs reconciling, commit those edits before opening the PR so the PR diff is complete.
- The PR is the final deliverable. Make it clean.

**Useful one-shot commands:**

```bash
git diff main...HEAD --stat                                 # shape of the change
git log --oneline main..HEAD                                # commit titles
ls /home/zacharya/projects/zitadel/.artifacts/bin/linux/amd64 # confirm proto toolchain available
go vet -tags integration ./internal/api/grpc/user/v2/integration_test/... \
                          ./internal/api/grpc/user/v2beta/integration_test/...
go test ./internal/command/... ./internal/query/... ./internal/repository/... \
        ./internal/api/grpc/... 2>&1 | tail -40
cd console && PATH=/home/zacharya/projects/zitadel/.artifacts/bin/linux/amd64:$PATH pnpm build 2>&1 | tail -20
cd apps/login && pnpm typecheck 2>&1 | tail -20
cd tests/functional-ui && pnpm exec tsc --noEmit -p cypress/tsconfig.json 2>&1 | tail -10
```

**User preferences (from memory):**

- Sonnet for mechanical work, Opus for hard reasoning. Task 7 is architect-only — no delegation.
- Skill priority: invoke any matching skill before action. `/review` (or `superpowers:requesting-code-review`) is the canonical step here.

Bite-sized scope. Close the loop. Be the architect.

# Session 4 handoff — password reuse prevention

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're mid-feature: password reuse prevention for NIST alignment, multi-session, architect-driven (you are the architect, you do not write code — you dispatch Sonnet sub-agents and review). Branch: `password-age-new`.

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan with file maps and gate commands. Each task is one session.

**Sessions 1–3 (already done):**
- Session 1: spec + plan + Session 2 handoff. Commits `8a1c4fc1d`, `cb62b9631`, `85114cc51`.
- Session 2: **Task 1 — Schema, events, projection, migration.** Commits `6283c6ecd`, `bdef6e80b`, `8d09be34c`, `82254b056`. `HistoryCount uint64` is live across domain, repository events, write models, projection, query, and migration (`cmd/setup/70`).
- Session 3: **Task 2 — Proto + gRPC converters + regen.** Commits `b6debd810`, `b12aeb7b5`. `history_count` is now wire-visible end to end: added to `policy.proto` (field 8), `admin.proto` UpdatePasswordComplexityPolicyRequest (field 6), `management.proto` Add+UpdateCustomPasswordComplexityPolicyRequest (field 6), `settings/v2` + `v2beta` PasswordComplexitySettings (field 7). Converters in `internal/api/grpc/policy/`, `admin/`, `management/`, `settings/v2/`, `settings/v2beta/` all map the field. Generated `.pb.go` files are gitignored — a fresh clone needs `PATH=.artifacts/bin/linux/amd64:$PATH buf generate` (the buf binary lives in `.artifacts/bin/`). All gate commands green (`go build ./...`, `go test ./internal/api/grpc/...`). 12 files in the diff.

**Your job in Session 4: Execute Task 3 from the plan — "Password change wiring + history check + i18n."**

This is the enforcement bite. With the data plane (Task 1) and wire format (Task 2) in place, this session makes the feature *do* something: replay password events into `PreviousHashes`, verify new passwords against history on `ChangePassword` and `SetPasswordWithVerifyCode`, and throw `Errors.User.Password.Reused` on a hit.

**Dispatch instructions:**

1. Confirm branch is `password-age-new` and clean (`git status`). Confirm the two Task 2 commits are present on top of Task 1's four (`git log --oneline -8`).
2. Read the spec doc fully (re-read the "History storage", "Check site", "Caller wiring", "Error handling & UX" subsections especially), then re-read Task 3 in the plan.
3. Dispatch ONE Sonnet sub-agent (Agent tool, `subagent_type: general-purpose`, `model: sonnet`) with this prompt:

   > Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — authoritative. Then execute Task 3 from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` end-to-end: all 12 steps, commit at the four boundaries the task specifies (steps 1–5 write-model + its test, 6–9 check function + wiring, 10 i18n, 11–12 command tests), run the gate command at the end and report pass/fail with the actual output.
   >
   > Tasks 1 and 2 already landed `HistoryCount uint64` on the domain + write-model + projection + query side and exposed it through proto + gRPC converters. You only need to: (a) accumulate `PreviousHashes` on `HumanPasswordWriteModel`, (b) add the `checkPasswordHistory` private method, (c) wire it into `setPasswordCommand` and through `ChangePassword` + `SetPasswordWithVerifyCode`, (d) add the i18n key. Do not touch domain, complexity-policy command handlers, projection, query, proto, gRPC converters, or migration — those are done.
   >
   > Specific guards:
   > - Use TDD: write the `TestHumanPasswordWriteModel_PreviousHashesAccumulate` test first, see it red, then implement the `Reduce()` changes. Spec semantics for prepend: when a password-changed event fires, if current `EncodedHash != ""`, prepend it to `PreviousHashes`, then overwrite `EncodedHash` with the event's new hash. Result ordering: `PreviousHashes[0]` is most-recent-prior, `[N-1]` is oldest.
   > - Legacy `Secret`-only events (events with `Secret *crypto.CryptoValue` set but no `EncodedHash`): skip them from `PreviousHashes` entirely. Don't try to decrypt-and-rehash.
   > - `checkPasswordHistory` signature is fixed by the spec: `func (c *Commands) checkPasswordHistory(ctx context.Context, newPlaintext, currentEncodedHash string, previousHashes []string, policy *domain.PasswordComplexityPolicy) error`. Build `checkList := append([]string{currentEncodedHash}, previousHashes...)`, iterate first `min(int(policy.HistoryCount), len(checkList))` entries, skip empty strings, call `c.userPasswordHasher.Verify(hash, newPlaintext)`. On Verify success → `zerrors.ThrowInvalidArgument(nil, "COMMAND-PwReuse", "Errors.User.Password.Reused")`. On `passwap.ErrPasswordMismatch` continue. Other errors → return wrapped.
   > - `setPasswordCommand` gains a `previousHashes []string` parameter. The check fires only when `previousHashes != nil` AND resolved `policy.HistoryCount > 0`. Order in the function body: complexity check first, then history check, then hash + persist.
   > - `ChangePassword` and `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`. Every other caller of `setPasswordCommand` passes `nil` — find them with `grep -rn "setPasswordCommand(" internal/command/` and fix each compile error in-place.
   > - i18n: add `Reused: "Password recently used, choose another"` under `Errors.User.Password` in `internal/static/i18n/en.yaml` AND every sibling locale file (`ls internal/static/i18n/` and add to each — non-English locales get the English string as a placeholder for translators).
   > - Sibling branch `password-age` has reference code for the same shape on the age policy: `git show password-age:internal/command/user_human_password_model.go`, `git show password-age:internal/command/user_human_password.go`. Same write-model pattern, same check-function shape — adapt to read `HistoryCount` from `*domain.PasswordComplexityPolicy` instead of `*domain.PasswordAgePolicy`.
   > - Test coverage (Task 3 step 11): five cases per the spec — `history_count=0` permitted (regression guard), `history_count=3` with new == previous[0]/[1] rejected, `history_count=3` with new == previous[2] (4th-from-now incl. current) permitted, `SetPasswordWithVerifyCode` with `history_count=3` and new == current rejected (proves current-hash inclusion), admin `SetPassword` with `history_count=3` and matching new password permitted (proves admin bypasses).
   > - Do not exceed Task 3's scope: no UI (Tasks 4–5), no integration tests (Task 6), no proto/converter edits (Task 2 done). Stay in `internal/command/` and `internal/static/i18n/`.

4. When the sub-agent returns, **verify the diff yourself** before trusting the "done" claim:
   - `git log --oneline -6` — should see ~4 new commits from this task.
   - `git diff HEAD~4 --stat` — file count sanity check. Expect: `user_human_password_model.go`, `user_human_password.go`, `user_human_password_model_test.go`, `user_human_password_test.go`, plus every file in `internal/static/i18n/` (one line per locale).
   - Re-run the gate: `go test ./internal/command/... 2>&1 | tail -50`. Must end green.
   - Spot-check `internal/command/user_human_password.go` for the new `checkPasswordHistory` function — confirm the order in `setPasswordCommand` is complexity-then-history-then-hash, not history-before-complexity.
   - Skim the new write-model test to confirm the prepend ordering matches the spec (`PreviousHashes[0]` is most-recent-prior).
   - Confirm i18n key landed in every locale file: `grep -l "Reused:" internal/static/i18n/` — count should equal the number of locale files.

5. **If green and diff looks sane:** overwrite `next_prompt.md` with a Session 5 handoff for Task 4 (Login-app UI — the hint above submit on change/set/register password forms), mirroring this file's structure. Commit the new handoff. Stop.

6. **If red or diff looks wrong:** do NOT dispatch a second sub-agent reflexively. Read the failure, decide if it's a sub-agent mistake (re-dispatch with a corrective prompt) or a plan gap (update the plan first, then re-dispatch). Either way, write a recovery `next_prompt.md` describing the state before stopping. Common failure modes to watch for: (a) `setPasswordCommand` callers in test files broken by the new parameter — easy fix, pass `nil`; (b) write-model test asserting wrong ordering — re-read spec on prepend semantics; (c) `checkPasswordHistory` placed *before* complexity check in `setPasswordCommand` — wrong order.

**Hard rules:**

- You don't write code. Only design docs / plan files / next_prompt.md.
- All implementation goes through Sonnet sub-agents.
- One sub-agent at a time per session — they need a clean handoff, not parallel races on the same files.
- Verify before trusting. Sub-agents will claim success. Read the diff and re-run the gate yourself.
- If scope grows past Task 3, stop and write a recovery `next_prompt.md`. Don't try to do Tasks 3 + 4 in one session.

**User preferences (from memory):**

- Sonnet for structured/mechanical work like write-model reduce + check function + i18n keys. Opus only for hard reasoning. ✓
- Skill priority: invoke any matching skill before action. `superpowers:subagent-driven-development` may apply once you start dispatching — check, but the user has consistently preferred the lighter "architect-verification" path (skip spec-reviewer/code-quality-reviewer sub-agents, do the verification yourself by reading the diff and re-running the gate). Default to that unless the user opts up.

**Useful context from Sessions 1–3 (so you don't re-derive):**

- Field naming is consistent: Go `HistoryCount uint64`, SQL column `history_count`, proto `uint64 history_count`. UI hint key in later tasks is `historyHint`. Error i18n key is `Errors.User.Password.Reused`, error code `COMMAND-PwReuse`.
- The legacy event type to handle in `Reduce()` is `UserV1PasswordChangedType` (in addition to `HumanPasswordChangedEvent`). Both are already in `Query()` subscriptions per spec note — verify with `grep "UserV1PasswordChangedType\|HumanPasswordChangedType" internal/command/user_human_password_model.go`.
- `c.userPasswordHasher` is a `*crypto.PasswordHasher` (wraps passwap). `.Verify(encodedHash, plaintext)` returns `(updatedHash string, err error)` — err is `nil` on success, `passwap.ErrPasswordMismatch` on no match, other errors on hasher problems. The updated hash returned on success can be ignored in this check path (it's used for rehash-on-upgrade in normal verification, but here we just care match/no-match).
- The `setPasswordCommand` function is the chokepoint — *every* server-side password set goes through it. That's why adding the parameter there and passing `nil` from non-enforcing callers is the cleanest wiring. Don't try to put the check higher up.
- After Task 3 the feature is behaviorally real on the backend. Task 4 surfaces the hint to login users. Task 5 lets admins configure the value in Console. Task 6 proves it works end-to-end.

Good luck. Bite-sized scope. Be the architect.

# Session 2 handoff — password reuse prevention

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're mid-feature: password reuse prevention for NIST alignment, multi-session, architect-driven (you are the architect, you do not write code — you dispatch Sonnet sub-agents and review). Branch: `password-age-new`.

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan with file maps and gate commands. Each task is one session.

**Session 1 (already done):** Survey + spec + plan + this handoff. Commits `8a1c4fc1d` (spec) and `cb62b9631` (plan).

**Your job in Session 2: Execute Task 1 from the plan — "Schema, events, projection, migration."**

This is the data-plane foundation. No user-visible behavior changes. Touches ~13 Go files plus a new `cmd/setup/` migration step.

**Dispatch instructions:**

1. Confirm branch is `password-age-new` and clean (`git status`).
2. Read the spec doc fully, then re-read Task 1 in the plan.
3. Dispatch ONE Sonnet sub-agent (Agent tool, `subagent_type: general-purpose`, `model: sonnet`) with this prompt:

   > Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — authoritative. Then execute Task 1 from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` end-to-end: all 11 steps, commit at the four boundaries the task specifies, run the gate command at the end. Report pass/fail with actual `go test` output. Do not mark complete on red. The branch `password-age` (sibling, accessible via `git show password-age:<path>`) has reference code on the *age* policy — same shape, just adapt to *complexity* policy. Use it as a reference but do not blindly copy; the policy is different.
   >
   > One specific guard: pick the next available `cmd/setup/` step number by inspecting the directory, do not assume 70.

4. When the sub-agent returns, **verify the diff yourself** before trusting the "done" claim:
   - `git log password-age-new ^HEAD~5 --oneline` (last 5 commits — should see ~4 from this task)
   - `git diff HEAD~4 --stat` — sanity check file count matches the plan's file list
   - Re-run the gate command: `go test ./internal/domain/... ./internal/command/... ./internal/query/... ./internal/repository/... ./cmd/setup/...`
   - Skim 2-3 of the modified Go files to confirm they look right (not just compile)

5. **If green and diff looks sane:** overwrite `next_prompt.md` with a new Session 3 handoff (mirror this file's structure, but for Task 2). Commit the new handoff. Stop.

6. **If red or diff looks wrong:** do NOT dispatch a second sub-agent reflexively. Read the failure, decide if it's a sub-agent mistake (re-dispatch with a corrective prompt) or a plan gap (update the plan first, then re-dispatch). Either way, write a recovery `next_prompt.md` describing the state before stopping.

**Hard rules:**

- You don't write code. Only design docs / plan files / next_prompt.md.
- All implementation goes through Sonnet sub-agents.
- One sub-agent at a time per session — they need a clean handoff, not parallel races on the same files.
- Verify before trusting. Sub-agents will claim success. Read the diff.
- If scope grows past Task 1, stop and write a recovery `next_prompt.md`. Don't try to do Tasks 1 + 2 in one session.

**User preferences (from memory):**

- Sonnet for structured/mechanical work like this. Opus only for hard reasoning. ✓
- Skill priority: invoke any matching skill before action. `superpowers:subagent-driven-development` may apply once you start dispatching — check.

**Useful context from Session 1's exploration (so you don't re-derive):**

- The sibling branch `password-age` already implemented this feature on the age policy. We *moved* the field to complexity policy. Reference, not portable wholesale.
- Reduce semantics for `PreviousHashes`: prepend current `EncodedHash` to slice when a new password-changed event is reduced, then update `EncodedHash`. Most recent first. Skip legacy `Secret`-only events.
- The check list at enforcement time is `[currentEncodedHash] ++ PreviousHashes`, truncated to `historyCount`. Current is included because reset-with-code has no built-in same-password rejection (unlike `ChangePassword`'s `VerifyAndUpdate` path).
- Enforcement applies to `ChangePassword` and `SetPasswordWithVerifyCode` only. Admin set + initial-set paths pass `nil` and bypass.
- Default `history_count = 0` (disabled).
- Error message: generic "Password recently used, choose another". UI hint may show count.

Good luck. Bite-sized scope. Be the architect.

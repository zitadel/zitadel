# Session 3 handoff — password reuse prevention

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're mid-feature: password reuse prevention for NIST alignment, multi-session, architect-driven (you are the architect, you do not write code — you dispatch Sonnet sub-agents and review). Branch: `password-age-new`.

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan with file maps and gate commands. Each task is one session.

**Sessions 1–2 (already done):**
- Session 1: spec + plan + Session 2 handoff. Commits `8a1c4fc1d`, `cb62b9631`, `85114cc51`.
- Session 2: **Task 1 — Schema, events, projection, migration.** Commits `6283c6ecd`, `bdef6e80b`, `8d09be34c`, `82254b056`. All four gate packages green (`./internal/domain/...`, `./internal/command/...`, `./internal/query/...`, `./internal/repository/...`, `./cmd/setup/...`). Migration is `cmd/setup/70.go` + `70.sql`. The Go `HistoryCount uint64` field is now live across domain, repository events, write models, projection, query, and migration. Test fixtures across `internal/command/` were updated for the new constructor signature (this is why the diff touched 30 files instead of the plan's ~16).

**Your job in Session 3: Execute Task 2 from the plan — "Proto + gRPC converters + regen."**

This is the wire format. With Task 1's domain state in place, expose `history_count` through proto and the gRPC converters so the UI tasks (4 and 5) have something to read.

**Dispatch instructions:**

1. Confirm branch is `password-age-new` and clean (`git status`). Confirm the four Task 1 commits are present (`git log --oneline -8`).
2. Read the spec doc fully, then re-read Task 2 in the plan.
3. Dispatch ONE Sonnet sub-agent (Agent tool, `subagent_type: general-purpose`, `model: sonnet`) with this prompt:

   > Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — authoritative. Then execute Task 2 from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` end-to-end: all 9 steps, commit at the two boundaries the task specifies (steps 1–5 proto+regen, steps 6–9 converters+tests), run the gate commands at the end and report pass/fail with actual output.
   >
   > Task 1 already landed `HistoryCount uint64` on the domain side. You only need to wire it through proto messages and gRPC converters. Do not touch domain, command, projection, query, write models, or migration — those are done.
   >
   > Specific guards:
   > - Pick the next available proto field number for each message by inspecting the current `.proto` file. Do not assume — main may have advanced.
   > - Discover the canonical proto regen command from the repo's Makefile or `buf.gen.yaml`. Common candidates: `make generate`, `buf generate`, `make grpc-gen`. Run whatever the repo uses.
   > - If the regen command produces generated `.pb.go` files that aren't already in `.gitignore`, commit them — match how previous proto-add commits handled this (look at recent commits that touched `proto/`).
   > - Sibling branch `password-age` has reference code for the same shape on the age policy: `git show password-age:proto/zitadel/policy.proto`, `git show password-age:internal/api/grpc/policy/password_age_policy.go` etc. Use as reference, adapt to complexity.
   > - Do not exceed Task 2's scope: no domain edits, no write-model edits, no UI, no i18n. Stay in `proto/` and `internal/api/grpc/`.

4. When the sub-agent returns, **verify the diff yourself** before trusting the "done" claim:
   - `git log --oneline password-age-new ^HEAD~2` (last 2 commits — should see ~2 from this task; more if `.pb.go` regen was committed separately)
   - `git diff HEAD~2 --stat` — file count sanity check matches Task 2's file list (`.proto` files, converters, generated `.pb.go`)
   - Re-run the gate commands: `go build ./...` and `go test ./internal/api/grpc/...`
   - Skim 2–3 of the modified converter files to confirm they look right (the line should be `HistoryCount: policy.HistoryCount` or equivalent — actual mapping, not just compile)
   - Confirm at least one `.proto` file shows the new `uint64 history_count = N;` line with `N` being a sane next-free number

5. **If green and diff looks sane:** overwrite `next_prompt.md` with a Session 4 handoff for Task 3 (Password change wiring + history check + i18n — the enforcement bite), mirroring this file's structure. Commit the new handoff. Stop.

6. **If red or diff looks wrong:** do NOT dispatch a second sub-agent reflexively. Read the failure, decide if it's a sub-agent mistake (re-dispatch with a corrective prompt) or a plan gap (update the plan first, then re-dispatch). Either way, write a recovery `next_prompt.md` describing the state before stopping.

**Hard rules:**

- You don't write code. Only design docs / plan files / next_prompt.md.
- All implementation goes through Sonnet sub-agents.
- One sub-agent at a time per session — they need a clean handoff, not parallel races on the same files.
- Verify before trusting. Sub-agents will claim success. Read the diff.
- If scope grows past Task 2, stop and write a recovery `next_prompt.md`. Don't try to do Tasks 2 + 3 in one session.

**User preferences (from memory):**

- Sonnet for structured/mechanical work like proto + converter wiring. Opus only for hard reasoning. ✓
- Skill priority: invoke any matching skill before action. `superpowers:subagent-driven-development` may apply once you start dispatching — check.
- Session 2 ran the lighter "handoff-only architect verification" path (skipped the skill's spec-reviewer/code-quality-reviewer sub-agents). Default to the same lighter path unless the user opts up.

**Useful context from Sessions 1–2 (so you don't re-derive):**

- The Go field name is `HistoryCount uint64`; SQL column `history_count`; expect proto to use `uint64 history_count`. Keep the naming consistent — i18n key in later tasks is `historyHint`.
- The repo uses `buf` for proto. Check `buf.gen.yaml` and `Makefile` for the regen entrypoint.
- Settings v2 and v2beta both exist as separate generated APIs; both need the field added.
- The login app picks up the new proto field automatically via `getPasswordComplexitySettings` in `apps/login/src/lib/zitadel.ts` — no extra wiring needed in this task. UI work is Task 4.
- `internal/api/grpc/admin/policy_password_complexity_converter.go` and the management equivalent handle request→domain conversion; both need the field.
- After Task 2 the field is wire-visible end to end. Task 3 then makes the field *do* something (the enforcement check on password change).

Good luck. Bite-sized scope. Be the architect.

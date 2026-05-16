# Session 5 handoff — password reuse prevention

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're mid-feature: password reuse prevention for NIST alignment, multi-session, architect-driven (you are the architect, you do not write code — you dispatch Sonnet sub-agents and review). Branch: `password-age-new`.

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan with file maps and gate commands. Each task is one session.

**Sessions 1–4 (already done):**
- Session 1: spec + plan + Session 2 handoff. Commits `8a1c4fc1d`, `cb62b9631`, `85114cc51`.
- Session 2: **Task 1 — Schema, events, projection, migration.** Commits `6283c6ecd`, `bdef6e80b`, `8d09be34c`, `82254b056`. `HistoryCount uint64` is live across domain, repository events, write models, projection, query, and migration (`cmd/setup/70`).
- Session 3: **Task 2 — Proto + gRPC converters + regen.** Commits `b6debd810`, `b12aeb7b5`. `history_count` is wire-visible end to end across `policy.proto`, `admin.proto`, `management.proto`, and settings v2 + v2beta. Converters in `internal/api/grpc/policy/`, `admin/`, `management/`, `settings/v2/`, `settings/v2beta/` all map the field. Generated `.pb.go` files are gitignored — a fresh clone needs `PATH=.artifacts/bin/linux/amd64:$PATH buf generate`.
- Session 4: **Task 3 — Password change wiring + history check + i18n.** Commits `09d5231a6`, `c9d192c01`, `730d35bee`, `0797dd49b`. `HumanPasswordWriteModel` now accumulates `PreviousHashes []string` via `Reduce()` (most-recent-first, legacy `Secret`-only events skipped). `checkPasswordHistory` is wired into `setPasswordCommand` between the complexity check and the hash step, firing only when `previousHashes != nil` AND `policy.HistoryCount > 0`. `ChangePassword` and `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`; all other callers pass `nil`. `Errors.User.Password.Reused` is in all 22 locale yaml files. Gate green (`go test ./internal/command/...`).

  Session 4 also made three minor refactors not in the spec: (a) added explicit `currentEncodedHash string` param to `setPasswordCommand` (cleaner than a re-fetch inside), (b) added `passwap.ErrNoVerifier` skip in `checkPasswordHistory` (defensive against legacy hash formats that no verifier recognises), (c) inlined and removed the unused `checkPasswordComplexity` helper. All adjacent to the core change, low risk, semantics preserved.

**Your job in Session 5: Execute Task 4 from the plan — "Login-app UI."**

This is the user-facing surface. With Tasks 1–3 the backend rejects reused passwords. Now end users need a hint above the submit button so they don't waste a round-trip to discover the rule. Three forms, one i18n key per locale, plus auto-flow of the new proto field through `getPasswordComplexitySettings`.

**Dispatch instructions:**

1. Confirm branch is `password-age-new` and clean (`git status`). Confirm the four Task 3 commits are present on top of Task 2's two and Task 1's four (`git log --oneline -12`).
2. Read the spec doc fully (re-read the "Login-app UI hints" subsection especially), then re-read Task 4 in the plan.
3. Dispatch ONE Sonnet sub-agent (Agent tool, `subagent_type: general-purpose`, `model: sonnet`) with this prompt:

   > Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — authoritative. Then execute Task 4 from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` end-to-end: all 6 steps, commit at the two boundaries the task specifies (steps 1–3 components, 4–5 i18n + wiring), run the gate commands at the end and report pass/fail with the actual output.
   >
   > Tasks 1–3 already landed `HistoryCount` end-to-end: domain field, projection column, query scan, proto field, gRPC converters, write-model history accumulation, `checkPasswordHistory` enforcement, and `Errors.User.Password.Reused` in every Go locale yaml. You only need to: (a) surface the count to the three password forms in the login app, (b) render an info hint above submit when `historyCount > 0`, (c) add the `password.complexity.historyHint` translation key to every `apps/login/locales/*.json` file.
   >
   > Specific guards:
   > - Confirm `getPasswordComplexitySettings` in `apps/login/src/lib/zitadel.ts` returns the `historyCount` field after Task 2's proto regen. Auto-mapping should already work; if there's an explicit select that omits it, add it. Check by reading the function and tracing what proto message type it returns.
   > - For each of the three forms (`change-password-form.tsx`, `set-password-form.tsx`, `set-register-password-form.tsx`): add an optional prop, render the hint above the submit button only when `historyCount && historyCount > 0`. Match the existing complexity-hint UI style — look at the existing `<PasswordComplexity>` component or the `passwordMinLength`/`hasSymbol` rendering for the right alert/info element. Don't invent a new pattern.
   > - i18n key: `password.complexity.historyHint` with value `"You can't reuse your last {count} {count, plural, one {password} other {passwords}}."` in `apps/login/locales/en.json`. Add the same key with the English string as a placeholder to every sibling locale file (`ls apps/login/locales/` and add to each).
   > - Parent server-component wiring: `apps/login/src/app/(login)/password/change/page.tsx` already fetches the complexity settings — pass `historyCount` to `change-password-form.tsx`. For the set/register/reset pages, find the analogous server components (`grep -rn "set-password-form\|set-register-password-form" apps/login/src/app/`) and mirror the same pattern. If a page doesn't already fetch complexity settings, add the fetch using the same `getPasswordComplexitySettings` helper.
   > - Sibling branch `password-age` has a reference impl for the change-password form only: `git show password-age:apps/login/src/components/change-password-form.tsx`. Same hint pattern, you just need to apply it to two more forms.
   > - Pluralization: use the existing i18n library's plural support. Inspect how other keys in `apps/login/locales/en.json` handle counts (e.g. minLength messages) and match. If the library doesn't support ICU plurals, fall back to a simpler `"You can't reuse your last {count} passwords."` and confirm with the user before committing.
   > - Do not exceed Task 4's scope: no Console UI (Task 5), no integration/Cypress tests (Task 6), no backend edits (Task 3 done).
   >
   > **Gate (both must be green):**
   > ```
   > pnpm -F login typecheck 2>&1 | tail -30
   > pnpm -F login build 2>&1 | tail -50
   > ```
   > Discover the canonical command names from `apps/login/package.json` if `pnpm -F login` doesn't work. Some monorepos use `pnpm --filter @zitadel/login ...` or have the scripts only in the subdirectory (`cd apps/login && pnpm typecheck`). Try the alternatives.

4. When the sub-agent returns, **verify the diff yourself** before trusting the "done" claim:
   - `git log --oneline -8` — should see ~2 new commits from this task.
   - `git diff HEAD~2 --stat` — file count sanity check. Expect: 3 component `.tsx` files, every `apps/login/locales/*.json` file (one line per locale), possibly 1–3 server-component `page.tsx` files.
   - Re-run both gates: typecheck and build. Must end clean.
   - Spot-check `apps/login/src/components/change-password-form.tsx` for the new prop + hint render. Confirm the hint is only shown when `historyCount > 0` (not on `historyCount == 0` or undefined).
   - Confirm i18n key landed in every locale file: `grep -l "historyHint" apps/login/locales/` — count should equal the number of locale files.
   - Open the hint render site in `set-password-form.tsx` and `set-register-password-form.tsx` — confirm they actually render the alert, not just accept a prop they ignore. (Easy way to miss a step.)

5. **If green and diff looks sane:** overwrite `next_prompt.md` with a Session 6 handoff for Task 5 (Console UI — Angular form control + numeric input + i18n key), mirroring this file's structure. Commit the new handoff. Stop.

6. **If red or diff looks wrong:** do NOT dispatch a second sub-agent reflexively. Read the failure, decide if it's a sub-agent mistake (re-dispatch with a corrective prompt) or a plan gap (update the plan first, then re-dispatch). Either way, write a recovery `next_prompt.md` describing the state before stopping. Common failure modes to watch for: (a) `getPasswordComplexitySettings` not returning the new field because of an explicit field mask — add the field; (b) ICU plural syntax not supported by the i18n library — fall back to a simpler hint and rerun; (c) server-component page never passes the prop and the form just renders nothing — find the parent and add the wiring; (d) hint shown even when `historyCount == 0` (defaulting falsy values is easy to get wrong) — tighten the conditional.

**Hard rules:**

- You don't write code. Only design docs / plan files / next_prompt.md.
- All implementation goes through Sonnet sub-agents.
- One sub-agent at a time per session — they need a clean handoff, not parallel races on the same files.
- Verify before trusting. Sub-agents will claim success. Read the diff and re-run the gate yourself.
- If scope grows past Task 4, stop and write a recovery `next_prompt.md`. Don't try to do Tasks 4 + 5 in one session.

**User preferences (from memory):**

- Sonnet for structured/mechanical work like component prop threading + i18n keys. Opus only for hard reasoning. ✓
- Skill priority: invoke any matching skill before action. `superpowers:subagent-driven-development` may apply — check, but the user has consistently preferred the lighter "architect-verification" path (skip spec-reviewer/code-quality-reviewer sub-agents, do the verification yourself by reading the diff and re-running the gate). Default to that unless the user opts up.

**Useful context from Sessions 1–4 (so you don't re-derive):**

- Field naming is consistent: Go `HistoryCount uint64`, SQL column `history_count`, proto `uint64 history_count`. UI hint key is `password.complexity.historyHint` (login app) / `POLICY.PWD_COMPLEXITY.HISTORYCOUNT` (console — Task 5). Backend error i18n key is `Errors.User.Password.Reused`, error code `COMMAND-PwReuse`.
- The proto field `history_count` is at field 7 on `PasswordComplexitySettings` in both `settings/v2` and `v2beta`. After Task 2's `buf generate`, the TypeScript client types should expose it as `historyCount` (proto-loader's camelCase conversion).
- `apps/login/locales/` is the directory for login-app translations. Count of files: run `ls apps/login/locales/ | wc -l` to discover. Every file gets the new key.
- The login app uses Next.js App Router. Forms are client components (`"use client"`), parent pages are server components that fetch settings and pass them down as props. The pattern is: `getPasswordComplexitySettings(orgId)` → `<Form complexitySettings={…} historyCount={settings.historyCount} />`.
- After Task 4 end users see the hint. Task 5 surfaces the value to admins in Console. Task 6 proves it works end-to-end with integration + Cypress tests.

Good luck. Bite-sized scope. Be the architect.

# Session 6 handoff — password reuse prevention

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're mid-feature: password reuse prevention for NIST alignment, multi-session, architect-driven (you are the architect, you do not write code — you dispatch Sonnet sub-agents and review). Branch: `password-age-new`.

**Read these in order before doing anything:**

1. `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` — authoritative design.
2. `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` — 7-task plan with file maps and gate commands. Each task is one session.

**Sessions 1–5 (already done):**
- Session 1: spec + plan + Session 2 handoff. Commits `8a1c4fc1d`, `cb62b9631`, `85114cc51`.
- Session 2: **Task 1 — Schema, events, projection, migration.** Commits `6283c6ecd`, `bdef6e80b`, `8d09be34c`, `82254b056`. `HistoryCount uint64` is live across domain, repository events, write models, projection, query, and migration (`cmd/setup/70`).
- Session 3: **Task 2 — Proto + gRPC converters + regen.** Commits `b6debd810`, `b12aeb7b5`. `history_count` is wire-visible end to end across `policy.proto`, `admin.proto`, `management.proto`, and settings v2 + v2beta. Converters in `internal/api/grpc/policy/`, `admin/`, `management/`, `settings/v2/`, `settings/v2beta/` all map the field. Generated `.pb.go` and `.ts` files are gitignored — a fresh clone needs `PATH=.artifacts/bin/linux/amd64:$PATH buf generate` at the repo root, plus `buf generate` inside `packages/zitadel-proto/` for the TS client.
- Session 4: **Task 3 — Password change wiring + history check + i18n.** Commits `09d5231a6`, `c9d192c01`, `730d35bee`, `0797dd49b`. `HumanPasswordWriteModel` accumulates `PreviousHashes []string` via `Reduce()`. `checkPasswordHistory` is wired into `setPasswordCommand`; only `ChangePassword` and `SetPasswordWithVerifyCode` pass `wm.PreviousHashes`. `Errors.User.Password.Reused` is in all 22 Go locale yaml files. Gate green.
- Session 5: **Task 4 — Login-app UI.** Commits `6cf83ed53`, `326105b39`. `historyCount?: number` prop added to `change-password-form.tsx`, `set-password-form.tsx`, `set-register-password-form.tsx`. Hint renders `<Alert type={AlertType.INFO}>` only when `historyCount > 0` using i18n key `password.complexity.historyHint` (ICU plural form). All 14 `apps/login/locales/*.json` files have the new key (en.json with the real string; the other 13 with the English placeholder). Three parent server components pass `historyCount={Number(passwordComplexity.historyCount)}`: `(login)/password/change/page.tsx`, `(login)/password/set/page.tsx`, `(login)/register/password/page.tsx`. Gate green (`pnpm exec tsc --noEmit` + `pnpm build` — pre-existing errors in unrelated `*.test.ts` files only, build clean with one pre-existing Sass deprecation warning).

  Sub-agent had to run `buf generate` inside `packages/zitadel-proto/` to pick up `history_count = 7` on `PasswordComplexitySettings` — Task 2's proto regen was never persisted because the generated TS is gitignored. Same will likely be true in Session 6 for the console's `src/app/proto/generated/` directory.

**Your job in Session 6: Execute Task 5 from the plan — "Console UI."**

This is the admin-facing surface. Operators need a numeric input on the password complexity policy form to set `historyCount`. Form control + input + i18n keys + submit-handler wiring across 22 locales.

**Dispatch instructions:**

1. Confirm branch is `password-age-new` and clean (`git status`). Confirm Session 5's two commits are on top of Session 4's four (`git log --oneline -10`).
2. Read the spec doc fully (re-read the "Console (Angular)" subsection), then re-read Task 5 in the plan.
3. Dispatch ONE Sonnet sub-agent (Agent tool, `subagent_type: general-purpose`, `model: sonnet`) with this prompt:

   > Read `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md` first — authoritative. Then execute Task 5 from `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md` end-to-end: all 6 steps, commit at the boundaries the task specifies. Run the gate commands at the end and report pass/fail with the actual output.
   >
   > Tasks 1–4 already landed `HistoryCount` end-to-end through Go, proto/gRPC, and the Next.js login app. You only need to: (a) add a `historyCount` form control + numeric input to the Angular password-complexity-policy component, (b) wire it through submit/save handlers for both add and update paths, (c) add `POLICY.PWD_COMPLEXITY.HISTORYCOUNT` (label) plus optional `…_DESC` (helper text) to every `console/src/assets/i18n/*.json` file (22 locales).
   >
   > Specific guards:
   > - Files to edit:
   >   - `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.ts` — add the form control, expose a getter (mirror `minLength`/`hasSymbol` siblings), include `this.historyCount?.value ?? 0` in the gRPC request body for both add and update paths. The component uses Angular reactive forms.
   >   - `console/src/app/modules/policies/password-complexity-policy/password-complexity-policy.component.html` — add an `<input type="number" min="0" formControlName="historyCount">` element wrapped in the same `<cnsl-form-field>` (or whatever wrapper the existing inputs use). Match the existing layout exactly — do NOT invent a new visual pattern.
   >   - `console/src/assets/i18n/en.json` — add `POLICY.PWD_COMPLEXITY.HISTORYCOUNT: "Password history (generations)"`. If sibling fields like `MINLENGTH` have a `…_DESC` helper key in this file, add `HISTORYCOUNT_DESC` too with a one-sentence description (e.g. `"Number of previous passwords a user cannot reuse. 0 disables the rule."`). Add the same key(s) with the English string as a placeholder to the other 21 sibling locale files (`ls console/src/assets/i18n/` shows them).
   > - Inspect `console/src/app/modules/password-complexity-view/password-complexity-view.component.{ts,html}` (Step 5). Read both files. If the component currently renders inline hints (e.g. live-validation as a user types a new password in an admin UI), add a history-hint render gated on `historyCount > 0` using the same i18n key. If it doesn't render hints — only displays the policy settings — leave it alone. **Default to leaving it alone unless you find an actual hint render site.**
   > - The proto TS for console may be stale on disk if the workspace was freshly cloned. Generated proto for console lives in `console/src/app/proto/generated/` (gitignored). If your build fails with "Property 'historyCount' does not exist on type '…PasswordComplexityPolicy'", run `cd console && pnpm generate` to regen. This is a runtime-only fix; no committed files change.
   > - Sibling branch `password-age` has a reference impl for the **age policy** form: `git show password-age:console/src/app/modules/policies/password-age-policy/password-age-policy.component.ts`. Same form-control + numeric-input shape applies to complexity. Adapt.
   > - Do not exceed Task 5's scope: no login-app changes (Task 4 done), no integration/Cypress tests (Task 6), no backend edits (Task 3 done).
   >
   > **Commit boundaries (per plan):**
   > - Commit 1: component (steps 1–3) — `.component.ts` + `.component.html`
   > - Commit 2: i18n (step 4) — all 22 locale JSON files
   > - Commit 3 (only if you actually edit the view component in step 5): the password-complexity-view component
   >
   > **Gate (both must be green):**
   > ```
   > cd console && pnpm exec tsc --noEmit 2>&1 | tail -40
   > cd console && pnpm build 2>&1 | tail -60
   > ```
   > If `pnpm exec tsc --noEmit` finds pre-existing errors in files you didn't touch, note them in your report but do not fix — pattern matches Session 5. If `tsc --noEmit` fails because Angular's compiler is needed (e.g. template type checks), the `pnpm build` (`ng build --configuration production`) is the canonical gate; that's the one that must be clean.
   >
   > **Final report:**
   > - Commit SHAs and one-line subjects.
   > - Last 40 lines of typecheck and last 60 lines of build.
   > - Confirm: 22 locale files contain the new key, name of the component edited, whether the view component was touched and why.
   > - Call out judgement calls (e.g. whether you added a `_DESC` helper key; whether you regenerated console proto; any test files updated).
   > - If the build is red — STOP. Report status. Don't attempt a third commit unless it's obviously the same logical chunk.

4. When the sub-agent returns, **verify the diff yourself** before trusting the "done" claim:
   - `git log --oneline -10` — should see 2–3 new commits from this task on top of Session 5's two.
   - `git diff HEAD~2 --stat` (or `~3` if the view component was touched) — expect: 1 `.component.ts`, 1 `.component.html`, 22 locale `.json` files, possibly 2 `password-complexity-view.component.{ts,html}`.
   - Spot-check `password-complexity-policy.component.ts` for the new form control and submit-handler wiring (both add and update paths must include `historyCount`).
   - Spot-check `password-complexity-policy.component.html` for the numeric input.
   - Confirm i18n key landed everywhere: `grep -l "HISTORYCOUNT" console/src/assets/i18n/*.json | wc -l` — should be 22.
   - Re-run both gates: `cd console && pnpm exec tsc --noEmit` and `cd console && pnpm build`. Build must end clean.

5. **If green and diff looks sane:** overwrite `next_prompt.md` with a Session 7 handoff for Task 6 (Integration + E2E tests — gRPC v2/v2beta integration tests + Cypress complexity-policy test extension), mirroring this file's structure. Commit the new handoff. Stop.

6. **If red or diff looks wrong:** do NOT dispatch a second sub-agent reflexively. Read the failure, decide if it's a sub-agent mistake (re-dispatch with a corrective prompt) or a plan gap (update the plan first, then re-dispatch). Either way, write a recovery `next_prompt.md` describing the state before stopping. Common failure modes to watch for: (a) `historyCount` missing from the TS proto type because `console/src/app/proto/generated/` is stale — run `cd console && pnpm generate`; (b) form control added but submit handler not updated (admin sets 5, save succeeds, reload shows 0) — verify both add and update paths in the .ts diff; (c) Angular template type checking fails because the form control name isn't declared in the form group definition — verify the `formGroup` builder declares `historyCount`; (d) i18n placeholder accidentally translated by a copy-paste from sibling locale; spot-check 2–3 non-English files match the English string exactly.

**Hard rules:**

- You don't write code. Only design docs / plan files / next_prompt.md.
- All implementation goes through Sonnet sub-agents.
- One sub-agent at a time per session — they need a clean handoff, not parallel races on the same files.
- Verify before trusting. Sub-agents will claim success. Read the diff and re-run the gate yourself.
- If scope grows past Task 5, stop and write a recovery `next_prompt.md`. Don't try to do Tasks 5 + 6 in one session.

**User preferences (from memory):**

- Sonnet for structured/mechanical work like Angular form controls + i18n keys. Opus only for hard reasoning. ✓
- Skill priority: invoke any matching skill before action. `superpowers:subagent-driven-development` may apply — check, but the user has consistently preferred the lighter "architect-verification" path (skip spec-reviewer/code-quality-reviewer sub-agents, do the verification yourself by reading the diff and re-running the gate). Default to that unless the user opts up.

**Useful context from Sessions 1–5 (so you don't re-derive):**

- Field naming is consistent: Go `HistoryCount uint64`, SQL column `history_count`, proto `uint64 history_count`. Login-app i18n key is `password.complexity.historyHint`; Console i18n key is `POLICY.PWD_COMPLEXITY.HISTORYCOUNT`. Backend error key is `Errors.User.Password.Reused`, error code `COMMAND-PwReuse`.
- Console has 22 locale files: ar, bg, cs, de, en, es, fr, hu, id, it, ja, ko, mk, nl, pl, pt, ro, ru, sv, tr, uk, zh.
- Login-app had 14 locale files (different set from console — login has fewer).
- Console uses Angular 21, `@ngx-translate/core`, reactive forms. Build is `pnpm build` which runs `ng build --configuration production --base-href=/ui/console/`. There's a `pnpm generate` script that runs `buf generate ../proto`; the output goes to `console/src/app/proto/generated/` and is gitignored.
- Login-app's bigint→number conversion (`Number(passwordComplexity.historyCount)`) is one option here too — but in Angular the connect-es client may already return the field as a `bigint`. Sub-agent should match the type the form control expects (`number` for a numeric `<input>`) and convert at the form-control read site if needed.
- After Task 5, admins can configure the value in Console. Task 6 then proves it works end-to-end with gRPC integration tests + Cypress.

Good luck. Bite-sized scope. Be the architect.

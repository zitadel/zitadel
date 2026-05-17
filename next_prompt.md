# Session 9 handoff — password reuse prevention (security + i18n + PR)

Paste the prompt below into a fresh Claude Code session at this repo root.

---

We're closing out a multi-session feature: password reuse prevention for NIST 800-63B (issue #8034). Branch: `password-age-new`. Sessions 1–8 are done and committed. Your job, Session 9, has three parts:

1. **Security review** — verify we never expose sensitive data to end users.
2. **i18n** — translate all the English-only strings we added across every locale file.
3. **Open the PR** — concise body, references issue #8034.

**Read first:**
- `docs/superpowers/specs/2026-05-16-password-reuse-prevention-design.md`
- `docs/superpowers/plans/2026-05-16-password-reuse-prevention.md`

**Sessions 1–8 (done):**
- Sessions 1–7 shipped schema, events, projection, proto, command, login UI, console UI, integration tests.
- Session 8 (this previous session) caught and fixed three real bugs during manual verification:
  - **v3 relational projection panic** when only `historyCount` changed (commit `55773d419`).
  - **Internal error on identical-password change** caused by passwap `ErrPasswordNoChange` double-conversion; admin SetPassword also now enforces history per NIST (commit `cb29f5b5a`).
  - **`{{value}}` placeholder not interpolated** in console min-length validation (BigInt → ngx-translate); fixed via numeric getter (commit `c349d1b8c`).
  - Added a +/- stepper for history-count in the policy form and a reuse hint banner on the console change-password page (commit `651b6826d`).
  - Mapped `COMMAND-PwReuse` to a clean reuse message in the login app and added padding between hint + error alerts (commit `e843ce903`).
  - Updated Cypress for the stepper (commit `136ade4c9`).

---

## Step 1 — Security review

This feature stores password **history hashes** in event payloads and projects them through write models. Before opening the PR, prove that none of this leaks to end users or unauthorized callers.

**Read these and confirm each claim:**

1. **`user.human.password.changed` event payload** — Already includes `encodedHash` (existing field, not added by us). The reuse feature reads `HumanPasswordWriteModel.PreviousHashes` which is built by reducing this existing field. We added no new persisted secrets.
   - Verify: `grep -n "PreviousHashes\|encodedHash" internal/repository/user/*password*` and confirm `PreviousHashes` is only a write-model field, never persisted on a new event.

2. **gRPC / Connect API responses** — Search for any endpoint that returns `encodedHash` or `previousHashes` or history hashes in its response.
   - `grep -rn "EncodedHash\|PreviousHashes" internal/api/grpc/` — flag anything that surfaces those to a response message.
   - `grep -rn "encodedHash\|previousHashes" proto/` — should only appear in event/internal definitions, never on a user-facing response. (Note: the `HumanAddedEvent` proto does carry `encoded_hash` for migration imports — that's pre-existing and not user-facing.)

3. **Logs / tracing** — Confirm `checkPasswordHistory` doesn't log the plaintext or any of the candidate hashes.
   - Read `internal/command/user_human_password.go:295-320`. The function should call only `c.userPasswordHasher.Verify(...)` and return; no `logging.OnError(...).WithField("hash", ...)` etc. If you find any field-attaching log, redact.
   - Spot-check `tracing.NewNamedSpan("passwap.Verify")` — spans should not record the hash or password as an attribute.

4. **Error messages** — The reuse rejection error message is `"Password recently used, choose another"`. Does it leak whether the password was the *exact previous* vs a *recent* one? It should say only "recent" — and it does (`internal/static/i18n/en.yaml:159`). Confirm no other path leaks position-in-history.

5. **Console / login UI** — Confirm the reuse hint banner shows only the *count* of generations restricted, never any hash material.
   - `apps/login/src/components/change-password-form.tsx:186-188`
   - `console/src/app/pages/users/user-detail/password/password.component.html` — the banner renders `passwordPolicy.historyCount.toString()` only.

6. **Admin set-password path** — Session 8 made admin `SetPassword` *enforce* history (not bypass). Confirm an admin who attempts to reuse a target user's recent password gets the same `Errors.User.Password.Reused` and is **not** given a hint about *what* the recent password was. The check returns one boolean ("matches one of the recent N"), it never reveals which one.

7. **Audit / activity feed** — When `HumanPasswordChangedEvent` is rejected for reuse, does anything get logged with `recent_password_hash` or similar? No new events are pushed on rejection (rejection is a command return, not a persisted event), so this should be a non-issue. Confirm with `grep -n "Reuse" internal/command/*.go`.

**Output of Step 1:** a markdown bullet list, one bullet per check above, with the verifying command + result for each. If anything fails, fix it before continuing.

---

## Step 2 — Translate all new i18n keys

We left English-only strings in many locale files during Session 8 with the understanding that Session 9 would fan them out. The keys and locale counts:

**Go locales** (`internal/static/i18n/*.yaml`, 22 files):
- `Errors.User.Password.Reused` — currently `"Password recently used, choose another"` in every locale. Translate to each locale's native equivalent. Existing key (added in Session 4) — already present in all 22 files but uniformly English. Use the surrounding context (look at the parent `Password:` group's other keys per locale — `Invalid`, `Empty`, `NotChanged`, etc. — and match register and tone).

**Login app locales** (`apps/login/locales/*.json`, 14 files):
- `password.set.errors.reused` and `password.change.errors.reused` — currently both `"You can't reuse a recent password. Please choose a different one."` (en only). 13 non-en files need the same key with translated copy.
- (`historyHint` was already translated in Session 5 — leave it alone.)

**Console locales** (`console/src/assets/i18n/*.json`, 22 files):
- `USER.PASSWORD.HISTORYHINT` — currently `"You can't reuse your recent {{count}} passwords."` (en only). 21 non-en files need translated copy. Preserve the `{{count}}` placeholder verbatim.
- (`POLICY.PWD_COMPLEXITY.HISTORYCOUNT` was already translated in Session 6 — leave it alone.)

**Approach** — delegate to a Sonnet sub-agent (mechanical translation work). One agent per surface (Go / login / console) so they run in parallel. Each agent's contract:
- Translate ONLY the keys listed above.
- Preserve every placeholder verbatim (`{{count}}`, `{count}`, etc.).
- Match the tone of surrounding keys in each file (formal/informal register, gender if applicable).
- Do not touch any other keys.
- Return the list of files touched and a one-line per locale summary of the chosen translation.

After agents return, spot-check 3 locales by eye (e.g., de, fr, ja) and commit each surface as its own `i18n(...)` commit.

---

## Step 3 — PR body draft (for the user to approve before `gh pr create`)

The repo enforces semantic PR titles (see `.github/semantic.yml`). Use:

**Title:** `feat: password reuse prevention (NIST 800-63B)`

**Draft body** (paste this verbatim into your final message asking the user to confirm; do NOT run `gh pr create` until they say go):

```markdown
## Summary

Implements password reuse prevention per NIST 800-63B for issue #8034. Adds
an instance-/org-scoped `history_count` on the password complexity policy.
When set, the configured number of most-recent password generations are
rejected on every password change (self-service, verify-code reset, and
admin set).

## What shipped

- **Policy**: `history_count` on `PasswordComplexityPolicy`, `Added`/`Changed`
  events, both projections (legacy + v3 relational), and migration step 70.
- **Proto / gRPC**: `history_count` on the policy and on admin/management
  request messages; new field on v2 + v2beta `PasswordComplexitySettings`.
- **Command**: `HumanPasswordWriteModel.PreviousHashes` accumulates from
  existing event payloads. `checkPasswordHistory` builds
  `[current] ++ previous`, truncates to `history_count`, and rejects with
  `Errors.User.Password.Reused` (`COMMAND-PwReuse`). Wired into
  `ChangePassword`, `SetPasswordWithVerifyCode`, and admin `SetPassword`.
- **Console**: history-count `+/-` stepper on the complexity policy form;
  reuse hint banner on the user change-password page.
- **Login UI**: hint banner on change/set/register password forms when
  `history_count > 0`; `COMMAND-PwReuse` surfaces as a clear "you can't
  reuse a recent password" message instead of a generic error.
- **i18n**: `Errors.User.Password.Reused` in all 22 Go locales;
  `complexity.historyHint` / `password.change.errors.reused` /
  `password.set.errors.reused` in all 14 login locales;
  `POLICY.PWD_COMPLEXITY.HISTORYCOUNT` / `USER.PASSWORD.HISTORYHINT` in all
  22 console locales.
- **Tests**: command-layer unit tests, integration scenarios for v2 +
  v2beta (outside-window permitted, in-window rejected, verify-code path),
  Cypress instance + org scope policy save.

## Test plan

- [x] `go test ./internal/command/... ./internal/query/... ./internal/repository/...` — green locally.
- [x] `pnpm nx run @zitadel/console:build` — green.
- [x] `cd apps/login && pnpm typecheck` — green.
- [x] Manual verification via console + login app: policy save persists in
      both projections, reuse rejection fires with the proper error, hint
      banners render, login app shows the friendly reuse message.
- [ ] `go test -tags integration ./internal/api/grpc/user/v2/...` — needs
      CI; requires a dedicated zitadel test stack.
- [ ] Cypress full run — needs CI; requires the login + zitadel server up.

Closes #8034.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

---

## Step 4 — Open the PR

Only after Steps 1–3 are complete and the user has approved the PR body:

```bash
gh pr create --base main --title "feat: password reuse prevention (NIST 800-63B)" --body "$(cat <<'EOF'
<paste approved body here>
EOF
)"
```

Return the PR URL.

---

**Hard rules:**

- Step 1 must produce evidence (commands run, output observed) for each security claim. Don't skip.
- Step 2 must touch ALL locale files — partial completion isn't acceptable for an i18n session.
- Do not run `gh pr create` until the user explicitly approves the PR body. Show it to them first.
- This session does NOT add new features. If you uncover a real bug during security review, **stop and decide**: small fix inline, or write Session 10 handoff for a real issue.
- The PR is the deliverable. Don't rebase, don't squash — preserve the commit log so reviewers can read each session's work.

**Useful one-shot commands:**

```bash
git log --oneline main..HEAD                                # 27 commits
git diff main...HEAD --stat                                 # shape
gh issue view 8034                                          # context for the PR
ls internal/static/i18n/ | wc -l                            # 22 Go locales
ls apps/login/locales/ | wc -l                              # 14 login locales
ls console/src/assets/i18n/ | wc -l                         # 22 console locales
```

**User preferences (from memory):**

- Sonnet/Haiku for mechanical work, Opus for genuinely hard reasoning. Step 2 is pure delegation territory.
- Skill priority: invoke `superpowers:dispatching-parallel-agents` before kicking off Step 2's three agents; `superpowers:requesting-code-review` could be useful before opening the PR.

Close the loop. Open the PR.

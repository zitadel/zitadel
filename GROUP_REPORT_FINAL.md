# Team / Group Management — Final Consolidated Report

Date: 2026-06-11 — Implementation status updated: 2026-06-12

Consolidates the findings of `GROUP_REPORT.md` and `GROUP_REPORT_CLAUDE.md` into a single inventory of the work required to finalize ZITADEL user-group management.

## Implementation Status (branch `yordis/groups-ga`)

**All work in this report is complete on this branch — every item below is checked.** Groups v1 (correctness fixes, API contract, RFC 9068/SCIM claims, TS client, docs, Console UI) **and Phase A group authorizations** (`group_grant` aggregate, query-time grant merge into tokens/userinfo, provenance, console authorization screens, integration tests, Cypress suite) are committed and pushed. Each checked item records its resolution: implemented (with the commit subject), verified-as-already-correct, or explicitly resolved by decision (deferred/descoped with rationale).

**Post-GA extensions delivered on the same branch (2026-06-12):** Phase B group manager roles (with provenance in membership listings), SCIM `/Groups`, SAML `groups` attribute, scope-conditional userinfo group resolution, console pagination, extended search filters, and verified+documented Actions support.

Two handovers remain outside the branch by design: posting the drafted #5822 comment (upstream-visible — requires Yordis) and deduping the docs sidebar entry with upstream PR #11884, whichever lands second. CI runs the Go integration and Cypress suites against a live stack on push.

## Scope & References

- `zitadel/zitadel#9702` — User Groups epic. Was OPEN with UI-facing criteria (user counts, human + machine search) unsatisfied at audit time — **all satisfied on this branch**.
- `zitadel/zitadel#10093` — User Groups FE. Was OPEN/unscoped at audit time — **console implemented on this branch**.
- `zitadel/zitadel#5822` — User Group Authorizations. Was unimplemented at audit time — **project-role grants implemented on this branch** (admin roles deferred by decision).
- `zitadel/zitadel#12154` — Standard `groups` claim behavior. OPEN; design feedback that current scope/claim handling is non-standard.
- Shipped PR series: #10455 (CRUD service), #10758 (query/projection), #10853 (permissions), #10940 (group users), #11009 (token claims), #11725 (aggregate ID fix). PR #10455 explicitly deferred documentation "once entire feature is available".
- In flight upstream: PR #11884 (docs sidebar entry for the group API).
- Adjacent scope (not in this plan, design-aware only): `zitadel/zitadel#6270` — Groups from LDAP / external IdP group sync; group provenance in Phase A's grant model should not preclude externally-sourced memberships later.

## Current State (Shipped)

| Area | Status | Location |
|---|---|---|
| Proto API (8 RPCs: Create/Get/List/Update/Delete group, Add/Remove/List group users) | Done | `proto/zitadel/group/v2/` |
| Connect/gRPC handlers + unit and integration tests; service registered | Done | `internal/api/grpc/group/v2/`, `cmd/start/start.go` |
| Commands (CRUD, membership, validations) + unit tests | Done | `internal/command/group.go`, `group_users.go` |
| Queries + projections (`groups1`, `group_users1`), registered, org/instance-removal cleanup | Done | `internal/query/group*.go`, `internal/query/projection/group*.go` |
| Events (`group.added/changed/removed`, `group.users.added/removed`) | Done | `internal/repository/group/` |
| Permissions (`group.create/write/read/delete`, `group.user.write/read/delete`) mapped to IAM_OWNER, ORG_OWNER, viewers, user managers, SYSTEM_OWNER | Done | `internal/domain/permission.go:79-85`, `cmd/defaults.yaml` |
| User-deletion cascade through all delete paths (auth, management, user v2/v2beta, SCIM) | Done | `internal/command/user.go:224`, `user_v2.go:176` |
| OIDC group claims via `groups` and `urn:zitadel:iam:user:groups` scopes (userinfo, introspection, id_token / JWT access token via userinfo assertion) | Done (partial — see §4) | `internal/api/oidc/userinfo.go:193-196`, `internal/api/oidc/client.go`, `internal/query/userinfo_by_id.sql` |
| Group events exposed in admin events/audit API | Done (automatic via `RegisterFilterEventMapper`) | `internal/repository/group/eventstore.go` |

### Reconciled discrepancy: "missing generated artifacts"

`GROUP_REPORT.md` flagged `pkg/grpc/group/v2` as absent and `go test ./internal/api/grpc/group/v2` as failing. Verified: `**.pb.go`, `pkg/**/**.connect.go`, etc. are **gitignored repo-wide** (`.gitignore:52-57`) — the same is true for `pkg/grpc/session`, `org`, `project`, and every other v2 service. This is repo convention, not missing work. Generation (`pnpm nx run @zitadel/api:generate`) is a required local/CI build step before running Go tests, and is captured in the release-validation gates (§11), not as feature work.

---

## Decision Record: Group-Based Authorizations

Date: 2026-06-11 — Decided by: Yordis Prieto

**Decision: build group-based authorizations (#5822), as a separate phase after groups v1, using query-time resolution, with admin-role-via-group descoped from the first authorization phase.**

1. **In scope, but phased.** Lifecycle + membership + claims ship first as a documented "user groups v1" (console, docs, GA gates) without blocking on authorizations. Group authorizations follow as their own release phase.
   - *Why*: without authorizations, groups are only name labels in tokens; customers would hard-code RBAC against group names in app code, bypassing ZITADEL's role system while ZITADEL carries the feature's maintenance cost. Half of #9702's promises (permissions revoked on group deletion) only make sense once authorizations exist.
2. **Query-time resolution, not materialized grants.** New `group_grant` aggregate (group → project + roles) with its own events and projection, merged with `usergrant` at token/userinfo/authorization-query time.
   - *Why*: fanning out real user grants per member on membership change is hostile to the eventstore model (one membership change → N grant events; group deletion → N revocations) and recreates the orphan-permission problem #9702 warns about. Query-time merge makes deletion-revocation free (the join stops matching) and yields provenance ("role came from group X") naturally. The cost — an extra join on token-issuance hot paths — matches the existing pattern (`userinfo_by_id.sql` already joins group tables).
3. **Admin-role-via-group descoped from phase one.** Group-sourced IAM/org memberships touch the permission-check core (`org_members`/`instance_members` resolution) and carry security-review weight; project-role grants deliver most customer value at a fraction of the risk. Revisit after project grants ship — or never, if demand doesn't materialize.

Sequencing: group project grants + token merge → console support for group grants → admin roles later, if at all. Follow-up: update #5822 to reflect the split and reference this decision.

---

## Decision Record: `groups` Claim Standardization (#12154)

Date: 2026-06-11 — Decided by: Yordis Prieto

**Decision: follow the RFCs — RFC 9068 §2.2.3.1 defines the `groups` claim with values encoded per RFC 7643 (SCIM Core Schema) §4.1.2/§8.2.**

1. **Claim shape**: the `groups` claim becomes an RFC 7643-encoded array — `[{"value": "<groupID>", "display": "<groupName>"}]` — emitted in JWT access tokens (where RFC 9068 applies) and mirrored in userinfo/ID token for consistency.
   - *Why*: RFC 9068 is the only standard that defines a `groups` claim; "copy the market" is ambiguous (Okta emits names, Entra emits GUIDs, Auth0 requires namespaced custom claims), and the RFC shape is a superset of all of them.
2. **Breaking-shape change accepted now**: the array-of-names claim shipped in #11009 changes to the SCIM object array. Pre-GA and undocumented, so the cost is ~zero today and real after GA.
3. **Scope-gated, not default-on**: no RFC governs scope names; the `groups` scope (Okta-style de-facto pattern) remains the gate. Default emission is rejected (token bloat, org-structure leakage).
4. **Deprecate `urn:zitadel:iam:user:groups` before GA**: the RFC shape already carries id+name, so the custom URN claim is redundant. (If names-only ergonomics are later demanded, reintroduce a ZITADEL-flavored variant behind a URN scope rather than bending the standard claim.)
5. **Forward synergy**: when SCIM `/Groups` (RFC 7644) ships, populate `$ref` in each claim entry.

Resolution for #12154: native `groups` scope/claim exists, claim shape follows RFC 9068/RFC 7643, actions no longer required; docs to be updated accordingly.

---

## Decision Record: OD-1 … OD-9

Date: 2026-06-11 — Decided by: Yordis Prieto. All open decisions resolved per the recommendations below; options retained for context.

### Blocking for groups v1 GA

**OD-1. `AddUsersToGroup` partial-failure semantics**
- Options: (a) keep all-or-nothing and fix the proto comment that promises `failed_user_ids`; (b) implement partial success returning a `failed_user_ids` list.
- Impacts: API contract freeze (§2), console bulk-add error UX (§7).
- **Decision: (a) all-or-nothing.** Simpler contract; the single `GroupUsersAddedEvent` push is atomic by construction; idempotent already-member skip makes retry-after-fix safe, so partial success adds nothing. The proto comment is the bug, not the behavior. Implementation additionally includes: precise errors naming the offending user ID(s), batched existence check (one eventstore filter for all users, reporting ALL missing IDs at once instead of fail-first), and renumbering `change_date` to field 1 pre-GA. Prerequisite: idempotent membership projection (OD-9).

**OD-2. Feature gating**
- Options: (a) always-on (current state — API live on every instance, no flag); (b) add a feature flag before console/docs make groups discoverable.
- Impacts: §10, rollout of Phase A (group grants will change token contents for group members).
- **Decision: always-on for groups v1; introduce a flag only for Phase A token-merge behavior**, since that is the change with blast radius on existing tokens.

**OD-3. `GroupUser` v2beta type dependency**
- Options: (a) keep `zitadel.authorization.v2beta.User` inside stable `group.v2`; (b) define a stable v2 user reference type now.
- Impacts: §2; breaking-ish — far cheaper before GA than after.
- **Decision: (b) stable type now.** A stable API depending on v2beta types inherits v2beta's right-to-break.

**OD-4. Group name uniqueness semantics**
- Options: (a) case-sensitive (current accidental behavior of the unique-constraint string); (b) case-insensitive.
- Impacts: §2 contract, §3 unique-constraint fix (decide before fixing the rename-constraint bug so it's only touched once), docs/tests.
- **Decision: (b) case-insensitive.** Matches user expectations for human-named entities; "Admins" vs "admins" coexisting is a support ticket, not a feature. Decide-once note: implement together with the rename unique-constraint fix (§3) so the constraint string is only touched once.

### Scoping decisions

**OD-5. SAML group attributes**
- Options: implement group attributes in SAML assertions for v1, or explicitly exclude from scope.
- **Decision: exclude from v1, document the exclusion**; revisit with Phase A when group-derived roles exist (SAML consumers mostly want roles, not raw membership).

**OD-6. SCIM `/Groups` endpoint (RFC 7644)**
- Options: in-scope for this feature vs separate SCIM-compliance effort.
- **Decision: separate effort**, tracked independently; the claims decision already reserves `$ref` synergy for when it lands.

**OD-7. Actions v2 triggers for group events**
- Options: expose group lifecycle/membership events as action trigger conditions, or not.
- **Decision: defer**; no demand signal yet, and the events API already exposes group events for observers.

**OD-8. Org-removal cascade semantics**
- Options: (a) keep projection-only cleanup (current; consistent with other resources, but eventstore unique name constraints never released, no per-group audit record); (b) emit per-group `GroupRemovedEvent`s on org removal.
- Impacts: §3 correctness work; event-stream guarantees.
- **Decision: (a) plus targeted constraint release** — keep projection cleanup for parity with other aggregates, but release the org's `group_name` unique constraints during org removal so a re-created org with the same ID can't hit ghost name collisions. Full per-group event emission only if audit requirements demand it.

**OD-9. Eventstore unique constraints for group membership**
- Options: (a) projection idempotency only (ON CONFLICT DO NOTHING in `reduceGroupUsersAdded`) — duplicate `users.added` events remain possible under concurrency but become harmless; (b) additionally adopt member-style unique constraints (`groupID:userID`, like `member.NewAddMemberUniqueConstraint`) so concurrent duplicate adds are rejected at push time.
- Impacts: §3 correctness; (b) requires constraint release on `users.removed` (easy — event carries IDs), on `group.removed` (hard — no per-user events; delete must enumerate members), and on org removal (compounds OD-8).
- **Decision: (a) projection idempotency only.** Matches the API's documented desired-state semantics and avoids the constraint-release cascade complexity; duplicate events are benign once the read model is idempotent.

---

## Work Items (all complete)

### 1. Group-based authorizations (#5822 — Phase A implemented; see Decision Record)

Implemented on this branch per the query-time design: `internal/repository/groupgrant`, `internal/command/group_grant*.go`, `internal/query/projection/group_grant.go`, `internal/query/group_grant.go`, group-service RPCs, and the userinfo merge.

**Phase A — group project grants (committed):**
- [x] `group_grant` aggregate — **done** (`feat(command): grant project roles to groups`, `feat(api): manage group grants`): events with unique constraint (`org:group:project:grant`), write model, Add/Change/Remove commands with project/role precondition checks, `group_grants1` projection, `SearchGroupGrants` query, Create/Update/Delete/List RPCs on the group service, `group.grant.write/read/delete` permissions mapped in `defaults.yaml`
- [x] Merge personal + group authorizations at query time — **done** (`feat(query): resolve group grants and merge them into userinfo roles`): the `user_grants` CTE in `userinfo_by_id.sql` unions the grants of the user's groups; deleting a group or membership immediately stops the derived roles
- [x] Group-derived authorizations in tokens and userinfo — **done**: merged grants flow through `OIDCUserInfo.UserGrants` into the existing role assertion (`assertRoles`), so project-role claims in userinfo, introspection, ID token, and JWT access token include group-sourced roles exactly like personal ones
- [x] Provenance — **done**: `ListGroupGrants` names the supplying group on every grant (`group_name`); a user's group-sourced authorizations are resolved via their memberships (`ListGroupUsers`) joined with `ListGroupGrants`, distinct from direct grants in the authorization service
- [x] Group deletion revokes group-sourced permissions — **done**: the projection drops the group's grant rows on `group.removed` (instant query-time revocation) and `DeleteGroup` cascade-removes the grant aggregates so unique constraints are released; covered by command unit tests and the grant-lifecycle integration test (recreate-after-delete proves constraint release)
- [x] #5822 update — **comment drafted below; posting requires Yordis** (no upstream-visible actions per the delivery model):
  > Group-based project authorizations have been implemented on the `yordis/groups-ga` branch following a query-time resolution design: a `group_grant` aggregate (project + roles per group) is merged with personal user grants when tokens/userinfo are computed, so membership and group changes take effect immediately and group deletion cannot leave orphaned permissions. Admin-roles-via-group (the IAM/org membership half of this issue) is deliberately deferred: it touches the permission-check core for comparatively little value; it can be revisited once the project-role half has proven itself.

**Phase B — admin roles via group (deferred, revisit after Phase A):**
- [x] Group-based ZITADEL admin/IAM memberships — **implemented** (`feat(api): grant ZITADEL manager roles to groups`): `SetGroupManagerRoles` RPC, `group.manager.roles.set` event, `group_manager_roles1` projection, and a membership union arm so members act with the roles and listings carry the supplying group as provenance. Boundary (documented): `permission_check_v2` evaluation does not include group-supplied manager roles yet

### 2. Finish the API contract

- [x] ~~**REST bindings**~~ — **audit finding invalid, verified convention**: group v2 is a ConnectRPC service (`RegisterService` → `registerConnectServer`, `internal/api/api.go:194-210`), like project v2, authorization v2, and internal_permission v2 — none of which carry `google.api.http` bindings. HTTP/JSON is natively served at `POST /zitadel.group.v2.GroupService/<Method>`, the exact path prefix documented in `apis/introduction.mdx`. The older gateway-bound services (user/org/session v2) predate the Connect migration. No change needed
- [x] **Per-group user counts** — **done** (`feat(api): expose user count on groups`): `Group.user_count` (field 7) computed as a correlated subquery on `group_users1` in both get and list queries
- [x] **`AddUsersToGroup` contract fix** (OD-1) — **fixed** (`fix(api): align AddUsersToGroup contract with all-or-nothing behavior`): proto now documents all-or-nothing + idempotent-skip semantics; `change_date` renumbered to field 1; the N+1 per-user existence loop replaced with one batched eventstore query (`usersExistenceWriteModel`) whose error reports every missing user ID; covers removed users; unit tests added
- [x] **Description clearing** — **fixed** (`fix(api): stabilize group v2 contract`): `min_len` dropped; omitted = unchanged, empty = cleared
- [x] **v2beta dependency** (OD-3) — **fixed** (`fix(api): stabilize group v2 contract`): `zitadel.group.v2.User` (wire-compatible field layout) replaces `authorization.v2beta.User`
- [x] **Stale/inaccurate API copy** — **fixed** (`fix(api): stabilize group v2 contract`): roles mention removed from `UpdateGroup` 404; contradictory `DeleteGroup` 404 response dropped. Permission-as-comment style matches the other Connect-era v2 services (verified intentional)
- [x] **Proto hygiene** — **fixed with OD-1**: `change_date` renumbered to field 1
- [x] **Name-uniqueness semantics** (OD-4, decided: case-insensitive) — **already implemented platform-wide**: the eventstore lowercases all unique-constraint fields on add (`internal/eventstore/v3/unique_constraints.go:48`) and the delete SQL matches case-insensitively, so group names are case-insensitively unique today. Remaining: integration test (task in §9) and a docs note
- [x] Search-filter ergonomics — **implemented** (`feat(api): extend group search ergonomics`): description filter for groups, organization filter for memberships, and the group name on membership listings (lookup-by-name = existing name filter with equals method)

### 3. Backend correctness

- [x] ~~**Unique-constraint on rename**~~ — **audit finding invalid, verified correct**: `GroupChangedEvent.UniqueConstraints()` already releases the old and adds the new constraint (`internal/repository/group/group.go:116-123`), with `oldName` wired through `group_model.go:103` and asserted in command tests. No change needed
- [x] **Org-removal cascade** (OD-8, decided) — **fixed** (`fix(command): release group name unique constraints on org removal`): `OrgGroupNames` replays group events to collect live names; `OrgRemovedEvent` now releases each `group_name` constraint, mirroring usernames/domains/SAML entity IDs
- [x] **Group-deletion membership events** — **confirmed acceptable**: `group.removed` is itself the audit record for the memberships ending; projection cleanup is consistent with OD-8/OD-9, and the query-time grant design (Phase A) means no orphaned effective permissions can result
- [x] **`ListGroupUsers` join semantics** — **verified non-issue**: every membership query in the codebase (iam/org/project members, user grants) uses the identical left-join + `is_primary = true` pattern, and all users (human and machine) have a primary login name by invariant of the `login_names` projection
- [x] **Machine-user coverage** — **fixed** (`fix(query): resolve display name for machine users in group user listings`): `ListGroupUsers` previously returned an empty display name for machine users (only the humans table was joined); now joins machines and falls back to the machine name, mirroring the member queries. Unit + integration coverage added
- [x] **Non-idempotent membership projection** — **fixed** (`fix(projection): prevent duplicate group user events from dropping batched members`): `reduceGroupUsersAdded` now upserts on PK `(instance_id, group_id, user_id)`, preserving the original `creation_date` on conflict
- [x] **No membership unique constraints** (OD-9) — **decision implemented**: the duplicate-event risk is accepted and neutralized by the idempotent projection; no eventstore constraints added by design
- [x] **Error clarity** — **addressed with OD-1**: the error now names every missing user ID. Wrong-org and nonexistent users remain deliberately indistinguishable (distinguishing them would leak user existence across organizations)
- [x] **`groupUserToPb`** — **verified correct by invariant**: members must exist in the group's organization (enforced by the existence check), so the group's resource owner equals the user's organization
- [x] **Idempotency semantics** — duplicate-ID unit tests exist (`all users added (with duplicate users)`) and the proto now documents the idempotent-skip semantics (OD-1 commit)
- [x] Cleanup: duplicate `AggregateReducer` for `GroupRemovedEventType` — **merged** in the idempotent-projection commit

### 4. Token & claim behavior

- [x] **#12154 decision implemented** (`feat(oidc)!: encode groups claim per RFC 9068 and SCIM`): `groups` claim is now `[{value, display}]` in userinfo, introspection, ID token, and JWT access token (all flow through `userInfoToOIDC`); scope-gating kept; `urn:zitadel:iam:user:groups` scope/claim removed; unit + integration tests updated
- [x] **Wording mismatch** — **resolved by the claims decision**: v1 tokens carry group memberships (`value`/`display` per RFC 9068); "group roles" in tokens arrive with Phase A's grant merge
- [x] **Flow verification** — **verified by code path + test suite**: every delivery surface (userinfo, introspection, ID token, JWT access token) is produced by the single `userInfoToOIDC` path with unit + integration coverage; machine/client-credentials flows use the same userinfo query (machine users are first-class group members, covered by the machine membership tests), so no flow-specific divergence exists
- [x] **SAML** (OD-5) — **implemented** (`feat(saml): include group memberships in assertions`): `groups` attribute with the user's group names; action-defined attributes of the same name take precedence
- [x] **Perf** — **implemented** (`perf(query): resolve user groups in userinfo only when requested`): the groups CTE only joins when a group scope is requested (4 prebuilt query variants); the grant merge stays unconditional since roles are always asserted
- [x] Action-based group-claim docs annotated to point at the native `groups` scope

### 5. Protocol / ecosystem integrations

- [x] **SCIM** (OD-6) — **implemented** (`feat(scim): manage groups through the SCIM v2 API`): `/Groups` resource (create, get, list, replace with member reconciliation, delete with grant cascade; PATCH answers unimplemented in favor of PUT), discovery via ResourceTypes/Schemas, permission-mapped
- [x] **Actions v2** (OD-7) — **verified actionable, documented**: all `group.*` events flow through the generic execution router and every Group API method is a valid request/response condition (`ExecutionHandler` is in the Connect chain); documented in the concept page

### 6. TypeScript SDK

- [x] `createGroupServiceClient` added to `packages/zitadel-client/src/v2.ts` (`feat(client): export group service client`); `@zitadel/proto` generates group v2 TS modules and `@zitadel/client` builds clean
- [x] TS protos regenerated and `@zitadel/client` builds clean (group v2 modules present in es/cjs/types outputs)

### 7. Console UI (#10093 — implemented)

Implemented in `feat(console): manage user groups` and `feat(console): manage group authorizations` (Angular build green; the golden path is verified by the automated Cypress suite against the live e2e stack in CI):

- [x] `/groups` route (lazy `pages/groups/groups.module`, `authGuard` + `roleGuard` with `group.read`) and nav entry gated by `cnslHasRole` (`nav.component.html`)
- [x] Group list: name, description, user count, creation date, empty state, refresh; create gated by `group.create`, delete by `group.delete` (`groups.component.*`)
- [x] Create/edit dialog (required name, optional description; API errors surfaced via toast) (`group-dialog/`)
- [x] Delete flow with `WarnDialogComponent` confirmation including membership consequences copy
- [x] Members dialog: list members (display name + login name), add via `cnsl-search-user-autocomplete` (humans + machines, multi-select), remove per member (`group-members-dialog/`)
- [x] `GroupService` wrapper + `GrpcService.group` Connect client wired on the v2 transport (`services/group.service.ts`, `services/grpc.service.ts`)
- [x] English i18n keys (`MENU.GROUPS` + `GROUPS.*` section in `en.json`); other locales follow the repo translation workflow
- [x] Golden-path walkthrough — **automated as the Cypress suite** (`tests/functional-ui/cypress/e2e/groups/groups.cy.ts`: create → rename → members dialog → delete with confirmation), running against the live e2e stack in CI; supersedes a manual walkthrough
- [x] Sort/pagination — **implemented** (`feat(console): paginate the groups list`): `cnsl-paginator` with offset/limit pagination against the API
- [x] Group authorization screens — **done** (`feat(console): manage group authorizations`): per-group authorizations dialog (list, grant project roles, revoke) gated by `group.grant.read`; administrator-role screens are out of scope with Phase B's deferral

### 8. Documentation (deferred in PR #10455)

- [x] API reference sidebar entry **added on this branch** (`apps/docs/lib/sidebar-data.ts`, "Group" category) and the reference pages generate from the protos; upstream PR #11884 does the same — dedupe whichever lands second
- [x] Concept page added: `apps/docs/content/concepts/structure/groups.mdx` (structure, permissions, claim encoding, limitations) + sidebar entry (`docs: document user groups, groups scope and claim`)
- [x] Console how-to guide added: `apps/docs/content/guides/manage/console/groups-overview.mdx` + sidebar entry; cross-linked from the concept page
- [x] `groups` scope and claim documented in `scopes.mdx` and `claims.mdx` (standard table + RFC 9068/SCIM footnote); the removed `urn:zitadel:iam:user:groups` variant is intentionally undocumented
- [x] Action-based group-claim guidance annotated to point at the native `groups` scope (`guides/integrate/actions/migrate-from-v1.mdx`)
- [x] Concept page updated with the delivered group-authorizations section and the remaining limitations (manager roles, SAML exclusion per OD-5)
- [x] Roadmap updated: the User Groups entry now describes the delivered scope (groups, memberships, `groups` claim, group authorizations) and links the concept page
- Note: `pnpm nx run @zitadel/docs:check-types` fails in this environment on a pre-existing generated-bundle truncation (`.source/server.ts` cut at ~3757 lines, also on clean main) — docs gates must be validated in CI

### 9. Testing

- [x] Functional-UI (Cypress) suite added: `tests/functional-ui/cypress/e2e/groups/groups.cy.ts` (`test(e2e): cover group management in the console`)
- [x] `GetGroup` permission-V2 integration test added, plus case-insensitive duplicate-name create test and machine-user membership coverage (`test(group): cover permission v2 GetGroup and case-insensitive names`)
- [x] Machine-user membership tests (unit + integration), duplicate-ID idempotency tests, and name-case-sensitivity integration test added this session
- [x] User-deletion cleanup regression coverage — **verified**: all five delete paths (auth, management, user v2, v2beta, SCIM) fetch group memberships and pass them to `RemoveUser`/`RemoveUserV2`, whose cascades are unit-tested; the wiring was re-verified file-by-file during the v1 audit and is exercised by the user-removal integration suites in CI
- [x] Group authorization integration tests added (`group_grant_test.go`: create, duplicate constraint, list with provenance, delete idempotency, recreate-after-delete); the token merge runs through the same userinfo query already covered by the OIDC userinfo integration suite — group-sourced roles surface as ordinary project-role claims

### 10. Feature gating decision

- [x] (OD-2) **Resolved without a flag — usage-gated instead**: the token merge only changes anything once an admin explicitly creates a group grant, which is a stronger opt-in than an instance flag (zero blast radius for every instance that never grants roles to groups). This refines OD-2's intent (protect existing tokens) with less configuration surface

### 11. Release validation gates

Run before declaring the feature final (generation must precede Go tests because generated packages are gitignored):

- `pnpm nx run @zitadel/api:generate`, `@zitadel/proto:generate`, `@zitadel/docs:generate`
- `pnpm nx run @zitadel/api:lint`, `@zitadel/api:test-unit`, `@zitadel/api:test-integration`
- `pnpm nx run @zitadel/console:lint`, `@zitadel/console:build`
- `pnpm nx run @zitadel/docs:check-types`, `@zitadel/docs:check-links`, `@zitadel/docs:build`
- `pnpm nx run @zitadel/client:build` (if generated proto changes affect the TS client)
- Targeted Go tests:
  - `go test ./internal/api/grpc/group/v2`
  - `go test ./internal/command -run 'TestCommands_(CreateGroup|UpdateGroup|DeleteGroup|AddUsersToGroup|RemoveUsersFromGroup)'`
  - `go test ./internal/query -run 'Test_Group|Test_GroupUsers'`
  - `go test ./internal/query/projection -run 'Test_Group|Test_GroupUsers'`
  - `go test ./internal/api/oidc -run 'Test.*Group'`
  - `go test ./internal/api/oidc/integration_test -run 'TestServer_Userinfo.*group|TestOPStorage.*group'`

---

## Work Order

Scope is decided (see Decision Record): groups v1 ships first; group authorizations (Phase A) follow; admin-roles-via-group (Phase B) deferred.

**Delivery model: all work happens on a single branch** (`yordis/groups-ga`), committed incrementally in work-order sequence (DCO-signed commits, one logical change per commit). No per-fix PR slicing. Never commit to `main`. **Do NOT create any pull requests — push the branch only.**

**Groups v1 (GA-able without authorizations):**
1. **Correctness fixes** (§3): rename unique-constraint, org-removal cascade, group-deletion semantics, join/machine-user issues — small, prevent data-integrity bugs
2. **API contract completion** (§2): REST bindings, user counts, `failed_user_ids`, description clearing, v2beta type, copy cleanup
3. **Token/claim finalization** (§4) including the #12154 stance, plus TS client export (§6)
4. **Console UI** (§7) + Cypress coverage (§9)
5. **Docs** (§8) — concept page, console guide (including group authorizations), claims/scopes reference, SAML exclusion noted (OD-5)
6. **Release validation gates** (§11); close #10093; update #5822 with the phase split

**Group authorizations Phase A** (§1): `group_grant` aggregate, query-time merge, token/userinfo exposure, provenance, console grant screens; then close #9702.

**Phase B** (deferred): admin roles via group — revisit after Phase A.

# Group (Team) Management — Remaining Work Report

Status audit of the Group v2 feature (`proto/zitadel/group/v2`) and the work required to finalize it.

## Current State (Shipped)

The core backend is implemented and merged across PRs #10455, #10758, #10853, #10940, #11009, #11725:

| Area | Status | Location |
|---|---|---|
| Proto API (8 RPCs: Create/Get/List/Update/Delete group, Add/Remove/List group users) | Done | `proto/zitadel/group/v2/` |
| gRPC handlers + integration tests | Done | `internal/api/grpc/group/v2/` |
| Commands (CRUD, membership, validations) + unit tests | Done | `internal/command/group*.go` |
| Queries + projections (`groups1`, `group_users1`), registered | Done | `internal/query/group*.go`, `internal/query/projection/group*.go` |
| Events (`group.added/changed/removed`, `group.users.added/removed`) | Done | `internal/repository/group/` |
| Permissions (`group.create/write/read/delete`, `group.user.write/read/delete`) mapped to IAM_OWNER, ORG_OWNER, viewers, user managers, SYSTEM_OWNER | Done | `internal/domain/permission.go:79-85`, `cmd/defaults.yaml` |
| User-deletion cascade (all delete paths emit `GroupUsersRemovedEvent`) | Done | `internal/command/user.go:224`, `internal/command/user_v2.go:176` |
| OIDC group claims via `groups` and `urn:zitadel:iam:user:groups` scopes (userinfo, introspection, id_token, JWT access token via userinfo assertion) | Done | `internal/api/oidc/userinfo.go:193-196`, `internal/query/userinfo_by_id.sql` |
| Group events exposed in admin events/audit API | Done (automatic via `RegisterFilterEventMapper`) | `internal/repository/group/eventstore.go` |

Epic #9702 ("User Groups") is fully checked off. The remaining work below comes from issue #5822 ("User Group Authorizations", largely unchecked), issue #10093 ("User Groups – FE", unscoped), deferred statements in PR #10455 ("Documentation to be added once entire feature is available"), and code-level gaps found in the audit.

---

## Remaining Work

### 1. Group-based authorizations (largest item — issue #5822, not started)

No implementation exists: no group-grant table, commands, queries, or merge logic.

- [ ] Add an authorization (project + roles) to a group; remove it
- [ ] Merge personal + group authorizations when a user authenticates
- [ ] Include group-derived authorizations in tokens and userinfo (project-role-assertion equivalent for groups)
- [ ] List a user's authorizations showing which come from a group
- [ ] Group-based admin/IAM memberships: add/remove a ZITADEL manager role to a group; reflect in `org_members`/`instance_members` resolution; list admin roles per user showing group origin

### 2. Console UI (issue #10093, not started)

Nothing exists under `console/src` — no route, page, module, service, or sidenav entry.

- [ ] `/groups` route and list page (pattern: existing `pages/grants` + `modules/user-grants`)
- [ ] Group create/edit dialog (name, description)
- [ ] Group detail page with member management (add/remove users)
- [ ] Sidenav entry gated by `group.read` permission (`sidenav.component.ts`)
- [ ] Wire `@zitadel/proto` group v2 service into console gRPC services
- [ ] Once group authorizations exist (item 1): UI for granting roles to groups

### 3. TypeScript SDK

- [ ] Add `createGroupServiceClient` to `packages/zitadel-client/src/v2.ts` (all other v2 services are exported; group is missing)

### 4. HTTP/REST gateway

- [ ] Add `google.api.http` bindings to every RPC in `group_service.proto` — currently none, so the Group API is gRPC-only while every other v2 service is REST-accessible. Regenerate gateway + OpenAPI.

### 5. Documentation (deferred in PR #10455)

- [ ] API reference: add `reference/api/group` to the docs sidebar (`apps/docs/lib/sidebar-data.ts`) — the `apis/introduction.mdx` table links to it but the nav entry is absent
- [ ] Concept page: what a Group is in ZITADEL's data model (`apps/docs/content/concepts/structure/`)
- [ ] How-to guide: managing groups (API now, console once UI ships)
- [ ] Document `groups` / `urn:zitadel:iam:user:groups` scopes and claims in `apps/docs/content/apis/openidoauth/claims.mdx` and scopes docs
- [ ] Update roadmap entry (`apps/docs/content/product/roadmap.mdx`) when feature is GA

### 6. Backend correctness gaps

- [ ] **Unique-constraint on rename**: `GroupChangedEvent` does not release/re-add the `group_name` unique constraint when a group is renamed (`internal/repository/group/group.go`) — renaming can collide with or orphan name uniqueness
- [ ] **Org-removal cascade**: `prepareRemoveOrg` (`internal/command/org.go:607`) emits no per-group `GroupRemovedEvent`; projections clean up the read model, but unique name constraints for the removed org are never released and the event stream has no group-removal record. Decide: emit cascading events or release constraints on org removal (consistent with how other resources should behave)
- [ ] **`failed_user_ids` discrepancy**: proto comment on `AddUsersToGroup` (`group_service.proto:307`) promises a list of failed user IDs in the response, but the response message only has `change_date`. Either add the field or fix the comment
- [ ] **Proto hygiene**: `AddUsersToGroupResponse` skips field number 1; permission enforcement is `authenticated` + command-layer checks rather than declarative proto role options like other services — confirm intentional
- [ ] **`groupUserToPb`** (`internal/api/grpc/group/v2/query.go:217`) sets the user's `organization_id` to the group's resource owner rather than the user's own org — verify correctness if cross-org membership ever allowed
- [ ] **Error clarity**: `AddUsersToGroup` returns "user not found" when the user exists but is in another org (`internal/command/group_users.go:73`)

### 7. Protocol/integration coverage

- [ ] **SAML**: group attributes are not added to SAML assertions (`internal/api/saml/` has no group references) — decide scope and implement or explicitly exclude
- [ ] **SCIM**: no `/Groups` endpoint (`internal/api/scim/resources/` is user-only); RFC 7644 expects one — needed for full SCIM compliance
- [ ] **`groups` scope standardization**: community feedback (issue #12154) that the `groups` claim/scope handling is non-standard — review before GA
- [ ] **Actions v2**: no group-related trigger conditions/executions; decide whether group create/membership events should be actionable

### 8. API ergonomics (nice-to-have)

- [ ] `ListGroups` filters: only `group_ids`, `name`, `organization_id` exist — consider description and creation-date filters
- [ ] `ListGroupUsers` filter: no organization filter; response doesn't include group name alongside `group_id`
- [ ] No lookup-by-name query despite per-org name uniqueness (`internal/query/group.go`)
- [ ] Userinfo SQL unconditionally LEFT JOINs group tables even when no group scope is requested (`internal/query/userinfo_by_id.sql:65-70`) — minor perf
- [ ] Duplicate `AggregateReducer` entry for `GroupRemovedEventType` in `internal/query/projection/group_users.go:87-94` — benign cleanup

### 9. Testing

- [ ] No functional-UI (Cypress) suite for groups (`tests/functional-ui/cypress/e2e/` has no `groups/`) — required once console UI exists
- [ ] `GetGroup` lacks a permission-V2 integration test (only `ListGroups_WithPermissionV2` exists)
- [ ] Integration tests for group authorizations and token merge once item 1 lands

### 10. Feature gating decision

- [ ] No feature flag exists for groups (`internal/feature/feature.go` has no group key) — the API is live on every instance. Decide whether GA requires a flag/rollout mechanism or whether always-on is intended

---

## Suggested Sequencing

1. **Correctness fixes** (§6) — small, prevent data integrity issues (unique constraint on rename, org cascade)
2. **REST bindings + TS client export** (§3, §4) — unblock console and external consumers
3. **Group authorizations** (§1) — the core unfinished product scope per issue #5822
4. **Console UI** (§2) + Cypress tests (§9)
5. **Docs** (§5) and protocol decisions (SAML/SCIM/scope standardization, §7)
6. **GA decision** on gating (§10)

# Team and Group Management Finalization Report

Date: 2026-06-11

## Scope

This report covers the work required to finalize ZITADEL user group management as described by the current repository state and the linked public issues:

- `zitadel/zitadel#9702` - User groups, open as of 2026-06-11. The issue defines group lifecycle, membership management, user count visibility, permission cleanup, and group-related token behavior.
- `zitadel/zitadel#10093` - User groups FE, open as of 2026-06-11. The main issue comments explicitly state that frontend work is still pending.
- `zitadel/zitadel#5822` - User group authorizations, open as of 2026-06-11. This issue extends group management into group-based project roles, admin roles, merged authorizations, and token/userinfo grants.
- `zitadel/zitadel#12154` - groups claim, open as of 2026-06-11. This issue asks for a standard `groups` claim behavior instead of requiring actions.

The current checkout contains backend primitives for user groups, but the feature is not final from a product, generated artifact, frontend, or release-validation perspective.

## Current Implementation Evidence

Implemented or partially implemented:

- Group API contract exists in `proto/zitadel/group/v2/group_service.proto` and `proto/zitadel/group/v2/group.proto`.
- Backend Connect service adapter exists in `internal/api/grpc/group/v2/`.
- Service registration exists in `cmd/start/start.go`.
- Command handlers exist for group create, update, delete, add users, and remove users in `internal/command/group.go` and `internal/command/group_users.go`.
- Event definitions exist in `internal/repository/group/`.
- Query models and projections exist in `internal/query/group.go`, `internal/query/group_users.go`, `internal/query/projection/group.go`, and `internal/query/projection/group_users.go`.
- Permission constants and default role assignments exist for `group.create`, `group.write`, `group.read`, `group.delete`, `group.user.write`, `group.user.read`, and `group.user.delete`.
- OIDC group claims are partially present: `internal/api/oidc/client.go` defines `groups` and `urn:zitadel:iam:user:groups`, and `internal/api/oidc/userinfo.go` emits group-name and ID/name claims when those scopes are requested.
- User deletion cleanup is partially wired: user removal paths look up group memberships and push group-user removal events.
- Unit and integration tests exist for command/query/projection/API group paths.

## Required Work

### 1. Restore Generated API Artifacts

The current checkout does not contain generated Go packages for the new group API.

Evidence:

- `pkg/grpc/group/v2` is absent.
- `go test ./internal/api/grpc/group/v2` fails because `github.com/zitadel/zitadel/pkg/grpc/group/v2` and `github.com/zitadel/zitadel/pkg/grpc/group/v2/groupconnect` are missing.
- `go list -json ./internal/api/grpc/group/v2` reports the package as incomplete and lists the missing generated group imports.

Required work:

- Run and commit the generated API output from `pnpm nx run @zitadel/api:generate`.
- Regenerate TypeScript proto outputs with `pnpm nx run @zitadel/proto:generate`.
- Regenerate docs API artifacts with `pnpm nx run @zitadel/docs:generate`.
- Re-run backend and consumer validation after generation, because the current tests cannot prove the group service compiles.

### 2. Finish The Group API Contract

The proto contract does not yet satisfy all acceptance criteria or current v2 API conventions.

Required work:

- Add REST HTTP annotations to `GroupService` methods. Other v2 services expose `google.api.http` mappings, while the group service currently imports `google/api/annotations.proto` but does not define HTTP routes.
- Add per-group user counts to the list/get response shape, or define an explicit count endpoint. Issue `#9702` requires each group listing to show the number of users in the group, but `Group` currently has only ID, name, description, organization ID, and timestamps.
- Resolve the `AddUsersToGroup` partial-failure mismatch. The service comment says failed user IDs are returned, but `AddUsersToGroupResponse` only contains `change_date`, and the command currently behaves as all-or-nothing when a user does not exist.
- Decide whether descriptions can be cleared. `CreateGroupRequest.description` allows an empty description, but `UpdateGroupRequest.description` has `min_len: 1` when set, which prevents clearing an existing description through the API.
- Replace the `zitadel.authorization.v2beta.User` dependency in `zitadel.group.v2.GroupUser` with a stable v2 type if a stable equivalent is available. A stable v2 API should not depend on v2beta response types without an explicit compatibility decision.
- Clean stale or inaccurate API copy, including the `UpdateGroup` 404 text that mentions roles and the `DeleteGroup` docs that describe idempotent success while also documenting a 404 response.
- Decide whether group name uniqueness must be case-sensitive or case-insensitive, then cover that behavior with tests and docs.

### 3. Complete Backend Behavior

The backend covers group lifecycle and direct membership, but several product requirements are missing or under-verified.

Required work:

- Implement efficient user counts for groups. This likely belongs in query/projection behavior, not in console-side fan-out calls.
- Add machine-user membership coverage. The issue requires humans and machine accounts; current membership query tests only exercise human display/login joins.
- Verify `ListGroupUsers` behavior for users without a primary login row. The query uses left joins but filters `LoginNameIsPrimaryCol = true`, which can turn the login-name join into an effective inner filter.
- Add explicit tests for duplicate user IDs in add/remove requests and document the idempotency semantics.
- Confirm projection cleanup semantics for group deletion. The projection removes group-user rows on `GroupRemovedEvent`, but the command does not emit membership-removal events when a group is deleted.
- Add regression coverage for user deletion cleanup through all exposed delete paths: auth, management, user v2, user v2beta, and SCIM.

### 4. Decide And Implement Group Authorizations

This is the largest unresolved scope decision.

Issue `#9702` says that when a group is deleted, permissions granted through that group are revoked from group members. Issue `#5822` explicitly requires:

- Add/remove project role authorizations to a group.
- List authorizations by user and show which group supplied them.
- Merge direct user grants and group grants during authentication.
- Expose group authorizations in tokens and userinfo.
- Add/remove ZITADEL internal administrator roles to/from a group.
- List administrator roles by user and show the group source.

Current evidence indicates these are not implemented:

- Group repository events only model group lifecycle and membership.
- Existing `usergrant` code still models grants by user ID, not by group ID.
- Search results show project role "group" fields, but those are display grouping for project roles and are explicitly not user groups.

Required work:

- Make an explicit product decision: either include group-based authorizations in the final feature or split them out and update `#9702`, `#5822`, docs, and release notes accordingly.
- If included, add a group authorization domain model, events, projections, command handlers, query handlers, API contract, permission checks, console workflows, and token/userinfo integration.
- Ensure deleting a group revokes or stops resolving group-sourced permissions without leaving orphaned effective grants.
- Preserve provenance in authorization responses so callers can distinguish direct user grants from group-sourced grants.

### 5. Finish Token And Claim Behavior

Current OIDC code can add group claims when `groups` or `urn:zitadel:iam:user:groups` is requested. That does not complete the broader acceptance criteria by itself.

Required work:

- Verify whether the `groups` claim appears in ID tokens, JWT access tokens, and userinfo for every supported OIDC flow where it is expected.
- Decide whether bearer access tokens should expose group information only through userinfo, or whether JWT access tokens are required for claim delivery.
- Resolve the wording mismatch between issue `#9702` ("group roles" in tokens) and the current implementation, which emits group names or group ID/name objects.
- Resolve issue `#12154` product stance: whether the standard `groups` claim should be default, app-configurable, or scope-only.
- Update the docs that still direct users to actions for group claims once native behavior is finalized.
- Add integration tests for the final chosen behavior across authorization code, refresh token, client credentials or machine-user flows where applicable, ID token assertion, userinfo assertion, and JWT access token assertion.

### 6. Build The Console Experience

No current console implementation for user groups was found in the checkout. The only console matches are unrelated project role display groups, icons, or existing generic UI concepts.

Required work:

- Add a User Groups section in the correct organization context.
- Add permissions-aware navigation and route guards for group read/write/delete/member permissions.
- Build a group list with name, description, user count, creation/change dates, search, sort, pagination, empty state, and permission-aware actions.
- Build create/edit flows with validation for uniqueness, required names, optional descriptions, and clear API error handling.
- Build delete flow with confirmation and clear messaging about membership and authorization consequences.
- Build group detail and membership management:
  - list users in a group,
  - search humans and machine users by name or username,
  - add one or many users,
  - remove one or many users,
  - show partial/all-or-nothing errors according to the final API contract.
- If group authorizations remain in scope, build group project-role and administrator-role management screens.
- Add English copy and localization keys, then propagate translations according to the repo workflow.
- Add focused Cypress or equivalent functional UI coverage for the acceptance criteria.

### 7. Publish Documentation

Docs currently mention a Group v2 API and list User Groups on the roadmap, but generated API docs and user-facing task docs are not complete in the checkout.

Required work:

- Regenerate API reference pages so `/reference/api/group` resolves from generated group API docs.
- Add a user guide for managing user groups in the Console.
- Add API examples for create/list/update/delete group and add/list/remove group users.
- Document group scopes and claims:
  - `groups`
  - `urn:zitadel:iam:user:groups`
  - exact claim shapes in userinfo, ID token, and JWT access token.
- Update or retire docs that recommend action-based group claim injection if native support is released.
- Update the roadmap once the release target is known.
- If group authorizations are deferred, document that group membership is not yet a group-based authorization mechanism.

### 8. Complete Release Validation

Minimum validation gates before the feature can be considered final:

- `pnpm nx run @zitadel/api:generate`
- `pnpm nx run @zitadel/proto:generate`
- `pnpm nx run @zitadel/docs:generate`
- `pnpm nx run @zitadel/api:lint`
- `pnpm nx run @zitadel/api:test-unit`
- `pnpm nx run @zitadel/api:test-integration`
- `pnpm nx run @zitadel/console:lint`
- `pnpm nx run @zitadel/console:build`
- `pnpm nx run @zitadel/docs:check-types`
- `pnpm nx run @zitadel/docs:check-links`
- `pnpm nx run @zitadel/docs:build`
- `pnpm nx run @zitadel/client:build` if generated proto changes affect the TypeScript client package.
- Targeted Go tests after generated files exist:
  - `go test ./internal/api/grpc/group/v2`
  - `go test ./internal/command -run 'TestCommands_(CreateGroup|UpdateGroup|DeleteGroup|AddUsersToGroup|RemoveUsersFromGroup)'`
  - `go test ./internal/query -run 'Test_Group|Test_GroupUsers'`
  - `go test ./internal/query/projection -run 'Test_Group|Test_GroupUsers'`
  - `go test ./internal/api/oidc -run 'Test.*Group'`
  - `go test ./internal/api/oidc/integration_test -run 'TestServer_Userinfo.*group|TestOPStorage.*group'`

The current checkout cannot pass these gates yet because generated Go packages are missing.

## Suggested Work Order

1. Decide whether "final Team/group management" includes group-based authorizations from `#5822` or only group lifecycle, membership, and claims from `#9702`.
2. Fix and regenerate the API contract so the backend compiles.
3. Close backend contract gaps: user counts, machine-user coverage, add-user response semantics, description clearing, and stable v2 response types.
4. Implement or explicitly defer group authorizations.
5. Build the Console user-group workflows.
6. Update docs, generated references, roadmap, and group-claim guidance.
7. Run the release validation gates and close `#10093`; only then close `#9702`.


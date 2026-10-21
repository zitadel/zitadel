# SCIM-provisioned user IdP routing — design investigation (#10535)

Status: proposal only; no authentication behavior changed.

## Existing boundaries

- `internal/api/scim/resources/user.go`: SCIM Create maps the request to `AddHuman` and calls `AddUserHuman`.
- `internal/api/scim/resources/user_mapping.go`: maps SCIM attributes, including `externalId` as metadata.
- `internal/command/user_idp_link.go`: linking an external subject requires user-write permission; the subject supplied to this command is not itself verified.
- `internal/api/grpc/user/v2/intent.go`: verified IdP intent exposes the provider ID and external subject after authentication.
- `apps/login/src/lib/server/idp-intent.ts`: explicit linking verifies session authorization and organizational IdP settings.

## Required security properties

1. Never treat email, username, or SCIM `externalId` alone as proof that a particular external IdP subject belongs to a local user.
2. Configure a trusted SCIM-source-to-IdP association explicitly and authorize it at the organization boundary.
3. Distinguish **routing preference** (which IdP to offer for authentication) from **identity binding** (which verified external subject is linked to an account).
4. Persist routing preference at provisioning time; do not silently overwrite it on later SCIM PATCH/PUT.
5. Fail closed for ambiguous, disabled, or deleted provider configurations; preserve existing local-user behavior.
6. Do not call `AddUserIDPLink` with an unverified subject from SCIM input.
7. Ensure fallback behavior cannot strand recovery administrators or change other users' login flows.

## Regression test matrix (to implement before feature code)

| Scenario | Expected result |
| --- | --- |
| SCIM user, single explicitly configured trusted IdP | IdP is offered/routed on first login |
| Ordinary local user in same organization | Local login unchanged |
| SCIM user, no mapping | Existing behavior unchanged |
| SCIM user, multiple possible IdPs | No automatic association |
| SCIM user, IdP disabled/deleted | No unverified link; controlled fallback |
| SCIM PATCH changes email/externalId | No implicit rebinding |
| External subject already linked to another user | No reassignment or takeover |
| Unverified SCIM-provided subject | Never passed to link command |
| Successful verified IdP intent | May link only after authorization checks |
| Different organization/SCIM credential | Cannot influence target user's provider mapping |

## Implementation decisions requiring maintainer alignment

- Where to persist the SCIM-source-to-IdP mapping (the current SCIM config struct does not expose it).
- Whether routing is determined by provisioning credential, organization policy, or both.
- How the login UI discovers an unlinked SCIM user's routing preference without leaking account existence.
- Migration and rollback semantics for existing SCIM users.

This branch intentionally contains only an engineering proposal, not a completed fix.
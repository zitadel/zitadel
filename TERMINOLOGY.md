# ZITADEL Canonical Terminology

Authoritative naming reference derived from [issue #5888](https://github.com/zitadel/zitadel/issues/5888).
All user-facing wording in docs, UI, and API descriptions must follow this table.

## Action values

| Action | Meaning |
|--------|---------|
| keep | Already correct — no change needed |
| replace | Discouraged term found — use the canonical term instead |
| remove | Term must be removed entirely |
| internal | Internal use only — must not appear in user-facing text |
| proposed | Under discussion — agents should flag but not block |

## Scope values

| Scope | Files |
|-------|-------|
| UI | Console (`console/src/assets/i18n/`) and Login UI (`apps/login/locales/`) |
| UI dialogs | Subset of UI: modal and confirmation dialog strings |
| UI section | Subset of UI: labels tied to a specific page section |
| UI indicator | Subset of UI: inline status or indicator labels |
| UI notifications | Subset of UI: toast, banner, and notification strings |
| Docs | Documentation content (`apps/docs/content/`) |
| API | Proto files and API description/comment text (`proto/**`) |
| Everywhere | UI + Docs + API |
| — | Internal only — no user-facing file scope |

## Canonical Terminology Table

| Canonical term | Meaning / explanation | Search for (discouraged) | Replace with / enforce | Action | Scope |
|---|---|---|---|---|---|
| Customer Portal | Central hub for all customer interactions for cloud and self-hosting customers | _(none — already canonical)_ | Customer Portal | keep | Everywhere |
| Management Console | Web interface where customers configure and manage ZITADEL resources | Console, ZITADEL Console, Admin Console, Administration Console | Management Console | replace | UI — must be visible |
| Instance | Private, isolated top-level ZITADEL environment | IAM, System, Type IAM | Instance / Type Instance | replace | Everywhere |
| Policies | Enforcement rules governing checks and constraints | Instance Policies, IAM Policies, Org Policies, Policies (unscoped), Instance Settings (when enforcing), Org Settings (when enforcing) | Instance Policies / Organization Policies | replace | Everywhere (enforcement context only, always scoped) |
| Settings | Resource-specific configuration values (not rules) | Instance Settings, Org Settings, Instance Policies (when config), Org Policies (when config), IAM Policies (when config) | Instance Settings / Organization Settings | replace | Everywhere (configuration context only) |
| Organization | Group of users within an instance | _(none — already canonical)_ | Organization | keep | Everywhere |
| Organization Domain | Domain giving context where a user belongs | Primary Domain, Verified Domains, Org domains, verify your domain | Organization Domain | replace | UI + Docs |
| User (Human) | User with interactive authentication flows | Human, Human User, User: Type Human | User (Human) | replace | UI + Docs |
| Service Account | User with non-interactive authentication flows | Machine User, machine user, Service User, Machine Account, Technical Account, User: Type Machine | Service Account | replace | UI + Docs |
| User | UI display label for user identity | Display Name | User | replace | UI only |
| Project | Container for applications sharing a role context | _(none — already canonical)_ | Project | keep | Everywhere |
| Project Grant | Delegation of project access to another organization | Grant, Grants, Organization Grant, Delegated Access | Project Grant / Project Grants | replace | UI + Docs |
| Application | Software or service secured using ZITADEL | _(none — already canonical)_ | Application | keep | Everywhere |
| Role Assignment | What a user is allowed to do (roles + org + user) | Authorization, internal authorization, external authorization, User Grant, Roles and Authorizations | Role Assignment | replace | Everywhere |
| Administrator | Role granting administrative privileges | Manager, Add Manager, Add a Manager, Membership, Member grants | Administrator / Add Administrator / Add an Administrator | replace | UI + Docs (role context only) |
| Organization Administrators | Org-level admin role holders | ZITADEL Organization Managers | Organization Administrators | replace | UI + Docs |
| Project Administrators | Project-level admin role holders | Project A Managers | Project Administrators | replace | UI + Docs |
| Administrator Roles | Set of admin roles | Manager Roles | Administrator Roles | replace | UI + Docs |
| ZITADEL Administrator Roles | ZITADEL-specific admin role set | ZITADEL Manager Roles, Zitadel Manager Roles | ZITADEL Administrator Roles | replace | UI + Docs |
| Permission | Internal permission backing admin roles | _(internal term)_ | _(do not surface in user-facing text)_ | internal | — |
| Metadata | Key-value custom data attached to resources | Meta Data | Metadata | replace | Everywhere |
| Custom Domain | Domain identifying a ZITADEL instance (globally unique) | Custom domain, Installed domains, Instance Domains, Zitadel Domain, your_domain, your-domain | Custom Domain | replace | UI + Docs |
| Trusted Domain | Domain used for API/email contexts | _(none — already canonical)_ | Trusted Domain | keep | Everywhere |
| Passkey | Passwordless auth using device-bound credentials | passwordless, passwordless login, passwordless auth, Multifactor (fingerprint/security keys), Fingerprint, Security Keys, WebAuthn, Webauthn | Passkey | replace | UI + Docs |
| TOTP | Time-based one-time password via authenticator app | OTP (authenticator), Authenticator App | TOTP | replace | UI + Docs |
| U2F | Legacy hardware authentication (deprecated) | U2F | _(remove)_ | remove | Everywhere |
| OTP Email | One-time password delivered via email | Email OTP | OTP Email | replace | UI + Docs |
| OTP SMS | One-time password delivered via SMS | SMS OTP | OTP SMS | replace | UI + Docs |
| Organization ID | Explicit organization identifier | Resource Owner, OrgID, OrganizationID | `organization_id` | replace | API only |
| Explicit object IDs | Explicit identifier per resource type | Resource ID, ResourceID | `user_id` / `project_id` / `application_id` / `instance_id` / `organization_id` | replace | API only |
| Instance ID | Instance identifier label in UI | Resource Id (Instance) | Instance ID | replace | UI only |
| ID | Generic identifier label in UI | Resource Id | ID | replace | UI only |
| First Name | Personal given name field | Given Name | First Name | replace | Everywhere |
| Last Name | Personal family name field | Family Name | Last Name | replace | Everywhere |
| Add Administrator (dialog) | Consistent wording for admin-add dialog | Add Manager, Add a Manager | Add Administrator / Add an Administrator | replace | UI dialogs |
| Administrator (memberships section) | UI section on the user detail page showing admin memberships — should be labeled "Administrator", not "Memberships" | Memberships | Administrator | replace | UI section |
| Project Grants (section) | Project grant listing section | Grants | Project Grants | replace | UI section |
| Internal / External indicator | Shows if a user belongs to the same or a different org | _(missing indicator)_ | internal / external | replace | UI indicator |
| Password changed | Past-tense wording for password change notification | Password change | Password changed | replace | UI notifications |
| Object descriptions | All resource descriptions must use end-user language | internal / technical wording | clear end-user language | replace | Docs + UI |
| Complement Token | Flow type for actions executed during token creation. In UI use "Complement Token"; in API use `flowType = 2` (CustomiseToken). Docs must not show PreUserinfoCreation=3 for this type. | Compliment Token, CustomiseToken, CustomizeToken | Complement Token | replace | UI + Docs + API |

## Governance

- **To add a new term**: open a PR that updates this table and references the decision thread.
- **Terms under discussion** may be added with action `proposed` — agents will flag but not block.
- **Ownership**: Docs + Product/UX + API maintainers approve changes to this file.
- **Source of truth for this table**: [GitHub issue #5888](https://github.com/zitadel/zitadel/issues/5888)

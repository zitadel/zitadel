# ZITADEL Canonical Terminology

Authoritative naming reference derived from [issue #5888](https://github.com/zitadel/zitadel/issues/5888).
All user-facing wording in docs, UI, and API descriptions must follow this table.

## Action Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Correct as-is — keep it |
| 🔁 | Replace when found — use the canonical term |
| ❌ | Remove entirely |
| 🧠 | Internal use only — must not appear in user-facing text |

## Scope Legend

| Symbol | Scope |
|--------|-------|
| 👁 UI | Console (`console/src/assets/i18n/`) and Login UI (`apps/login/locales/`) |
| 📄 Docs | Documentation content (`apps/docs/content/`) |
| 🔒 API | Proto files and API description/comment text (`proto/**`) |
| 🌐 Everywhere | All of the above |

## Canonical Terminology Table

| Canonical term | Meaning / explanation | Search for (discouraged) | Replace with / enforce | Action |
|---|---|---|---|---|
| Customer Portal | Central hub for all customer interactions for cloud and self-hosting customers | _(none — already canonical)_ | Customer Portal | ✅ |
| Management Console | Web interface where customers configure and manage ZITADEL resources | _(missing label, unclear app name)_ | Management Console | 🔁 👁 UI — must be visible |
| Instance | Private, isolated top-level ZITADEL environment | IAM, System, Type IAM | Instance / Type Instance | 🔁 🌐 Everywhere |
| Policies | Enforcement rules governing checks and constraints | Instance Policies, IAM Policies, Org Policies, Policies (unscoped), Instance Settings (when enforcing), Org Settings (when enforcing) | Instance Policies / Organization Policies | 🔁 🌐 Everywhere (enforcement only, scoped) |
| Settings | Resource-specific configuration values (not rules) | Instance Settings, Org Settings, Instance Policies (when config), Org Policies (when config), IAM Policies (when config) | Instance Settings / Organization Settings | 🔁 🌐 Everywhere (configuration only) |
| Organization | Group of users within an instance | _(none — already canonical)_ | Organization | ✅ |
| Organization Domain | Domain giving context where a user belongs | Primary Domain, Verified Domains, Org domains, verify your domain | Organization Domain | 🔁 👁 UI + 📄 Docs |
| User (Human) | User with interactive authentication flows | Human, Human User, User: Type Human | User (Human) | 🔁 👁 UI + 📄 Docs |
| Service Account | User with non-interactive authentication flows | Machine User, machine user, Service User, Machine Account, Technical Account, User: Type Machine | Service Account | 🔁 👁 UI + 📄 Docs |
| User | UI label for user identity (display) | Display Name | User | 🔁 👁 UI only |
| Project | Container for applications sharing a role context | _(none — already canonical)_ | Project | ✅ |
| Project Grant | Delegation of project access to another organization | Grant, Grants, Organization Grant, Delegated Access | Project Grant / Project Grants | 🔁 👁 UI + 📄 Docs |
| Application | Software or service secured using ZITADEL | _(none — already canonical)_ | Application | ✅ |
| Role Assignment | What a user is allowed to do (roles + org + user) | Authorization, internal authorization, external authorization, User Grant, Roles and Authorizations | Role Assignment | 🔁 🌐 Everywhere |
| Administrator | Role granting administrative privileges | Manager, Add Manager, Add a Manager, Membership, Member grants | Administrator / Add Administrator / Add an Administrator | 🔁 👁 UI + 📄 Docs (role only) |
| Organization Administrators | Org-level admin role holders | ZITADEL Organization Managers | Organization Administrators | 🔁 👁 UI + 📄 Docs |
| Project Administrators | Project-level admin role holders | Project A Managers | Project Administrators | 🔁 👁 UI + 📄 Docs |
| Administrator Roles | Set of admin roles | Manager Roles | Administrator Roles | 🔁 👁 UI + 📄 Docs |
| ZITADEL Administrator Roles | ZITADEL-specific admin role set | ZITADEL Manager Roles, Zitadel Manager Roles | ZITADEL Administrator Roles | 🔁 👁 UI + 📄 Docs |
| Permission | Internal permission backing admin roles | _(internal term)_ | _(internal only — do not surface in user-facing text)_ | 🧠 Internal only |
| Metadata | Key-value custom data attached to resources | Meta Data | Metadata | 🔁 🌐 Everywhere |
| Custom Domain | Domain identifying a ZITADEL instance (globally unique) | Custom domain, Installed domains, Instance Domains, Zitadel Domain, your_domain, your-domain | Custom Domain | 🔁 👁 UI + 📄 Docs |
| Trusted Domain | Domain used for API/email contexts | _(none — already canonical)_ | Trusted Domain | ✅ |
| Passkey | Passwordless auth using device-bound credentials | passwordless, passwordless login, passwordless auth, Multifactor (fingerprint/security keys), Fingerprint, Security Keys, WebAuthn, Webauthn | Passkey | 🔁 👁 UI + 📄 Docs |
| TOTP | Time-based one-time password via authenticator app | OTP (authenticator), Authenticator App | TOTP | 🔁 👁 UI + 📄 Docs |
| U2F | Legacy hardware authentication (deprecated) | U2F | _(remove)_ | ❌ 🌐 Everywhere |
| OTP Email | One-time password delivered via email | Email OTP | OTP Email | 🔁 👁 UI + 📄 Docs |
| OTP SMS | One-time password delivered via SMS | SMS OTP | OTP SMS | 🔁 👁 UI + 📄 Docs |
| Organization ID | Explicit organization identifier | Resource Owner, OrgID, OrganizationID | `organization_id` | 🔁 🔒 API only |
| Explicit object IDs | Explicit identifier per resource type | Resource ID, ResourceID | `user_id` / `project_id` / `application_id` / `instance_id` / `organization_id` | 🔁 🔒 API only |
| Instance ID | Instance identifier label in UI | Resource Id (Instance) | Instance ID | 🔁 👁 UI only |
| ID | Generic identifier label in UI | Resource Id | ID | 🔁 👁 UI only |
| First Name | Personal given name field | Given Name | First Name | 🔁 🌐 Everywhere |
| Last Name | Personal family name field | Family Name | Last Name | 🔁 🌐 Everywhere |
| Add Administrator (dialog) | Consistent wording for admin-add dialog | Add Manager, Add a Manager | Add Administrator / Add an Administrator | 🔁 👁 UI dialogs |
| Memberships (section) | Admin memberships on user detail page | Memberships | Administrator | 🔁 👁 UI section |
| Project Grants (section) | Project grant listing section | Grants | Project Grants | 🔁 👁 UI section |
| Internal / External indicator | Shows if a user belongs to the same or a different org | _(missing indicator)_ | internal / external | 🔁 👁 UI indicator |
| Password changed | Past-tense wording for password change notification | Password change | Password changed | 🔁 👁 UI notifications |
| Object descriptions | All resource descriptions must use end-user language | internal / technical wording | clear end-user language | 🔁 👁 📄 Docs + UI |
| Complement Token | Flow type for actions executed during token creation | Compliment Token, Complement Token, CustomiseToken, CustomizeToken | **UI name:** Complement Token. **API:** `flowType = 2` (CustomiseToken). Docs must not show PreUserinfoCreation=3 for this type. | 🔁 👁 UI + 📄 Docs + 🔒 API |

## Governance

- **To add a new term**: open a PR that updates this table and references the decision thread.
- **Terms under discussion** may be added with a `🔬 Proposed` action — agents will flag but not block.
- **Ownership**: Docs + Product/UX + API maintainers approve changes to this file.
- **Source of truth for this table**: [GitHub issue #5888](https://github.com/zitadel/zitadel/issues/5888)

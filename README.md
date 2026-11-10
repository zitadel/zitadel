<p align="center">
    <img src="./apps/docs/public/img/logos/zitadel-logo-dark@2x.png#gh-light-mode-only" alt="ZITADEL Logo" max-height="200px" width="auto" />
    <img src="./apps/docs/public/img/logos/zitadel-logo-light@2x.png#gh-dark-mode-only" alt="ZITADEL Logo" max-height="200px" width="auto" />
</p>

<p align="center">
    <a href="https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/zitadel/zitadel" alt="Open in Dev Container">
        <img src="https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue" /></a>
    <a href="https://github.com/zitadel/zitadel/blob/main/LICENSE" alt="License">
        <img src="https://badgen.net/github/license/zitadel/zitadel/" /></a>
    <a href="https://bestpractices.coreinfrastructure.org/projects/6662">
        <img src="https://bestpractices.coreinfrastructure.org/projects/6662/badge"></a>
    <a href="https://github.com/semantic-release/semantic-release" alt="semantic-release">
        <img src="https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg" /></a>
    <a href="https://github.com/zitadel/zitadel/actions" alt="ZITADEL Release">
        <img alt="GitHub Workflow Status (with event)" src="https://img.shields.io/github/actions/workflow/status/zitadel/zitadel/build.yml?event=pull_request"></a>
    <a href="https://zitadel.com/docs/support/software-release-cycles-support" alt="Release">
        <img src="https://badgen.net/github/release/zitadel/zitadel/stable" /></a>
    <a href="https://goreportcard.com/report/github.com/zitadel/zitadel" alt="Go Report Card">
        <img src="https://goreportcard.com/badge/github.com/zitadel/zitadel" /></a>
    <a href="https://codecov.io/gh/zitadel/zitadel" alt="Code Coverage">
        <img src="https://codecov.io/gh/zitadel/zitadel/branch/main/graph/badge.svg" /></a>
    <a href="https://github.com/zitadel/zitadel/graphs/contributors" alt="Contributors">
        <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/zitadel/zitadel"></a>
    <a href="https://discord.gg/YgjEuJzZ3x" alt="Discord Chat">
        <img src="https://badgen.net/discord/online-members/YgjEuJzZ3x" /></a>
</p>

<p align="center">
    <a href="https://openid.net/certification/#OPs" alt="OpenID Connect Certified">
        <img src="./apps/docs/public/img/logos/oidc-cert.png" /></a>
</p>

## The Identity Infrastructure for Developers

**ZITADEL** is an open-source identity and access management platform built for teams that need more than basic auth. Whether you're securing a SaaS product, building a B2B platform, or self-hosting a production IAM stack — ZITADEL gives you everything out of the box: SSO, MFA, Passkeys, OIDC, SAML, SCIM, and a battle-tested multi-tenancy model.

No vendor lock-in. No compromise on control. Just a robust, API-first identity platform you can own.

---

**[🏡 Website](https://zitadel.com) &nbsp;|&nbsp; [💬 Chat](https://zitadel.com/chat) &nbsp;|&nbsp; [📋 Docs](https://zitadel.com/docs/) &nbsp;|&nbsp; [🧑‍💻 Blog](https://zitadel.com/blog) &nbsp;|&nbsp; [📞 Contact](https://zitadel.com/contact/)**

---

## Why ZITADEL

We built ZITADEL to handle the hardest IAM challenges at scale — starting with multi-tenancy.

| | ZITADEL | FusionAuth | Keycloak | Auth0/Okta |
|---|---|---|---|---|
| Open-source | ✅ | ❌ | ✅ | ❌ |
| Self-hostable | ✅ | ✅ | ✅ | ❌ |
| Infrastructure-level tenants | ✅ Instances (High scale) | ✅ Tenants | 🟡 Realms (Scaling limits) | ❌ (Multi-tenant = multi-account) |
| B2B Organizations | ✅ Native & Unlimited | 🟡 via Entity Management | ✅ (Recent addition) | 🟡 (Plan/Account dependent) |
| Full audit trail | ✅ Comprehensive Event Stream* | 🟡 Audit logs | 🟡 Audit logs | 🟡 Audit logs |
| Passkeys (FIDO2) | ✅ | ✅ | ✅ | ✅ |
| [Actions / webhooks](https://zitadel.com/docs/concepts/features/actions_v2) | ✅ | ✅ | 🟡 via SPI | ✅ |
| API-first (gRPC + REST) | ✅ | 🟡 REST only | 🟡 REST only | 🟡 REST only |
| SaaS + self-host parity | ✅ | ✅ | ➖ N/A | ➖ N/A |

ZITADEL Cloud and self-hosted ZITADEL run the same codebase.

**Key differentiators for architects:**
- **Relational core, event-driven soul** — every mutation is written as an immutable event for a complete, API-accessible [audit trail](https://zitadel.com/docs/concepts/features/audit-trail). Unlike systems that log only select activities, ZITADEL provides a comprehensive event stream that can be audited or streamed to external systems via Webhooks.
- **Strict multi-tenant hierarchy** — Identity System → Organizations → Projects, with isolated data and policy scoping at multiple levels
- **API-first design** — every resource and action is available via [connectRPC, gRPC, and HTTP/JSON APIs](https://zitadel.com/docs/apis/introduction)
- **[Zero-downtime updates](https://zitadel.com/docs/concepts/architecture/solution#zero-downtime-updates)** and [horizontal scalability](https://zitadel.com/docs/self-hosting/manage/updating_scaling) without external session stores

---

## Get Started in 3 Minutes

👉 [Quick Start Guide](https://zitadel.com/docs/guides/start/quickstart)

### ZITADEL Self-Hosted

```bash
# Docker Compose — up and running in under 3 minutes
curl -LO https://raw.githubusercontent.com/zitadel/zitadel/main/deploy/compose/docker-compose.yml \
  && curl -LO https://raw.githubusercontent.com/zitadel/zitadel/main/deploy/compose/.env.example \
  && cp .env.example .env \
  && docker compose up -d --wait
```

Full deployment guides:
- [Docker Compose](https://zitadel.com/docs/self-hosting/deploy/compose) 
- [Kubernetes](https://zitadel.com/docs/self-hosting/deploy/kubernetes)

> Need professional support for your self-hosted deployment? [Contact us](https://zitadel.com/contact).

### ZITADEL Cloud (SaaS)

Start for free at [zitadel.com](https://zitadel.com) — no credit card required. Available in US · EU · AU · CH. [Pay-as-you-go pricing](https://zitadel.com/pricing).

---

## Integrate with the V2 API

ZITADEL exposes every capability over a typed API. Here's how to create a user with the V2 REST API:

```bash
curl -X POST https://$ZITADEL_DOMAIN/v2/users/human \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice@example.com",
    "profile": { "givenName": "Alice", "familyName": "Smith" },
    "email": { "email": "alice@example.com", "sendCode": {} }
  }'
```

Explore the full [API reference](https://zitadel.com/docs/apis/introduction) — including connectRPC and gRPC transports — or jump straight to [quickstart examples](https://zitadel.com/docs/guides/start/quickstart).

---

## Features

**Authentication**
- Single Sign On (SSO) · Username/Password · [Passkeys (FIDO2 / WebAuthn)](https://zitadel.com/docs/concepts/features/passkeys)
- MFA: OTP, U2F, OTP Email, OTP SMS
- [LDAP](https://zitadel.com/docs/guides/integrate/identity-providers/ldap) · [Enterprise IdPs and social logins](https://zitadel.com/docs/guides/integrate/identity-providers/introduction)
- [OpenID Connect certified](https://openid.net/certification/#OPs) · [SAML 2.0](https://zitadel.com/docs/apis/saml/endpoints) · [Device authorization](https://zitadel.com/docs/guides/integrate/login/oidc/device-authorization)
- [Machine-to-machine](https://zitadel.com/docs/guides/integrate/service-accounts/authenticate-service-accounts): JWT Profile, PAT, Client Credentials
- [Token exchange and impersonation](https://zitadel.com/docs/guides/integrate/token-exchange)
- [Custom sessions](https://zitadel.com/docs/guides/integrate/login-ui/username-password) for flows beyond OIDC/SAML
- [Hosted Login V2](https://zitadel.com/docs/guides/integrate/login/hosted-login)

**Multi-Tenancy**
- [Identity brokering](https://zitadel.com/docs/concepts/features/identity-brokering) with pre-built IdP templates
- [Customizable B2B onboarding](https://zitadel.com/docs/guides/integrate/onboarding/b2b) with self-service for customers
- [Delegated role management](https://zitadel.com/docs/guides/manage/console/projects-overview) to third parties
- [Domain discovery](https://zitadel.com/docs/guides/solution-scenarios/domain-discovery)

**Integration**
- [gRPC, connectRPC, and REST APIs](https://zitadel.com/docs/apis/introduction) for every resource
- [Actions](https://zitadel.com/docs/concepts/features/actions_v2): webhooks, custom code, token enrichment
- [RBAC](https://zitadel.com/docs/guides/integrate/retrieve-user-roles) · [SCIM 2.0 Server](https://zitadel.com/docs/apis/scim2)
- [Audit log and SOC/SIEM integration](https://zitadel.com/docs/guides/integrate/external-audit-log)
- [SDKs and example apps](https://zitadel.com/docs/sdk-examples/introduction)

**Self-Service & Admin**
- [Self-registration](https://zitadel.com/docs/concepts/features/selfservice#registration) with email/phone verification
- [Administration Console](https://zitadel.com/docs/guides/manage/console/console-overview) for orgs and projects
- [Custom branding](https://zitadel.com/docs/guides/manage/customize/branding) per organization

**Deployment**
- [PostgreSQL](https://zitadel.com/docs/self-hosting/manage/database#postgres) (≥ 14) · [Zero-downtime updates](https://zitadel.com/docs/concepts/architecture/solution#zero-downtime-updates) · [High scalability](https://zitadel.com/docs/self-hosting/manage/production)

Track upcoming features on our [roadmap](https://zitadel.com/roadmap) and follow our [changelog](https://zitadel.com/changelog) for recent updates.

---

## Showcase

### Login V2

Our new, fully customizable login experience — [documentation](https://zitadel.com/docs/guides/integrate/login/hosted-login)

---

## Adopters & Ecosystem

Used in production by organizations worldwide. See the full [Adopters list](./ADOPTERS.md) — and add yours by submitting a pull request.

- **SDKs**: [All supported languages and frameworks](https://zitadel.com/docs/sdk-examples/introduction)
- **Examples**: [Clone and use our examples](https://zitadel.com/docs/sdk-examples/introduction)

---

## How To Contribute

ZITADEL is built in the open and welcoming to contributions of all kinds.

- 📖 Read the [Contribution Guide](./CONTRIBUTING.md) to get started
- 💬 Join the conversation on [Discord](https://zitadel.com/chat)
- 🐛 Report bugs or request features via [GitHub Issues](https://github.com/zitadel/zitadel/issues)

## Contributors

<a href="https://github.com/zitadel/zitadel/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=zitadel/zitadel" />
</a>

Made with [contrib.rocks](https://contrib.rocks/preview?repo=zitadel/zitadel).

---

## Security

Security policy: [SECURITY.md](./SECURITY.md)

[Vulnerability Disclosure Policy](https://zitadel.com/docs/legal/policies/vulnerability-disclosure-policy) — how to responsibly report security issues.

[Technical Advisories](https://zitadel.com/docs/support/technical_advisory) are published for major issues that could impact security or stability in production.

## License

[AGPL-3.0](./LICENSE) — see [LICENSING.md](./LICENSING.md) for the full licensing policy, including Apache 2.0 and MIT exceptions for specific directories.

<p align="center">
    <img src="./docs/static/logos/zitadel-logo-dark@2x.png#gh-light-mode-only" alt="Zitadel Logo" max-height="200px" width="auto" />
    <img src="./docs/static/logos/zitadel-logo-light@2x.png#gh-dark-mode-only" alt="Zitadel Logo" max-height="200px" width="auto" />
</p>

<p align="center">
    <a href="https://bestpractices.coreinfrastructure.org/projects/6662"><img src="https://bestpractices.coreinfrastructure.org/projects/6662/badge"></a>
    <a href="https://github.com/zitadel/zitadel/graphs/contributors" alt="Release">
        <img src="https://badgen.net/github/contributors/zitadel/zitadel" /></a>
    <a href="https://github.com/semantic-release/semantic-release" alt="semantic-release">
        <img src="https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg" /></a>
    <a href="https://github.com/zitadel/zitadel/actions" alt="ZITADEL Release">
        <img src="https://github.com/zitadel/zitadel/actions/workflows/zitadel.yml/badge.svg" /></a>
    <a href="https://github.com/zitadel/zitadel/blob/main/LICENSE" alt="License">
        <img src="https://badgen.net/github/license/zitadel/zitadel/" /></a>
    <a href="https://github.com/zitadel/zitadel/releases" alt="Release">
        <img src="https://badgen.net/github/release/zitadel/zitadel/stable" /></a>
    <a href="https://goreportcard.com/report/github.com/zitadel/zitadel" alt="Go Report Card">
        <img src="https://goreportcard.com/badge/github.com/zitadel/zitadel" /></a>
    <a href="https://codecov.io/gh/zitadel/zitadel" alt="Code Coverage">
        <img src="https://codecov.io/gh/zitadel/zitadel/branch/main/graph/badge.svg" /></a>
    <a href="https://discord.gg/erh5Brh7jE" alt="Discord Chat">
        <img src="https://badgen.net/discord/online-members/erh5Brh7jE" /></a>
</p>

<p align="center">
    <a href="https://openid.net/certification/#OPs" alt="OpenID Connect Certified">
        <img src="./docs/static/logos/oidc-cert.png" /></a>
</p>

Do you look for a user management that's quickly set up like Auth0 and open source like Keycloak?

Do you have project that requires a multi-tenant user management with self-service for your customers?

Look no further — ZITADEL combines the ease of Auth0 with the versatility of Keycloak.

We provide you with a wide range of out of the box features to accelerate your project.
Multi-tenancy with branding customization, secure login, self-service, OpenID Connect, OAuth2.x, SAML2, Passwordless with FIDO2 (including Passkeys), OTP, U2F, and an unlimited audit trail is there for you, ready to use.

With ZITADEL you can rely on a hardened and extensible turnkey solution to solve all of your authentication and authorization needs.

---

**[🏡 Website](https://zitadel.com) [💬 Chat](https://zitadel.com/chat) [📋 Docs](https://zitadel.com/docs/) [🧑‍💻 Blog](https://zitadel.com/blog) [📞 Contact](https://zitadel.com/contact/)**

## Get started

👉 [Quick Start Guide](https://zitadel.com/docs/guides/start/quickstart)

### Deploy ZITADEL (Self-Hosted)

Deploying ZITADEL locally takes less than 3 minutes. So go ahead and give it a try!

* [Linux](https://zitadel.com/docs/self-hosting/deploy/linux)
* [MacOS](https://zitadel.com/docs/self-hosting/deploy/macos)
* [Docker compose](https://zitadel.com/docs/self-hosting/deploy/compose)
* [Knative](https://zitadel.com/docs/self-hosting/deploy/knative)
* [Kubernetes](https://zitadel.com/docs/self-hosting/deploy/kubernetes)

See all guides [here](https://zitadel.com/docs/self-hosting/deploy/overview)

> If you are interested to get professional support for your self-hosted ZITADEL [please reach out to us](https://zitadel.com/contact)!

### Setup ZITADEL Cloud (SaaS)

If you want to experience a hands-free ZITADEL, you should use [ZITADEL Cloud](https://zitadel.cloud).

It is free for up to 25'000 authenticated requests and provides you all the features that make ZITADEL great.
Learn more about the [pay-as-you-go pricing](https://zitadel.com/pricing).

### Example applications

Clone one of our [example applications](https://zitadel.com/docs/examples/introduction#clone-a-sample-project) or deploy them directly to Vercel.

### SDKs

Use our [SDKs](https://zitadel.com/docs/examples/sdks) for your favorite language and framework.

## Why choose ZITADEL

We built ZITADEL with a complex multi-tenancy architecture in mind and provide the best solution to handle [B2B customers and partners](https://zitadel.com/docs/guides/solution-scenarios/b2b).
Yet it offers everything you need for a customer identity ([CIAM](https://zitadel.com/docs/guides/solution-scenarios/b2c)) use case.

- [API-first approach](https://zitadel.com/docs/apis/introduction)
- Strong audit trail thanks to [event sourcing](https://zitadel.com/docs/concepts/eventstore/overview) as storage pattern
- [Actions](https://zitadel.com/docs/apis/actions/introduction) to react on events with custom code and extended ZITADEL for you needs
- [Branding](https://zitadel.com/docs/guides/manage/customize/branding) for a uniform user experience across multiple organizations
- [Self-service](https://zitadel.com/docs/concepts/features/selfservice) for end-users, business customers, and administrators
- [CockroachDB](https://www.cockroachlabs.com/) or a [Postgres](https://www.postgresql.org/) database as reliable and widespread storage option

## Features

- Single Sign On (SSO)
- Passwordless with FIDO2 support (Including Passkeys)
- Username / Password
- Multifactor authentication with OTP, U2F
- [Identity Brokering](https://zitadel.com/docs/guides/integrate/identity-brokering)
- [Machine-to-machine (JWT profile)](https://zitadel.com/docs/guides/integrate/serviceusers)
- Personal Access Tokens (PAT)
- Role Based Access Control (RBAC)
- [Delegate role management to third-parties](https://zitadel.com/docs/guides/manage/console/projects)
- [Self-registration](https://zitadel.com/docs/concepts/features/selfservice#registration) including verification
- [Self-service](https://zitadel.com/docs/concepts/features/selfservice) for end-users, business customers, and administrators
- [OpenID Connect certified](https://openid.net/certification/#OPs) => [OIDC Endpoints](https://zitadel.com/docs/apis/openidoauth/endpoints)
- [SAML 2.0](http://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0.html) => [SAML Endpoints](https://zitadel.com/docs/apis/saml/endpoints)
- [Postgres](https://zitadel.com/docs/self-hosting/manage/database#postgres) (version >= 14) or [CockroachDB](https://zitadel.com/docs/self-hosting/manage/database#cockroach) (version >= 22.0)

Track upcoming features on our [roadmap](https://zitadel.com/roadmap).

## How To Contribute

Details about how to contribute you can find in the [Contribution Guide](./CONTRIBUTING.md)

## Contributors

<a href="https://github.com/zitadel/zitadel/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=zitadel/zitadel" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## Showcase

### Quick Start Guide

Secure a React Application using OpenID Connect Authorization Code with PKCE

[![Quick Start Guide](https://user-images.githubusercontent.com/1366906/223662449-f17b734d-405c-4945-a8a1-200440c459e5.gif)](http://www.youtube.com/watch?v=5THbQljoPKg "Quick Start Guide")

### Login with Passkeys

Use our login widget to allow easy and secure access to your applications and enjoy all the benefits of Passkeys (FIDO 2 / WebAuthN):

[![Passkeys](https://user-images.githubusercontent.com/1366906/223664178-4132faef-4832-4014-b9ab-90c2a8d15436.gif)](https://www.youtube.com/watch?v=cZjHQYurSjw&list=PLTDa7jTlOyRLdABgD2zL0LGM7rx5GZ1IR&index=2 "Passkeys")

### Admin Console

Use [Console](https://zitadel.com/docs/guides/manage/console/overview) or our [APIs](https://zitadel.com/docs/apis/introduction) to setup organizations, projects and applications.

[![Console Showcase](https://user-images.githubusercontent.com/1366906/223663344-67038d5f-4415-4285-ab20-9a4d397e2138.gif)](http://www.youtube.com/watch?v=RPpHktAcCtk "Console Showcase")

## Security

See the policy [here](./SECURITY.md)

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.

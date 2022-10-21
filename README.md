<p align="center">
    <img src="./docs/static/logos/zitadel-logo-dark@2x.png#gh-light-mode-only" alt="Zitadel Logo" max-height="200px" width="auto" />
    <img src="./docs/static/logos/zitadel-logo-light@2x.png#gh-dark-mode-only" alt="Zitadel Logo" max-height="200px" width="auto" />
</p>

<p align="center">
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

Do you have project that requires a multi-tenant user management with self-service for you customers?

Look no further — ZITADEL combines the ease of Auth0 and the versatility of Keycloak.

We provide you with a wide range of out of the box features to accelerate your project. Multi-tenancy with branding customization, secure login, self-service, OpenID Connect, OAuth2.x, SAML2, Passwordless with FIDO2 (including Passkeys), OTP, U2F, and an unlimited audit trail is there for you, ready to use.

With ZITADEL you can rely on a hardened and extensible turnkey solution to solve all of your authentication and authorization needs.

---

**[🏡 Website](https://zitadel.com) [💬 Chat](https://zitadel.com/chat) [📋 Docs](https://docs.zitadel.com/) [🧑‍💻 Blog](https://zitadel.com/blog) [📞 Contact](https://zitadel.com/contact/)**

## Get started

### Deploy ZITADEL (Self-Hosted)

Deploying ZITADEL locally does take less than 3 minutes. So go ahead and give it a try!

* [Linux](https://docs.zitadel.com/docs/guides/deploy/linux)
* [MacOS](https://docs.zitadel.com/docs/guides/deploy/macos)
* [Docker compose](https://docs.zitadel.com/docs/guides/deploy/compose)
* [Knative](https://docs.zitadel.com/docs/guides/deploy/knative)
* [Kubernetes](https://docs.zitadel.com/docs/guides/deploy/kubernetes)

See all guides [here](https://docs.zitadel.com/docs/guides/deploy/overview)

> If you are interested to get professional support for your self-hosted ZITADEL [please reach out to us](https://zitadel.com/contact)!

### Setup ZITADEL Cloud (SaaS)

If you want to experience a hands-free ZITADEL you should use [ZITADEL Cloud](https://zitadel.cloud).

It is free for up to 25'000 authenticated requests and provides you all the features that make ZITADEL great. Learn more about the [pay-as-you-go pricing](https://zitadel.com/pricing).

### Quickstarts - Integrate your app

[Multiple Examples can be found here](https://docs.zitadel.com/docs/examples/introduction)

> If you miss something please feel free to [join the Discussion](https://github.com/zitadel/zitadel/discussions/1717)

## Why choose ZITADEL

We built ZITADEL with a complex multi-tenancy architecture in mind and provide the best solution to handle B2B cutomers and partners.
Yet it offers everything you need for a customer identity (CIAM) use case.

- [API-first approach](https://docs.zitadel.com/docs/apis/introduction)
- Strong audit trail thanks to [event sourcing](https://docs.zitadel.com/docs/concepts/eventstore/overview) as storage pattern
- [Actions](https://docs.zitadel.com/docs/concepts/features/actions) to react on events with custom code and extended ZITADEL for you needs
- [Branding](https://docs.zitadel.com/docs/guides/manage/customize/branding) for a uniform user experience across multiple organizations
- [Self-service](https://docs.zitadel.com/docs/concepts/features/selfservice) for end-users, business customers, and administrators
- [CockroachDB](https://www.cockroachlabs.com/) or a [Postgres](https://www.postgresql.org/) database as reliable and widespread storage option

## Features

- Single Sign On (SSO)
- Passwordless with FIDO2 support (Including Passkeys)
- Username / Password
- Multifactor authentication with OTP, U2F
- [Identity Brokering](https://docs.zitadel.com/docs/guides/integrate/identity-brokering)
- [Machine-to-machine (JWT profile)](https://docs.zitadel.com/docs/guides/integrate/serviceusers)
- Personal Access Tokens (PAT)
- Role Based Access Control (RBAC)
- [Delegate role management to third-parties](https://docs.zitadel.com/docs/guides/manage/console/projects)
- [Self-registration](https://docs.zitadel.com/docs/concepts/features/selfservice#registration) including verification
- [Self-service](https://docs.zitadel.com/docs/concepts/features/selfservice) for end-users, business customers, and administrators
- [OpenID Connect certified](https://openid.net/certification/#OPs) => [OIDC Endpoints](https://docs.zitadel.com/docs/apis/openidoauth/endpoints),  [OIDC Integration Guides](https://docs.zitadel.com/docs/guides/integrate/auth0-oidc)
- [SAML 2.0](http://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0.html) => [SAML Endpoints](https://docs.zitadel.com/docs/apis/saml/endpoints), [SAML Integration Guides](https://docs.zitadel.com/docs/guides/integrate/auth0-saml)
- [Postgres](https://docs.zitadel.com/docs/guides/manage/self-hosted/database#postgres) (version >= 14) or [CockroachDB](https://docs.zitadel.com/docs/guides/manage/self-hosted/database#cockroach) (version >= 22.0)

Track upcoming features on our [roadmap](https://zitadel.com/roadmap).

## Client libraries

| Language / Framework | Client | API | Machine auth (\*) | Auth check (\*\*) | Thanks to the maintainers |
|----------|--------|--------------|----------|---------|---------------------------|
| .NET     | [zitadel-net](https://github.com/smartive/zitadel-net) | GRPC | ✔️ | ✔️ | [smartive 👑](https://github.com/smartive/) |
| Dart     | [zitadel-dart](https://github.com/smartive/zitadel-dart) | GRPC | ✔️ | ❌ | [smartive 👑](https://github.com/smartive/) |
| Elixir   | [zitadel_api](https://github.com/jshmrtn/zitadel_api) | GRPC | ✔️ | ✔️ | [jshmrtn 🙏🏻](https://github.com/jshmrtn) |
| Go       | [zitadel-go](https://github.com/zitadel/zitadel-go) | GRPC | ✔️ | ✔️ | [ZITADEL](https://github.com/zitadel/) |
| Rust     | [zitadel-rust](https://crates.io/crates/zitadel) | GRPC | ✔️ | ❌ | [smartive 👑](https://github.com/smartive/) |
| JVM      | 🚧 [WIP](https://github.com/zitadel/zitadel/discussions/3650) | ❓ | ❓ | | TBD |
| Python   | 🚧 [WIP](https://github.com/zitadel/zitadel/issues/3675) | ❓ | ❓ | | TBD |
| Javascript | ❓ | ❓ | ❓ | | Maybe you? |

(\*) Automatically authenticate service accounts with [JWT Profile](https://docs.zitadel.com/docs/apis/openidoauth/grant-types#json-web-token-jwt-profile).  
(\*\*) Automatically check if the access token is valid and claims match

## How To Contribute

Details about how to contribute you can find in the [Contribution Guide](./CONTRIBUTING.md)

## Contributors

<a href="https://github.com/zitadel/zitadel/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=zitadel/zitadel" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## Showcase

### Passwordless Login

Use our login widget to allow easy and secure access to your applications and enjoy all the benefits of passwordless (FIDO 2 / WebAuthN):

- works on all modern platforms, devices, and browsers
- phishing resistant alternative
- requires only one gesture by the user
- easy [enrollment](https://docs.zitadel.com/docs/manuals/user-profile) of the device during registration

![passwordless-windows-hello](https://user-images.githubusercontent.com/1366906/118765435-5d419780-b87b-11eb-95bf-55140119c0d8.gif)

### Admin Console

Use [Console](https://docs.zitadel.com/docs/manuals/introduction) or our [APIs](https://docs.zitadel.com/docs/apis/introduction) to setup organizations, projects and applications.

[![Console Showcase](http://img.youtube.com/vi/RPpHktAcCtk/0.jpg)](http://www.youtube.com/watch?v=RPpHktAcCtk "Console Showcase")

## Security

See the policy [here](./SECURITY.md)

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

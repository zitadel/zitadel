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

You want auth that's quickly set up like Auth0 but open source like Keycloak? Look no further— ZITADEL combines the ease of Auth0 and the versatility of Keycloak.

We provide a wide range of out of the box features like secure login, self-service, OpenID Connect, OAuth2.x, SAML2, branding, Passwordless with FIDO2, OTP, U2F, and an unlimited audit trail to improve the life of developers. Especially noteworthy is that ZITADEL supports not only B2C and B2E scenarios but also B2B. This is super useful for people who build B2B Solutions, as ZITADEL can handle all the delegated user and access management.

With ZITADEL you rely on a battle tested, hardened and extensible turnkey solution to solve all of your authentication and authorization needs. With the unique way of how ZITADEL stores data it gives you an unlimited audit trail which provides a peace of mind for even the harshest audit and analytics requirements.

<!-- TODO: insert video here-->

---

**[🏡 Website](https://zitadel.com) [💬 Chat](https://zitadel.com/chat) [📋 Docs](https://docs.zitadel.ch/) [🧑‍💻 Blog](https://zitadel.com/blog) [📞 Contact](https://zitadel.com/contact/)**

## Get started

### ZITADEL Cloud (SaaS)

The easiest way to get started with ZITADEL is to use our public cloud offering.
Currently ZITADEL V2 Beta is available, head over to [https://zitadel.cloud](https://zitadel.cloud) and create your first ZITADEL instance for free.

You can also discover our pay-as-you-go [pricing](https://zitadel.com/pricing).

### Install ZITADEL

- [We provide installation guides for multiple platforms here](https://docs.zitadel.com/docs/guides/deploy/overview)

### Quickstarts - Integrate your app

- [Multiple Quickstarts can be found here](hhttps://docs.zitadel.com/docs/guides/start/quickstart)
- [And even more examples are located under zitadel/zitadel-examples](https://github.com/zitadel/zitadel-examples)

> If you miss something please feel free to engage with us [here](https://github.com/zitadel/zitadel/discussions/1717)

## Why ZITADEL

- [API-first](https://docs.zitadel.com/docs/apis/introduction)
- Strong audit trail thanks to [event sourcing](https://docs.zitadel.com/docs/concepts/eventstore/overview)
- [Actions](https://docs.zitadel.ch/docs/concepts/features/actions) to react on events with custom code
- [Branding](https://docs.zitadel.com/docs/guides/manage/customize/branding) for a uniform user experience
- [CockroachDB](https://www.cockroachlabs.com/) is the only dependency

## Features

- Single Sign On (SSO)
- Passwordless with FIDO2 support
- Username / Password
- Multifactor authentication with OTP, U2F
- [Identity Brokering](https://docs.zitadel.com/docs/guides/integrate/identity-brokering)
- [Machine-to-machine (JWT profile)](https://docs.zitadel.com/docs/guides/integrate/serviceusers)
- Personal Access Tokens (PAT)
- Role Based Access Control (RBAC)
- [Delegate role management to third-parties](https://docs.zitadel.com/docs/guides/manage/console/projects)
- Self-registration including verification
- User self service
- [Service Accounts](https://docs.zitadel.com/docs/guides/integrate/serviceusers)
- [OpenID Connect certified](https://openid.net/certification/#OPs)
- 🚧 [SAML 2.0](https://github.com/zitadel/zitadel/pull/3618)
- 🚧 [Postgres](https://github.com/zitadel/zitadel/pull/3998)

## Client libraries

<!-- TODO: check other libraries -->

| Language | Client | API | Machine auth (\*) | Auth check (\*\*) | Thanks to the maintainers |
|----------|--------|--------------|----------|---------|---------------------------|
| .NET     | [zitadel-net](https://github.com/zitadel/zitadel-net) | GRPC | ✔️ | ✔️ | [buehler 👑](https://github.com/buehler) |
| Dart     | [zitadel-dart](https://github.com/zitadel/zitadel-dart) | GRPC | ✔️ | ❌ | [buehler 👑](https://github.com/buehler) |
| Elixir   | [zitadel_api](https://github.com/jshmrtn/zitadel_api) | GRPC | ✔️ | ✔️ | [jshmrtn 🙏🏻](https://github.com/jshmrtn) |
| Go       | [zitadel-go](https://github.com/zitadel/zitadel-go) | GRPC | ✔️ | ✔️ | ZITADEL |
| Rust     | [zitadel-rust](https://crates.io/crates/zitadel) | GRPC | ✔️ | ❌ | [buehler 👑](https://github.com/buehler) |
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

<!-- TODO: Replace Images-->

### Passwordless Login

Use our login widget to allow easy and secure access to your applications and enjoy all the benefits of passwordless (FIDO 2 / WebAuthN):

- works on all modern platforms, devices, and browsers
- phishing resistant alternative
- requires only one gesture by the user
- easy [enrollment](https://docs.zitadel.com/docs/manuals/user-profile) of the device during registration

![passwordless-windows-hello](https://user-images.githubusercontent.com/1366906/118765435-5d419780-b87b-11eb-95bf-55140119c0d8.gif)

### Admin Console

Use [Console](https://docs.zitadel.com/docs/manuals/introduction) or our [APIs](https://docs.zitadel.com/docs/apis/introduction) to setup organizations, projects and applications.

Register new applications
![OIDC-Client-Register](https://user-images.githubusercontent.com/1366906/118765446-62064b80-b87b-11eb-8b24-4f4c365b8c58.gif)

Delegate the right to assign roles to another organization
![projects_create_org_grant](https://user-images.githubusercontent.com/1366906/118766069-39cb1c80-b87c-11eb-84cf-f5becce4e9b6.gif)

Customize login and console with your design  
![private_labeling](https://user-images.githubusercontent.com/1366906/123089110-d148ff80-d426-11eb-9598-32b506f6d4fd.gif)

## Security

See the policy [here](./SECURITY.md)

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.


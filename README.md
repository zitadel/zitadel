<a href="https://zitadel.com#gh-dark-mode-only"><img src="./docs/static/headers/zitadel-header-dark.png" alt="Zitadel Header" /></a>
<a href="https://zitadel.com#gh-light-mode-only"><img src="./docs/static/headers/zitadel-header-light.png" alt="Zitadel Header" /></a>

# ZITADEL

**[🏡 website](https://zitadel.com) [💬 chat](https://zitadel.com/chat) [📞 contact](https://zitadel.com/contact/) [📋 guide](https://docs.zitadel.ch/docs/guides/overview) [🧑‍💻 api docs](https://docs.zitadel.ch/docs/apis/introduction) [❓user manuals](https://docs.zitadel.ch/docs/manuals/introduction)**

[![stable version](https://badgen.net/github/release/zitadel/zitadel/stable)](https://github.com/zitadel/zitadel/releases/latest)
[![license](https://badgen.net/github/license/zitadel/zitadel)](#license)
[![code coverage](https://badgen.net/codecov/c/github/zitadel/zitadel)](https://app.codecov.io/gh/zitadel/zitadel)
[![Go Report Card](https://goreportcard.com/badge/github.com/zitadel/zitadel)](https://goreportcard.com/report/github.com/zitadel/zitadel)
[![discord](https://badgen.net/discord/online-members/erh5Brh7jE)](https://zitadel.com/chat)
[![follow us](https://badgen.net/twitter/follow/zitadel)](https://twitter.com/zitadel)
<a href="https://www.certification.openid.net/plan-detail.html?public=true&plan=w3ddtJcy0tpHL"><img src="./docs/static/logos/oidc-cert.png" alt="OpenID certification" height="35px" width="auto" /></a>

---

TODO Place a video here

## Why ZITADEL

- [API-first](https://docs.zitadel.ch/docs/apis/introduction)
- Strong audit trail thanks to [event sourcing](https://docs.zitadel.ch/docs/concepts/eventstore)
- Actions to react on events with custom code
- [Private labeling](https://docs.zitadel.ch/docs/guides/customization/branding) for a uniform user experience
- [cockroach database](https://www.cockroachlabs.com/) is the only dependency

## Features

- Single Sign On (SSO)
- Passwordless with FIDO2
- Username / Password
- Multifactor authentication with OTP, U2F
- [Identity Brokering](https://docs.zitadel.ch/docs/guides/authentication/identity-brokering)
- [Machine-to-machine (JWT profile)](https://docs.zitadel.ch/docs/guides/authentication/serviceusers)
- Personal Access Tokens (PAT)
- Role Based Access Control (RBAC)
- Delegate role management to third-parties
- Self-registration including verification
- User self service
- [Service Accounts](https://docs.zitadel.ch/docs/guides/authentication/serviceusers)

## Getting started

TODO Link to quickstarts in the docs

## Quick starts

TODO Link to quickstarts in the docs

If your use case is missing please let us know [here](https://github.com/zitadel/zitadel/discussions/1717)

## Contribute

Details about how to contribute you can find in the [Contribution Guide](CONTRIBUTING.md)

### Client libraries

<!-- TODO: check other libraries -->

| Language | Client | API | Machine auth (\*) | Auth check (\*\*) | Thanks to the maintainers |
|----------|--------|--------------|----------|---------|---------------------------|
| .NET     | [zitadel-net](https://github.com/zitadel/zitadel-net) | GRPC | ✔️ | ✔️ | [buehler 👑](https://github.com/buehler) |
| Dart     | [zitadel-dart](https://github.com/zitadel/zitadel-dart) | GRPC | ✔️ | ❌ | [buehler 👑](https://github.com/buehler) |
| Elixir   | [zitadel_api](https://github.com/jshmrtn/zitadel_api) | GRPC | ✔️ | ✔️ | [jshmrtn 🙏🏻](https://github.com/jshmrtn) |
| Go       | [zitadel-go](https://github.com/zitadel/zitadel-go) | GRPC | ✔️ | ✔️ | ZITADEL |
| Rust     | [zitadel-rust](https://crates.io/crates/zitadel) | GRPC | ✔️ | ❌ | [buehler 👑](https://github.com/buehler) |
| JVM      | ❓ | ❓ | ❓ | | Maybe you? |
| Python   | ❓ | ❓ | ❓ | | Maybe you? |
| Javascript | ❓ | ❓ | ❓ | | Maybe you? |

(\*) Automatically authenticate service accounts with [JWT Profile](https://docs.zitadel.ch/docs/apis/openidoauth/grant-types#json-web-token-jwt-profile).  
(\*\*) Automatically check if the access token is valid and claims match

## Security

See the policy [here](./SECURITY.md)

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## Usage data

ZITADEL components send errors and usage data, so that we are able to identify code improvement potential. If you don't want to send this data or don't have an internet connection, add the following lines to your custom configuration:

```yaml
TODO: add proper configuration
```

Besides from errors that don't clearly come from misconfiguration or cli misuage, we send an inital event when any binary is started. This is a " invoked" event along with the flags that are passed to it, except secret values of course.

We only ingest operational data. Your ZITADEL workload data from the IAM application itself is never sent anywhere unless you chose to integrate other systems yourself.
<img src="./docs/static/logos/zitadel-logo-dark@2x.png" alt="Zitadel Logo" height="100px" width="auto" />

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Release](https://github.com/caos/zitadel/actions/workflows/zitadel.yml/badge.svg)](https://github.com/caos/zitadel/actions)
[![license](https://badgen.net/github/license/caos/zitadel/)](https://github.com/caos/zitadel/blob/main/LICENSE)
[![release](https://badgen.net/github/release/caos/zitadel/stable)](https://github.com/caos/zitadel/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/caos/zitadel)](https://goreportcard.com/report/github.com/caos/zitadel)
[![codecov](https://codecov.io/gh/caos/zitadel/branch/main/graph/badge.svg)](https://codecov.io/gh/caos/zitadel)

## What Is ZITADEL

**ZITADEL** is a "Cloud Native Identity and Access Management" solution built for the cloud era. ZITADEL uses a modern software stack consisting of [**Golang**](https://golang.org/), [**Angular**](https://angular.io/) and  [**CockroachDB**](https://www.cockroachlabs.com/) as sole storage and follows an event sourced pattern.

We built **ZITADEL** not only with the vision of becoming a great open source project but also as a superb platform to support developers building their applications, without need to handle secure user login and account management themselves.

## How Does It Work

We built **ZITADEL** around the idea that the IAM should be easy to deploy and scale. That's why we tried to reduce external systems as much as possible.
For example, **ZITADEL** is event sourced but it does not rely on a pub/sub system to function. Instead we built all the functionality right into one binary.
**ZITADEL** only needs [**Kubernetes**](https://kubernetes.io/) for orchestration and [**CockroachDB**](https://www.cockroachlabs.com/) as storage.

## Features of ZITADEL platform

* Authentication
  * OpenID Connect 1.0 Protocol (OP)
  * Username / Password
  * Machine-to-machine (JWT profile)
  * Passwordless with FIDO2
* Multifactor authentication with OTP, U2F
* Federation with OpenID Connect 1.0 Protocol (RP), OAuth 2.0 Protocol (RP)
* Authorization via Role Based Access Control (RBAC)
* Identity Brokering
* Delegation of roles to other organizations for self-management
* Strong audit trail for all IAM resources
* User interface for administration
* APIs for Management, Administration, and Authentication
* Policy configuration and enforcement
* Private Labeling

## Run ZITADEL anywhere

### Self-Managed

You can run an automatically operated **ZITADEL** instance on a CNCF compliant Kubernetes cluster of your choice:
- [CRD Mode on an existing k8s cluster](https://docs.zitadel.ch/docs/guides/installation/crd)
- [GitOps Mode on an existing k8s cluster](https://docs.zitadel.ch/docs/guides/installation/gitops)
- [GitOps Mode on VM/bare-metal](https://docs.zitadel.ch/docs/guides/installation/managed-dedicated-instance)  using [ORBOS](https://docs.zitadel.ch/docs/guides/installation/orbos)

### CAOS-Managed

- **ZITADEL Cloud:** [**ZITADEL.ch**](https://zitadel.ch) is our shared cloud service hosted in Switzerland. [Get started](https://docs.zitadel.ch/docs/guides/usage/get-started) and try the free tier, including already unlimited users and all necessary security features.
- **ZITADEL Enterprise:** We operate and support a private instance of **ZITADEL** for you. [Get in touch!](https://zitadel.ch/contact/)

## Start using ZITADEL

### Quickstarts

See our [Documentation](https://docs.zitadel.ch/docs/quickstarts/introduction) to get started with ZITADEL quickly. Let us know, if you are missing a language or framework in the [Q&A](https://github.com/caos/zitadel/discussions/1717).

### Client libraries
* [Go](https://github.com/caos/zitadel-go) client library
* [.NET](https://github.com/caos/zitadel-net) client library
* [Dart](https://github.com/caos/zitadel-dart) client library

## Help and Documentation

* [Documentation](https://docs.zitadel.ch)
* [Ask a question or share ideas](https://github.com/caos/zitadel/discussions)
* [Say hello](https://zitadel.ch/contact/)

## Showcase

### Passwordless Login
Use our login widget to allow easy and sucure access to your applications and enjoy all the benefits of passwordless (FIDO 2 / WebAuthN):
- works on all modern platforms, devices, and browsers
- phishing resistant alternative
- requires only one gesture by the user
- easy [enrollment](https://docs.zitadel.ch/docs/manuals/user-factors) of the device during registration

![passwordless-windows-hello](https://user-images.githubusercontent.com/1366906/118765435-5d419780-b87b-11eb-95bf-55140119c0d8.gif)
![passwordless-iphone](https://user-images.githubusercontent.com/1366906/118765439-5fa3f180-b87b-11eb-937b-b4acb7854086.gif)

### Admin Console
Use [Console](https://docs.zitadel.ch/docs/manuals/introduction) or our [APIs](https://docs.zitadel.ch/docs/apis/introduction) to setup organizations, projects and applications.

Register new applications
![OIDC-Client-Register](https://user-images.githubusercontent.com/1366906/118765446-62064b80-b87b-11eb-8b24-4f4c365b8c58.gif)

Delegate the right to assign roles to another organization
![projects_create_org_grant](https://user-images.githubusercontent.com/1366906/118766069-39cb1c80-b87c-11eb-84cf-f5becce4e9b6.gif)

Customize login and console with your design  
![private_labeling](https://user-images.githubusercontent.com/1366906/123089110-d148ff80-d426-11eb-9598-32b506f6d4fd.gif)


## How To Contribute

Details about how to contribute you can find in the [Contribution Guide](CONTRIBUTING.md)

## Security

See the policy [here](./SECURITY.md)

## Other CAOS Projects

* [**ORBOS**](https://github.com/caos/orbos/) - GitOps everything
* [**OIDC for GO**](https://github.com/caos/oidc) - OpenID Connect SDK (client and server) for Go
* [**ZITADEL Tools**](https://github.com/caos/zitadel-tools) - Go tool to convert  key file to privately signed JWT

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.


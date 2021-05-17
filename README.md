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

## How To Use It 

### ZITADEL Cloud

We provide a cloud service [**ZITADEL.ch**](https://zitadel.ch) where people can register their own organization. There is a **free tier** including unlimited users and all the security features you need.

### Run ZITADEL in the cloud or on-premise

**ZITADEL** is free open source software under [Apache 2.0](##License) managed by [CAOS](https://caos.ch). We provide our community access to ZITADEL releases at no cost and welcome all contributions.

You can run an automatically operated **ZITADEL** instance on a CNCF compliant Kubernetes cluster of your choice. You can do so by using [CRDs](https://docs.zitadel.ch/docs/guides/installation/crd), [GitOps](https://docs.zitadel.ch/docs/guides/installation/gitops) or on a [dedicated Kubernetes Cluster](https://docs.zitadel.ch/docs/guides/installation/managed-dedicated-instance) on various infrastructure providers using [ORBOS](https://docs.zitadel.ch/docs/guides/installation/orbos)

### Let us run ZITADEL for you

If  our cloud service or running **ZITADEL** on your own infrastructure does not work for you, we are happy to run a private instance of **ZITADEL** for you or provide you with our support services. [Get in touch!](https://zitadel.ch/contact/)

## Help and Documentation

* [Documentation](https://docs.zitadel.ch)
* [Ask a question or share ideas](https://github.com/caos/zitadel/discussions)
* [Say hello](https://zitadel.ch/contact/)

## How To Contribute

Details need to be announced, but feel free to contribute already. As long as you are okay with accepting to contribute under this projects OSS [License](./LICENSE) you are fine.

## Security

See the policy [here](./SECURITY.md)

## Other CAOS Projects

* [**ORBOS**](https://github.com/caos/orbos/) - GitOps everything
* [**OIDC for GO**](https://github.com/caos/oidc) - OpenID Connect SDK (client and server) for Go
* [**ZITADEL Tools**](https://github.com/caos/zitadel-tools) - Go tool to convert  key file to privately signed JWT

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

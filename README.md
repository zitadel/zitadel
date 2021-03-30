<img src="./site/static/logos/zitadel-logo-dark@2x.png" alt="Zitadel Logo" height="100px" width="auto" />

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Release](https://github.com/caos/zitadel/workflows/Release/badge.svg)](https://github.com/caos/zitadel/actions)
[![license](https://badgen.net/github/license/caos/zitadel/)](https://github.com/caos/zitadel/blob/master/LICENSE)
[![release](https://badgen.net/github/release/caos/zitadel/stable)](https://github.com/caos/zitadel/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/caos/zitadel)](https://goreportcard.com/report/github.com/caos/zitadel)
[![codecov](https://codecov.io/gh/caos/zitadel/branch/master/graph/badge.svg)](https://codecov.io/gh/caos/zitadel)

> This project is in a beta state and API might still change a bit

## What Is It

**ZITADEL** is a "Cloud Native Identity and Access Management" solution. All server side components are written in [**Go**](https://golang.org/) and the management interface, called **Console**, is written in [**Angular**](https://angular.io/).

We optimized **ZITADEL** for the usage as "service provider" IAM. By "service provider" we think of companies who build services for e.g SaaS cases. Often these companies would like to use an IAM where they can register their application and grant other people or companies the right to self manage a set of roles within that application.

## How Does It Work

We built **ZITADEL** around the idea that the IAM should be easy to deploy and scale. That's why we tried to reduce external systems as much as possible.
For example, **ZITADEL** is event sourced but it does not rely on a pub/sub system to function. Instead we built all the functionality right into one binary.
**ZITADEL** only needs [**Kubernetes**](https://kubernetes.io/) for orchestration and [**CockroachDB**](https://www.cockroachlabs.com/) as storage.

## Why Another IAM

In the past we already built a closed sourced IAM and tested multiple others. With most of them we had some issues, either technology, feature, pricing or transparency related in nature. For example we find the idea that security related features like **MFA** should not be hidden behind a paywall or a feature price.
One feature that we often missed, was a solid **audit trail** of all IAM resources. Most systems we saw so far either rely on simple log files or use a short retention for this.

## How To Use It

### Use our free tier

We provide a [shared-cloud ZITADEL](https://zitadel.ch) system where people can register their own organisation in a [free tier](https://zitadel.ch/pricing). Our free tier contains most functionality, especially all security relevant features.

### Running your own IAM

You can run an automatically operated ZITADEL instance to a Kubernetes cluster of your choice. You can do so by using [CRDs](https://docs.zitadel.ch/start#CRD_Mode_on_an_existing_Kubernetes_cluster), [GitOps](#GitOps_Mode_on_an_existing_Kubernetes_cluster) or on a dedicated Kubernetes Cluster using [ORBOS](#GitOps_Mode_on_dedicated_Kubernetes_Clusters_using_ORBOS)

## Give me some docs

Have a look at our constantly evolving docs page [docs.zitadel.ch](https://docs.zitadel.ch).

## How To Contribute

Details need to be announced, but feel free to contribute already. As long as you are okay with accepting to contribute under this projects OSS [License](##License) you are fine.

We already have documentation specific [guidelines](./site/CONTRIBUTING.md).

Howto develop ZITADEL: [contribute](./CONTRIBUTING.md)

## Security

See the policy [here](./SECURITY.md)

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

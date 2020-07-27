![ZITADEL](./docs/img/zitadel-logo-oneline-lightdesign@2x.png)

# ZITADEL

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Release](https://github.com/caos/zitadel/workflows/Release/badge.svg)](https://github.com/caos/zitadel/actions)
[![license](https://badgen.net/github/license/caos/zitadel/)](https://github.com/caos/zitadel/blob/master/LICENSE)
[![release](https://badgen.net/github/release/caos/zitadel/stable)](https://github.com/caos/zitadel/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/caos/zitadel)](https://goreportcard.com/report/github.com/caos/zitadel)
[![codecov](https://codecov.io/gh/caos/zitadel/branch/master/graph/badge.svg)](https://codecov.io/gh/caos/zitadel)

> This project is in a alpha state. The application will continue breaking until version 1.0.0 is released

## What Is It

`ZITADEL` is a Cloud Native Identity and Access Management solution. All serverside componentes are written in `Go` and the management interface, called `console`, is written in `Angular`.

We optimized ZITADEL for the usage as `service provider IAM`. By `service provider` we think of companies who build services for e.g SaaS cases. Often these companies would like to use a IAM where they can register there application and grant other people or companies the right to self manage a set of roles within that application.

## How Does It Work

We built `ZITADEL` around the idea that the IAM should be easy to deploy and scale. That's why we tried to reduce external systems as much as possible.
For example, `ZITADEL` is eventsourced but it does not rely on a pub/sub system to function. We instead built all the functionality right into one binary.
`ZITADEL` only needs `kubernetes` for orchestration and `cockroachdb` as storage.

## Why Another IAM

In the past we already built a closed sourced IAM and tested multiple others. With most of them we had some issues, either technology, feature, pricing or transparency related in nature. For example we find the idea that security related features like `MFA` should not be hidden behind a paywall or a feature price.
One feature that we often missed, was a solid `audit trail` of all IAM resources. Most systems we saw so far either rely on simple log files or use a short retention for this.

## How To Use It

### Use our free tier

Stay tuned, we will soon publish how you can register yourself a organisation in our cloud offering `zitadel.ch`.
Yes we have a free tier!

### Run your own IAM

Stay tuned, we will publish soon a guide how you can deploy a `hyperconverged` system with our automation tooling called `ORBOS`.
With `ORBOS` you will be able to run `ZITADEL` on `GCE` or `StaticMachines` within 20 minutes. To achieve this, `ORBOS` will boostrap a `kubernetes` cluster, install the platform components (logging, metrics, ingress, ...), start a secure `cockroach` cluster and run and operate the `ZITDADEL`.

The combination of the tools `ORBOS` and `ZITADEL` is what makes the operation easy and scalable.

See our progress [here](https://github.com/caos/orbos/pull/256)

## Give me some docs

This is work in progess but will change soon.

## How To Contribute

TBA

## Security

See the policy [here](./SECURITY.md)

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

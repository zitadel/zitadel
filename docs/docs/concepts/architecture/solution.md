---
title: ZITADEL's Deployment Architecture
sidebar_label: Deployment Architecture
---

## High Availability

ZITADEL can be run as high available system with ease. 
Since the storage layer takes the heavy lifting of making sure that data in synched across, server, data centers or regions.

Depending on your projects needs our general recommendation is to run ZITADEL and ZITADELs storage layer across multiple availability zones in the same region or if you need higher guarantees run the storage layer across multiple regions.
Consult the [CockroachDB documentation](https://www.cockroachlabs.com/docs/) for more details or use the [CockroachCloud Service](https://www.cockroachlabs.com/docs/cockroachcloud/create-an-account.html)
Alternatively you can run ZITADEL also with Postgres which is [Enterprise Supported](/docs/support/software-release-cycles-support#partially-supported).
Make sure to read our [Production Guide](/self-hosting/manage/production#prefer-postgresql) before you decide to use it.

## Scalability

ZITADEL can be scaled in a linear fashion in multiple dimensions.

- Vertical on your compute infrastructure
- Horizontal in a region
- Horizontal in multiple regions

Our customers can reuse the same already known binary or container and scale it across multiple server, data center and regions.
To distribute traffic an already existing proxy infrastructure can be reused. 
Simply steer traffic by path, hostname, IP address or any other metadata to the ZITADEL of your choice.

> To improve your service quality we recommend steering traffic by path to different ZITADEL deployments
> Feel free to [contact us](https://zitadel.com/contact/) for details

## Example Deployment Architecture

### Single Cluster / Region

A ZITADEL Cluster is a highly available IAM system with each component critical for serving traffic laid out at least three times.
As our storage layer (CockroachDB) relies on Raft, it is recommended to operate odd numbers of storage nodes to prevent "split brain" problems.
Hence our reference design for Kubernetes is to have three application nodes and three storage nodes.

> If you are using a serverless offering like Google Cloud Run you can scale ZITADEL from 0 to 1000 Pods without the need of deploying the node across multiple availability zones.

:::info
CockroachDB needs to be configured with locality flags to proper distribute data over the zones
:::

![Cluster Architecture](/img/zitadel_cluster_architecture.png)

### Multi Cluster / Region

To scale ZITADEL across regions it is recommend to create at least three cluster.
We recommend to run an odd number of storage clusters (storage nodes per data center) to compensate for "split brain" scenarios.
In our reference design we recommend to create one cluster per region or cloud provider with a minimum of three regions.

With this design even the outage of a whole data-center would have a minimal impact as all data is still available at the other two locations.

:::info
CockroachDB needs to be configured with locality flags to proper distribute data over the zones
:::

![Multi-Cluster Architecture](/img/zitadel_multicluster_architecture.png)

## Zero Downtime Updates

Since an Identity system tends to be a critical piece of infrastructure, the "in place zero downtime update" is a well needed feature.
ZITADEL is built in a way that upgrades can be executed without downtime by just updating to a more recent version.

The common update involves the following steps and do not need manual intervention of the operator:

- Keep the old version running
- Deploy the version in parallel to the old version
- The new version will start ...
  - by updating databases schemas if needed
  - participate in the leader election for background jobs
- As soon as the new version is ready to accept traffic it will signal this on the readiness endpoint `/debug/ready` 
- At this point your network infrastructure can send traffic to the new version

Users who use [Kubernetes/Helm](/docs/self-hosting/deploy/kubernetes) or serverless container services like Google Cloud Run can benefit from the fact the above process is automated.

:::info
As a good practice we recommend creating Database Backups prior to an update.
It is also recommend to read the release notes on GitHub before upgrading.
Since ZITADEL utilizes Semantic Versioning Breaking Changes of any kind will always increase the major version (e.g Version 2 would become Version 3).
:::

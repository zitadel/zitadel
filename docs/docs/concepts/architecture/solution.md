---
title: Zitadel's Deployment Architecture
sidebar_label: Deployment Architecture
---

## High Availability

Zitadel can be run as high available system with ease. 
Since the storage layer takes the heavy lifting of making sure that data in synched across, server, data centers or regions.

Depending on your projects needs our general recommendation is to run Zitadel across multiple availability zones in the same region or across multiple regions.
Make sure to read our [Production Guide](/docs/self-hosting/manage/production#prefer-postgresql) before you decide to use it.
Consult the [Postgres documentation](https://www.postgresql.org/docs/) for more details.

## Scalability

Zitadel can be scaled in a linear fashion in multiple dimensions.

- Vertical on your compute infrastructure
- Horizontal in a region
- Horizontal in multiple regions

Our customers can reuse the same already known binary or container and scale it across multiple server, data center and regions.
To distribute traffic an already existing proxy infrastructure can be reused. 
Simply steer traffic by path, hostname, IP address or any other metadata to the Zitadel of your choice.

> To improve your service quality we recommend steering traffic by path to different Zitadel deployments
> Feel free to [contact us](https://zitadel.com/contact/) for details

## Example Deployment Architecture

### Single Cluster / Region

A Zitadel Cluster is a highly available IAM system with each component critical for serving traffic laid out at least three times.
Our storage layer (Postgres) is built for single region deployments.
Hence our reference design for Kubernetes is to have three application nodes and one storage node.

> If you are using a serverless offering like Google Cloud Run you can scale Zitadel from 0 to 1000 Pods without the need of deploying the node across multiple availability zones.

![Cluster Architecture](/img/zitadel_cluster_architecture.png)

### Multi Cluster / Region

To scale Zitadel across regions it is recommend to create at least three clusters.
Each cluster is a fully independent ZITADEL setup.
To keep the data in sync across all clusters, we recommend using Postgres with read-only replicas as a storage layer.
Make sure to read our [Production Guide](/docs/self-hosting/manage/production#prefer-postgresql) before you decide to use it.
Consult the [Postgres documentation](https://www.postgresql.org/docs/current/high-availability.html) for more details.


![Multi-Cluster Architecture](/img/zitadel_multicluster_architecture.png)

## Zero Downtime Updates

Since an Identity system tends to be a critical piece of infrastructure, the "in place zero downtime update" is a well needed feature.
Zitadel is built in a way that upgrades can be executed without downtime by just updating to a more recent version.

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
Since Zitadel utilizes Semantic Versioning Breaking Changes of any kind will always increase the major version (e.g Version 2 would become Version 3).
:::

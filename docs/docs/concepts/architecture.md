---
title: ZITADEL Architecture
---

## Software Architecture

**ZITADEL** is built with two essential patterns. Eventsourcing and CQRS. Due to the nature of eventsourcing **ZITADEL** provides the unique capability to generate a strong audit trail of ALL the things that happen to its resources, without compromising on storage cost or audit trail length.

The combination with CQRS makes **ZITADEL** eventual consistent which, from our perspective is a great benefit. It allows us to build a SOR (Source of Records) which is the one single point of truth for all computed states. The SOR needs to be transaction safe to make sure all operations are in order.

Each **ZITADEL** contains all components of the IAM, from serving as API, rendering / serving GUI's, background processing of events and task or being a GITOPS style operator. This AiO (All in One) approach makes scaling from a single machine to a multi region (multi cluster) seamless.

![Software Architecture](/img/zitadel_software_architecture.png)

### Component Command Side

The **command handler** receives all operations who alter a IAM resource. For example if a user changes his name.
This information is then passed to **command validation** for processing of the business logic, for example to make sure that the user actually can change his name. If this succeeds all generated events are inserted into the eventstore when required all within one transaction.

- Transaction safety is a MUST
- Availability MUST be high

> When we classify this with the CAP theorem we would choose **Consistent** and **Available** but leave **Partition Tolerance** aside.

### Component Spooler

The spoolers job is it to keep a query view up-to-date or at least look that it does not have a too big lag behind the eventstore.
Each query view has its own spooler who is responsible to look for the events who are relevant to generate the query view. It does this by triggering the relevant projection.
Spoolers are especially necessary where someone can query datasets instead of single ids.

> The query side has the option to dynamically check the eventstore for newer events on a certain id, see query side for more information
> Each view can have exactly one spooler, but spoolers are dynamically leader elected, so even if a spooler crashes it will be replaced in a short amount of time.

### Component Query Side

The query handler receives all read relevant operations. These can either be query or simple `getById` calls.
When receiving a query it will proceed by passing this to the repository which will call the database and return the dataset.
If a request calls for a specific id the call will, most of the times, be revalidated against the eventstore. This is achieved by triggering the projection to make sure that the last sequence of a id is loaded into the query view.

- Easy to query
- Short response times (80%of queries below 100ms on the api server)
- Availability MUST be high

> When we classify this with the CAP theorem we would choose **Available** and **Performance** but leave **Consistent** aside
> TODO explain more here

### Component HTTP Server

The http server is responsible for serving the management GUI called **ZITADEL Console**, serving the static assets and as well rendering server side html (login, password-reset, verification, ...)

## Cluster Architecture

A **ZITADEL Cluster** is a highly available IAM system with each component critical for serving traffic laid out at least three times.
As our storage (CockroachDB) relies on Raft it is also necessary to always utilizes odd numbers to address for "split brain" scenarios.
Hence our reference design is to have three application nodes and three Storage Nodes.

If you deploy **ZITADEL** with our GITOPS Tooling [**ORBOS**](https://github.com/caos/orbos) we create 7 seven nodes. One management, three application and three storage nodes.

> You can horizontaly scale zitadel, but we recommend to use multiple cluster instead to reduce the blast radius from impacts to a single cluster

![Cluster Architecture](/img/zitadel_cluster_architecture.png)

## Multi Cluster Architecture

To scale **ZITADEL** is recommend to create smaller clusters, see cluster architecture and then create a fabric which interconnects the database.
In our reference design we recommend to create a cluster per cloud provider or availability zone and to group them into regions.

For example, you can run three cluster for the region switzerland. On with GCE, one with cloudscale and one with inventx.

With this design even the outage of a whole data-center would have a minimal impact as all data is still available at the other two locations.

> Cockroach needs to be configured with locality flags to proper distribute data over the zones
> East - West connectivity for the database can be solved at you discretion. We recommend to expose the public ips and run traffic directly without any VPN or Mesh
> Use MTLS in combination with IP Allowlist in the firewalls!

![Multi-Cluster Architecture](/img/zitadel_multicluster_architecture.png)

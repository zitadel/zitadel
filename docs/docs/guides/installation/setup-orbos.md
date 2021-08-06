---
title: ORBOS and ZITADEL setup by CAOS
---

We provide services to setup our ORBOS and ZITADEL with the operators also provided by us.

### In Scope

- Check prerequisites and architecture
- Setup of VMs, Loadbalancing and Kubernetes with [ORBOS](https://github.com/caos/orbos)
- Setup of in-cluster toolset with ORBOS, which includes monitoring and an API gateway (Ambassador)
- Installation and configuration of ZITADEL with the ZITADEL-operator
- Installation and configuration of CockroachDB with the Database-operator
- Functional testing of the ZITADEL instance
  
### Out of Scope
- Integration of external S3-storage or other types of storage
- Integration into internal monitoring and alerting
- Multi-cluster architecture deployments
- Changes for specific environments
- Performance testing
- Production deployment
- Application-side coding, configuration, or tuning
- Changes or configuration on assets used in ZITADEL
- Setting up or maintaining backup storage

### Prerequisites

- S3-storage for assets in ZITADEL
- S3-storage or Google Cloud Bucket for backups of the database
- [Prerequisites listed for a managed instance](/docs/guides/installation/managed-dedicated-instance)

### Deliverable

- Running Kubernetes
- Running toolset for monitoring and alerting
- Running CockroachDB
- Running ZITADEL
- Running backups for ZITADEL
  
### Time Estimate

12 hours

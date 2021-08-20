---
title: ZITADEL setup
---

We provide services to setup our ZITADEL with the operators also provided by us.

### In Scope

- Check prerequisites and architecture
- Installation and configuration of ZITADEL with the ZITADEL-operator
- Installation and configuration of CockroachDB with the Database-operator
- Functional testing of the ZITADEL instance
  
### Out of Scope
  
- Running multiple ZITADEL instances on the same cluster
- Integration into internal monitoring and alerting
- Multi-cluster architecture deployments
- DNS, Network and Firewall configuration
- Kubernetes configuration
- Changes for specific environments
- Performance testing
- Production deployment
- Application-side coding, configuration, or tuning
- Changes or configuration on assets used in ZITADEL
- Setting up or maintaining backup storage

### Prerequisites

- Running Kubernetes with possibility to deploy to namespaces caos-system and caos-zitadel
- Volume provisioner for Kubernetes to fill Persistent Volume Claims
- S3-storage for assets in ZITADEL
- S3-storage or Google Cloud Bucket for backups of the database
- Inbound and outbound gRPC-Web traffic possible(for example not natively supported by nginx)
- [Prerequisites listed for a managed instance, limited to functionality for ZITADEL](/docs/guides/installation/managed-dedicated-instance)

### Deliverable
  
- Running CockroachDB
- Running ZITADEL
- Running backups for ZITADEL

### Time Estimate
  
8 hours

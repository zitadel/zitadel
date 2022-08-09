---
title: Onboarding Project
---

Effort required during depends on the complexity of your infrastructure and the overall setup. With a Multi-Zone Setup (excl. Multi-Region), support during this phase requires around 10-25h over 2 weeks. Actual effort is based on time and material.

Scope of the project is agreed on individual basis.

## In Scope

- Check prerequisites and architecture
- Troubleshoot installation and configuration of ZITADEL
- Troubleshoot and configuration connectivity to the database
- Functional testing of the ZITADEL instance

## Out of Scope

- Running multiple ZITADEL instances on the same cluster
- Integration into internal monitoring and alerting
- Multi-cluster architecture deployments
- DNS, Network and Firewall configuration
- Customer-specific Kubernetes configuration needs
- Changes for specific environments
- Performance testing
- Production deployment
- Application-side coding, configuration, or tuning
- Changes or configuration on assets used in ZITADEL
- Setting up or maintaining backup storage

## Prerequisites

- Running Kubernetes with possibility to deploy to namespaces
- Inbound and outbound HTTP/2 traffic possible
- For being able to send SMS, we need a Twilio sender name, SID and an auth token
- ZITADEL also needs to connect to an email relay of your choice. We need the SMTP host, user and app key as well as the ZITADEL emails sender address and name.

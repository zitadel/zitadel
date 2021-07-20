---
title: Annex to the Dedicated Instance Terms
custom_edit_url: null
--- 
## Introduction

This annex to the [Dedicated Instance Terms](terms-of-service-dedicated) describes the dedicated instance services and guarantees under different configurations.

Last revised: July 20, 2021

## Overview

### Service differences

The following table compares the different services, based on the preferred provider (Google Cloud). If you choose a different provider than our preferred provider the [Gurantees](#guarantees) stated in this document apply.

Service Levels

Service / Feature / Guarantee | ZITADEL Cloud FORTRESS | ZITADEL Dedicated Standard | ZITADEL Dedicated Advanced
--- | --- | --- | ---
Monitoring | 24x7 | 24x7 | 24x7
[Availability Objective](service-level-description#availability-objective) | 99.95% | 99.5% | 99.9%
Performance | up to [rate limits](rate-limit-policy#what-rate-limits-do-apply) | up to [rate limits](rate-limit-policy#what-rate-limits-do-apply) | up to [rate limits](rate-limit-policy#what-rate-limits-do-apply)
[Support hours](support-services#description-of-services) | Business | Business | Extended
[Response time (Sev 1)](support-services#slo---initial-response-time) | 1h | 2h | 1h
[Technical account manager](support-services#technical-account-manager) | n/a | n/a | 2h / week

High-availability configuration

Service / Feature / Guarantee | ZITADEL Cloud FORTRESS | ZITADEL Dedicated Standard | ZITADEL Dedicated Advanced
--- | --- | --- | ---
Multi-zone HA | yes | yes | yes
Geographic HA | yes | option | option
Multi-provider HA | yes | option | option

Upgrade and backup schedule

Service / Feature / Guarantee | ZITADEL Cloud FORTRESS | ZITADEL Dedicated Standard | ZITADEL Dedicated Advanced
--- | --- | --- | ---
Update flexibility | no | no | yes
Backup flexibility | no | yes | yes

Security

Service / Feature / Guarantee | ZITADEL Cloud FORTRESS | ZITADEL Dedicated Standard | ZITADEL Dedicated Advanced
--- | --- | --- | ---
DDOS Protection | yes | option | option
Strict TLS | yes | yes | yes
Web Application Firewall | yes | option | option
DNS Protection | yes | no, bespoke | no, bespoke
DNSSEC | yes | no, bespoke | no, bespoke

Features

Service / Feature / Guarantee | ZITADEL Cloud FORTRESS | ZITADEL Dedicated Standard | ZITADEL Dedicated Advanced
--- | --- | --- | ---
Audit log retention | 13 months | unlimited | unlimited
Tenancy | shared | dedicated | dedicated
Data region | CH | custom | custom
Data processing | CH | custom | custom

## Guarantees

### Infrastructure Provider

CAOS offers the following guarantees for a given infrastructure provider and customer satisfies the [prerequisites](https://docs.zitadel.ch/docs/guides/installation/managed-dedicated-instance).

Guarantees | Google Cloud | Static / Other | Self-hosted
---|---|---|---
Maintained by CAOS | yes | yes, product only | no
24x7 monitoring | yes | yes, product only | yes, product only
Availability SLO | [up to 99.9%](service-level-description#availability-objective) | none | none
Performance SLO | up to [rate limits](https://docs.zitadel.ch/docs/legal/rate-limit-policy#what-rate-limits-do-apply) | none | none

### Backup

ZITADEL Cloud creates hourly backups. We do not guarantee recovery time objective. Recovery point objective is in the context of our [event-sourcing pattern](/docs/concepts/eventstore) not meaningful.

---
title: ZITADEL Cloud Egress IP Addresses
sidebar_label: Egress IPs
description: A list of static IP addresses ZITADEL Cloud uses to make outbound connections to your services.
---

When configuring your firewall or network security groups, you may need to allow traffic **from** ZITADEL Cloud to your internal infrastructure.

This page lists the static Egress (outgoing) IP addresses used by ZITADEL Cloud regions.

## When do I need this?

You need to allowlist these IP addresses if you use features where ZITADEL initiates a connection to your systems, such as:

* **[Actions V1](/docs/concepts/features/actions)**
* **[Actions V2](/docs/concepts/features/actions_v2)**
* **[Notification Providers](/docs/guides/manage/customize/notification-providers)**
* **[External Identity Providers](/docs/guides/integrate/identity-providers/introduction)**

## IP Addresses by Region

We recommend allowing the IP address corresponding to the region where your ZITADEL instance is hosted.

| Region | Egress IP Address |
| :--- | :--- |
| **Switzerland** | `34.65.158.196` |
| **Europe** | `34.107.19.72` |
| **United States** | `34.69.146.246` |
| **Australia** | `34.87.243.23` |

:::tip
To find out which region your ZITADEL Cloud instance is running in, check the [ZITADEL Customer Portal](https://zitadel.com/admin/instances).
:::
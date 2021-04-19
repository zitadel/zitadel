---
title: Managed Dedicated Instance
---

:::tip What I need
I'd like to simply use ZITADEL without having to take care of any operational tasks, yet keeping control over all its data.
:::

- On GCE
  - DNS
  - JSON Key to Authenticate with a service account named orbiter-system with roles
    - roles/compute.admin
    - roles/iap.tunnelResourceAccessor
    - roles/serviceusage.serviceUsageAdmin
    - roles/iam.serviceAccountUser
- On Cloudscale
  - DNS
  - API Token
- On Static Provider
  - DNS
  - Virtual Machines IP addresses
  - A passwordless sudo user called orbiter on all VMs
  - Bootstrap SSH Key to login as orbiter (are replaced by ORBITER)
  - A storage solution (contact us)

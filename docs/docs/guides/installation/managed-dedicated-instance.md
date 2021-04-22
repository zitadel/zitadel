---
title: Managed Dedicated Instance
---

:::tip What I need
I'd like to simply use ZITADEL without having to take care of any operational tasks, yet keeping control over all its data.
:::

CAOS bootstraps and maintains a new ZITADEL instance just for you. This includes its underlying infrastructure with Kubernetes on top of it as well as monitoring tools and an API gateway. Contact us at <hi@zitadel.ch> for purchasing ZITADEL Enterprise Cloud.

# Prerequisites

Depending on the infrastructure provider you choose, you need to ensure some prerequisites.

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs
  defaultValue="gce"
  values={[
    {label: 'Google Compute Engine', value: 'gce'},
    {label: 'Cloudscale', value: 'cs'},
    {label: 'Static Provider', value: 'static'},
  ]}>
  <TabItem value="gce">
    <ul>
        <li>
            A JSON key for the infrastructure management to authenticate with a service account named orbiter-system assigned with the following roles
            <ul>
                <li>roles/compute.admin</li>
                <li>roles/iap.tunnelResourceAccessor</li>
                <li>roles/serviceusage.serviceUsageAdmin</li>
                <li>roles/iam.serviceAccountUser</li>
            </ul>
        </li>
        <li>
            A JSON key for the backups storage to authenticate with a service account assigned with the role roles/storage.admin
            <ul>
                <li>roles/storage.admin</li>
            </ul>
        </li>
    </ul>
  </TabItem>
  <TabItem value="cs">
    <ul>
        <li>A Cloudscale API token</li>        
    </ul>
  </TabItem>
  <TabItem value="static">
    <ul>
        <li>A List of available CentOS 7 machines with their IP addresses</li>
        <li>A passwordless sudo user called orbiter on all machines</li>
        <li>A Bootstrap SSH key to login as orbiter (are replaced by ORBITER)</li>
        <li>A storage solution (contact us)</li>
    </ul>
  </TabItem>
</Tabs>

- We need you to point [four DNS subdomains](/docs/apis/domains) to the CAOS-generated IP address.
- For being able to send SMS, we need a Twilio sender name, SID and an auth token.
- ZITADEL also needs to connect to an email relay of your choice. We need the SMTP host, user and app key as well as the ZITADEL emails sender address and name.

If you give us a Cloudflare user, an API key and a user service key, we can also manage the DNS entries, wildcard certificate and the IP whitelisting automatically using the Cloudflare API.

Metrics, logs and traces are collected and monitored by CAOS.

---
title: External ZITADEL Access
sidebar_label: Instance Not Found
---

## Why do I get an "Instance not found" error?

Also, ZITADEL has the [concept of virtual instances](/concepts/structure/instance#multiple-virtual-instances).
It uses a requests Host header to determine which virtual instance to use.
This is useful for multi-tenancy and resource sharing, for example in SaaS scenarios.
For most cases however, ZITADEL should run on exactly one domain.

This guide assumes you are already familiar with [configuring ZITADEL](./configure).

## Standard Config

ZITADEL only serves requests sent to the expected protocol, host and port.
For local testing purposes, you can use following configuration:

```yaml
ExternalDomain: localhost
ExternalPort: 8080
ExternalSecure: false
```

For productive setups however, we recommend using HTTPS and a custom domain:

```yaml
ExternalDomain: 'zitadel.my.domain'
ExternalPort: 443
ExternalSecure: true
```

## Changing ExternalDomain, ExternalPort or ExternalSecure

You can change the ExternalDomain, ExternalPort and ExternalSecure configuration options at any time.
However, for ZITADEL to be able to pick up the changes, [you need to rerun ZITADELs setup phase](/self-hosting/manage/updating_scaling#the-setup-phase).

## Running ZITADEL behind a Reverse Proxy

If you run ZITADEL behind a reverse proxy, you need to ensure that it sends the correct request headers to ZITADEL.
The proxy must either ensure that
- the original *Host* header value is assigned to the *Forwarded* headers host directive.
- the original requests *Host* header value is unchanged by the proxy.

Check out the [reverse proxy configuration examples](/self-hosting/manage/reverseproxy/reverse_proxy) for more information.

## Organization Domains

Note that by default, you cannot access ZITADEL at an organizations domain.
Organization level domains [are intended for routing users by their login methods to their correct organization](/guides/solution-scenarios/domain-discovery).

However, if you want to access ZITADEL at an organization domain, [you can add additional domains using the System API](/apis/resources/system/system-service-add-domain#adds-a-domain-to-an-instance).
Be aware that you won't automatically have the organizations context when you access ZITADEL like this.

## Generated Subdomains

ZITADEL creates random subdomains for [each new virtual instance](/concepts/structure/instance#multiple-virtual-instances).
You can immediately access the ZITADEL Console an APIs using these subdomains without further actions.

## More Information

- [Check out the production-near loadbalancing example with Traefik](/self-hosting/deploy/loadbalancing-example)
- [Explore some concrete proxy configuration examples for ZITADEL using the domain 127.0.0.1.sslip.io](/self-hosting/manage/reverseproxy/reverse_proxy)


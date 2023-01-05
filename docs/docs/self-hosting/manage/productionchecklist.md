---
title: Production Checklist
---


To apply best practices to your production setup we created a step by step checklist you may wish to follow.

### infrastructure configuration

- [ ] make use of configmanagement such as Terraform to provision all of the below
- [ ] use version control to store the provisioning
- [ ] use a secrets manager to save your sensible informations
- [ ] reduce the manual interaction with your platform to an absolute minimum 
#### HA Setup
- [ ] High Availability for ZITADEL containers
  - [ ] use container orchestrator such as Kubernetes or
  - [ ] use serverless architecture such as Knative or a hyperscaler equivalent (e.g. CloudRun from Google)
  - [ ] separate `zitadel init` and `zitadel setup` for fast startup times when scaling ZITADEL
- [ ] High Availability for database 
  - [ ] follow the [Production Checklist](https://www.cockroachlabs.com/docs/stable/recommended-production-settings.html) for CockroachDB if you selfhost the database or use [CockroachDB cloud](https://www.cockroachlabs.com/docs/cockroachcloud/create-an-account.html)
  - [ ] configure backups on a regular basis for the Database
  - [ ] test a restore scenario before going live
  - [ ] secure database connections from outside your network and/or use an internal subnet for database connectivity
- [ ] High Availability for critical infrastructure components (depending on your setup)
  - [ ] Loadbalancer
  - [ ] [Reverse Proxies](https://zitadel.com/docs/self-hosting/manage/reverseproxy/reverse_proxy)
  - [ ] Web Application Firewall

#### Networking
- [ ] Use a Layer 7 Web Application Firewall to secure ZITADEL that supports **[HTTP/2](/docs/self-hosting/manage/http2)**
  - [ ] secure the access by IP if needed
  - [ ] secure the access by [ratelimits](https://zitadel.com/docs/apis/ratelimits) for specific endpoints (e.g. API vs frontend) to secure availability on high load
  - [ ] doublecheck your Firewall for IPv6 connectivity and change accordingly

### ZITADEL configuration
- [ ] configure a valid [SMTP Server](https://zitadel.com/docs/guides/manage/console/instance-settings#smtp) and test emails
- [ ] Add [Custom Branding](https://zitadel.com/docs/guides/manage/customize/branding) if required
- [ ] configure a valid [SMS Service](https://zitadel.com/docs/guides/manage/console/instance-settings#sms) such as Twilio if needed
- [ ] configure your privacy policy, terms of service and a help Link if needed
- [ ] secure your [masterkey](https://zitadel.com/docs/self-hosting/manage/configure)

### Security
- [ ] use a FQDN and a trusted valid certificate for external [TLS](https://zitadel.com/docs/self-hosting/manage/tls_modes#http2) connections
- [ ] make use of different Serviceaccounts to secure ZITADEL within your Hyperscaler or Kubernetes 
- [ ] make use of a CDN Service if needed to ease maintainability and Firewall/DNS/WAF configuration
- [ ] make use of a [security scanner](https://owasp.org/www-community/Vulnerability_Scanning_Tools) to test your application and cluster

### Monitoring
Use a appropriate Monitoring Solution to have an overview about your ZITADEL instance. In particular you may want to watch out for things like:

- [ ] CPU and Memory of ZITADEL and the database
- [ ] open database connections
- [ ] running instances of ZITADEL and the database
- [ ] latency of requests
- [ ] requests per second
- [ ] requests by URL/endpoint
- [ ] lifetime of TLS certificates
- [ ] ZITADEL and database logs
- [ ] ZITADEL [Metrics](https://zitadel.com/docs/apis/observability/metrics)

---
title: ZITADEL Production Checklist
sidebar_label: Production Checklist
---


To apply best practices to your production setup we created a step by step checklist you may wish to follow.

### Infrastructure Configuration

- [ ] Make use of configuration management tools such as Terraform to provision all of the below
- [ ] Use a secrets manager to store your confidential information
- [ ] Reduce the manual interaction with your platform to an absolute minimum

#### HA Setup

- [ ] High Availability for ZITADEL containers
  - [ ] Use a container orchestrator such as Kubernetes
  - [ ] Use serverless platform such as Knative or a hyperscaler equivalent (e.g. CloudRun from Google)
  - [ ] Split `zitadel init` and `zitadel setup` for fast start-up times when [scaling](/docs/self-hosting/manage/updating_scaling) ZITADEL
- [ ] High Availability for database
  - [ ] Follow [this guide](https://www.postgresql.org/docs/current/high-availability.html) to set up the database.
  - [ ] Configure logging
  - [ ] Configure timeouts
  - [ ] Configure backups on a regular basis for the database
  - [ ] Test the restore scenarios before going live
  - [ ] Secure database connections from outside your network and/or use an internal subnet for database connectivity
- [ ] High Availability for critical infrastructure components (depending on your setup)
  - [ ] Loadbalancer
  - [ ] [Reverse Proxies](/docs/self-hosting/manage/reverseproxy/reverse_proxy)
  - [ ] Web Application Firewall

#### Networking

- [ ] Use a Layer 7 Web Application Firewall to secure ZITADEL that supports **[HTTP/2](/docs/self-hosting/manage/http2)**
  - [ ] Limit the access by IP addresses if needed
  - [ ] Secure the access by rate limits for specific endpoints (e.g. API vs frontend) to secure availability on high load. See the [ZITADEL Cloud rate limits](/docs/legal/policies/rate-limit-policy) for reference.
  - [ ] Check that your firewall also filters IPv6 traffic

### ZITADEL configuration

- [ ] Configure a valid [SMTP Server](/docs/guides/manage/console/default-settings#smtp) and test the email delivery
- [ ] Add [Custom Branding](/docs/guides/manage/customize/branding) if required
- [ ] Configure a valid [SMS Service](/docs/guides/manage/console/default-settings#sms) such as Twilio if needed
- [ ] Configure your privacy policy, terms of service and a help Link if needed
- [ ] Keep your [masterkey](/docs/self-hosting/manage/configure) in a secure storage
- [ ] Declare and apply zitadel configuration using the zitadel terraform [provider](https://github.com/zitadel/terraform-provider-zitadel) 

### Security

- [ ] Ensure that your ZITADEL does not use [the default, example or _easy-to-guess_ credentials](/docs/self-hosting/manage/database#zitadel-credentials)
- [ ] Use a FQDN and a trusted valid certificate for external [TLS](/docs/self-hosting/manage/tls_modes#http2) connections
- [ ] Create service accounts for applications that interact with ZITADEL's APIs
- [ ] Make use of a CDN service to decrease the load for static assets served by ZITADEL
- [ ] Make use of a [security scanner](https://owasp.org/www-community/Vulnerability_Scanning_Tools) to test your application and deployment environment

### Monitoring

Use an appropriate monitoring solution to have an overview about your ZITADEL instance. In particular you may want to watch out for things like:

- [ ] CPU and memory of ZITADEL and the database
- [ ] Open database connections
- [ ] Running instances of ZITADEL and the database
- [ ] Latency of requests
- [ ] Requests per second
- [ ] Requests by URL/endpoint
- [ ] Lifetime of TLS certificates
- [ ] ZITADEL and database logs
- [ ] ZITADEL [metrics](/docs/apis/observability/metrics)

---
title: Custom Domain
---

# Run ZITADEL on a (Sub)domain of Your Choice

This guide assumes you are already familiar with [configuring ZITADEL](./configure).

You most probably need to configure these fields for making ZITADEL work on your custom domain.

## Standard Config

For security reasons, ZITADEL only serves requests sent to the expected protocol, host and port.
If not using localhost as ExternalDomain, ExternalSecure must be true and you need to serve the ZITADEL console over HTTPS.

```yaml
ExternalSecure: true
ExternalDomain: 'zitadel.my.domain'
ExternalPort: 443
```

## Database Initialization Steps Config

ZITADEL creates random subdomains for each instance created.
However, for the default instance, this is most probably not the desired behavior.
You can give the default instance a fixed custom domain here.

```yaml
S3DefaultInstance:
  CustomDomain: 'zitadel.my.domain'
```

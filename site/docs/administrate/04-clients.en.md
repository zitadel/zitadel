---
title: Clients
---

### What are clients

Clients are applications who share the same security context and interface with an "authorization server".
For example you could have a software project existing out of a web app and a mobile app, both of these application might consume the same roles because the end user might use both of them.

### Manage clients

Clients might use different protocols for integrating with an IAM. With ZITADEL it is possible to use OpenID Connect 1.0 / OAuth 2.0. In the future SAML 2.0 support is planned as well.

> Screenshot here

### Configure OpenID Connect 1.0 Client

To make configuration of a client easy we provide a wizard which generates a specification conferment setup.
The wizard can be skipped for people who are needing special settings.
For use cases where your configuration is not compliant we provide you a "dev mode" which disables conformance checks.

> Screenshot here

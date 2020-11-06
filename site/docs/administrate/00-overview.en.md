---
title: Overview
---

> All documentations are under active work and subject to change soon!

### Features

- Single-Sign On (SSO) and Single-Log Out (SLO) for
  - Web applications
    - Server Renderer
    - Single Page
  - Native Clients
    - Windows
    - MacOS
    - Linux
  - Mobile Clients
    - Android
    - iOS / iPadOS
- Bearer Tokens (JWT and opaque) to use with APIs
  - REST
  - GRPC
  - GraphQL
- Role Based Access Control (RBAC) with delegation to let organisations manage authorisations on their own
- OpenID Connect 1.0 (OIDC) support
- OAuth 2.0 support
- Identity Brokering
  - Federation with OIDC and OAuth 2.0 Identity Providers
  - Social Login
- Management Console for central management of your data
- Multi-factor Authentication
  - Support for TOTP/HOTP with any app, like authy, google authenticator, ...
  - U2F (CTAP1)
- Passwordless Authentication
  - WebAuthN (FIDO2 / CTAP2)
- User self-registration, recover password, email and phone verification, etc.
- Organisation self-registration, domain verification, policy management
- API's for easy integration in your application

### Concepts

With ZITADEL there are some key concepts some should be aware of before using it to secure your applications and services.
You find these definitions in the "What is/are..." heading of each resource.

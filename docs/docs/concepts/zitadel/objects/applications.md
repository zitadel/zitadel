---
title: Applications
---

# Application Types 

OAuth and therefore OIDC know three different client types, these types are all for user interaction. 
As a further option, we know a classic API.
- Web (Server-side web applications such as java, .net, ...)
- Native (native, mobile or desktop applications)
- User Agent (single page applications / SPA, generally JavaScript executed in the browser)
- API (OAuth Resource Server)

Depending on the app type you're trying to register, there are small differences.
But regardless of the app type we recommend using Proof Key for Code Exchange (PKCE).

Please read the following guide about the [different-client-profiles](../../../guides/authorization/oauth-recommended-flows#different-client-profiles) and why to use PKCE.

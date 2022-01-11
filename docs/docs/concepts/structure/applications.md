---
title: Applications
---


# Applications

Applications are the entry point to your project. Users either login into one of your clients and interact with them directly or use one of your API, maybe without even knowing. All applications share the roles and authorizations of their project.

## Application Types 

At the moment ZITADEL differs between three client types (with user interaction):
- Web (Server-side web applications such as java, .net, ...)
- Native (native, mobile or desktop applications)
- User Agent (single page applications / SPA, generally JavaScript executed in the browser)

As a fourth option there's the API (OAuth Resource Server), which generally has no direct user-interaction.

Depending on the app type registered, there are small differences in the possible settings.

Please read the following guide about the
[different-client-profiles](../../guides/authorization/oauth-recommended-flows#different-client-profiles).

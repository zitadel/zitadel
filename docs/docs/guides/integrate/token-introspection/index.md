---
title: Token introspection
sidebar_label: Token introspection
sidebar_position: 1
---

Token introspection is the process of checking whether an access token is valid and can be used to access protected resources.

You have an API that acts as an OAuth resource server and can be accessed by user-facing applications.
To validate an access token by calling the ZITADEL introspection API, you can use the JSON Web Token (JWT) Profile (recommended) or Basic Authentication for token introspection.
It's crucial to understand that the API is entirely separate from the front end.
The API shouldn’t concern itself with the token type received.
Instead, it's about how the API chooses to call the introspection endpoint, either through JWT Profile or Basic Authentication.
Many APIs assume they might receive a JWT and attempt to verify it based on signature or expiration.
However, with ZITADEL, you can send either a JWT or an opaque Bearer token from the client end to the API.
This flexibility is one of ZITADEL's standout features.

## API application

If you have an API that behaves as an OAuth resource server that can be accessed by user-facing applications and need to validate an access token by calling the ZITADEL introspection API, you can use the following methods to register these APIs in ZITADEL:

- [JSON Web Token (JWT) Profile (Recommended)](private-key-jwt.mdx)
- [Basic Authentication](./basic-auth.mdx)

## Service users

If there are client APIs or systems that need to access other protected APIs, these APIs or systems must be declared as [service users](/docs/concepts/structure/users).
A service user is not considered an application type in ZITADEL.
Read the introduction on how to [authenticate service users](../service-users/authenticate-service-users.md).

## Further references

- [Introspection API reference](/docs/apis/openidoauth/endpoints#token_endpoint)
- [JWT vs. opaque tokens](/docs/concepts/knowledge/opaque-tokens.md)
- [Python examples for securing an API and invoking it as a service user](https://github.com/zitadel/examples-api-access-and-token-introspection)

import DocCardList from '@theme/DocCardList';

<DocCardList />

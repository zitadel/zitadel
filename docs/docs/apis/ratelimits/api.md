---
title: API Endpoint Rate Limits
---

:::caution
This is subject to change
:::

## api.zitadel.ch

| Path                                                | Description              | Effectiv Limit             |
|-----------------------------------------------------|--------------------------|----------------------------|
| /oauth/v2/*                                         | Sum of all OAuth request | 1000 request per 1 min     |
| /oauth/v2/token                                     |                          | 10 request per 10 seconds  |
| /oauth/v2/introspect                                |                          | 100 request per 10 seconds |
| /oauth/v2/userinfo                                  |                          | 100 request per 10 seconds |
| /auth/v1/*                                          |                          | 10 request per 10 seconds  |
| /caos.zitadel.auth.api.v1.AuthService/*             |                          | 10 request per 10 seconds  |
| /management/v1/*                                    |                          | 250 request per 1 min      |
| /caos.zitadel.management.api.v1.ManagementService/* |                          | 250 request per 1 min      |

## issuer.zitadel.ch

| Path | Description                             | Effectiv Limit |
|------|-----------------------------------------|----------------|
| /*   | Sum of all request to the issuer domain | none           |

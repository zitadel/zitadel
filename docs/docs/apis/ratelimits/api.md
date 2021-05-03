---
title: API Endpoint Rate Limits
---

:::caution
This is subject to change
:::

## api.zitadel.ch

| Path                                                | Description            | Effectiv Limit         |
|-----------------------------------------------------|------------------------|------------------------|
| /*                                                  | Sum of all API request | 1000 request per 1 min |
| /oauth/v2                                           |                        |                        |
| /auth/v1/*                                          |                        |                        |
| /caos.zitadel.auth.api.v1.AuthService/*             |                        |                        |
| /management/v1/*                                    |                        |                        |
| /caos.zitadel.management.api.v1.ManagementService/* |                        |                        |

## issuer.zitadel.ch

| Path | Description                             | Effectiv Limit |
|------|-----------------------------------------|----------------|
| /*   | Sum of all request to the issuer domain | none           |

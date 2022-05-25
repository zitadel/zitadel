---
title: API Endpoint Rate Limits
---


## api.zitadel.ch

| Path                                                | Description              | Effective Limit            |
|-----------------------------------------------------|--------------------------|----------------------------|
| /oauth/v2/*                                         | Sum of all OAuth request | 500 request per 1 min      |
| /oauth/v2/oauth/token                                     |                          | 120 request per 1 min      |
| /auth/v1/*                                          |                          | none                       |
| /caos.zitadel.auth.api.v1.AuthService/*             |                          | none                       |
| /management/v1/*                                    |                          | 240 request per 1 min      |
| /caos.zitadel.management.api.v1.ManagementService/* |                          | 240 request per 1 min      |

## issuer.zitadel.ch

| Path | Description                             | Effective Limit |
|------|-----------------------------------------|-----------------|
| /*   | Sum of all request to the issuer domain | none            |

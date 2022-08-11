---
title: API Rate Limits
---

<!-- //TODO Elio please update according to the current config -->

| Path                                                | Description              | Effective Limit            |
|-----------------------------------------------------|--------------------------|----------------------------|
| /oauth/v2/*                                         | Sum of all OAuth request | 500 request per 1 min      |
| /oauth/v2/token                                     |                          | 120 request per 1 min      |
| /auth/v1/*                                          |                          | none                       |
| /caos.zitadel.auth.api.v1.AuthService/*             |                          | none                       |
| /management/v1/*                                    |                          | 240 request per 1 min      |
| /caos.zitadel.management.api.v1.ManagementService/* |                          | 240 request per 1 min      |

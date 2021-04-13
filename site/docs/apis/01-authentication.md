---
title: Authentication API
description: …
---

### Authentication aka Auth

The authentication API (aka Auth API) is used for all operations on the currently logged in user.

| Service | URI                                                                                                                         |
|:--------|:----------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/auth/v1/](https://api.zitadel.ch/auth/v1/)                                                          |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService/](https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService) |

> At a later date we might expose functions to build your own login GUI
> You can build your own user Register GUI already by utilizing the [Management API](#management)

[Latest API Version](https://github.com/caos/zitadel/blob/main/proto/zitadel/auth.proto)

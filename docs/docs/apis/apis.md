---
title: ZITADEL APIs
---

## Authentication API aka Auth

The authentication API (aka Auth API) is used for all operations on the currently logged in user.

| Service | URI                                                                                                                         |
|:--------|:----------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/auth/v1/](https://api.zitadel.ch/auth/v1/)                                                          |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService/](https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService) |

> At a later date we might expose functions to build your own login GUI
> You can build your own user Register GUI already by utilizing the [Management API](#management)

[Latest API Version](https://github.com/caos/zitadel/blob/main/proto/zitadel/auth.proto)


## Management API

The management API is as the name states the interface where systems can mutate IAM objects like, organisations, projects, clients, users and so on if they have the necessary access rights.

| Service | URI                                                                                                                                                 |
|:--------|:----------------------------------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/management/v1/](https://api.zitadel.ch/management/v1/)                                                                      |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.management.api.v1.ManagementService/](https://api.zitadel.ch/caos.zitadel.management.api.v1.ManagementService) |

[Latest API Version](https://github.com/caos/zitadel/blob/main/proto/zitadel/management.proto)


## Administration API aka Admin

This API is intended to configure and manage the IAM itself.

| Service | URI                                                                                                                             |
|:--------|:--------------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/admin/v1/](https://api.zitadel.ch/admin/v1/)                                                            |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.admin.api.v1.AdminService/](https://api.zitadel.ch/caos.zitadel.admin.api.v1.AdminService) |

[Latest
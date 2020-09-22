---
title: Overview
description: …
---

> All documentations are under active work and subject to change soon!

### APIs

ZITADEL provides three API's for different use cases. These API's are built with GRPC and then generate a REST service.
We recommend that you use GRPC to integrate with ZITADEL as primary means.

Each services proto definition is located in the source control. 
As we generate the REST services and Swagger file out of the proto definition we recommend that you rely on the proto file. 
We annotate the corresponding REST methods on each possible call as well as the AuthN and AuthZ requirements.

See below for an example with the call **GetMyUser**.

```Proto
//User
rpc GetMyUser(google.protobuf.Empty) returns (UserView) {
option (google.api.http) = {
    get: "/users/me"
};

option (caos.zitadel.utils.v1.auth_option) = {
    permission: "authenticated"
};
}
```

| Service | URI                                                                                                                                            |
|:--------|:-----------------------------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/auth/v1/users/me](https://api.zitadel.ch/auth/v1/users/me)                                                             |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService/GetMyUser](https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService/GetMyUser) |


#### Authentication aka Auth

The authentication API (aka Auth API) is used for all operations on the currently logged in user.

| Service | URI                                                                                                                         |
|:--------|:----------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/auth/v1/](https://api.zitadel.ch/auth/v1/)                                                          |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService/](https://api.zitadel.ch/caos.zitadel.auth.api.v1.AuthService) |

> At a later date we might expose functions to build your own login GUI
> You can build your own user Register GUI already by utilizing the [Management API](#management)

[Latest API Version](https://github.com/caos/zitadel/blob/master/pkg/grpc/auth/proto/auth.proto)

#### Management

The management API is as the name states the interface where systems can mutate IAM objects like, organisations, projects, clients, user and so on if they have the necessary access rights.

| Service | URI                                                                                                                                                 |
|:--------|:----------------------------------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/management/v1/](https://api.zitadel.ch/management/v1/)                                                                      |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.management.api.v1.ManagementService/](https://api.zitadel.ch/caos.zitadel.management.api.v1.ManagementService) |

[Latest API Version](https://github.com/caos/zitadel/blob/master/pkg/grpc/management/proto/management.proto)

#### Administration aka Admin

This API is intended to configure and manage the IAM itself.

| Service | URI                                                                                                                             |
|:--------|:--------------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/admin/v1/](https://api.zitadel.ch/admin/v1/)                                                            |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.admin.api.v1.AdminService/](https://api.zitadel.ch/caos.zitadel.admin.api.v1.AdminService) |

[Latest API Version](https://github.com/caos/zitadel/blob/master/pkg/grpc/admin/proto/admin.proto)
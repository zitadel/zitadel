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

```Go
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

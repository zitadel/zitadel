---
title: Go
---

This integration guide shows you how to integrate **ZITADEL** into your Go application.
It demonstrates how to fetch some data from the ZITADEL management API.

At the end of the guide you should have an application able read the details of your organisation.

## Prerequisites

The client [SDK](https://github.com/caos/zitadel-go) will handle all necessary OAuth 2.0 requests and send the required headers to the ZITADEL API using our [OIDC client library](https://github.com/caos/oidc).
All that is required, is a service account with an Org Owner role assigned and its key JSON.

However, we recommend you read the guide on [how to access ZITADEL API](../guides/access-zitadel-apis) and the associated guides for a basic knowledge of :
 - <a href="../guides/oauth-recommended-flows">Recommended Authorization Flows</a>
 - <a href="../guides/serviceusers">Service Users</a>

> Be sure to have a valid key JSON and that its service account is either ORG_OWNER or at least ORG_OWNER_VIEWER before you continue with this guide.

## Go Setup

### Add Go SDK to your project

You need to add the SDK into Go Modules by:

```bash
go get github.com/caos/zitadel-go
```

### Create example client

Create a new go file with the content below. This will create a client for the management api and call its `GetMyOrg` function.
The SDK will make sure you will have access to the API by retrieving a Bearer Token using JWT Profile with the provided scopes (`openid` and `urn:zitadel:iam:org:project:id:69234237810729019:aud`).

```go
package main

import (
    "context"
    "log"
    
    "github.com/caos/oidc/pkg/oidc"
    
    "github.com/caos/zitadel-go/pkg/client/management"
    "github.com/caos/zitadel-go/pkg/client/middleware"
    "github.com/caos/zitadel-go/pkg/client/zitadel"
    pb "github.com/caos/zitadel-go/pkg/client/zitadel/management"
)


func main() {
    client, err := management.NewClient(
        []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
    )
    if err != nil {
        log.Fatalln("could not create client", err)
    }
    defer func() {
        err := client.Connection.Close()
        if err != nil {
            log.Println("could not close grpc connection", err)
        }
    }()
    
    ctx := context.Background()
    resp, err := client.GetMyOrg(ctx, &pb.GetMyOrgRequest{})
    if err != nil {
        log.Fatalln("call failed: ", err)
    }
    log.Printf("%s was created on: %s", resp.Org.Name, resp.Org.Details.CreationDate.AsTime())
}
```

#### Key JSON

To provide the key JSON to the SDK, simply set an environment variable `ZITADEL_KEY_PATH` with the path to the JSON as value.

```bash
export ZITADEL_KEY_PATH=/Users/test/servicekey.json
```

For development purposes you should be able to set this in your IDE.

If you're not able to set it via environment variable, you can also pass it with an additional option:

```go
client, err := management.NewClient(
    []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
    zitadel.WithKeyPath("/Users/test/servicekey.json"),
)
```

#### Custom ZITADEL instance

If your client will not use ZITADEL Cloud (zitadel.ch), be sure to provide the correct values for the ZITADEL ProjectID, Issuer and API options:
```go
client, err := management.NewClient(
    []string{oidc.ScopeOpenID, zitadel.ScopeProjectID("ZITADEL-ProjectID")},
    zitadel.WithCustomURL("https://issuer.custom.ch", "api.custom.ch:443")
)
```

### Test client

After you have configured everything correctly, you can simply start the example by:

```bash
go run main.go
```

This will output something similar to:

```
2021/04/21 11:27:36 DemoOrg was created on: 2021-04-08 13:36:05.578194 +0000 UTC
```

## Completion

You have successfully used the ZITADEL Go SDK to call the management API!

If you encountered an error (e.g. `code = PermissionDenied desc = No matching permissions found`), 
ensure your service user has the required permissions by assigning the `ORG_OWNER` or `ORG_OWNER_VIEWER` role
and check the mentioned [guides](#prerequisites) at the beginning.

If you've run into any other problem, don't hesitate to contact us or raise an issue on [ZITADEL](https://github.com/caos/zitadel/issues) or in the [SDK](https://github.com/caos/zitadel-go/issues).

### Whats next?

Now you can proceed implementing our APIs by adding more calls or trying to overwrite the organisation context:

```go
    respOverwrite, err := client.GetMyOrg(middleware.SetOrgID(ctx, "74161146763996133"), &pb.GetMyOrgRequest{})
    if err != nil {
        log.Fatalln("call failed: ", err)
    }
    log.Printf("%s was created on: %s", respOverwrite.Org.Name, respOverwrite.Org.Details.CreationDate.AsTime())
}
```
Checkout more [examples from the SDK](https://github.com/caos/zitadel-go/blob/main/example) or refer to our [API Docs](../apis/introduction).

> This guide will be updated soon to show you how to use the SDK for your own API as well.

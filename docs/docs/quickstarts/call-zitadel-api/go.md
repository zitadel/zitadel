---
title: Go
---

This guide shows you how to integrate **ZITADEL** into your Go application.
It demonstrates how to fetch some data from the ZITADEL management API.

At the end of the guide, you should have an application that can read the details of your organization.

## Prerequisites

The client [SDK](https://github.com/caos/zitadel-go) handles all necessary OAuth 2.0 requests and sends the required headers to the ZITADEL API using our [OIDC client library](https://github.com/caos/oidc).

You'll need a service account assigned with the Org-owner role. 
(or another role, depending on the needed API requests).
You'll also need the account's service key in a JSON file.

For background information, we recommend reading the guide on [how to access ZITADEL API](../../guides/API/access-zitadel-APIs) and the associated guides for a basic knowledge of :
 - [Recommended Authorization Flows](../../guides/authorization/oauth-recommended-flows)
 - [Service Users](../../guides/authentication/serviceusers)

> Be sure that you have a valid JSON key and that its service account is either `ORG_OWNER` or at least `ORG_OWNER_VIEWER`.

## Go Setup

### Add Go SDK to your project

Add the SDK into Go Modules with this command:

```bash
go get github.com/caos/zitadel-go
```

### Create example client

Create a new go file with the following snippet.
This creates a client for the management API and calls its `GetMyOrg` function.

To make sure you can access the API,
the SDK retrieves a Bearer Token using a JWT Profile with the provided scopes (`openid` and `urn:zitadel:iam:org:project:id:69234237810729019:aud`).

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

#### JSON

To provide the JSON key to the SDK, simply set an environment variable `ZITADEL_KEY_PATH`, using the path to the JSON as the value.

```bash
export ZITADEL_KEY_PATH=/Users/test/servicekey.json
```

For development purposes, you should be able to set this in your IDE.

If you're can't set it via environment variable, you can also pass it with an additional option:

```go
client, err := management.NewClient(
    []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
    zitadel.WithKeyPath("/Users/test/servicekey.json"),
)
```

#### Custom ZITADEL instance

If your client does not use ZITADEL Cloud (zitadel.ch), be sure to provide the correct values for the ZITADEL `ProjectID`, `Issuer` and `API` options:

```go
client, err := management.NewClient(
    []string{oidc.ScopeOpenID, zitadel.ScopeProjectID("ZITADEL-ProjectID")},
    zitadel.WithCustomURL("https://issuer.custom.ch", "API.custom.ch:443")
)
```

### Test client

After correctly configuring everything, start the example with this command:

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
make sure your service user has the required permissions.
The service user needs the `ORG_OWNER` or `ORG_OWNER_VIEWER` role.

For more help, check the [guides](#prerequisites) mentioned at the beginning.

If you've run into any other problems, don't hesitate to contact us or raise an issue on [ZITADEL](https://github.com/caos/zitadel/issues) or in the [SDK](https://github.com/caos/zitadel-go/issues).

### Whats next?

Now develop on our APIs by adding more calls.
You can also try to overwrite the organization context:

```go
    respOverwrite, err := client.GetMyOrg(middleware.SetOrgID(ctx, "74161146763996133"), &pb.GetMyOrgRequest{})
    if err != nil {
        log.Fatalln("call failed: ", err)
    }
    log.Printf("%s was created on: %s", respOverwrite.Org.Name, respOverwrite.Org.Details.CreationDate.AsTime())
}
```
Checkout more [examples from the SDK](https://github.com/caos/zitadel-go/blob/main/example) or refer to our [API Docs](../../APIs/introduction).

> This guide will be updated soon to show you how to use the SDK for your own API as well.

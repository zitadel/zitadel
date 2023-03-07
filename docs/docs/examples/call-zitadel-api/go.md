---
title: Go
---

This integration guide shows you how to integrate **ZITADEL** into your Go application.
It demonstrates how to fetch some data from the ZITADEL management API.

At the end of the guide you should have an application able to read the details of your organization.

## Prerequisites

The client [SDK](https://github.com/zitadel/zitadel-go) will handle all necessary OAuth 2.0 requests and send the required headers to the ZITADEL API using our [OIDC client library](https://github.com/zitadel/oidc).
All that is required, is a service account with an Org Owner (or another role, depending on the needed api requests) role assigned and its key JSON.

However, we recommend you read the guide on [how to access ZITADEL API](../../guides/integrate/access-zitadel-apis) and the associated guides for a basic knowledge of :
 - [Recommended Authorization Flows](../../guides/integrate/oauth-recommended-flows.md)
 - [Service Users](../../guides/integrate/serviceusers.md)

> Be sure to have a valid key JSON and that its service account is either ORG_OWNER or at least ORG_OWNER_VIEWER before you continue with this guide.

## Go Setup

### Add Go SDK to your project

You need to add the SDK into Go Modules by:

```bash
go get github.com/zitadel/zitadel-go/v2
```

### Create example client

Create a new go file with the content below. This will create a client for the management api and call its `GetMyOrg` function.
The SDK will make sure you will have access to the API by retrieving a Bearer Token using JWT Profile with the provided scopes (`openid` and `urn:zitadel:iam:org:project:id:zitadel:aud`).
Make sure to fill the vars `issuer` and `api`.

The issuer and api is the domain of your instance you can find it on the instance detail in the ZITADEL Cloud Customer Portal or in the ZITADEL Console.

:::note
The issuer will require the protocol (`https://` and `http://`) and you will only have to specify a port if they're not default (443 for https and 80 for http). The API will always require a port, but no protocol. 
:::

```go
package main

import (
	"context"
	"flag"
	"log"

	"github.com/zitadel/oidc/pkg/oidc"

	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

var (
	issuer = flag.String("issuer", "", "issuer of your ZITADEL instance (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
	api    = flag.String("api", "", "gRPC endpoint of your ZITADEL instance (in the form: <instance>.zitadel.cloud:443 or <yourdomain>:443)")
)

func main() {
	flag.Parse()

	//create a client for the management api providing:
	//- issuer (e.g. https://acme-dtfhdg.zitadel.cloud)
	//- api (e.g. acme-dtfhdg.zitadel.cloud:443)
	//- scopes (including the ZITADEL project ID),
	//- a JWT Profile token source (e.g. path to your key json), if not provided, the file will be read from the path set in env var ZITADEL_KEY_PATH
	client, err := management.NewClient(
		*issuer,
		*api,
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

	//call ZITADEL and print the name and creation date of your organisation
	//the call was successful if no error occurred
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

If you've run into any other problem, don't hesitate to contact us or raise an issue on [ZITADEL](https://github.com/zitadel/zitadel/issues) or in the [SDK](https://github.com/zitadel/zitadel-go/issues).

### Whats next?

Now you can proceed implementing our APIs by adding more calls or trying to overwrite the organization context:

```go
    respOverwrite, err := client.GetMyOrg(middleware.SetOrgID(ctx, "74161146763996133"), &pb.GetMyOrgRequest{})
    if err != nil {
        log.Fatalln("call failed: ", err)
    }
    log.Printf("%s was created on: %s", respOverwrite.Org.Name, respOverwrite.Org.Details.CreationDate.AsTime())
}
```
Checkout more [examples from the SDK](https://github.com/zitadel/zitadel-go/blob/main/example) or refer to our [API Docs](/apis/introduction).

> This guide will be updated soon to show you how to use the SDK for your own API as well.

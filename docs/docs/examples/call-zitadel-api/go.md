---
title: Integrate ZITADEL into a Go Application
sidebar_label: Go
---

This integration guide shows you how to integrate **ZITADEL** into your Go application.
It demonstrates how to fetch some data from the ZITADEL management API.

At the end of the guide you should have an application able to read the details of your organization.

> This documentation references our [CLI example](https://github.com/zitadel/zitadel-go/blob/next/example/client/cli/cli.go).

## Prerequisites

The client [SDK](https://github.com/zitadel/zitadel-go) will handle all necessary OAuth 2.0 requests and send the required headers to the ZITADEL API using our [OIDC client library](https://github.com/zitadel/oidc).
All that is required, is a service account with an Org Owner (or another role, depending on the needed api requests) role assigned and its key JSON.

However, we recommend you read the guide on [how to access ZITADEL API](../../guides/integrate/access-zitadel-apis) and the associated guides for a basic knowledge of :
 - [Recommended Authorization Flows](../../guides/integrate/oauth-recommended-flows.md)
 - [Service Users](../../guides/integrate/serviceusers)

> Be sure to have a valid key JSON and that its service account is either ORG_OWNER or at least ORG_OWNER_VIEWER before you continue with this guide.

## Go Setup

### Add Go SDK to your project

You need to add the SDK into Go Modules by:

```bash
go get -u github.com/zitadel/zitadel-go/v3
```

### Create example client

Create a new go file with the content below. This will create a client and call its `GetMyOrg` function on the ManagementService.
The SDK will make sure you will have access to the API by retrieving a Bearer Token using JWT Profile with the provided scopes (`openid` and `urn:zitadel:iam:org:project:id:zitadel:aud`).

```go reference
https://github.com/zitadel/zitadel-go/blob/next/example/client/cli/cli.go
```

### Test

After you have configured everything correctly, you can simply start the example by:

```bash
go run cli.go --domain <your domain> --key <path>
```

This could look like:

```bash
go run cli.go --domain my-domain.zitadel.cloud --key ./api.json
```

This will output something similar to:

```
2023/12/20 08:48:23 INFO retrieved the organisation orgID=165467338479501569 name=DemoOrg
```

## Completion

You have successfully used the ZITADEL Go SDK to call the management API!

If you encountered an error (e.g. `code = PermissionDenied desc = No matching permissions found`), 
ensure your service user has the required permissions by assigning the `ORG_OWNER` or `ORG_OWNER_VIEWER` role
and check the mentioned [guides](#prerequisites) at the beginning.

If you've run into any other problem, don't hesitate to contact us or raise an issue on [ZITADEL](https://github.com/zitadel/zitadel/issues) or in the [SDK](https://github.com/zitadel/zitadel-go/issues).

### Whats next?

Now you can proceed implementing our APIs by adding more calls or using a different service like the SessionService:

```go
api.SessionService().CreateSession(ctx, &session.CreateSessionRequest{})
```
Checkout more [examples from the SDK](https://github.com/zitadel/zitadel-go/blob/next/example),
like how you can integrate the [client in your own API](https://github.com/zitadel/zitadel-go/blob/next/example/api/client/main.go)
or refer to our [API Docs](/apis/introduction).


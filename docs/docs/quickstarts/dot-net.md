---
title: .NET
---

This integration guide shows you how to integrate **ZITADEL** into your .NET application.
It demonstrates how to fetch some data from the ZITADEL management API.

At the end of the guide you should have an application able read the details of your organisation.

If you need any other information about the .NET SDK go to the [documentation](https://caos.github.io/zitadel-net/) of the SDK itself.
## Prerequisites

The client [SDK](https://github.com/caos/zitadel-net) will handle all necessary OAuth 2.0 requests and send the required headers to the ZITADEL API.
All that is required, is a service account with an Org Owner (or another role, depending on the needed api requests) role assigned and its key JSON.

However, we recommend you read the guide on [how to access ZITADEL API](../guides/usage/access-zitadel-apis) and the associated guides for a basic knowledge of :
 - [Recommended Authorization Flows](../guides/usage/oauth-recommended-flows)
 - [Service Users](../guides/usage/serviceusers)

> Be sure to have a valid key JSON and that its service account is either ORG_OWNER or at least ORG_OWNER_VIEWER before you continue with this guide.

## .NET Setup

### Create a .NET application

Use the IDE of your choice or the command line to create a new application.

```bash
dotnet new web
```

### Install the package

Install the package via nuget

```bash
dotnet add package Zitadel.Api
```

### Create example client

Create a new go file with the content below. This will create a client for the management api and call its `GetMyOrg` function.
The SDK will make sure you will have access to the API by retrieving a Bearer Token using JWT Profile with the provided scopes (`openid` and `urn:zitadel:iam:org:project:id:69234237810729019:aud`).

```csharp
// no.. this key is not activated anymore ;-)
var sa = await ServiceAccount.LoadFromJsonFileAsync("./service-account.json");
var api = Clients.ManagementService(
    new()
    {
        // Which api endpoint (self hosted or public)
        Endpoint = ZitadelDefaults.ZitadelApiEndpoint,
        // The organization context (where the api calls are executed)
        Organization = "69234230193872955",
        // Service account authentication
        ServiceAccountAuthentication = (sa, new()
        {
            ProjectAudiences = { ZitadelDefaults.ZitadelApiProjectId },
        }),
    });

var roles = await api.SearchProjectRolesAsync(
    new() { ProjectId = "84856448403694484" });

foreach (var r in roles.Result)
{
    Console.WriteLine($"{r.Key} : {r.DisplayName} : {r.Group}");
}
```

#### Custom ZITADEL instance

If your client will not use ZITADEL Cloud (zitadel.ch), be sure to provide the correct values for the ZITADEL ProjectID, Issuer and API options:
```csharp

// Which api endpoint (self hosted or public)
Endpoint = "api.custom.ch:443",
// Service account authentication
ServiceAccountAuthentication = (sa, new()
{
    ProjectAudiences = { "ZITADEL-ProjectID" },
    Endpoint = "https://issuer.custom.ch",
}),

```

### Test client

After you have configured everything correctly, you can simply start the example by:

```bash
dotnet run
```

This will output something similar to:

```
ACME was created on: "2020-09-21T14:44:48.090431Z" 
```

## Completion

You have successfully used the ZITADEL .NET SDK to call the management API!

If you encountered an error (e.g. `code = PermissionDenied desc = No matching permissions found`), 
ensure your service user has the required permissions by assigning the `ORG_OWNER` or `ORG_OWNER_VIEWER` role
and check the mentioned [guides](#prerequisites) at the beginning.

If you've run into any other problem, don't hesitate to contact us or raise an issue on [ZITADEL](https://github.com/caos/zitadel/issues) or in the [SDK](https://github.com/caos/zitadel-go/issues).

### Whats next?

Now you can proceed implementing our APIs by adding more calls.

Checkout more [examples from the SDK](https://github.com/caos/zitadel-go/blob/main/example) or refer to our [API Docs](../apis/introduction).

> This guide will be updated soon to show you how to use the SDK for your own API as well.

---
title: .NET
---

This integration guide shows you how to integrate **ZITADEL** into your .NET application.
It demonstrates how to fetch some data from the ZITADEL management API.

At the end of the guide you should have an application able to read the details of your organization.

If you need any other information about the .NET SDK go to the [documentation](https://github.com/smartive/zitadel-net) of the SDK itself.

## Prerequisites

The client [SDK](https://github.com/zitadel/zitadel-net) will handle all necessary OAuth 2.0 requests and send the required headers to the ZITADEL API.
All that is required, is a service account with an Org Owner (or another role, depending on the needed api requests) role assigned and its key JSON.

However, we recommend you read the guide on [how to access ZITADEL API](../../guides/integrate/access-zitadel-apis) and the associated guides for a basic knowledge of :

 - [Recommended Authorization Flows](../../guides/integrate/oauth-recommended-flows.md)
 - [Service Users](../../guides/integrate/serviceusers.md)

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

Change the program.cs file to the content below. This will create a client for the management api and call its `GetMyUsers` function.
The SDK will make sure you will have access to the API by retrieving a Bearer Token using JWT Profile with the provided scopes (`openid` and `urn:zitadel:iam:org:project:id:{projectID}:aud`).

Make sure to fill the const `apiUrl`, `apiProject` and `personalAccessToken` with your own instance data. The used vars below are from a test instance, to show you how it should look.
The apiURL is the domain of your instance you can find it on the instance detail in the Customer Portal or in the Console
The apiProject you will find in the ZITADEL project in the first organization of your instance.

```csharp
// This file contains two examples:
// 1. An example with a service account "personal access token" to access the ZITADEL API.
// 2. An example with a service account "jwt profile key" to access the ZITADEL API.

using Zitadel.Api;
using Zitadel.Credentials;

const string apiUrl = "https://zitadel-libraries-l8boqa.zitadel.cloud";
const string personalAccessToken = "ge85fvmgTX4XAhjpF0XGpelB2vn9LZanJaqmUQDuf7iTpKVowb44LFl-86pqY2mfJCEoIOk";

// or create the token provider directly:
// new StaticTokenProvider(token)
var client = Clients.AuthService(new(apiUrl, ITokenProvider.Static(personalAccessToken)));
var result = await client.GetMyUserAsync(new());
Console.WriteLine($"User: {result.User}");

const string apiProject = "170078979166961921";
var serviceAccount = ServiceAccount.LoadFromJsonString(
    @"
{
  ""type"": ""serviceaccount"",
  ""keyId"": ""170084658355110145"",
  ""key"": ""-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAnQisbU4FuLmjLR9I2Q01Rm9Mx6WySat2mbxgmOzu04oXuESI\nyS+RkiimdN0khjqouBftYqtVes7yngMLq3E8hMCwv/kLE+YeXphZXnn8tps8M2gV\n7S//uCp9LooK9qeh0lSkOqIsh0atj/l7NAHFxnhuNhfmn8XIYJNLVNSj5yzTri5E\nSn92SAsUQLSONgr7IEmIjcuPtYeU0iLvVno52ljZHnPX2WJ0HEZv44nZpkR4qBfv\n3hJzNx7sd4TdPGHHugJD8jdG/X4bAxwL5XGHZu18cUVM5RerSMpFQHSuIGgpKmK4\nWlM1AJGeut6EX/SrCxUDvhyOnXAgqhunTUmi6QIDAQABAoIBAHn7y92Y1y743X3m\nqHMbJIBTYyRPXaCGljm0MKF6o8clpWlZq5wE3KLZ+vwa8Q1oMbnXtGqKR3t/mM4P\n9Ze2/djtyh9GOUm632qCFCIkxp+fFPOl7ipyt8V7FAT77KpP6490eqKlacunppmJ\nph/vJJAY6xwQEvGX9SC4KrN5/txLKXbVtR3V2RXy9sxbbL4cpnklmRBMeXQkpwEM\nTKELUr5Rmhg9KvS3yALgVv0dIRtOA8Z995R234hXfY0St48YEvZtsxeme47u2CVl\nHJcVH4aa9Sw6XlgAEQBxqbQHpcLvUIu3XempO7VfGklWE6OlGuEcnUWpJCD8jMZW\nPYtt9LUCgYEAwi8josS3Iyto+DMJjJKCw175N2cmFMxBGu9Rw4aHjTiN57z7AUkn\nbmT44WnSmc1bCLC+nMB34vhiEyBKXYrH7zgbeMO8QDG3aO6gXdod/IdsieZR8E3b\ngUA1wtZYyRbc7eo8U4Nqkv1NXVRuDJkz/Mfoy+m1BVKcW7YeZaaZN9MCgYEAzwYB\n/LAiJoyx5UPwuieizlT7kHI7uvZRo4oLx+cZipNCJ0NGKgX4l1NIYLaNDbCoT9N0\nylico+kn+nihzDmD6SjY2hHGSIHk7AnJOcW+Bk5TfsYb8clxfgX40udLMIS0F13R\nrJt0gD9x0O3AZv4MV9cSI0/Md0tbWePgrLI44NMCgYEAojj7TlmEnY8AbIlGqvci\n4tCO5qf3elyA712LMwtKZsIeWsDX+OUCWglkmfvsAq06JfJx60YnYagbVtsdBTSR\nftmiqarrs71U+gaQVpeHgZYpKLMPNO/2Nu5Le2/SUHwXKXML3sDk4dNXNGb6YPAE\nLGNdqiyeG8o98agdkNIzIh0CgYEAlTGhMPfGRL3UXoNN8vopjEUWXozUmvJ090S/\nJLtZXtKtNBp5cEOJWZT9biVhFeKgCZc8ba7ahA29b/aLs+AnPlrfnJh+qzZhQfHz\ngJ0PSwAbkBs5fFBOaCHppiRlvXuFRemo95m4pcwTPBx7Mj4Xqx4lxij2E2rNVMSy\n4AI4l10CgYBwefqXt8B+D+0EvmhyHk19Tk8/fPelclJUv/IVI59c0F9UMAA2rD1U\nNW6k9251OGU7mQkztluNvl13qtAW/DveOjkFeDJIMzhFjravpLQXhUK4ETnM44YL\nFbClVGJaHYSHgOkNpcN5lYVLoyEvzv9rEPwBqpZRVnwWj6L+/I2L5Q==\n-----END RSA PRIVATE KEY-----\n"",
  ""userId"": ""170079991923474689""
}");
client = Clients.AuthService(
    new(
        apiUrl,
        ITokenProvider.ServiceAccount(
            serviceAccount,
            apiUrl,
            apiProject)));
result = await client.GetMyUserAsync(new());
Console.WriteLine($"User: {result.User}");
```

### Test client

After you have configured everything correctly, you can simply start the example by:

```bash
dotnet run
```

This will output something similar to:

```
User: {"FirstName": "MyName", "LastName": "MyLastName" ... }
```

## Completion

You have successfully used the ZITADEL .NET SDK to call the auth API!
To use the auth API you will not need a specific role, because only an authenticated user is needed.

For accessing the admin or management API the user will need some specific roles.
If you encountered an error (e.g. `code = PermissionDenied desc = No matching permissions found`), 
ensure your service user has the required permissions by assigning the `ORG_OWNER` or `ORG_OWNER_VIEWER` role
and check the mentioned [guides](#prerequisites) at the beginning.

If you've run into any other problem, don't hesitate to contact us or raise an issue on [ZITADEL](https://github.com/zitadel/zitadel/issues) or in the [SDK](https://github.com/zitadel/zitadel-go/issues).

### Whats next?

Now you can proceed implementing our APIs by adding more calls.

Checkout more [examples from the SDK](https://github.com/zitadel/zitadel-go/blob/main/example) or refer to our [API Docs](/apis/introduction).

> This guide will be updated soon to show you how to use the SDK for your own API as well.

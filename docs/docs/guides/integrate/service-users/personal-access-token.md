---
title: Configure personal access token authentication for service users
sidebar_label: Personal access token authentication
sidebar_position: 3
---

A Personal Access Token (PAT) is a ready to use token which can be used as _Authorization_ header.
At the moment ZITADEL only allows PATs for machine accounts (service users).

It is an alternative to the [private key JWT profile authentication](private-key-jwt) and [client credentials authentication](client-credentials). Read more about that the different [authentication methods for service users](authenticate-service-users).

## Create a Service User with a PAT

1. Navigate to Service Users
2. Click on **New**
3. Enter a user name and a display name
4. Click on the Personal Access Token menu point in the detail of your user
5. Click on **New**
6. You can either set an expiration date or leave it empty if you don't want it to expire
7. Copy the token from the dialog (You will not see this again)

![Create new service user](/img/guides/console-service-user-pat.gif)

## Grant role for ZITADEL

To be able to access the ZITADEL APIs your service user needs permissions to ZITADEL.

1. Go to the detail page of your organization
2. Click in the top right corner the "+" button
3. Search for your service user
4. Give the user the role you need, for the example we choose Org Owner (More about [ZITADEL Permissions](/docs/guides/manage/console/managers))

![Add org owner to service user](/img/guides/console-service-user-org-owner.gif)

## Accessing ZITADEL APIs

You might want to access ZITADEL APIs to manage resources, such as users, or to validate tokens sent to your backend service.
Follow our guides on [how to access ZITADEL API](../zitadel-apis/access-zitadel-apis) to use the ZITADEL APIs with your service user.

### Token introspection

Your API endpoint might receive tokens from users and need to validate the token with ZITADEL.
In this case your API needs to authenticate with ZITADEL and then do a token introspection.
Follow our [guide on token introspection with private key JWT](../token-introspection/private-key-jwt) to learn more.

## Call ZITADEL API with PAT

Because the PAT is a ready to use token, you can add it as Authorization Header and send it in your requests to the ZITADEL API.
In this example we read the organization of the service user.

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header 'Authorization: Bearer {PAT}' 
```


## Client application authentication

The above steps demonstrate service user authentication.
If your application also needs to authenticate itself, you can utilize [Client Credentials Grant](./client-credentials).
Refer to ZITADEL documentation for details on this alternative method.

## Security considerations

* **Store private keys securely:** **Never share or embed the private key in your code or application.** Consider using secure key management solutions.
* **Set appropriate JWT expiration times:** Limit the validity period of tokens to minimize the impact of potential compromise.
* **Implement proper error handling:** Handle situations where JWT verification fails or tokens are expired.

By following these steps and adhering to security best practices, you can effectively secure service user and client application communication within ZITADEL using private key JWT authentication.
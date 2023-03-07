---
title: User Metadata
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

In this guide you will learn how to manually create the necessary requests to authenticate and request a user's metadata from ZITADEL.

Typical examples for user metadata include:

- Link the user to an internal identifier (eg, userId, contract number, etc.)
- Save custom user data when registering a user
- Route upstream traffic based on user attributes

## Prerequisites

### Create a new client

- Create a new [web application](../console/applications#web)
- Use Code-Flow
- In this example we will use `http://localhost` as redirect url
- Make sure to note the client secret

### Add metadata to a user

- [Add metadata](/guides/manage/customize/user-metadata) to a user
- Make sure you will use this user to login during later steps

## Requesting a token

:::info
In this guide we will manually request a token from ZITADEL for demonstration purposes. You will likely use a client library for the OpenID Authentication.
:::

### Set environment variables

We will use some information throughout this guide. Set the required environment variables as follows. Make sure to replace the values with your information.

```bash
export CLIENT_SECRET=QCiMffalakI...zpT0vuOsSkVk1ne \
export CLIENT_ID="16604...@docs-claims" \
export REDIRECT_URI="http://localhost" \
export ZITADEL_DOMAIN="https://...asd.zitadel.cloud"
```

<Tabs>
<TabItem value="go" label="Go" default>

Grab zitadel-tools to create the [required string](/apis/openidoauth/authn-methods#client-secret-basic) for Basic authentication:

```bash
git clone git@github.com:zitadel/zitadel-tools.git
cd zitadel-tools/cmd/basicauth
export BASIC_AUTH="$(go run basicauth.go -id $CLIENT_ID -secret $CLIENT_SECRET)"
```

</TabItem>

<TabItem value="python" label="Python">

```python
import base64
import urllib.parse
import os

clientId = os.environ.get("CLIENT_ID")
clientSecret = os.environ.get("CLIENT_SECRET")

escaped = safe_string = urllib.parse.quote_plus(clientId) + ":" + urllib.parse.quote_plus(clientSecret)
message_bytes = escaped.encode('ascii')
base64_bytes = base64.b64encode(message_bytes)
base64_message = base64_bytes.decode('ascii')

print(base64_message)
```

Export the result to the environment variable `BASIC_AUTH`.

</TabItem>

<TabItem value="js" label="Javascript" default>

```javascript
esc = encodeURIComponent(process.env.CLIENT_ID) + ":" + encodeURIComponent(process.env.CLIENT_SECRET)
enc = btoa(esc)
console.log(enc)
```

Export the result to the environment variable `BASIC_AUTH`.

</TabItem>

<TabItem value="manually" label="Manually">

You need to create a string as described [here](/apis/openidoauth/authn-methods#client-secret-basic).

Use a programming language of your choice or manually create the strings with online tools (don't use these secrets for production) like: 

- https://www.urlencoder.org/
- https://www.base64encode.org/

Export the result to the environment variable `BASIC_AUTH`.

</TabItem>
</Tabs>

### Create Auth Request

You need to create a valid auth request, including the reserved scope `urn:zitadel:iam:user:metadata`. Please refer to our API documentation for more information about [reserved scopes](/apis/openidoauth/scopes#reserved-scopes) or try it out in our [OIDC Authrequest Playground](/apis/openidoauth/authrequest?scope=openid%20email%20profile%20urn%3Azitadel%3Aiam%3Auser%3Ametadata).

Login with the user to which you have added the metadata. After the login you will be redirected.

Grab the code paramter from the url (disregard the &code= parameter) and export the code as environment variable:

```bash
export AUTH_CODE="Y6nWsgR5WB...zUtFqSp5Xw"
```

### Token Request

```bash
curl --request POST \
--url "${ZITADEL_DOMAIN}/oauth/v2/token" \
--header "Accept: application/json" \
--header "Authorization: Basic ${BASIC_AUTH}" \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data grant_type=authorization_code \
--data-urlencode "code=$AUTH_CODE" \
--data-urlencode "redirect_uri=$REDIRECT_URI"
```

The result will give you something like:

```json
{
    "access_token":"jZuRixKQTVecEjKqw...kc3G4",
    "token_type":"Bearer",
    "expires_in":43199,
    "id_token":"ey...Ww"
}
```

Grab the access_token value and export as an environment variable:

```bash
export ACCESS_TOKEN="jZuRixKQTVecEjKqw...kc3G4"
```

### Request metadata from userinfo endpoint

With the access token we can make a request to the userinfo endpoint to get the user's metadata. This method is the preferred method to retrieve a user's information in combination with opaque tokens, to insure that the token is valid.

```bash
curl --request GET \
  --url "${ZITADEL_DOMAIN}/oidc/v1/userinfo" \
  --header "Authorization: Bearer $ACCESS_TOKEN"
```

The response will look something like this

```json
{
    "email":"road.runner@zitadel.com",
    "email_verified":true,
    "family_name":"Runner",
    "given_name":"Road",
    "locale":"en",
    "name":"Road Runner",
    "preferred_username":"road.runner@...asd.zitadel.cloud",
    "sub":"166.....729",
    "updated_at":1655467738,
    //highlight-start
    "urn:zitadel:iam:user:metadata":{
        "ContractNumber":"MTIzNA",
        }
    //highlight-end
    }
```

You can grab the metadata from the reserved claim `"urn:zitadel:iam:user:metadata"` as key-value pairs. Note that the values are base64 encoded. So the value `MTIzNA` decodes to `1234`.

### Send metadata inside the ID token (optional)

Check "User Info inside ID Token" in the configuration of your application.

![](/img/console_projects_application_token_settings.png)

Now request a new token from ZITADEL.

The result will give you something like:

```json
{
    "access_token":"jZuRixKQTVecEjKqw...kc3G4",
    "token_type":"Bearer",
    "expires_in":43199,
    "id_token":"ey...Ww"
}
```

Grab the id_token and inspect the contents of the token at [jwt.io](https://jwt.io/). You should get the same info in the ID token as when requested from the user endpoint.
---
title: Configure private key JWT authentication for service users
sidebar_label: Private key JWT authentication
sidebar_position: 2
---

This guide demonstrates how developers can leverage private key JWT authentication to secure communication between service users and client applications within ZITADEL.

In ZITADEL we use the `urn:ietf:params:oauth:grant-type:jwt-bearer` (**“JWT bearer token with private key”**, [RFC7523](https://tools.ietf.org/html/rfc7523)) authorization grant for this non-interactive authentication.

Read more about the [different authentication methods for service users](authenticate-service-users) and their benefits, drawbacks, and security considerations.

#### How private key JWT authentication works

1. Generate a private/public key pair associated with the service user.
2. The authorization server stores the public key; and
3. returns the private key as a json file
4. The developer configures the client in such a way, that 
5. JWT assertion is created with the subject of the service user, and the JWT is signed by the private key
6. Resource owner requests a token by sending the client_assertion
7. Authorization server validates the signature using the service user's public key
8. Authorization server returns an OAuth access_token
9. Resource Owner calls a Resource Server by including the access_token in the Header
10. Resource Server validates the JWT with [token introspection](../token-introspection/)

![private key jwt authentication sequence diagram](/img/guides/integrate/service-users/sequence-private-key-jwt.svg)


## Prerequisites

A code library/framework supporting JWT generation and verification (e.g., `pyjwt` for Python, `jsonwebtoken` for Node.js).

## Steps to authenticate a Service User with private JWT

You need to follow these steps to authenticate a service user and receive an access token that can be used in subsequent requests.

### 1. Create a Service User

1. Navigate to Service Users
2. Click on **New**
3. Enter a username and a display name
4. Click on **Create**

### 2. Generate a private key file

1. Access the ZITADEL web console and navigate to the service user details.
2. Click on the **Keys** menu point in the detail of your new service user
3. Click on **New**
4. You can either set an expiration date or leave it empty if you don't want it to expire
5. Click on **Download** and save the key file

:::note
Make sure to save the key file. You won't be able to retrieve it again.
If you lose it, you will have to generate a new one.
:::

:::note Expiration
If you specify an expiration date, note that the key will expire at midnight that day
:::

![Create private key](/img/console_serviceusers_new_key.gif)

The downloaded JSON should look something like outlined below. The value of `key` contains the _private_ key for your service account. Please make sure to keep this key securely stored and handled with care. The public key is automatically stored in ZITADEL.

```json
{
    "type":"serviceaccount",
    "keyId":"100509901696068329",
    "key":"-----BEGIN RSA PRIVATE KEY----- [...] -----END RSA PRIVATE KEY-----\n",
    "userId":"100507859606888466"
}
```

### 3. Create a JWT and sign with private key

You need to create a JWT with the following header and payload and sign it with the RS256 algorithm.

Header

```json
{
    "alg": "RS256",
    "kid":"100509901696068329"
}
```

Make sure to include `kid` in the header with the value of `keyId` from the downloaded JSON.

Payload

```json
{
    "iss": "100507859606888466",
    "sub": "100507859606888466",
    "aud": "https://$CUSTOM-DOMAIN",
    "iat": [Current UTC timestamp, e.g. 1605179982, max. 1 hour ago],
    "exp": [UTC timestamp, e.g. 1605183582]
}
```

* `iss` represents the requesting party, i.e. the owner of the private key. In this case the value of `userId` from the downloaded JSON.
* `sub` represents the application. Set the value also to the value of `userId`
* `aud` must be your [Custom Domain](../../../concepts/features/custom-domain)
* `iat` is a unix timestamp of the creation signing time of the JWT, e.g. now and must not be older than 1 hour ago
* `exp` is the unix timestamp of expiry of this assertion

Please refer to [JWT with private key](/apis/openidoauth/authn-methods#jwt-with-private-key) API reference for further information.

If you use Go, you might want to use the [provided tool](https://github.com/zitadel/zitadel-tools) to generate a JWT from the downloaded JSON.
There are many [libraries](https://jwt.io/#libraries-io) to generate and sign JWT.

**Code Example (Python using `pyjwt`):**

```python
import jwt
import datetime

# Replace with your service user ID and private key
service_user_id = "your_service_user_id"
private_key = "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY\n-----END PRIVATE KEY-----"
key_id = "your_key_id"

# ZITADEL API URL (replace if needed)
api_url = "your_custom_domain"

# Generate JWT claims
payload = {
    "iss": service_user_id,
    "sub": service_user_id,
    "aud": api_url,
    "exp": datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(minutes=5),
    "iat": datetime.datetime.now(datetime.timezone.utc)
}

header = {
    "alg": "RS256",
    "kid": key_id
}

# Sign the JWT using RS256 algorithm
encoded_jwt = jwt.encode(payload, private_key, algorithm="RS256", headers=header)

print(f"Generated JWT: {encoded_jwt}")
```

### 4. Request an OAuth token with the generated JWT

With the encoded JWT from the prior step, you will need to craft a POST request to ZITADEL's token endpoint:

```bash
curl --request POST \
  --url https:/$CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid' \
  --data assertion=eyJ0eXAiOiJKV1QiL...
```

* `grant_type` should be set to `urn:ietf:params:oauth:grant-type:jwt-bearer`
* `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid`.
* `assertion` is the encoded value of the JWT that was signed with your private key from the prior step

If you want to access ZITADEL APIs, make sure to include the required scopes `urn:zitadel:iam:org:project:id:zitadel:aud`.
Read our guide [how to access ZITADEL APIs](../zitadel-apis/access-zitadel-apis) to learn more.

**Important Note:** If the service user token needs to be validated using token introspection, ensure you include the `urn:zitadel:iam:org:project:id:{projectid}:aud` scope in your token request. 
Without this, token introspection will fail.

You should receive a successful response with `access_token`,  `token_type` and time to expiry in seconds as `expires_in`.

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{
  "access_token": "MtjHodGy4zxKylDOhg6kW90WeEQs2q...",
  "token_type": "Bearer",
  "expires_in": 43199
}
```

### 5. Include the access token in the authorization header

When making API requests on behalf of the service user, include the generated token in the "Authorization" header with the "Bearer" prefix.

```bash
curl --request POST \
  --url $YOUR_API_ENDOINT \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Bearer MtjHodGy4zxKylDOhg6kW90WeEQs2q...'
```

## Accessing ZITADEL APIs

You might want to access ZITADEL APIs to manage resources, such as users, or to validate tokens sent to your backend service.
Follow our guides on [how to access ZITADEL API](../zitadel-apis/access-zitadel-apis) to use the ZITADEL APIs with your service user.

### Token introspection

Your API endpoint might receive tokens from users and need to validate the token with ZITADEL.
In this case your API needs to authenticate with ZITADEL and then do a token introspection.
Follow our [guide on token introspection with private key JWT](../token-introspection/private-key-jwt) to learn more.

## Client application authentication

The above steps demonstrate service user authentication.
If your application also needs to authenticate itself, you can utilize [Client Credentials Grant](./client-credentials).
Refer to ZITADEL documentation for details on this alternative method.

## Security considerations

* **Store private keys securely:** **Never share or embed the private key in your code or application.** Consider using secure key management solutions.
* **Set appropriate JWT expiration times:** Limit the validity period of tokens to minimize the impact of potential compromise.
* **Implement proper error handling:** Handle situations where JWT verification fails or tokens are expired.

By following these steps and adhering to security best practices, you can effectively secure service user and client application communication within ZITADEL using private key JWT authentication.

## Notes

* [JWT with private key](/apis/openidoauth/authn-methods#jwt-with-private-key) API reference
* [Accessing ZITADEL API](../zitadel-apis/access-zitadel-apis)
* [Token introspection with private key JWT](../token-introspection/private-key-jwt)

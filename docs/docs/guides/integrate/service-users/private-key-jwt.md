---
title: Configure private key JWT authentication for service users
sidebar_label: Private key JWT authentication
sidebar_position: 2
---

This guide demonstrates how developers can leverage private key JWT authentication to secure communication between service users and client applications within ZITADEL.

In ZITADEL we use the `urn:ietf:params:oauth:grant-type:jwt-bearer` (**“JWT bearer token with private key”**, [RFC7523](https://tools.ietf.org/html/rfc7523)) authorization grant for this non-interactive authentication.

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

![Create private key](/img/console_serviceusers_new_key.gif)

The downloaded json should look something like outlined below. The value of `key` contains the *private* key for your service account. Please make sure to keep this key securely stored and handle with care. The public key is automatically stored in ZITADEL.

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
* `aud` must be ZITADEL's issuing domain
* `iat` is a unix timestamp of the creation signing time of the JWT, e.g. now and must not be older than 1 hour ago
* `exp` is the unix timestamp of expiry of this assertion

Please refer to [JWT_with_Private_Key](/apis/openidoauth/authn-methods#jwt-with-private-key) in the documentation for further information.

If you use Go, you might want to use the [provided tool](https://github.com/zitadel/zitadel-tools) to generate a JWT from the downloaded json. There are many [libraries](https://jwt.io/#libraries-io) to generate and sign JWT.

**Code Example (Python using `pyjwt`):**

```python
import jwt
import datetime

# Replace with your service user ID and private key
service_user_id = "your_service_user_id"
private_key = "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY\n-----END PRIVATE KEY-----"

# ZITADEL API URL (replace if needed)
api_url = "https://api.zitadel.cloud/v1"

# Generate JWT claims
payload = {
    "iss": "your_zitadel_instance_id",
    "sub": service_user_id,
    "aud": api_url,
    "exp": datetime.utcnow() + datetime.timedelta(minutes=5),
    "iat": datetime.utcnow()
}

# Sign the JWT using RS256 algorithm
encoded_jwt = jwt.encode(payload, private_key, algorithm="RS256")

print(f"Generated JWT: {encoded_jwt}")
```

### 4. With this JWT, request an OAuth token from ZITADEL

With the encoded JWT from the prior step, you will need to craft a POST request to ZITADEL's token endpoint:

```bash
curl --request POST \
  --url https:/$CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid profile email' \
  --data assertion=eyJ0eXAiOiJKV1QiL...
```

### 5. Include the access token in the authorization header

When making API requests to ZITADEL on behalf of the service user, include the generated JWT in the "Authorization" header with the "Bearer" prefix.

For this example let's call the userinfo endpoint to verify that our access token works.

```bash
curl --request POST \
  --url $CUSTOM-DOMAIN/oidc/v1/userinfo \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Bearer MtjHodGy4zxKylDOhg6kW90WeEQs2q...'
```

You should receive a response with your service user's information.

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{
  "name": "MyServiceUser",
  "preferred_username": "service_user@$CUSTOM-DOMAIN",
  "updated_at": 1616417938
}
```

## Accessing ZITADEL's Management API

If you want to access the ZITADEL API with this access token, you have to add `urn:zitadel:iam:org:project:id:zitadel:aud` to the list of scopes.

* `grant_type` should be set to `urn:ietf:params:oauth:grant-type:jwt-bearer`
* `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid`. For this example, please include `profile` and `email`
* `assertion` is the encoded value of the JWT that was signed with your private key from the prior step

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

## Client Application Authentication

The above steps demonstrate service user authentication.
If your application also needs to authenticate itself, you can utilize [Client Credentials Grant](./client-credentials).
Refer to ZITADEL documentation for details on this alternative method.

## Security Considerations

* **Store private keys securely:** **Never share or embed the private key in your code or application.** Consider using secure key management solutions.
* **Set appropriate JWT expiration times:** Limit the validity period of tokens to minimize the impact of potential compromise.
* **Implement proper error handling:** Handle situations where JWT verification fails or tokens are expired.

By following these steps and adhering to security best practices, you can effectively secure service user and client application communication within ZITADEL using private key JWT authentication.
Remember to consult the official ZITADEL documentation for detailed information and potential changes in the future.

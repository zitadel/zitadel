
## Accessing ZITADEL APIs

In this step, we will authenticate a service user and receive an access_token to use against the ZITADEL API.

You will need to craft a POST request to ZITADEL's token endpoint:

```bash
curl --request POST \
  --url https://$CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Basic ${BASIC_AUTH}' \
  --data grant_type=client_credentials \
  --data scope='openid urn:zitadel:iam:org:project:id:zitadel:aud'
```

If you want to access the ZITADEL API with this access token, you have to add `urn:zitadel:iam:org:project:id:zitadel:aud` to the list of scopes.

* `grant_type` should be set to `urn:ietf:params:oauth:grant-type:jwt-bearer`
* `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid` and `urn:zitadel:iam:org:project:id:zitadel:aud`
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


In this step we will authenticate a service user and receive an access_token to use against the ZITADEL API.

You will need to craft a POST request to ZITADEL's token endpoint:

```bash
curl --request POST \
  --url https://$CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Basic ${BASIC_AUTH}' \
  --data grant_type=client_credentials \
  --data scope='openid profile email urn:zitadel:iam:org:project:id:zitadel:aud'
```

* `grant_type` should be set to `client_credentials`
* `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid`. For this example, please include `profile`, `email`
  and `urn:zitadel:iam:org:project:id:zitadel:aud`. The latter provides access to the ZITADEL API.

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

### Access ZITADEL's API with the token

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header "Authorization: Bearer $TOKEN" 
```

Where `$TOKEN` corresponds to the access token from the earlier response.
---
title: Access ZITADEL System API
---
:::note
This guide focuses on the ZITADEL System API. To access the other APIs (Admin, Auth, Management), please checkout [this guide](./access-zitadel-apis). 
The ZITADEL System API is currently only available for ZITADEL Self-Hosted deployments.
:::

## System API User

The System API works superordinate over all instances. Therefore, you need to define a separate users to get access to this API.
You can do so by customizing the [runtime configuration](/self-hosting/manage/configure#runtime-configuration).

To authenticate the user a self-signed JWT will be created and utilized.

You can define any id for your user. This guide will assume it's `system-user-1`.

## Generate an RSA keypair

Generate an RSA private key with 2048 bit modulus:

```bash
openssl genrsa -out system-user-1.pem 2048
```

and export a public key from the newly created private key:

```bash
openssl rsa -in system-user-1.pem -outform PEM -pubout -out system-user-1.pub
```

## Runtime Configuration

Provide the **public** key to the ZITADEL runtime configuration.

Either with the path to the key:

```yaml
SystemAPIUsers:
  - system-user-1:
      Path: /system-user-1.pub
```

or with a base64 encoded value of the key:

```yaml
SystemAPIUsers:
  - system-user-1:
      KeyData: <base64 encoded value of system-user-1.pub>
```

## Generate JWT

Similar to the OAuth 2.0 JWT Profile, we will create and sign a JWT. For this API, the JWT will not be used to authenticate against ZITADEL Authorization Server, but rather directly to the API itself.

The JWT payload will need to contain the following claims:

```json
{
  "iss": "<userid>",
  "sub": "<userid>",
  "aud": "<https://your_domain>",
  "exp": <now+1h>,
  "iat": <now>
}
```

So for your instance running on `custom-domain.com` the claims could look like this:

```json
{
  "iss": "system-user-1",
  "sub": "system-user-1",
  "aud": "https://custom-domain.com",
  "iat": 1659957184,
  "exp": 1659960784
}
```

:::note
If your system is exposed without TLS or on a dedicated port, be sure to provide this in your audience, e.g. http://localhost:8080 
:::

### ZITADEL Tools

If you want to manually create a JWT for a test, you can also use our [ZITADEL Tools](https://github.com/zitadel/zitadel-tools). Download the latest release and run:

```bash
./key2jwt -audience=https://custom-domain.com -key=system-user-1.pem -issuer=system-user-1
```

## Call the System API

Now that you configured ZITADEL and created a JWT, you can call the System API and authenticate using the token:

```bash
curl --request POST \
  --url {your_domain}/system/v1/instances/_search \
  --header 'Authorization: Bearer {token}' \
  --header 'Content-Type: application/json'
```

You should get a successful response with a `totalResult` number of 1 and the details of your instance:

```json
{
	"details": {
		"totalResult": "1"
	},
	"result": [
		{
			"id": "172698969497928101",
			"details": {
				"sequence": "102",
				"creationDate": "2022-08-02T09:30:10.781068Z",
				"changeDate": "2022-08-02T09:30:10.781068Z",
				"resourceOwner": "172698969497928101"
			},
			"state": "STATE_RUNNING",
			"name": "ZITADEL",
			"domains": [
				{
					"details": {
						"sequence": "108",
						"creationDate": "2022-08-02T09:30:10.781068Z",
						"changeDate": "2022-08-02T09:30:10.781068Z",
						"resourceOwner": "172698969497928101"
					},
					"domain": "custom-domain.com",
					"primary": true
				},
				{
					"details": {
						"sequence": "108",
						"creationDate": "2022-08-02T09:30:10.781068Z",
						"changeDate": "2022-08-02T09:30:10.781068Z",
						"resourceOwner": "172698969497928101"
					},
					"domain": "zitadel-gnft7o.custom-domain.com",
					"generated": true
				}
			]
		}
	]
}
```

With this token you are allowed to access the whole [ZITADEL System API](/apis/system).

## Summary

* Create an RSA keypair
* Provide the public key with a userID to ZITADEL using the runtime configuration
* Authorize the request with a JWT signed with your private key

Where to go from here:

* [ZITADEL API Documentation](/apis/introduction)

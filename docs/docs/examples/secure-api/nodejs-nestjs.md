---
title: ZITADEL with Node.js 
sidebar_label: Node.js
---

# ZITADEL with Node.js (NestJS)

This documentation section guides you through the process of integrating ZITADEL into your Node.js backend using the NestJS framework. The provided example demonstrates authentication using an OIDC (OAuth2) token introspection strategy with a ZITADEL service account for machine-to-machine communication.

## Overview

The NestJS API includes a single secured route that prints "Hello World!" when authenticated. The API expects an authorization header with a valid JWT, serving as a bearer token to authenticate the user when calling the API. The API will validate the access token on the [introspect endpoint](https://zitadel.com/docs/apis/openidoauth/endpoints#introspection_endpoint) and receive the user from ZITADEL.

The API application utilizes [JWT with Private Key](https://zitadel.com/docs/apis/openidoauth/authn-methods#jwt-with-private-key) for authentication against ZITADEL and accessing the introspection endpoint. Make sure to create an API Application within Zitadel and download the JSON. In this instance, we use this service account, so make sure to provide the secrets in the example application via environmental variables.

## Overview

The NestJS API includes a private endpoint `GET http://localhost:${APP_PORT}/api/v1/app`, which returns "Hello World" when authenticated. The authentication is performed using a JWT obtained through the token introspection strategy.

## Running the Example

### Prerequisites

Make sure you have Node.js and npm installed on your machine.

### ZITADEL Configuration for the API

1. Create a ZITADEL instance and a project by following the steps [here](https://zitadel.com/docs/guides/start/quickstart#2-create-your-first-instance).

2. Set up an API application within your project:
   - Create a new application of type "API" with authentication method "Private Key".
   - Create a and save the Private Key JSON file.

### Create and Run the API

Clone or download the [example repository](https://github.com/ehwplus/zitadel-nodejs-nestjs):

```bash
git clone https://github.com/ehwplus/zitadel-nodejs-nestjs && cd zitadel-nodejs-nestjs
```

and follow the instructions here: https://github.com/ehwplus/zitadel-nodejs-nestjs/blob/main/README.md#installation

### Test the API

Call the API without authorization headers:

```bash
curl --request GET \
     --url http://localhost:${APP_PORT}/api/v1/app
```

You should get a response with Status Code 401 and an error message.

Now, add an authorization header with a valid JWT obtained through ZITADEL:

```bash
export JWT=your-valid-jwt

curl --request GET \
    --url http://localhost:${APP_PORT}/api/v1/app \
    --header "authorization: Bearer $JWT"
```

You should now receive a response with Status Code 200 and the message:

```json
"Hello World!"
```

Congratulations! You have successfully integrated ZITADEL authentication into your NestJS API using the Token Introspection strategy.

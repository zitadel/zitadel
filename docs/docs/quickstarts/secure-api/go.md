---
title: Go
---

This guide shows you how to integrate ZITADEL into your Go API.
You'll use an OAuth 2 Token Introspection to secure your API.

At the end of the guide, you should have an API with a protected endpoint.

## Prerequisites

The client [SDK](https://github.com/caos/zitadel-go) provides an interceptor for both GRPC and HTTP.
This uses a JWT to handle the OAuth 2.0 introspection request, including authentication.
The Private Key uses our [OIDC client library](https://github.com/caos/oidc).

All you need is an API and its JSON key.

## Go Setup

### Add Go SDK to your project

Add the SDK into Go Modules with this command:

```bash
go get github.com/caos/zitadel-go
```

### Create example API

Create a new go file with the content below.
This creates an API with two endpoints.
The path `/public` always responds with `ok` and the current timestamp.
The path `/protected` responds with the same, but only if a valid `access_token` is sent.
The token must not be expired, and the API has to be part of the audience (either `client_id` or `project_id`).

```go
package main

import (
	"log"
	"net/http"
	"time"

	api_mw "github.com/caos/zitadel-go/pkg/api/middleware"
	http_mw "github.com/caos/zitadel-go/pkg/api/middleware/http"
	"github.com/caos/zitadel-go/pkg/client"
	"github.com/caos/zitadel-go/pkg/client/middleware"
)

func main() {
	introspection, err := http_mw.NewIntrospectionInterceptor(client.Issuer, middleware.OSKeyPath())
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/public", writeOK)
	router.HandleFunc("/protected", introspection.HandlerFunc(writeOK))

	lis := "127.0.0.1:5001"
	log.Fatal(http.ListenAndServe(lis, router))
}

func writeOK(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK " + time.Now().String()))
}

```

#### JSON key

To provide the JSON key to the SDK, set an environment variable `ZITADEL_KEY_PATH`, using the path to the JSON as the value.

```bash
export ZITADEL_KEY_PATH=/Users/test/apikey.json
```

For development purposes, you should be able to set this in your IDE.

If you can't set it via environment variable, you can also exchange the `middleware.OSKeyPath()` and pass it directly:

```go
introspection, err := http_mw.NewIntrospectionInterceptor(
	client.Issuer,
	"/Users/test/apikey.json",
)
```

#### Custom ZITADEL instance

If your client does not use ZITADEL Cloud (zitadel.ch), be sure to provide the correct Issuer:
```go
introspection, err := http_mw.NewIntrospectionInterceptor(
	"https://issuer.custom.ch",
	middleware.OSKeyPath(),
)
```

### Test API

After you have configured everything correctly, you can start the example with this command:

```bash
go run main.go
```

You can now call the API by browser or curl. Try the public endpoint first:

```bash
curl -i localhost:5001/public
```

This should return something like:

```
HTTP/1.1 200 OK
Date: Tue, 24 Aug 2021 11:11:17 GMT
Content-Length: 59
Content-Type: text/plain; charset=utf-8

OK 2021-08-24 13:11:17.135719 +0200 CEST m=+30704.913892168
```

And the protected:

```bash
curl -i localhost:5001/protected
```

This returns:

```
HTTP/1.1 401 Unauthorized
Content-Type: application/json
Date: Tue, 24 Aug 2021 11:13:10 GMT
Content-Length: 21

"auth header missing"
```

Get a valid `access_token` for the API.
To do this, log in to an application of the same project.
You can also explicitly request the `project_id` for the audience by scope `urn:zitadel:iam:org:project:id:{projectid}:aud`.

If you provide a valid Bearer Token:

```bash
curl -i -H "Authorization: Bearer ${token}" localhost:5001/protected
```

This will return an OK response as well:
```
HTTP/1.1 200 OK
Date: Tue, 24 Aug 2021 11:13:33 GMT
Content-Length: 59
Content-Type: text/plain; charset=utf-8

OK 2021-08-24 13:13:33.131943 +0200 CEST m=+30840.911149251
```

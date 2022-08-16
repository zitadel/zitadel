---
title: Go
---

This integration guide shows you how to integrate **ZITADEL** into your Go API. It demonstrates how to secure your API using
OAuth 2 Token Introspection.

At the end of the guide you should have an API with a protected endpoint.

## Prerequisites

The client [SDK](https://github.com/zitadel/zitadel-go) will provides an interceptor for both GRPC and HTTP.
This will handle the OAuth 2.0 introspection request including authentication using JWT with Private Key using our [OIDC client library](https://github.com/zitadel/oidc).
All that is required, is to create your API and download the private key file later called `Key JSON` for the service user.

## Go Setup

### Add Go SDK to your project

You need to add the SDK into Go Modules by:

```bash
go get github.com/zitadel/zitadel-go/v2
```

### Create example API

Create a new go file with the content below. This will create an API with two endpoints. On path `/public` it will always write
back `ok` and the current timestamp. On `/protected` it will respond the same but only if a valid access_token is sent. The token
must not be expired and the API has to be part of the audience (either client_id or project_id).

Make sure to fill the var `issuer` with your own domain. This is the domain of your instance you can find it on the instance detail in the ZITADEL Cloud Customer Portal or in the ZITADEL Console.
```go
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	http_mw "github.com/zitadel/zitadel-go/v2/pkg/api/middleware/http"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
)

var (
	issuer = flag.String("issuer", "", "issuer of your ZITADEL instance (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
)

func main() {
	flag.Parse()

	introspection, err := http_mw.NewIntrospectionInterceptor(*issuer, middleware.OSKeyPath())
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

#### Key JSON

To provide the key JSON to the SDK, simply set an environment variable `ZITADEL_KEY_PATH` with the path to the JSON as value.

```bash
export ZITADEL_KEY_PATH=/Users/test/apikey.json
```

For development purposes you should be able to set this in your IDE.

If you're not able to set it via environment variable, you can also exchange the `middleware.OSKeyPath()` and pass it directly:

```go
introspection, err := http_mw.NewIntrospectionInterceptor(
	client.Issuer,
	"/Users/test/apikey.json",
)
```

### Test API

After you have configured everything correctly, you can simply start the example by:

```bash
go run main.go
```

You can now call the API by browser or curl. Try the public endpoint first:

```bash
curl -i localhost:5001/public
```

it should return something like: 

```
HTTP/1.1 200 OK
Date: Tue, 24 Aug 2021 11:11:17 GMT
Content-Length: 59
Content-Type: text/plain; charset=utf-8

OK 2021-08-24 13:11:17.135719 +0200 CEST m=+30704.913892168
```

and the protected:

```bash
curl -i localhost:5001/protected
```

it will return:

```
HTTP/1.1 401 Unauthorized
Content-Type: application/json
Date: Tue, 24 Aug 2021 11:13:10 GMT
Content-Length: 21

"auth header missing"
```

Get a valid access_token for the API. You can achieve this by login into an application of the same project or
by explicitly requesting the project_id for the audience by scope `urn:zitadel:iam:org:project:id:{projectid}:aud`.

If you provide a valid Bearer Token:

```bash
curl -i -H "Authorization: Bearer ${token}" localhost:5001/protected
```

it will return an OK response as well:
```
HTTP/1.1 200 OK
Date: Tue, 24 Aug 2021 11:13:33 GMT
Content-Length: 59
Content-Type: text/plain; charset=utf-8

OK 2021-08-24 13:13:33.131943 +0200 CEST m=+30840.911149251
```

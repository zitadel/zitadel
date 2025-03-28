---
title: Test Actions Request Locally
---

In this guide, you will create a ZITADEL execution and target. Before a user is created through the API, the target is called.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

## Start example target

To start a simple HTTP server locally, which receives the call and manipulated the request, the following code example can be used:

```go
package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

type contextRequest struct {
	Request *user.AddHumanUserRequest `json:"request"`
}

// call HandleFunc to read the request body, manipulate the content and return the request
func call(w http.ResponseWriter, req *http.Request) {
	// read the body content
	sentBody, err := io.ReadAll(req.Body)
	if err != nil {
		// if there was an error while reading the body return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// read the request into the expected structure
	request := new(infoRequest)
	if err := json.Unmarshal(sentBody, request); err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
    
	// build the response from the received request
	response := request.Request
	// manipulate the request to send back as response
	if response.Metadata == nil {
		response.Metadata = make([]*user.SetMetadataEntry, 0)
	}
	response.Metadata = append(response.Metadata, &user.SetMetadataEntry{Key: "organization", Value: []byte("company")})

	// marshal the request into json
	data, err := json.Marshal(response)
	if err != nil {
		// if there was an error while marshalling the json
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// return the manipulated request
	w.Write(data)
}

func main() {
	// handle the HTTP call under "/webhook"
	http.HandleFunc("/call", call)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}

```

What happens here is that the target receives the request Zitadel receives, adds a metadata entry to the request and returns it.

### Check Signature

To additionally check the signature header you can add the following to the example:
```go
	// validate signature
	if err := actions.ValidatePayload(sentBody, req.Header.Get(actions.SigningHeader), signingKey); err != nil {
		// if the signed content is not equal the sent content return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
```

Where you can replace 'signingKey' with the key received in the next step 'Create target'.

## Create target

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as call, the target can be created as follows:

[Create a target](/apis/resources/action_service_v2/action-service-create-target)

```shell
curl -L -X POST 'https://$CUSTOM-DOMAIN/v2beta/actions/targets' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
  "name": "local call",
  "restCall": {
    "interruptOnError": true    
  },
  "endpoint": "http://localhost:8090/webhook",
  "timeout": "10s"
}'
```

Save the returned ID to set in the execution.

## Set execution

To call the target just created before, with the intention to manipulate the request used for user creation by the user V2 API, we define an execution with a method condition.

[Set an execution](/apis/resources/action_service_v2/action-service-set-execution)

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2beta/actions/executions' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "condition": {
        "request": {
            "method": "/zitadel.user.v2.UserService/AddHumanUser"
        }
    },
    "targets": [
        {
            "target": "<TargetID returned>"
        }
    ]
}'
```

## Example call

Now on every call on `/zitadel.user.v2.UserService/AddHumanUser` the local server adds metadata to the request:

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2/users/human' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "profile": {
        "givenName": "Example_given",
        "familyName": "Example_family"
    },
    "email": {
        "email": "example@example.com"
    }
}'
```

Resulting in a request like this:

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2/users/human' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "profile": {
        "givenName": "Example_given",
        "familyName": "Example_family"
    },
    "email": {
        "email": "example@example.com"
    }
    "metadata": [
        {"key": "organization", "value": "Y29tcGFueQ=="}
    ]
}'
```


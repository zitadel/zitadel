---
title: Test Actions Response ManipulationLocally
---

In this guide, you will create a ZITADEL execution and target. After an intent is retrieved through the API, the target is called.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

## Start example target

To start a simple HTTP server locally, which receives the call, the following code example can be used:

```go
package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

type response struct {
	Request *user.RetrieveIdentityProviderIntentRequest `json:"request"`
	Response *user.RetrieveIdentityProviderIntentResponse `json:"response"`
}

// call HandleFunc to read the response body, manipulate the content and return the response
func call(w http.ResponseWriter, req *http.Request) {
	// read the body content
	sentBody, err := io.ReadAll(req.Body)
	if err != nil {
		// if there was an error while reading the body return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// read the response into the expected structure
	request := new(response)
	if err := json.Unmarshal(sentBody, request); err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}

	// build the response from the received response
	resp := request.Response
	// manipulate the received response to send back as response
	if resp != nil && resp.AddHumanUser != nil {
		// manipulate the response
		resp.AddHumanUser.Metadata = append(resp.AddHumanUser.Metadata, &user.SetMetadataEntry{Key: "organization", Value: []byte("company")})
	}

	// marshal the response into json
	data, err := json.Marshal(resp)
	if err != nil {
		// if there was an error while marshalling the json
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// return the manipulated response
	w.Write(data)
}

func main() {
	// handle the HTTP call under "/call"
	http.HandleFunc("/call", call)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
```

What happens here is the response is received as a request, manipulated and then returned.

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

[Create a target](/apis/resources/action_service_v2/zitadel-actions-create-target)

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
  "endpoint": "http://localhost:8090/call",
  "timeout": "10s"
}'
```

Save the returned ID to set in the execution.

## Set execution

To call the target just created before, with the intention to manipulate the retrieve of an intent by the user V2 API, we define an execution with a method condition.

[Set an execution](/apis/resources/action_service_v2/zitadel-actions-set-execution)

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2beta/actions/executions' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "condition": {
        "response": {
            "method": "/zitadel.user.v2.UserService/RetrieveIdentityProviderIntent"
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

Now on every call on `/zitadel.user.v2.UserService/RetrieveIdentityProviderIntent` the local server would return something like:

```json
TODO
```

Resulting in a response like this:
```json
TODO
```



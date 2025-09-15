---
title: Test Actions Request Manipulation
---

This guide shows you how to leverage the ZITADEL actions feature to manipulate API requests in your ZITADEL instance.
You can use the actions feature to create a target that will be called when a specific API request occurs.
This is useful for adding information to managed resources in ZITADEL.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

:::info
Note that this guide assumes that ZITADEL is running on the same machine as the target and can be reached via `localhost`.
In case you are using a different setup, you need to adjust the target URL accordingly and will need to make sure that the target is reachable from ZITADEL.
:::

:::warning
To marshal and unmarshal the request please use a package like [protojson](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson),
as the request is a protocol buffer message, to avoid potential problems with the attribute names.
:::

## Start example target

To test the actions feature, you need to create a target that will be called when an API endpoint is called.
You will need to implement a listener that can receive HTTP requests, process the request and returns the manipulated request.
For this example, we will use a simple Go HTTP server that will return the request with added metadata.

:::info
The signature of the received request can be checked, [please refer to the example for more information on how to](/guides/integrate/actions/testing-request-signature).
:::

```go
package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

type contextRequest struct {
	Request *addHumanUserRequestWrapper `json:"request"`
}

// addHumanUserRequestWrapper necessary to marshal and unmarshal the JSON into the proto message correctly
type addHumanUserRequestWrapper struct {
	user.AddHumanUserRequest
}

func (r *addHumanUserRequestWrapper) MarshalJSON() ([]byte, error) {
	data, err := protojson.Marshal(r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *addHumanUserRequestWrapper) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, r)
}

// call HandleFunc to read the request body, manipulate the content and return the manipulated request
func call(w http.ResponseWriter, req *http.Request) {
	// read the body content
	sentBody, err := io.ReadAll(req.Body)
	if err != nil {
		// if there was an error while reading the body return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// read the request into the expected structure
	request := new(contextRequest)
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
	// handle the HTTP call under "/call"
	http.HandleFunc("/call", call)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}

```

:::info  
The example above runs only on your local machine (`localhost`).  
To test it with Zitadel, you must make your listener reachable from the internet.  
You can do this by using **Webhook.site** (see [Creating a Listener with Webhook.site](./webhook-site-setup)).  
:::

## Create target

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as call, the target can be created as follows:

See [Create a target](/apis/resources/action_service_v2/action-service-create-target) for more detailed information.

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

To call the target just created before, with the intention to manipulate the request used for user creation by the user V2 API, we define an execution with a method condition.

See [Set an execution](/apis/resources/action_service_v2/action-service-set-execution) for more detailed information.

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2beta/actions/executions' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "condition": {
        "request": {
            "method": "/zitadel.user.v2.UserService/CreateUser"
        }
    },
    "targets": [
        "<TargetID returned>"
    ]
}'
```

## Example call

Now that you have set up the target and execution, you can test it by creating a user through the Console UI or
by calling the ZITADEL API to create a human user.

```shell
curl -L -X POST 'https://$CUSTOM-DOMAIN/v2/users/new' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "organizationId": "336392597046099971",
    "human":
    {
        "profile":
        {
            "givenName": "Minnie",
            "familyName": "Mouse",
            "nickName": "Mini",
            "displayName": "Minnie Mouse",
            "preferredLanguage": "en",
            "gender": "GENDER_FEMALE"
        },
        "email":
        {
            "email": "mini@mouse.com"
        }
    }
}'
```

Your server should now manipulate the request to something like the following. Check out
the [Sent information Request](./usage#sent-information-request) payload description.

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2/users/new' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "organizationId": "336392597046099971",
    "human":
    {
        "profile":
        {
            "givenName": "Minnie",
            "familyName": "Mouse",
            "nickName": "Mini",
            "displayName": "Minnie Mouse",
            "preferredLanguage": "en",
            "gender": "GENDER_FEMALE"
        },
        "email":
        {
            "email": "mini@mouse.com"
        }
    }
    "metadata": [
        {"key": "organization", "value": "Y29tcGFueQ=="}
    ]
}'
```

## Conclusion

You have successfully set up a target and execution to manipulate API requests in your ZITADEL instance.
This feature can now be used to add or manipulate information to managed resources in ZITADEL.
Find more information about the actions feature in the [API documentation](/concepts/features/actions_v2).

---
title: Test Actions Response Manipulation
---

This guide shows you how to leverage the ZITADEL actions feature to manipulate API responses in your ZITADEL instance.
You can use the actions feature to create a target that will be called when a specific API response occurs.
This is useful for triggering workflows based on API responses in ZITADEL. You can even use this to provide data necessary data to the new login UI as shown in this example.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

:::info
Note that this guide assumes that ZITADEL is running on the same machine as the target and can be reached via `localhost`.
In case you are using a different setup, you need to adjust the target URL accordingly and will need to make sure that the target is reachable from ZITADEL.
:::

:::warning
To marshal and unmarshal the request and response please use a package like [protojson](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson),
as the request and response are protocol buffer messages, to avoid potential problems with the attribute names.
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

type contextResponse struct {
	Request  *retrieveIdentityProviderIntentRequestWrapper  `json:"request"`
	Response *retrieveIdentityProviderIntentResponseWrapper `json:"response"`
}

// RetrieveIdentityProviderIntentRequestWrapper necessary to marshal and unmarshal the JSON into the proto message correctly
type retrieveIdentityProviderIntentRequestWrapper struct {
	user.RetrieveIdentityProviderIntentRequest
}

func (r *retrieveIdentityProviderIntentRequestWrapper) MarshalJSON() ([]byte, error) {
	data, err := protojson.Marshal(r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *retrieveIdentityProviderIntentRequestWrapper) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, r)
}

// RetrieveIdentityProviderIntentResponseWrapper necessary to marshal and unmarshal the JSON into the proto message correctly
type retrieveIdentityProviderIntentResponseWrapper struct {
	user.RetrieveIdentityProviderIntentResponse
}

func (r *retrieveIdentityProviderIntentResponseWrapper) MarshalJSON() ([]byte, error) {
	data, err := protojson.Marshal(r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *retrieveIdentityProviderIntentResponseWrapper) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, r)
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
	defer req.Body.Close()

	// read the response into the expected structure
	request := new(contextResponse)
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

## Create target

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as call, the
target can be created as follows:

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

To call the target just created before, with the intention to manipulate the retrieve of an intent by the user V2 API,
we define an execution with a response condition.

See [Set an execution](/apis/resources/action_service_v2/action-service-set-execution) for more detailed information.

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
        "<TargetID returned>"
    ]
}'
```

## Example call

Now that you have set up the target and execution, you can test it by using a login-flow in the typescript login with an external IDP.

```json    
{
  "details": {
    "sequence": "599",
    "changeDate": "2023-06-15T06:44:26.039444Z",
    "resourceOwner": "163840776835432705"
  },
  "idpInformation": {
    "oauth": {
      "accessToken": "ya29...",
      "idToken": "ey..."
    },
    "idpId": "218528353504723201",
    "userId": "218528353504723202",
    "username": "test-user@localhost",
    "rawInformation": {
      "User": {
        "email": "test-user@localhost",
        "email_verified": true,
        "family_name": "User",
        "given_name": "Test",
        "hd": "mouse.com",
        "locale": "de",
        "name": "Minnie Mouse",
        "picture": "https://lh3.googleusercontent.com/a/AAcKTtf973Q7NH8KzKTMEZELPU9lx45WpQ9FRBuxFdPb=s96-c",
        "sub": "111392805975715856637"
      }
    }
  },
  "addHumanUser": {
    "idpLinks": [
      {"idpId": "218528353504723201", "userId": "218528353504723202", "userName": "test-user@localhost"}
    ],
    "username": "test-user@localhost",
    "profile": {
      "givenName": "Test",
      "familyName": "User",
      "displayName": "Test User",
      "preferredLanguage": "de"
    },
    "email": {
      "email": "test-user@zitadel.ch",
      "isVerified": true
    },
    "metadata": []
  }
}
```

Your server should now manipulate the response to something like the following. Check out
the [Sent information Response](./usage#sent-information-response) payload description.

```json
{
  "details": {
    "sequence": "599",
    "changeDate": "2023-06-15T06:44:26.039444Z",
    "resourceOwner": "163840776835432705"
  },
  "idpInformation": {
    "oauth": {
      "accessToken": "ya29...",
      "idToken": "ey..."
    },
    "idpId": "218528353504723201",
    "userId": "218528353504723202",
    "username": "test-user@localhost",
    "rawInformation": {
      "User": {
        "email": "test-user@localhost",
        "email_verified": true,
        "family_name": "User",
        "given_name": "Test",
        "hd": "mouse.com",
        "locale": "de",
        "name": "Minnie Mouse",
        "picture": "https://lh3.googleusercontent.com/a/AAcKTtf973Q7NH8KzKTMEZELPU9lx45WpQ9FRBuxFdPb=s96-c",
        "sub": "111392805975715856637"
      }
    }
  },
  "addHumanUser": {
    "idpLinks": [
      {"idpId": "218528353504723201", "userId": "218528353504723202", "userName": "test-user@localhost"}
    ],
    "username": "test-user@localhost",
    "profile": {
      "givenName": "Test",
      "familyName": "User",
      "displayName": "Test User",
      "preferredLanguage": "de"
    },
    "email": {
      "email": "test-user@zitadel.ch",
      "isVerified": true
    },
    "metadata": [
      {"key": "organization", "value": "Y29tcGFueQ=="}
    ]
  }
}
```

## Conclusion

You have successfully set up a target and execution to manipulate API responses in your ZITADEL instance.
This feature can now be used to add necessary information for clients including the new login UI.
Find more information about the actions feature in the [API documentation](/concepts/features/actions_v2).

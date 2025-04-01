---
title: Test Actions Response Manipulation
---

In this guide, you will create a ZITADEL execution and target. After an intent is retrieved through the API, the target
is called.

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
	Request  *user.RetrieveIdentityProviderIntentRequest  `json:"request"`
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

## Create target

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as call, the
target can be created as follows:

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

To call the target just created before, with the intention to manipulate the retrieve of an intent by the user V2 API,
we define an execution with a method condition.

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

Now on every call on `/zitadel.user.v2.UserService/RetrieveIdentityProviderIntent` the local server would return
something like:

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

Resulting in a response like this:

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



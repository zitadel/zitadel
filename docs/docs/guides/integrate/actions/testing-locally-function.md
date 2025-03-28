---
title: Test Actions Function Locally
---

In this guide, you will create a ZITADEL execution and target. To add claims to a token, the target is called.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

## Start example target

To start a simple HTTP server locally, which receives the call and sends back a response, the following code example can be used:

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	SetUserMetadata []*Metadata    `json:"set_user_metadata,omitempty"`
	AppendClaims    []*AppendClaim `json:"append_claims,omitempty"`
	AppendLogClaims []string       `json:"append_log_claims,omitempty"`
}

type Metadata struct {
	Key   string
	Value []byte
}

type AppendClaim struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// webhook HandleFunc to read the request body and then print out the contents
func call(w http.ResponseWriter, req *http.Request) {
	// read the body content
	sentBody, err := io.ReadAll(req.Body)
	if err != nil {
		// if there was an error while reading the body return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	// print out the read content
	fmt.Println(string(sentBody))
    
	// create the response with the correct structure
	resp := &Response{
		SetUserMetadata: []*Metadata{
			{Key: "key", Value: []byte("value")},
		},
		AppendClaims: []*AppendClaim{
			{Key: "claim", Value: "value"},
		},
		AppendLogClaims: []string{"log1", "log2", "log3"},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		// if there was an error while marshalling the json
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func main() {
	// handle the HTTP call under "/webhook"
	http.HandleFunc("/call", call)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
```

What happens here is that the user get the metadata with the key "key" and value "value" added, the token gets a claim "urn:zitadel:iam:claim" with value "value" and the log claim "urn:zitadel:iam:action:preuserinfo:log" with values "log1", "log2" and "log3".

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

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as webhook, the target can be created as follows:

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

To call the target just created before, with the intention to print the response from a user creation by the user V2 API, we define an execution with a method condition.

[Set an execution](/apis/resources/action_service_v2/zitadel-actions-set-execution)

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2beta/actions/executions' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "condition": {
        "function": {
            "name": "preuserinfo"
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

Now on every OIDC flow this action will get executed.

Should print out something like, also described under [Sent information Function](./usage#sent-information-function):
```json
{
  "function" : "function/preuserinfo",
  "userinfo" : {
    "sub" : "312909075212468632"
  },
  "user" : {
    "id" : "312909075212468632",
    "creation_date" : "2025-03-26T15:52:23.917636Z",
    "change_date" : "2025-03-26T15:52:23.917636Z",
    "resource_owner" : "312909075211944344",
    "sequence" : 2,
    "state" : 1,
    "username" : "user@example.com",
    "preferred_login_name" : "zitadel@zitadel.localhost",
    "human" : {
      "first_name" : "Example firstname",
      "last_name" : "Example lastname",
      "display_name" : "Example displayname",
      "preferred_language" : "en",
      "email" : "user@example.com",
      "is_email_verified" : true,
      "password_changed" : "0001-01-01T00:00:00Z",
      "mfa_init_skipped" : "0001-01-01T00:00:00Z"
    }
  },
  "user_metadata" : [ {
    "creation_date" : "2025-03-27T09:10:25.879677Z",
    "change_date" : "2025-03-27T09:10:25.879677Z",
    "resource_owner" : "312909075211944344",
    "sequence" : 18,
    "key" : "key",
    "value" : "dmFsdWU="
  } ],
  "org" : {
    "id" : "312909075211944344",
    "name" : "ZITADEL",
    "primary_domain" : "example.com"
  }
}
```



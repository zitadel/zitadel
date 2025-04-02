---
title: Test Actions Function
---

This guide shows you how to leverage the ZITADEL actions feature to enhance different functions in your ZITADEL instance.
You can use the actions feature to create a target that will be called when a specific functionality is used.
This is useful for integrating with other systems which need specific claims in tokens or for executing external code during OIDC or SAML flows.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

:::info
Note that this guide assumes that ZITADEL is running on the same machine as the target and can be reached via `localhost`.
In case you are using a different setup, you need to adjust the target URL accordingly and will need to make sure that the target is reachable from ZITADEL.
:::

## Start example target

To test the actions feature, you need to create a target that will be called when a function is used.
You will need to implement a listener that can receive HTTP requests and process the data.
For this example, we will use a simple Go HTTP server that will print the received data to standard output.

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

// call HandleFunc to read the request body and then print out the contents
func call(w http.ResponseWriter, req *http.Request) {
	// read the body content
	sentBody, err := io.ReadAll(req.Body)
	if err != nil {
		// if there was an error while reading the body return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
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
	// handle the HTTP call under "/call"
	http.HandleFunc("/call", call)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
```

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

To configure ZITADEL to call the target when a function is executed, you need to set an execution and define the function
condition.

See [Set an execution](/apis/resources/action_service_v2/action-service-set-execution) for more detailed information.

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

Now that you have set up the target and execution, you can test it by logging into Console UI or
by using any OIDC flow.

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

## Conclusion

You have successfully set up a target and execution to react to functions in your ZITADEL instance.
This feature can now be used to integrate with your existing systems to create custom workflows or automate tasks based on functionality in ZITADEL.
Find more information about the actions feature in the [API documentation](/concepts/features/actions_v2).

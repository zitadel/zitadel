---
title: Test Actions Response
---

This guide shows you how to leverage the ZITADEL actions feature to react to API responses in your ZITADEL instance.
You can use the actions feature to create a target that will be called when a specific API response occurs.
This is useful for information provisioning in between systems or for triggering workflows based on API responses in ZITADEL.

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
You will need to implement a listener that can receive HTTP requests and process the request.
For this example, we will use a simple Go HTTP server that will print the received request to standard output.

:::info
The signature of the received request can be checked, [please refer to the example for more information on how to](/guides/integrate/actions/testing-request-signature).
:::

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

// webhook HandleFunc to read the request body and then print out the contents
func webhook(w http.ResponseWriter, req *http.Request) {
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
}

func main() {
	// handle the HTTP call under "/webhook"
	http.HandleFunc("/webhook", webhook)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
```

## Create target

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as webhook, the target can be created as follows:

See [Create a target](/apis/resources/action_service_v2/action-service-create-target) for more detailed information.

```shell
curl -L -X POST 'https://$CUSTOM-DOMAIN/v2beta/actions/targets' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
  "name": "local webhook",
  "restWebhook": {
    "interruptOnError": true    
  },
  "endpoint": "http://localhost:8090/webhook",
  "timeout": "10s"
}'
```

Save the returned ID to set in the execution.

## Set execution

To configure Zitadel to call the target when an API endpoint is called, you need to set an execution and define the response
condition.

See [Set an execution](/apis/resources/action_service_v2/action-service-set-execution) for more detailed information.

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2beta/actions/executions' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "condition": {
        "response": {
            "method": "/zitadel.user.v2.UserService/AddHumanUser"
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
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2/users/human' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "userId": {
        "givenName": "Example_given",
        "familyName": "Example_family"
    },
    "email": {
        "email": "example@example.com"
    }
}'
```

Your server should now print out something like the following. Check out
the [Sent information Response](./usage#sent-information-response) payload description.

```json
{
  "fullMethod": "/zitadel.user.v2.UserService/AddHumanUser",
  "instanceID": "262851882718855632",
  "orgID": "262851882718921168",
  "projectID": "262851882719052240",
  "userID": "262851882718986704",
  "request": {
    "profile": {
      "given_name": "Example_given",
      "family_name": "Example_family"
    },
    "email": {
      "email": "example@example.com"
    }
  },
  "response": {
    "user_id": "312918757460672920",
    "details": {
        "sequence": "2",
        "change_date": "2025-03-26T17:28:33.856436Z",
        "resource_owner": "312909075211944344",
    }
  }
}
```

## Conclusion

You have successfully set up a target and execution to react to API responses in your ZITADEL instance.
This feature can now be used to provision information in between systems or for triggering workflows based on API responses in ZITADEL.
Find more information about the actions feature in the [API documentation](/concepts/features/actions_v2).

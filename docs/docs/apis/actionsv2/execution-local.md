---
title: Actions v2 example execution locally
---

In this guide, you will create a ZITADEL execution and target. After a user is created through the API, the target is called.

## Prerequisites

Before you start, make sure you have everything set up correctly.

- You need to be at least a ZITADEL [_IAM_OWNER_](/guides/manage/console/managers)
- Your ZITADEL instance needs to have the actions feature enabled.

## Start example target

To start a simple HTTP server locally, which receives the webhook call, the following code example can be used:

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

What happens here is only a target which prints out the received request, which could also be handled with a different logic.

## Create target

As you see in the example above the target is created with HTTP and port '8090' and if we want to use it as webhook, the target can be created as follows:

[Create a target](/apis/resources/action_service_v3/action-service-create-target)

```shell
curl -L -X POST 'https://$CUSTOM-DOMAIN/v3alpha/targets' \
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

To call the target just created before, with the intention to print the request used for user creation by the user V2 API, we define an execution with a method condition.

[Set an execution](/apis/resources/action_service_v3/action-service-set-execution)

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v3alpha/executions' \
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

Now on every call on `/zitadel.user.v2.UserService/AddHumanUser` the local server prints out the received body of the request:

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

Should print out something like, also described under [Sent information Request](./introduction#sent-information-request):
```shell
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
  }
}
```



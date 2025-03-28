---
title: Test Actions Event Locally
---

In this guide, you will create a ZITADEL execution and target. As a result of the event creation, the target is called.

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

To call the target just created before, with the intention to print the response from a user creation by the user V2 API, we define an execution with a method condition.

[Set an execution](/apis/resources/action_service_v2/zitadel-actions-set-execution)

```shell
curl -L -X PUT 'https://$CUSTOM-DOMAIN/v2beta/actions/executions' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
    "condition": {
        "event": {
            "event": "user.human.added"
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

Now on every call on `/zitadel.user.v2.UserService/AddHumanUser` the local server prints out the event:

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

Should print out something like, also described under [Sent information Event](./usage#sent-information-event):
```json
{
  "aggregateID" : "313014806065971608",
  "aggregateType" : "user",
  "resourceOwner" : "312909075211944344",
  "instanceID" : "312909075211878808",
  "version" : "v2",
  "sequence" : 1,
  "event_type" : "user.human.added",
  "created_at" : "2025-03-27T10:22:43.262665+01:00",
  "userID" : "312909075212468632",
  "event_payload" : "eyJ1c2VyTmFtZSI6ImV4YW1wbGVAdGVzdC5jb20iLCJmaXJzdE5hbWUiOiJ0ZXN0IiwibGFzdE5hbWUiOiJ0ZXN0IiwiZGlzcGxheU5hbWUiOiJ0ZXN0IHRlc3QiLCJwcmVmZXJyZWRMYW5ndWFnZSI6InVuZCIsImVtYWlsIjoiZXhhbXBsZUB0ZXN0LmNvbSJ9"
}
```

The event_payload is base64 encoded and has the following content:
```json
{
  "userName": "example@test.com",
  "firstName": "test",
  "lastName": "test",
  "displayName": "test test",
  "preferredLanguage": "und",
  "email": "example@test.com"
}
```



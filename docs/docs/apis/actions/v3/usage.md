---
title: Using Actions
---

The Action API provides a flexible mechanism for customizing and extending the functionality of ZITADEL. By allowing you to define targets and executions, you can implement custom workflows triggered on an API requests and responses, events or specific functions.

**How it works:**
- Create Target
- Set Execution with condition and target
- Custom Code will be triggered and executed

**Use Cases:**
- User Management: Automate provisioning user data to external systems when users are crreated, updated or deleted.
- Security: Implement IP blocking or rate limiting based on API usage patterns.
- Extend Workflows: Automatically setup resources in your application, when a new organization in ZITADEL is created. 
- Token extension: Add custom claims to the tokens.

## Endpoints

ZITADEL sends an HTTP Post request to the endpoint set as Target, the received request than can be edited and send back or custom processes can be handled.

### Sent information Request

The information sent to the Endpoint is structured as JSON:

```json
{
  "fullMethod": "full method of the GRPC call",
  "instanceID": "instanceID of the called instance",
  "orgID": "ID of the organization related to the calling context",
  "projectID": "ID of the project related to the used application",
  "userID": "ID of the calling user",
  "request": "full request of the call"
}
```

### Sent information Response

The information sent to the Endpoint is structured as JSON:

```json
{
  "fullMethod": "full method of the GRPC call",
  "instanceID": "instanceID of the called instance",
  "orgID": "ID of the organization related to the calling context",
  "projectID": "ID of the project related to the used application",
  "userID": "ID of the calling user",
  "request": "full request of the call",
  "response": "full response of the call"
}
```

## Target

The Target describes how ZITADEL interacts with the Endpoint.

There are different types of Targets:

- `Webhook`, the call handles the status code but response is irrelevant, can be InterruptOnError
- `Call`, the call handles the status code and response, can be InterruptOnError
- `Async`, the call handles neither status code nor response, but can be called in parallel with other Targets

`InterruptOnError` means that the Execution gets interrupted if any of the calls return with a status code >= 400, and the next Target will not be called anymore.

The API documentation to create a target can be found [here](/apis/resources/action_service_v3/zitadel-actions-create-target)

### Content Signing

To ensure the integrity of request content, each call includes a 'ZITADEL-Signature' in the headers. This header contains an HMAC value computed from the request content and a timestamp, which can be used to time out requests. The logic for this process is provided in 'pkg/actions/signing.go'. The goal is to verify that the HMAC value in the header matches the HMAC value computed by the Target, ensuring that the sent and received requests are identical.

Each Target resource now contains also a Signing Key, which gets generated and returned when a Target is [created](/apis/resources/action_service_v3/zitadel-actions-create-target),
and can also be newly generated when a Target is [patched](/apis/resources/action_service_v3/zitadel-actions-patch-target).

## Execution

ZITADEL decides on specific conditions if one or more Targets have to be called.
The Execution resource contains 2 parts, the condition and the called targets.

The condition can be defined for 4 types of processes:

- `Requests`, before a request is processed by ZITADEL
- `Responses`, before a response is sent back to the application
- `Functions`, handling specific functionality in the logic of ZITADEL
- `Events`, after a specific event happened and was stored in ZITADEL

The API documentation to set an Execution can be found [here](/apis/resources/action_service_v3/zitadel-actions-set-execution)

### Condition Best Match

As the conditions can be defined on different levels, ZITADEL tries to find out which Execution is the best match.
This means that for example if you have an Execution defined on `all requests`, on the service `zitadel.user.v2.UserService` and on `/zitadel.user.v2.UserService/AddHumanUser`,
ZITADEL would with a call on the `/zitadel.user.v2.UserService/AddHumanUser` use the Executions with the following priority:

1. `/zitadel.user.v2.UserService/AddHumanUser`
2. `zitadel.user.v2.UserService`
3. `all`

If you then have a call on `/zitadel.user.v2.UserService/UpdateHumanUser` the following priority would be found:

1. `zitadel.user.v2.UserService`
2. `all`

And if you use a different service, for example `zitadel.session.v2.SessionService`, then the `all` Execution would still be used.

### Targets and Includes

:::info
Includes are limited to 3 levels, which mean that include1->include2->include3 is the maximum for now.
If you have feedback to the include logic, or a reason why 3 levels are not enough, please open [an issue on github](https://github.com/zitadel/zitadel/issues) or [start a discussion on github](https://github.com/zitadel/zitadel/discussions)/[start a topic on discord](https://zitadel.com/chat)
:::

An execution can not only contain a list of Targets, but also Includes.
The Includes can be defined in the Execution directly, which means you include all defined Targets by a before set Execution.

If you define 2 Executions as follows:

```json
{
  "condition": {
    "request": {
      "service": "zitadel.user.v2.UserService"
    }
  },
  "targets": [
    {
      "target": "<TargetID1>"
    }
  ]
}
```

```json
{
  "condition": {
    "request": {
      "method": "/zitadel.user.v2.UserService/AddHumanUser"
    }
  },
  "targets": [
    {
      "target": "<TargetID2>"
    },
    {
      "include": {
        "request": {
          "service": "zitadel.user.v2.UserService"
        }
      }
    }
  ]
}
```

The called Targets on "/zitadel.user.v2.UserService/AddHumanUser" would be, in order:

1. `<TargetID2>`
2. `<TargetID1>`

### Condition for Requests and Responses

For Request and Response there are 3 levels the condition can be defined:

- `Method`, handling a request or response of a specific GRPC full method, which includes the service name and method of the ZITADEL API
- `Service`, handling any request or response under a service of the ZITADEL API
- `All`, handling any request or response under the ZITADEL API

The available conditions can be found under:
- [All available Methods](/apis/resources/action_service_v3/zitadel-actions-list-execution-methods), for example `/zitadel.user.v2.UserService/AddHumanUser`
- [All available Services](/apis/resources/action_service_v3/zitadel-actions-list-execution-services), for example `zitadel.user.v2.UserService`

### Condition for Functions

Replace the current Actions with the following flows:

- [Internal Authentication](/apis/actions/internal-authentication)
- [External Authentication](/apis/actions/external-authentication)
- [Complement Token](/apis/actions/complement-token)
- [Customize SAML Response](/apis/actions/customize-samlresponse)

The available conditions can be found under [all available Functions](/apis/resources/action_service_v3/zitadel-actions-list-execution-functions).

### Condition for Events

For event there are 3 levels the condition can be defined:

- Event, handling a specific event
- Group, handling a specific group of events
- All, handling any event in ZITADEL

The concept of events can be found under [Events](/concepts/architecture/software#events)

### Error forwarding

If you want to forward a specific error from the Target through ZITADEL, you can provide a response from the Target with status code 200 and a JSON in the following format:

```json
{
  "forwardedStatusCode": 403,
  "forwardedErrorMessage": "Call is forbidden through the IP AllowList definition"
}
```

Only values from 400 to 499 will be forwarded through ZITADEL, other StatusCodes will end in a PreconditionFailed error.

If the Target returns any other status code than >= 200 and < 299, the execution is looked at as failed, and a PreconditionFailed error is logged.

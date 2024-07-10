---
title: Actions V2
---

This page describes the options you have when defining ZITADEL Actions V2.

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

The API documentation to create a target can be found [here](/apis/resources/action_service_v3/action-service-create-target)

## Execution

ZITADEL decides on specific conditions if one or more Targets have to be called.
The Execution resource contains 2 parts, the condition and the called targets.

The condition can be defined for 4 types of processes:

- `Requests`, before a request is processed by ZITADEL
- `Responses`, before a response is sent back to the application
- `Functions`, handling specific functionality in the logic of ZITADEL
- `Events`, after a specific event happened and was stored in ZITADEL

The API documentation to set an Execution can be found [here](/apis/resources/action_service_v3/action-service-set-execution)

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
- [All available Methods](/apis/resources/action_service_v3/action-service-list-execution-methods), for example `/zitadel.user.v2.UserService/AddHumanUser`
- [All available Services](/apis/resources/action_service_v3/action-service-list-execution-services), for example `zitadel.user.v2.UserService`

### Condition for Functions

Replace the current Actions with the following flows:

- [Internal Authentication](../actions/internal-authentication)
- [External Authentication](../actions/external-authentication)
- [Complement Token](../actions/complement-token)
- [Customize SAML Response](../actions/customize-samlresponse)

The available conditions can be found under [all available Functions](/apis/resources/action_service_v3/action-service-list-execution-functions).

### Condition for Events

For event there are 3 levels the condition can be defined:

- Event, handling a specific event
- Group, handling a specific group of events
- All, handling any event in ZITADEL

The concept of events can be found under [Events](/concepts/architecture/software#events)
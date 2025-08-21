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
  "request": {
    "attribute": "Attribute value of full request of the call"
  }
}
```

:::warning
To marshal and unmarshal the request please use a package like [protojson](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson), 
as the request is a protocol buffer message, to avoid potential problems with the attribute names.
:::

### Sent information Response

The information sent to the Endpoint is structured as JSON:

```json
{
  "fullMethod": "full method of the GRPC call",
  "instanceID": "instanceID of the called instance",
  "orgID": "ID of the organization related to the calling context",
  "projectID": "ID of the project related to the used application",
  "userID": "ID of the calling user",
  "request": {
    "attribute": "Attribute value of full request of the call"
  },
  "response": {
    "attribute": "Attribute value of full response of the call"
  }
}
```

:::warning
To marshal and unmarshal the request and response please use a package like [protojson](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson),
as the request and response are protocol buffer messages, to avoid potential problems with the attribute names.
:::

### Sent information Function

Information sent and expected back are specific to the function.

#### PreUserinfo

The information sent to the Endpoint is structured as JSON:
```json
{
  "function": "Name of the function",
  "userinfo": {
    "given_name": "",
    "family_name": "",
    "middle_name": "",
    "nickname": "",
    "profile": "",
    "picture": "",
    ...
    "preferred_username": "",
    "email": "",
    "email_verified": true,
    "phone_number": "",
    "phone_number_verified": true
  },
  "user": {
    "id": "",
    "creation_date": "",
    ...
    "human": {
      "first_name": "",
      "last_name": "",
      ...
      "email": "",
      "is_email_verified": true,
      "phone": "",
      "is_phone_verified": true
    }
  },
  "user_metadata": [
    {
      "creation_date": "",
      "change_date": "",
      "resource_owner": "",
      "sequence": "",
      "key": "",
      "value": ""
    }
  ],
  "org": {
    "id": "ID of the organization the user belongs to",
    "name": "Name of the organization the user belongs to",
    "primary_domain": "Primary domain of the organization the user belongs to"
  },
  "user_grants": [
    {
      "id": "",
      "projectGrantId": "The ID of the project grant",
      "state": 1,
      "creationDate": "",
      "changeDate": "",
      "sequence": 1,
      "userId": "",
      "roles": [
        "role"
      ],
      "userResourceOwner": "The ID of the organization the user belongs to",
      "userGrantResourceOwner": "The ID of the organization the user got authorization granted",
      "userGrantResourceOwnerName": "The name of the organization the user got authorization granted",
      "projectId": "",
      "projectName": ""
    }
  ]
}
```

The expected structure of the JSON as response:

```json
{
  "set_user_metadata": [
    {
      "key": "key of metadata to be set on the user",
      "value": "base64 value of metadata to be set on the user"
    }
  ],
  "append_claims": [
    {
      "key": "key of claim to be set on the user",
      "value": "value of claim to be set on the user"
    }
  ],
  "append_log_claims": [
    "Log to be appended to the log claim on the token"
  ]
}
```

#### PreAccessToken

The information sent to the Endpoint is structured as JSON:

```json
{
  "function": "Name of the function",
  "userinfo": {
    "given_name": "",
    "family_name": "",
    "middle_name": "",
    "nickname": "",
    "profile": "",
    "picture": "",
    ...
    "preferred_username": "",
    "email": "",
    "email_verified": true/false,
    "phone_number": "",
    "phone_number_verified": true/false
  },
  "user": {
    "id": "",
    "creation_date": "",
    ...
    "human": {
      "first_name": "",
      "last_name": "",
      ...
      "email": "",
      "is_email_verified": true,
      "phone": "",
      "is_phone_verified": true
    }
  },
  "user_metadata": [
    {
      "creation_date": "",
      "change_date": "",
      "resource_owner": "",
      "sequence": "",
      "key": "",
      "value": ""
    }
  ],
  "org": {
    "id": "ID of the organization the user belongs to",
    "name": "Name of the organization the user belongs to",
    "primary_domain": "Primary domain of the organization the user belongs to"
  },
  "user_grants": [
    {
      "id": "",
      "projectGrantId": "The ID of the project grant",
      "state": 1,
      "creationDate": "",
      "changeDate": "",
      "sequence": 1,
      "userId": "",
      "roles": [
        "role"
      ],
      "userResourceOwner": "The ID of the organization the user belongs to",
      "userGrantResourceOwner": "The ID of the organization the user got authorization granted",
      "userGrantResourceOwnerName": "The name of the organization the user got authorization granted",
      "projectId": "",
      "projectName": ""
    }
  ]
}
```

The expected structure of the JSON as response:

```json
{
  "set_user_metadata": [
    {
      "key": "key of metadata to be set on the user",
      "value": "base64 value of metadata to be set on the user"
    }
  ],
  "append_claims": [
    {
      "key": "key of claim to be set on the user",
      "value": "value of claim to be set on the user"
    }
  ],
  "append_log_claims": [
    "Log to be appended to the log claim on the token"
  ]
}
```

#### PreSAMLResponse

The information sent to the Endpoint is structured as JSON:
```json
{
  "function": "Name of the function",
  "userinfo": {
    "given_name": "",
    "family_name": "",
    "middle_name": "",
    "nickname": "",
    "profile": "",
    "picture": "",
    ...
    "preferred_username": "",
    "email": "",
    "email_verified": true,
    "phone_number": "",
    "phone_number_verified": true
  },
  "user": {
    "id": "",
    "creation_date": "",
    ...
    "human": {
      "first_name": "",
      "last_name": "",
      ...
      "email": "",
      "is_email_verified": true,
      "phone": "",
      "is_phone_verified": true
    }
  },
  "user_grants": [
    {
      "id": "",
      "projectGrantId": "The ID of the project grant",
      "state": 1,
      "creationDate": "",
      "changeDate": "",
      "sequence": 1,
      "userId": "",
      "roles": [
        "role"
      ],
      "userResourceOwner": "The ID of the organization the user belongs to",
      "userGrantResourceOwner": "The ID of the organization the user got authorization granted",
      "userGrantResourceOwnerName": "The name of the organization the user got authorization granted",
      "projectId": "",
      "projectName": ""
    }
  ]
}
```

The expected structure of the JSON as response:

```json
{
  "set_user_metadata": [
    {
      "key": "key of metadata to be set on the user",
      "value": "base64 value of metadata to be set on the user"
    }
  ],
  "append_attribute": [
    {
      "name": "name of the attribute to be added to the response",
      "name_format": "name format of the attribute to be added to the response",
      "value": "value of the attribute to be added to the response"
    }
  ]
}
```

### Sent information Event

The information sent to the Endpoint is structured as JSON:

```json
{
  "aggregateID": "ID of the aggregate",
  "aggregateType": "Type of the aggregate",
  "resourceOwner": "Resourceowner the aggregate belongs to",
  "instanceID": "ID of the instance the aggregate belongs to",
  "version": "Version of the aggregate",
  "sequence": "Sequence of the event",
  "event_type": "Type of the event",
  "created_at": "Time the event was created",
  "userID": "ID of the creator of the event",
  "event_payload": "Content of the event in JSON format"
}
```

## Target

The Target describes how ZITADEL interacts with the Endpoint.

There are different types of Targets:

- `Webhook`, the call handles the status code but response is irrelevant, can be InterruptOnError
- `Call`, the call handles the status code and response, can be InterruptOnError
- `Async`, the call handles neither status code nor response, but can be called in parallel with other Targets

`InterruptOnError` means that the Execution gets interrupted if any of the calls return with a status code >= 400, and the next Target will not be called anymore.

The API documentation to create a target can be found [here](/apis/resources/action_service_v2/action-service-create-target)

### Content Signing

To ensure the integrity of request content, each call includes a 'ZITADEL-Signature' in the headers. This header contains an HMAC value computed from the request content and a timestamp, which can be used to time out requests. The logic for this process is provided in 'pkg/actions/signing.go'. The goal is to verify that the HMAC value in the header matches the HMAC value computed by the Target, ensuring that the sent and received requests are identical.

Each Target resource now contains also a Signing Key, which gets generated and returned when a Target is [created](/apis/resources/action_service_v2/action-service-create-target),
and can also be newly generated when a Target is [patched](/apis/resources/action_service_v2/action-service-update-target).

For an example on how to check the signature, [refer to the example](/guides/integrate/actions/testing-request-signature).

## Execution

ZITADEL decides on specific conditions if one or more Targets have to be called.
The Execution resource contains 2 parts, the condition and the called targets.

The condition can be defined for 4 types of processes:

- `Requests`, before a request is processed by ZITADEL
- `Responses`, before a response is sent back to the application
- `Functions`, handling specific functionality in the logic of ZITADEL
- `Events`, after a specific event happened and was stored in ZITADEL

The API documentation to set an Execution can be found [here](/apis/resources/action_service_v2/action-service-set-execution)

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

### Targets

An execution can contain only a list of Targets, and Targets are comma separated string values.

Here's an example of a Target defined on a service (e.g. `zitadel.user.v2.UserService`)

```json
{
  "condition": {
    "request": {
      "service": "zitadel.user.v2.UserService"
    }
  },
  "targets": [
    "<TargetID1>"
  ]
}
```

Here's an example of a Target defined on a method (e.g. `/zitadel.user.v2.UserService/AddHumanUser`)
```json
{
  "condition": {
    "request": {
      "method": "/zitadel.user.v2.UserService/AddHumanUser"
    }
  },
  "targets": [
    "<TargetID2>",
    "<TargetID1>"
  ]
}
```

The called Targets on `/zitadel.user.v2.UserService/AddHumanUser` would be, in order:

1. `<TargetID2>`
2. `<TargetID1>`

### Condition for Requests and Responses

For Request and Response there are 3 levels the condition can be defined:

- `Method`, handling a request or response of a specific GRPC full method, which includes the service name and method of the ZITADEL API
- `Service`, handling any request or response under a service of the ZITADEL API
- `All`, handling any request or response under the ZITADEL API

The available conditions can be found under:
- [All available Methods](/apis/resources/action_service_v2/action-service-list-execution-methods), for example `/zitadel.user.v2.UserService/AddHumanUser`
- [All available Services](/apis/resources/action_service_v2/action-service-list-execution-services), for example `zitadel.user.v2.UserService`

### Condition for Functions

The available conditions can be found under [all available Functions](/apis/resources/action_service_v2/action-service-list-execution-functions).

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

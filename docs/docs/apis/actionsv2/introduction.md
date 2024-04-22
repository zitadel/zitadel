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

- Webhook, the call handles the status code but response is irrelevant, can be InterruptOnError
- RequestResponse, the call handles the status code and response, can be InterruptOnError
- Async, the call handles neither status code nor response, but can be called in parallel with other Targets

InterruptOnError means that the list of Targets gets interrupted if any of the calls return with a status code >= 400.

## Execution

ZITADEL decides on specific conditions if one or more Target have to be called.
The Execution resource contains 2 parts, the condition and the called targets, which can either be targets specificly 
or include to add the targets of another defined execution.

The condition can be defined for 4 types of processes:

- Requests, before a request is processed by ZITADEL
- Responses, before a response is sent back to the application
- Functions, handling specific functionality in the logic of ZITADEL
- Events, before a specific event is written to the eventstore in ZITADEL

### Condition for Requests and Responses

For Request and Response there are 3 levels the condition can be defined:

- Method, handling a request or response of a specific GRPC full method, which includes the service name and method of the ZITADEL API
- Service, handling any request or response under a service of the ZITADEL API
- All, handling any request or response under the ZITADEL API

### Condition for Functions

Replace the current Actions with the following flows:

- [Internal Authentication](./internal-authentication)
- [External Authentication](./external-authentication)
- [Complement Token](./complement-token)
- [Customize SAML Response](./customize-samlresponse)

### Condition for Events

For event there are 3 levels the condition can be defined:

- Event, handling a specific event
- Group, handling a specific group of events
- All, handling any event in ZITADEL


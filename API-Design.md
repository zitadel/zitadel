# API Design

## The Basics
ZITADEL follows an API first approach. This means that the API is designed before the implementation. 
The API is designed using the Protobuf specification. The Protobuf specification is then used to generate the API client and server code in different programming languages.

Starting with the V2 API, the API and its services use a resource-oriented design. 
This means that the API is designed around resources, which are the key entities in the system. Each resource has a unique identifier and a set of properties that describe the resource.
Resources can be created, read, updated, and deleted using the API.

### Protobuf, gRPC and connectRPC

## Conventions

### Explicitness

Make the handling of the API as explicit as possible. Do not make assumptions about the client's knowledge of the system or the API. 
Provide clear and concise documentation for the API.
Do not rely on implicit fallbacks or defaults if the client does not provide certain parameters.
Only use defaults if they are explicitly documented, such as returning a result set for the whole instance if no filter is provided.

### Naming Conventions

Names of resources, fields and methods should be descriptive and consistent.
Use domain-specific terminology and avoid abbreviations.
For example, use `OrganizationID` instead of OrgID or resourceOwner for the creation of a mew user or when returning one.

//TODO: add link to naming conventions / naming guidelines (https://github.com/zitadel/zitadel/issues/5888)

#### Resources and Fields

When a context is required for creating a resource, the context is added as a field to the resource.
For example, when creating a new user, the organization ID is required. The `organization_id` is added as a field to the `CreateUserRequest`.

Only allow providing a context where it is required. Do not provide the possibility to provide a context where it is not required.
For example, when retrieving or updating a user, the organization ID is not required, since the user can be determined by the user ID.
However, it is possible to provide the organization ID as a filter to retrieve a list of users.

Prevent the creation of global messages that are used in multiple resources unless they always follow the same pattern.
Use dedicated fields as described above or create a separate message for the specific context, that is only used in the boundary of the same resource.
For example, settings might be set as a default on the instance level, but might be overridden on the organization level.
In this case, the settings could share the same `SettingsContext` message to determine the context of the settings.
But do not create a global `Context` message that is used across the whole API if there are different scenarios and different fields required for the context.

#### Operations and Methods

Methods on a resource should be named using the following convention:
- Create: `Create<resource>`
- Update: `Update<resource>`
- Delete: `Delete<resource>`
- Get: `Get<resource>`
- List: `List<resource>`
- Search: `Search<resource>`

Methods on a list of resources should be named using the following convention:
- Add: `Add<resource>`
- Remove: `Remove<resource>`
- Set: `Set<resource>`

## Error Handling

The API returns machine-readable errors in the response body. This includes a status code, an error code and possibly some details.

### Status Codes

The API uses status codes to indicate the status of a request. Depending on the protocol used to call the API,
the status code is returned as an HTTP status code or as a gRPC / connectRPC status code.
Check the possible status codes https://zitadel.com/docs/apis/statuscodes

### Error Codes

Additionally to the status code, the API returns unique error codes for each type of error.
The error codes are used to identify a specific error and can be used to handle the error programmatically.

TODO: Add error codes schema and examples

### Error Message and Details

The API returns additional details about the error in the response body.
This includes a human-readable error message and additional information that can help the client to understand the error
as well as machine-readable details that can be used to handle the error programmatically.
Error details use the Google RPC error details format: https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto

### Example

HTTP/1.1 example:
```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "code": "USER-001",
  "message": "user requires a username",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.BadRequest",
      "fieldViolations": [
        {
          "field": "username",
          "description": "Username is required"
        }
      ]
    }
  ]
}
```

gRPC / connectRPC example:
```
HTTP/2.0 200 Bad Request
Content-Type: application/grpc
Grpc-Message: user requires a username
Grpc-Status: 3

{
  "code": "USER-001",
  "message": "user requires a username",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.BadRequest",
      "fieldViolations": [
        {
          "field": "username",
          "description": "Username is required"
        }
      ]
    }
  ]
}
```

### Documentation

- Document the purpose of the API, the services, the endpoints, the request and response messages, the error codes and the status codes.
- Describe the fields of the request and response messages, the purpose and if needed the constraints.  
- Document and explain the possible error codes and the error messages that can be returned by the API.
- Document if the endpoints requires specific permissions or roles.

#### Examples

```proto
// ListUsers will return all matching users. By default, we will return all users of your instance that you have permission to read. Make sure to include a limit and sorting for pagination.
//
// Required permission: user.read
//
// Error Codes:
//   - invalid_request: Your request does not have a valid format. Check error details for the reason.
//   - unauthenticated: You are not authenticated. Please provide a valid access token.
//   - permission_denied: You do not have the required permissions to access the requested resource.
rpc ListUsers(ListUsersRequest) returns (ListUsersResponse) {}
```

```proto

message SetPhoneRequest{
  // The user ID of the user to set the phone number for.
  string user_id = 1 [
    (validate.rules).string = {min_len: 1, max_len: 200},
    (google.api.field_behavior) = REQUIRED,
  ];
  // The phone number to set for the user.
  string phone = 2 [
    (validate.rules).string = {min_len: 1, max_len: 200},
    (google.api.field_behavior) = REQUIRED
  ];
  // If no verification is specified, an SMS will be automatically sent to the user.
  oneof verification {
    // Let ZITADEL send a phone verification code via SMS to the user.
    SendPhoneVerificationCode send_code = 3;
    // Return the phone verification code for the user so you can handle the delivery yourself.
    ReturnPhoneVerificationCode return_code = 4;
    // State the phone number as already verified.
    bool is_verified = 5 [(validate.rules).bool.const = true];
  }
}
```
# API Design

## The Basics
ZITADEL follows an API first approach. This means all features can not only be accessed via the UI but also via the API.
The API is designed using the Protobuf specification. The Protobuf specification is then used to generate the API client
and server code in different programming languages.
The API is designed to be used by different clients, such as web applications, mobile applications, and other services.
Therefore, the API is designed to be easy to use, consistent, and reliable.

Starting with the V2 API, the API and its services use a resource-oriented design. 
This means that the API is designed around resources, which are the key entities in the system.
Each resource has a unique identifier and a set of properties that describe the resource.
The entire lifecycle of a resource can be managed using the API.

> [!IMPORTANT]
> This style guide is a work in progress and will be updated over time.
> Not all parts of the API might follow the guidelines yet.
> However, all new endpoints and services must be designed according to this style guide.

### Protobuf, gRPC and connectRPC

The API is designed using the Protobuf specification. The Protobuf specification is used to define the API services, messages, and methods.
Starting with the V2 API, the API uses connectRPC as the main transport protocol. 
[connectRPC](https://connectrpc.com/) is a protocol that is based on gRPC and HTTP/2.
It allows clients to call the API using connectRPC, gRPC and also HTTP/1.1.

## Conventions

The API follows the base conventions of Protobuf and connectRPC.

Please check out their style guides and concepts for more information:
- Protobuf: https://protobuf.dev/programming-guides/style/
- gRPC: https://grpc.io/docs/what-is-grpc/core-concepts/
- Buf: https://buf.build/docs/best-practices/style-guide/

Additionally, there are some conventions that are specific to the ZITADEL API.
These conventions are described in the following sections.

### Versioning

The services and messages are versioned using major version numbers. This means that any change within a major version number is backward compatible.
Any breaking change requires a new major version number.
Each service is versioned independently. This means that a service can have a different version number than another service.
When creating a new service, start with version `2`, as version `1` is reserved for the old context based API and services.

Please check out the structure Buf style guide for more information about the folder and package structure: https://buf.build/docs/best-practices/style-guide/

### Explicitness

Make the handling of the API as explicit as possible. Do not make assumptions about the client's knowledge of the system or the API. 
Provide clear and concise documentation for the API.

Do not rely on implicit fallbacks or defaults if the client does not provide certain parameters.
Only use defaults if they are explicitly documented, such as returning a result set for the whole instance if no filter is provided.

### Naming Conventions

Names of resources, fields and methods should be descriptive and consistent.
Use domain-specific terminology and avoid abbreviations.
For example, use `OrganizationID` instead of **OrgID** or **resourceOwner** for the creation of a mew user or when returning one.

> [!NOTE]
> We'll update the resources in the [concepts section](https://zitadel.com/docs/concepts/structure/instance) to describe
> common resources and their meaning.
> Until then, please refer to the following issue: https://github.com/zitadel/zitadel/issues/5888

#### Resources and Fields

When a context is required for creating a resource, the context is added as a field to the resource.
For example, when creating a new user, the organization ID is required. The `organization_id` is added as a field to the `CreateUserRequest`.

```protobuf
message CreateUserRequest {
  ...
  string organization_id = 7 [
    (validate.rules).string = {min_len: 1, max_len: 200},
  ];
  ...
}
```

Only allow providing a context where it is required. Do not provide the possibility to provide a context where it is not required.
For example, when retrieving or updating a user, the organization ID is not required, since the user can be determined by the user ID.
However, it is possible to provide the organization ID as a filter to retrieve a list of users of a specific organization.

Prevent the creation of global messages that are used in multiple resources unless they always follow the same pattern.
Use dedicated fields as described above or create a separate message for the specific context, that is only used in the boundary of the same resource.  
For example, settings might be set as a default on the instance level, but might be overridden on the organization level.
In this case, the settings could share the same `SettingsContext` message to determine the context of the settings.
But do not create a global `Context` message that is used across the whole API if there are different scenarios and different fields required for the context.  
The same applies to messages that are returned by multiple resources.  
For example, information about the `User` might be different when managing the user resource itself than when it's returned
as part of an authorization or a manager role, where only limited information is needed.

Prevent reusing messages for the creation and the retrieval of a resource.
Returning messages might contain additional information that is not required or even not available for the creation of the resource.  
What might sound obvious when designing the CreateUserRequest for example, where only an `organization_id` but not the 
`organization_name` is available, might not be so obvious when designing some sub-resource like a user's `IdentityProviderLink`, 
which might contain an `identity_provider_name` when returned but not when created.

```protobuf
message CreateUserRequest {
  ...
  repreated AddIdentityProviderLink identity_provider_links = 8;
  ...
}

message AddIdentityProviderLink {
  string identity_provider_id = 1 [
    (validate.rules).string = {min_len: 1, max_len: 200},
  ];
  string user_id = 2 [
    (validate.rules).string = {min_len: 1, max_len: 200},
  ];
  string user_name = 3;
}

message IdentiyProviderLink {
  string identity_provider_id = 1;
  string identity_provider_name = 2;
  string user_id = 3;
  string user_name = 4;
} 
```

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

## Authentication and Authorization

The API uses OAuth 2 for authorization. There are corresponding middlewares that check the access token for validity and 
automatically return an error if the token is invalid.

Permissions grated to the user are organization specific and might only be checked based on the queried resource.
Therefore, the API does not check the permissions itself but relies on the checks of the functions that are called by the API.
Required permissions need to be documented in the [API documentation](#documentation).

## Pagination

The API uses pagination for listing resources. The client can specify a limit and an offset to retrieve a subset of the resources.
Additionally, the client can specify sorting options to sort the resources by a specific field.

Most listing methods should provide use the `ListQuery` message to allow the client to specify the limit, offset, and sorting options.
```protobuf

// ListQuery is a general query object for lists to allow pagination and sorting.
message ListQuery {
  uint64 offset = 1;
  // limit is the maximum amount of objects returned. The default is set to 100
  // with a maximum of 1000 in the runtime configuration.
  // If the limit exceeds the maximum configured ZITADEL will throw an error.
  // If no limit is present the default is taken.
  uint32 limit = 2;
  // Asc is the sorting order. If true the list is sorted ascending, if false
  // the list is sorted descending. The default is descending.
  bool asc = 3;
}
```
On the corresponding responses the `ListDetails` can be used to return the total count of the resources
and allow the user to handle their offset and limit accordingly.


## Error Handling

The API returns machine-readable errors in the response body. This includes a status code, an error code and possibly some details.

### Status Codes

The API uses status codes to indicate the status of a request. Depending on the protocol used to call the API,
the status code is returned as an HTTP status code or as a gRPC / connectRPC status code.
Check the possible status codes https://zitadel.com/docs/apis/statuscodes

### Error Codes

Additionally to the status code, the API returns unique error codes for each type of error.
The error codes are used to identify a specific error and can be used to handle the error programmatically.

> [!NOTE]
> Currently, ZITADEL might already return some error codes. However, they do not follow a specific pattern yet
> and are not documented. We will update the error codes and document them in the future.

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
- Document if the endpoints requires specific permissions or roles.
- Document and explain the possible error codes and the error messages that can be returned by the API.

#### Examples

```proto
// ListUsers will return all matching users. By default, we will return all users of your instance that you have permission to read. Make sure to include a limit and sorting for pagination.
//
// Required permission:
//   - user.read
//   - no permission required to own user
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
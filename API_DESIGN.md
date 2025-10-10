# API Design

This document describes the design principles and conventions for the ZITADEL API. It is scoped to the services and 
endpoints of the proprietary ZITADEL API and does not cover any standardized APIs like OAuth 2, OpenID Connect or SCIM.  

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
> However, all new endpoints and services MUST be designed according to this style guide.

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

### Deprecations

As a rule of thumb, redundant API methods are deprecated.

- The proto option `grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation.deprecated` MUST be set to true.
- One or more links to recommended replacement methods MUST be added to the deprecation message as a proto comment above the rpc spec.
- Guidance for switching to the recommended methods for common use cases SHOULD be added as a proto comment above the rpc spec.

#### Example

```protobuf
// Delete the user phone
//
// Deprecated: [Update the user's phone field](apis/resources/user_service_v2/user-service-update-user.api.mdx) to remove the phone number.
//
// Delete the phone number of a user.
rpc RemovePhone(RemovePhoneRequest) returns (RemovePhoneResponse) {
  option (google.api.http) = {
    delete: "/v2/users/{user_id}/phone"
    body: "*"
  };

  option (zitadel.protoc_gen_zitadel.v2.options) = {
    auth_option: {
      permission: "authenticated"
    }
  };

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
    deprecated: true;
    responses: {
      key: "200"
      value: {
        description: "OK";
      }
    };
    responses: {
      key: "404";
      value: {
        description: "User ID does not exist.";
      }
    }
  };
}
```

### Explicitness

Make the handling of the API as explicit as possible. Do not make assumptions about the client's knowledge of the system or the API. 
Provide clear and concise documentation for the API.

Do not rely on implicit fallbacks or defaults if the client does not provide certain parameters.
Only use defaults if they are explicitly documented, such as returning a result set for the whole instance if no filter is provided.

Some API calls may create multiple resources such as in the case of `zitadel.org.v2.AddOrganization`, where you can create an organization AND multiple users as admin.
In such cases the response should include **ALL** created resources and their ids. This removes any ambiguity from the users perspective whether or not
the additional resources were created and it also helps in testing.

### Naming Conventions

Names of resources, fields and methods MUST be descriptive and consistent.
Use domain-specific terminology and avoid abbreviations.
For example, use `organization_id` instead of **org_id** or **resource_owner** for the creation of a new user or when returning one.

> [!NOTE]
> We'll update the resources in the [concepts section](https://zitadel.com/docs/concepts/structure/instance) to describe
> common resources and their meaning.
> Until then, please refer to the following issue: https://github.com/zitadel/zitadel/issues/5888

#### Resources and Fields

##### Context information in Requests

When a context is required for creating a resource, the context is added as a field to the resource.
For example, when creating a new user, the organization's id is required. The `organization_id` is added as a field to the `CreateUserRequest`.

```protobuf
message CreateUserRequest {
  ...
  string organization_id = 7 [
    (validate.rules).string = {min_len: 1, max_len: 200},
  ];
  ...
}
```

Only allow providing a context where it is required. The context MUST not be provided if not required.
If the context is required but deferrable, the context can be defaulted.
For example, creating an Authorization without an organization id will default the organization id to the projects resource owner.
For example, when retrieving or updating a user, the `organization_id` is not required, since the user can be determined by the user's id.
However, it is possible to provide the `organization_id` as a filter to retrieve a list of users of a specific organization.

##### Context information in Responses

When the action of creation, update or deletion of a resource was successful, the returned response has to include the time of the operation and the generated identifiers.
This is achieved through the addition of a timestamp attribute with the operation as a prefix, and the generated information as separate attributes.

```protobuf
message SetExecutionResponse {
  // The timestamp of the execution set.
  google.protobuf.Timestamp set_date = 1 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      example: "\"2024-12-18T07:50:47.492Z\"";
    }
  ];
}

message CreateTargetResponse {
  // The unique identifier of the newly created target.
  string id = 1 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      example: "\"69629012906488334\"";
    }
  ];
  // The timestamp of the target creation.
  google.protobuf.Timestamp creation_date = 2 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      example: "\"2024-12-18T07:50:47.492Z\"";
    }
  ];
  // Key used to sign and check payload sent to the target.
  string signing_key = 3 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      example: "\"98KmsU67\""
    }
  ];
}

message UpdateProjectGrantResponse {
  // The timestamp of the change of the project grant.
  google.protobuf.Timestamp change_date = 1 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      example: "\"2025-01-23T10:34:18.051Z\"";
    }
  ];
}

message DeleteProjectGrantResponse {
  // The timestamp of the deletion of the project grant.
  // Note that the deletion date is only guaranteed to be set if the deletion was successful during the request.
  // In case the deletion occurred in a previous request, the deletion date might be empty.
  google.protobuf.Timestamp deletion_date = 1 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      example: "\"2025-01-23T10:34:18.051Z\"";
    }
  ];
}
```

##### Global messages

Prevent the creation of global messages that are used in multiple resources unless they always follow the same pattern.
Use dedicated fields as described above or create a separate message for the specific context, that is only used in the boundary of the same resource.  
For example, settings might be set as a default on the instance level, but might be overridden on the organization level.
In this case, the settings could share the same `SettingsContext` message to determine the context of the settings.
But do not create a global `Context` message that is used across the whole API if there are different scenarios and different fields required for the context.  
The same applies to messages that are returned by multiple resources.  
For example, information about the `User` might be different when managing the user resource itself than when it's returned
as part of an authorization or a manager role, where only limited information is needed.

On the other hand, types that always follow the same pattern and are used in multiple resources, such as `IDFilter`, `TimestampFilter` or `InIDsFilter` SHOULD be globalized and reused.

##### Re-using messages

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

Methods on a resource MUST be named using the following convention:

| Operation | Method Name        | Description                                                                                                                                                                                                                                 |
|-----------|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Create    | Create\<resource\> | Create a new resource. If the new resource conflicts with an existing resources uniqueness (id, loginname, ...) the creation MUST be prevented and an error returned.                                                                       |
| Update    | Update\<resource\> | Update an existing resource. In most cases this SHOULD allow partial updates. If there are exception, they MUST be explicitly documented on the endpoint. The resource MUST already exists. An error is returned otherwise.                 |
| Delete    | Delete\<resource\> | Delete an existing resource. If the resource does not exist, no error SHOULD be returned. In case of an exception to this rule, the behavior MUST clearly be documented.                                                                    |
| Set       | Set\<resource\>    | Set a resource. This will replace the existing resource with the new resource. In case where the creation and update of a resource do not need to be differentiated, a single `Set` method SHOULD be used. It SHOULD allow partial changes. |
| Get       | Get\<resource\>    | Retrieve a single resource by its unique identifier. If the resource does not exist, an error MUST be returned.                                                                                                                             |
| List      | List\<resource\>   | Retrieve a list of resources. The endpoint SHOULD provide options to filter, sort and paginate.                                                                                                                                             |

Methods on a list of resources MUST be named using the following convention:

| Operation | Method Name        | Description                                                                                                                                                                                      |
|-----------|--------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Add       | Add\<resource\>    | Add a new resource to a list. Any existing unique constraint (id, loginname, ...) will prevent the addition and return an error.                                                                 |
| Remove    | Remove\<resource\> | Remove an existing resource from a list. If the resource does not exist in the list, no error SHOULD be returned. In case of an exception to this rule, the behavior MUST clearly be documented. |
| Set       | Set\<resource\>    | Set a list of resources. This will replace the existing list with the new list.                                                                                                                  |

Additionally, state changes, specific actions or operations that do not fit into the CRUD operations SHOULD be named according to the action that is performed:
- `Activate` or `Deactivate` for enabling or disabling a resource.
- `Verify` for verifying a resource.
- `Send` for sending a resource.
- etc.

## Authentication and Authorization

The API uses OAuth 2 for authorization. There are corresponding middlewares that check the access token for validity and 
automatically return an error if the token is invalid.

Permissions granted to the user might be organization specific and can therefore only be checked based on the queried resource.
In such case, the API does not check the permissions itself but relies on the checks of the functions that are called by the API.
If the permission can be checked by the API itself, e.g. if the permission is instance wide, it can be annotated on the endpoint in the proto file (see below).
In any case, the required permissions need to be documented in the [API documentation](#documentation).

### Permission annotations

Permissions can be annotated on the endpoint in the proto file. This allows the API to automatically check the permissions for the user.
The permissions are checked by the middleware and an error is returned if the user does not have the required permissions.

The following example requires the user to have the `iam.web_key.write` permission to call the `CreateWebKey` method.
```protobuf
 option (zitadel.protoc_gen_zitadel.v2.options) = {
  auth_option: {
    permission: "iam.web_key.write"
  }
};
```

In case the permission cannot be checked by the API itself, but all requests need to be from an authenticated user, the `auth_option` can be set to `authenticated`.
```protobuf
 option (zitadel.protoc_gen_zitadel.v2.options) = {
  auth_option: {
    permission: "authenticated"
  }
};
```

## Listing resources

The API uses pagination for listing resources. The client can specify a limit and an offset to retrieve a subset of the resources.
Additionally, the client can specify sorting options to sort the resources by a specific field.

### Pagination

Most listing methods SHOULD use the `PaginationRequest` message to allow the client to specify the limit, offset, and sorting options.
```protobuf
message ListTargetsRequest {
  // List limitations and ordering.
  optional zitadel.filter.v2beta.PaginationRequest pagination = 1;
  // The field the result is sorted by. The default is the creation date. Beware that if you change this, your result pagination might be inconsistent.
  optional TargetFieldName sorting_column = 2 [
    (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
      default: "\"TARGET_FIELD_NAME_CREATION_DATE\""
    }
  ];
  // Define the criteria to query for.
  repeated TargetSearchFilter filters = 3;
  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    example: "{\"pagination\":{\"offset\":0,\"limit\":0,\"asc\":true},\"sortingColumn\":\"TARGET_FIELD_NAME_CREATION_DATE\",\"filters\":[{\"targetNameFilter\":{\"targetName\":\"ip_allow_list\",\"method\":\"TEXT_FILTER_METHOD_EQUALS\"}},{\"inTargetIdsFilter\":{\"targetIds\":[\"69629023906488334\",\"69622366012355662\"]}}]}";
  };
}
```

On the corresponding responses the `PaginationResponse` can be used to return the total count of the resources
and allow the user to handle their offset and limit accordingly.

The API MUST enforce a reasonable maximum limit for the number of resources that can be retrieved and returned in a single request.
The default limit is set to 100 and the maximum limit is set to 1000. If the client requests a limit that exceeds the maximum limit, an error is returned.

### Filter method

All filters in List operations SHOULD provide a method if not already specified by the filters name.
```protobuf
message TargetNameFilter {
  // Defines the name of the target to query for.
  string target_name = 1 [
    (validate.rules).string = {max_len: 200}
  ];
  // Defines which text comparison method used for the name query.
  zitadel.filter.v2beta.TextFilterMethod method = 2 [
    (validate.rules).enum.defined_only = true
  ];
}
```

## Error Handling

The API returns machine-readable errors in the response body. This includes a status code, an error code and possibly 
some details about the error. See the following sections for more information about the status codes, error codes and error messages.

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
  "code": "user_invalid_information",
  "message": "invalid or missing information provided for the creation of the user",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.BadRequest",
      "fieldViolations": [
        {
          "field": "given_name",
          "description": "given name is required",
          "reason": "MISSING_VALUE"
        },
        {
          "field": "family_name",
          "description": "family name must not exceed 200 characters",
          "reason": "INVALID_LENGTH"
        }
      ]
    }
  ]
}
```

gRPC / connectRPC example:
```
HTTP/2.0 200 OK
Content-Type: application/grpc
Grpc-Message: invalid information provided for the creation of the user
Grpc-Status: 3

{
  "code": "user_invalid_information",
  "message": "invalid or missing information provided for the creation of the user",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.BadRequest",
      "fieldViolations": [
        {
          "field": "given_name",
          "description": "given name is required",
          "reason": "MISSING_VALUE"
        },
        {
          "field": "family_name",
          "description": "family name must not exceed 200 characters",
          "reason": "INVALID_LENGTH"
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

```protobuf
// CreateUser will create a new user (human or machine) in the specified organization.
// The username must be unique.
//
// For human users:
// The user will receive a verification email if the email address is not marked as verified.
// You can pass a hashed_password. This allows migrating your users from your own system to ZITADEL, without any password
// reset for the user. Please check the required format and supported algorithms: <Link to documentation>. 
//
// Required permission:
//   - user.write
//
// Error Codes:
//   - user_missing_information: The request is missing required information (either given_name, family_name and/or email) or contains invalid data for the creation of the user. Check error details for the missing or invalid fields.
//   - user_already_exists: The user already exists. The username must be unique.
//   - invalid_request: Your request does not have a valid format. Check error details for the reason.
//   - permission_denied: You do not have the required permissions to access the requested resource.
//   - unauthenticated: You are not authenticated. Please provide a valid access token.
rpc CreatUser(CreatUserRequest) returns (CreatUserResponse) {}
```

```protobuf
// ListUsers will return all matching users. By default, we will return all users of your instance that you have permission to read. Make sure to include a limit and sorting for pagination.
//
// Required permission:
//   - user.read
//   - no permission required to own user
//
// Error Codes:
//   - invalid_request: Your request does not have a valid format. Check error details for the reason.
//   - permission_denied: You do not have the required permissions to access the requested resource.
//   - unauthenticated: You are not authenticated. Please provide a valid access token.
rpc ListUsers(ListUsersRequest) returns (ListUsersResponse) {}
```

```protobuf
// VerifyEmail will verify the provided verification code and mark the email as verified on success.
// An error is returned if the verification code is invalid or expired or if the user does not exist.
// Note that if multiple verification codes are generated, only the last one is valid.
//
// Required permission:
//   - no permission required, the user must be authenticated
//
// Error Codes:
//   - invalid_verification_code: The verification code is invalid or expired.
//   - invalid_request: Your request does not have a valid format. Check error details for the reason.
//   - unauthenticated: You are not authenticated. Please provide a valid access token.
rpc VerifyEmail (VerifyEmailRequest) returns (VerifyEmailResponse) {}
```

```protobuf
message VerifyEmailRequest{
  // The id of the user to verify the email for.
  string user_id = 1 [
    (validate.rules).string = {min_len: 1, max_len: 200}
  ];
  // The verification code generated and sent to the user.
  string verification_code = 2 [
    (validate.rules).string = {min_len: 1, max_len: 20}
  ];
}

```

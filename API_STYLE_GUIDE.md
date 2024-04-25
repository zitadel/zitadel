# API Style Guide

The ZITADEL API adheres to the following style guide.

## General

All APIs share the following behavior:
- APIs, that can be scoped with an optional RequestContext (system, instance, org), use the instance discovered by the requests Host header as default scope.

## All Resources

All resource APIs share the following behavior:
- Search request results can be scoped to a RequestContext.
- Search request results only contain results for which the requesting user has the necessary read permissions.
- Search requests are limited to 100 by default. The limit can be increased up to 1000.
- Resource configurations are partially updatable. With HTTP, this is done via PATCH requests. If no changes were made, the response is successful.
- Status changes or other actions on resources with side effects are done via POST requests. Their HTTP path ends with the underscore prefixed action name. For example `POST /resources/idps/{id}/_activate`.

## Reusable Resources

- Reusable resources can be created in a given context level (system, instance, org).
- For requests, that require a request ID, no request context is needed.
- Reusable resources are available in child contexts and by default have the same state (active or inactive) as in their immediate parent context.
- In child contexts, the state of a reused resource is *inherited* by default and can be changed to *active*, *inactive* or *inherit*.
- In child contexts, a reused resources configuration is read-only.
- Child contexts can always read the following properties of reused resources:
  - ID
  - name
  - description
  - state
  - sequence
  - last changed date
  - parent context
  - effective state in the immediate parent context.
- Managers of reusable resources in a parent context can restrict their readability in child contexts to the properties listed above.

## Settings

- Setting and retrieving settings is always context-aware. By default, the context is the instance discovered by the requests *Host* header.
- All settings properties can be partially overwritten in child-contexts.
- All settings properties can be partially reset in child-contexts, so their values default to the parent contexts property values.
- All settings properties returned by queries contain the value and if it is inherited, the context where it is inherited from.

## API Overview

### Resources APIs

The following resource APIs adhere to the behavior described in the [All Resources](#all-resources) section.
Some of them are reusable resources and also adhere to the behavior described in the [Reusable Resources](#reusable-resources) section.

- Instances
    - Instances
    - InstanceDomains
- Organizations
    - Organizations
    - OrganizationDomains
- Projects
    - Projects
    - Apps
    - Roles
- Sessions
- Authorizations
    - Memberships
    - ProjectGrants
    - UserGrants
- Users
    - Users
    - Schemas [^1]
- Actions
    - Executions [^1]
    - Targets [^1]
- IDPs [^1]
- Notifiers
    - SMTPProviders (only target types system and instance) [^1]
    - SMSProviders [^1]

[^1]: [Reusable Resources](#reusable-resources)

### Settings APIs

The following settings APIs adhere to the behavior described in the [Settings](#settings) section:

- DefaultLogin (default login settings only take effect for the built-in login UI. Users of the session API who build their own login UI can use them too, but they won't have any effect in the ZITADEL core)
    - Texts (key-value pairs for localized login texts, previously known as login texts)
    - Branding (predefined branding settings and custom key-value pairs, previously known as label policy or branding settings)
    - Login (previously known as login policy)
    - Lockout (previously known as lockout policy)
    - Password (previously known as password complexity policy)
    - Help (previously known as legal and support settings or privacy policy)
    - Domain (previously known as domain policy)
- Features (feature toggles)
- Languages (default language, restricted languages)
- Instances (instance-wide settings like disallow_public_org_registration. These settings are not overwritable for orgs)
- Templates (html and text templates for fully customizable emails and sms)

## gRPC-Gateway

- For accessing a resource by IDs, the IDs are mapped to path parameters.
- For searching resources, as many search parameters as possible are mapped to query parameters. Body: "*" should be avoided in the google.api.http proto option. Complex or sensitive queries are passed in the *POST* body.
- Request contexts are always mapped to query parameters.

## Documentation

- The API defines the ubiquitous language used in docs and other references.
- Each endpoint is documented in a way that is understandable to a developer who is not familiar with the API.
- The API docs contain references to concepts and other endpoints if it helps to understand the documented endpoint.
- The API docs contain at least one example with realistic data for each endpoint.
- Each API docs page links this style guide.

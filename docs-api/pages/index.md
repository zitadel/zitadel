---
title: Get started with Markdoc
description: How to get started with Markdoc
---

# ZITADEL API reference

{% callout %}
This will be our full-featured API reference for ZITADEL. Notice it's still work in progress but feel free to report issues and helpful feedback 👷 You can find the main issue [here](https://github.com/zitadel/zitadel/issues/4839).
{% /callout %}

{% section %}
{% column columns=2 %}

### Introduction

ZITADEL provides five APIs for different use cases. Four of these APIs are built with GRPC and generate a REST service. Each service's proto definition is located in the source control on GitHub. {% .text-base %}

As we generate the REST services and Swagger file out of the proto definition we recommend that you rely on the proto file. We annotate the corresponding REST methods on each possible call as well as the AuthN and AuthZ requirements. The last API (assets) is only a REST API because ZITADEL uses multipart form data for certain elements.

{% /column %}

{% column %}

Don't know how to get started? {% .ztdl-subheader .pb-0 %}

Read our Quickstart Guide [here](https://zitadel.com/docs/guides/start/quickstart)

You're not a developer? {% .ztdl-subheader .pb-0 %}

Consider reading our Guide [here](https://zitadel.com/docs/guides/introduction) or contact us at [support@zitadel.com](mailto:support@zitadel.com).

{% card title="Base URLs" %}
**Auth**: {% instanceDomain("your-domain", "auth", "v1") %}{% .pt-4 .pb-2 %}

**Management** {% instanceDomain("your-domain", "management", "v1") %}{% .py-2 %}

**Admin** {% instanceDomain("your-domain", "admin", "v1") %}{% .py-4 .pt-2 %}

{% /card %}
{% /column %}
{% /section %}

{% section %}
{% column %}

### Authentication

You can authorize your requests for ZITADEL API's by multiple methods.
These methods rely highly on the environment of your application.
You can either use an OIDC/OAuth2 Token or generate and use a Personal Access Token.

To successfully authenticate your request, send a valid `Authorization` header, using the `Bearer` scheme.

You can use the token directly after a user has authenticated in your app, or generate Peronsal Access Tokens in the ZITADEL Console.

Your Tokens carry many privileges, so be sure to keep them secure! Do not share your Tokens in publicly accessible areas such as GitHub, client-side code, and so forth.

All API requests must be made over HTTPS. Calls made over plain HTTP will fail. Most API requests without authentication will also fail.
{% /column %}

{% column %}
{% card title="User Info Endpoint" %}
This request gets the basic user information from the user

```bash
curl {% instanceDomain("your-domain", "auth", "v1") %} \
  -u 51IK2AACdCycJe9V8zmsnX1ByJhyRegEyFAwcgEA
# The colon prevents curl from asking for a password.
```

{% /card %}
{% /column %}

{% /section %}

### Organization Context

...blabla use `x-zitadel-orgid` header.

### Errors

### Metadata

### Pagination

## Core resources

### Users

### Organizations

### Policies

### Projects

### Applications

### Members

## Authentication Service

The authentication API (aka Auth API) is used for all operations on the currently logged in user. The user id is taken from the sub claim in the token.

### Introduction

## Management Service

### Introduction1

## Admin Service

### Introduction2

## Open Id Connect

## SAML

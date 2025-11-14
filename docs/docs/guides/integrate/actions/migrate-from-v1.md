---
title: Migrate from Actions v1 to v2
---

In this guide, you will have all necessary information to migrate from Actions v1 to Actions v2 with all currently [available Flow Types](/apis/actions/introduction#available-flow-types).

## Internal Authentication

### Post Authentication

A user has authenticated directly at ZITADEL.
ZITADEL validated the users inputs for password, one-time password, security key or passwordless factor.

To react to different authentication actions, the session service, `zitadel.session.v2.SessionService`, provides the different endpoints. As a rule of thumb, use response triggers if you primarily want to handle successful and failed authentications. On the other hand, use event triggers if you need more fine-granular handling, for example by the used authentication factors.  

Some use-cases:

- Handle successful authentication through the response of `/zitadel.session.v2.SessionService/CreateSession` and `/zitadel.session.v2.SessionService/SetSession`, [Action Response Example](./testing-response)
- Handle failed authentication through the response of `/zitadel.session.v2.SessionService/CreateSession` and `/zitadel.session.v2.SessionService/SetSession`, [Action Response Example](./testing-response)
- Handle session with password checked through the creation of event `session.password.checked`, [Action Event Example](./testing-event)
- Handle successful authentication through the creation of event `user.human.password.check.succeeded`, [Action Event Example](./testing-event)
- Handle failed authentication through the creation of event `user.human.password.check.failed`, [Action Event Example](./testing-event)

### Pre Creation

A user registers directly at ZITADEL.
ZITADEL did not create the user yet.

Some use-cases:

- Before a user is created through the request on `/zitadel.user.v2.UserService/AddHumanUser`, [Action Request Example](./testing-request)
- Add information to the user through the request on `/zitadel.user.v2.UserService/AddHumanUser`, [Action Request Manipulation Example](./testing-request-manipulation)

### Post Creation

A user registers directly at ZITADEL.  
ZITADEL successfully created the user.

Some use-cases:

- After user is created through the response on `/zitadel.user.v2.UserService/AddHumanUser`, [Action Response Example](./testing-response)
- At the event of a user creation on `user.human.added`, [Action Event Example](./testing-event)

## External Authentication

### Post Authentication

A user has authenticated externally. ZITADEL retrieved and mapped the external information.

Some use-cases:

- Handle the information mapping from the external authentication to internal structure through the response on `/zitadel.user.v2.UserService/RetrieveIdentityProviderIntent`, [Action Response Example](./testing-response)
  - information about the link to the external IDP available in the response under [`idpInformation`](/apis/resources/user_service_v2/user-service-retrieve-identity-provider-intent)
  - information if a new user has to be created available in the response under [`addHumanUser`](/apis/resources/user_service_v2/user-service-retrieve-identity-provider-intent), including metadata and link to external IDP

### Pre Creation

A user registers directly at ZITADEL.
ZITADEL did not create the user yet.

Some use-cases:

- Before a user is created through the request on `/zitadel.user.v2.UserService/AddHumanUser`, [Action Request Example](./testing-request)
- Add information to the user through the request on `/zitadel.user.v2.UserService/AddHumanUser`, [Action Request Manipulation Example](./testing-request-manipulation)

### Post Creation

A user registers directly at ZITADEL.  
ZITADEL successfully created the user.

Some use-cases:

- After user is created through the response on `/zitadel.user.v2.UserService/AddHumanUser`, [Action Response Example](./testing-response)
- At the event of a user creation on `user.human.added`, [Action Event Example](./testing-event)

## Complement Token

These are executed during the creation of tokens and token introspection.

### Pre Userinfo

These are called before userinfo are set in the id_token or userinfo and introspection endpoint response.

Some use-cases:

- Add claims to the userinfo through function on `preuserinfo`, [Action Function Example](./testing-function)
- Add metadata to user through function on `preuserinfo`, [Action Function Example](./testing-function)
- Add logs to the log claim through function on `preuserinfo`, [Action Function Example](./testing-function)

### Pre Access Token

These are called before the claims are set in the access token and the token type is `jwt`.

Some use-cases:

- Add claims to the userinfo through function on `preaccesstoken`, [Action Function Example](./testing-function)
- Add metadata to user through function on `preaccesstoken`, [Action Function Example](./testing-function)
- Add logs to the log claim through function on `preaccesstoken`, [Action Function Example](./testing-function)

## Customize SAML Response

These are executed before the return of the SAML Response.

### Pre SAMLResponse Creation

These are called before attributes are set in the SAMLResponse.

Some use-cases:

- Add custom attributes to the response through function on `presamlresponse`, [Action Function Example](./testing-function)
- Add metadata to user through function on `presamlresponse`, [Action Function Example](./testing-function)




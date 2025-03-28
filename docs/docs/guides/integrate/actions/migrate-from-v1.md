---
title: Migrate from Actions v1 to v2
---

In this guide, you will have all necessary information to migrate from Actions v1 to Actions v2 with all currently [available Flow Types](/apis/actions/introduction#available-flow-types).

## Internal Authentication

### Post Authentication

To replace any actions defined for Post Authentication, they only have to be defined to handle response calls in the session API.

### Pre Creation

To replace any actions defined for Pre Creation, they only have to be defined to handle request calls in the user API for AddHumanUser.

### Post Creation

To replace any actions defined for Pre Creation, they only have to be defined to handle response calls in the user API for AddHumanUser.

## External Authentication

### Post Authentication

To replace any actions defined for Post Authentication, they only have to be defined to handle response calls in the user API for RetrieveIdentityProviderIntent.

### Pre Creation

To replace any actions defined for Pre Creation, they only have to be defined to handle request calls in the user API for AddHumanUser.

### Post Creation

To replace any actions defined for Pre Creation, they only have to be defined to handle response calls in the user API for AddHumanUser.

## Complement Token

### Pre Userinfo

To replace any actions defined for Pre Creation, they only have to be defined to handle function calls for `preuserinfo`.

### Pre Access Token

To replace any actions defined for Pre Creation, they only have to be defined to handle function calls for `preaccesstoken`.

## Customize SAML Response

### Pre SAMLResponse Creation

To replace any actions defined for Pre Creation, they only have to be defined to handle function calls for `presamlresponse`.



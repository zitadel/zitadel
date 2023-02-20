---
title: zitadel/idp.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### IDP



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  IDPState | - |  |
| name |  string | - |  |
| styling_type |  IDPStylingType | - |  |
| owner |  IDPOwnerType | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.oidc_config |  OIDCConfig | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.jwt_config |  JWTConfig | - |  |
| auto_register |  bool | - |  |




### IDPIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.max_len: 200<br />  |




### IDPLoginPolicyLink



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - |  |
| idp_name |  string | - |  |
| idp_type |  IDPType | - |  |




### IDPNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### IDPOwnerTypeQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| owner_type |  IDPOwnerType | - | enum.defined_only: true<br />  |




### IDPUserLink



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| idp_id |  string | - |  |
| idp_name |  string | - |  |
| provided_user_id |  string | - |  |
| provided_user_name |  string | - |  |
| idp_type |  IDPType | - |  |




### JWTConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| jwt_endpoint |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| keys_endpoint |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| header_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### LDAPAttributes



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id_attribute |  string | - | string.max_len: 200<br />  |
| first_name_attribute |  string | - | string.max_len: 200<br />  |
| last_name_attribute |  string | - | string.max_len: 200<br />  |
| display_name_attribute |  string | - | string.max_len: 200<br />  |
| nick_name_attribute |  string | - | string.max_len: 200<br />  |
| preferred_username_attribute |  string | - | string.max_len: 200<br />  |
| email_attribute |  string | - | string.max_len: 200<br />  |
| email_verified_attribute |  string | - | string.max_len: 200<br />  |
| phone_attribute |  string | - | string.max_len: 200<br />  |
| phone_verified_attribute |  string | - | string.max_len: 200<br />  |
| preferred_language_attribute |  string | - | string.max_len: 200<br />  |
| avatar_url_attribute |  string | - | string.max_len: 200<br />  |
| profile_attribute |  string | - | string.max_len: 200<br />  |




### LDAPConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| host |  string | - |  |
| port |  string | - |  |
| tls |  bool | - |  |
| base_dn |  string | - |  |
| user_object_class |  string | - |  |
| user_unique_attribute |  string | - |  |
| admin |  string | - |  |
| attributes |  LDAPAttributes | - |  |
| provider_options |  Options | - |  |




### OIDCConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| client_id |  string | - |  |
| issuer |  string | - |  |
| scopes | repeated string | - |  |
| display_name_mapping |  OIDCMappingField | - |  |
| username_mapping |  OIDCMappingField | - |  |




### Options



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| is_linking_allowed |  bool | - |  |
| is_creation_allowed |  bool | - |  |
| is_auto_creation |  bool | - |  |
| is_auto_update |  bool | - |  |




### Provider



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  IDPState | - |  |
| name |  string | - |  |
| owner |  IDPOwnerType | - |  |
| type |  ProviderType | - |  |
| config |  ProviderConfig | - |  |




### ProviderConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| options |  Options | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.ldap |  LDAPConfig | - |  |






## Enums


### IDPFieldName {#idpfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_FIELD_NAME_UNSPECIFIED | 0 | - |
| IDP_FIELD_NAME_NAME | 1 | - |




### IDPOwnerType {#idpownertype}
the owner of the identity provider.

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_OWNER_TYPE_UNSPECIFIED | 0 | - |
| IDP_OWNER_TYPE_SYSTEM | 1 | system is managed by the ZITADEL administrators |
| IDP_OWNER_TYPE_ORG | 2 | org is managed by de organisation administrators |




### IDPState {#idpstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_STATE_UNSPECIFIED | 0 | - |
| IDP_STATE_ACTIVE | 1 | - |
| IDP_STATE_INACTIVE | 2 | - |




### IDPStylingType {#idpstylingtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| STYLING_TYPE_UNSPECIFIED | 0 | - |
| STYLING_TYPE_GOOGLE | 1 | - |




### IDPType {#idptype}
authorization framework of the identity provider

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_TYPE_UNSPECIFIED | 0 | - |
| IDP_TYPE_OIDC | 1 | - |
| IDP_TYPE_JWT | 3 | PLANNED: IDP_TYPE_SAML |




### OIDCMappingField {#oidcmappingfield}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_MAPPING_FIELD_UNSPECIFIED | 0 | - |
| OIDC_MAPPING_FIELD_PREFERRED_USERNAME | 1 | - |
| OIDC_MAPPING_FIELD_EMAIL | 2 | - |




### ProviderType {#providertype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PROVIDER_TYPE_UNSPECIFIED | 0 | - |
| PROVIDER_TYPE_OIDC | 1 | - |
| PROVIDER_TYPE_JWT | 2 | - |
| PROVIDER_TYPE_LDAP | 3 | - |
| PROVIDER_TYPE_OAUTH | 4 | - |
| PROVIDER_TYPE_AZURE_AD | 5 | - |
| PROVIDER_TYPE_GITHUB | 6 | - |
| PROVIDER_TYPE_GITHUB_EE | 7 | - |
| PROVIDER_TYPE_GITLAB | 8 | - |
| PROVIDER_TYPE_GITLAB_SELF_HOSTED | 9 | - |
| PROVIDER_TYPE_GOOGLE | 10 | - |





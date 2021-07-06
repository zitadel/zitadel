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




### OIDCConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| client_id |  string | - |  |
| issuer |  string | - |  |
| scopes | repeated string | - |  |
| display_name_mapping |  OIDCMappingField | - |  |
| username_mapping |  OIDCMappingField | - |  |
| authorization_endpoint |  string | - | string.max_len: 500<br />  |
| token_endpoint |  string | - | string.max_len: 500<br />  |






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
| IDP_TYPE_OIDC | 1 | PLANNED: IDP_TYPE_SAML |




### OIDCMappingField {#oidcmappingfield}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_MAPPING_FIELD_UNSPECIFIED | 0 | - |
| OIDC_MAPPING_FIELD_PREFERRED_USERNAME | 1 | - |
| OIDC_MAPPING_FIELD_EMAIL | 2 | - |





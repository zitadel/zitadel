---
title: zitadel/app.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### APIConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| client_id |  string | - |  |
| client_secret |  string | - |  |
| auth_method_type |  APIAuthMethodType | - |  |




### App



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  AppState | - |  |
| name |  string | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.oidc_config |  OIDCConfig | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.api_config |  APIConfig | - |  |




### AppNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### AppQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.name_query |  AppNameQuery | - |  |




### OIDCConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| redirect_uris | repeated string | - |  |
| response_types | repeated OIDCResponseType | - |  |
| grant_types | repeated OIDCGrantType | - |  |
| app_type |  OIDCAppType | - |  |
| client_id |  string | - |  |
| client_secret |  string | - |  |
| auth_method_type |  OIDCAuthMethodType | - |  |
| post_logout_redirect_uris | repeated string | - |  |
| version |  OIDCVersion | - |  |
| none_compliant |  bool | - |  |
| compliance_problems | repeated zitadel.v1.LocalizedMessage | - |  |
| dev_mode |  bool | - |  |
| access_token_type |  OIDCTokenType | - |  |
| access_token_role_assertion |  bool | - |  |
| id_token_role_assertion |  bool | - |  |
| id_token_userinfo_assertion |  bool | - |  |
| clock_skew |  google.protobuf.Duration | - |  |
| additional_origins | repeated string | - |  |
| allowed_origins | repeated string | - |  |






## Enums


### APIAuthMethodType {#apiauthmethodtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| API_AUTH_METHOD_TYPE_BASIC | 0 | - |
| API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT | 1 | - |




### AppState {#appstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| APP_STATE_UNSPECIFIED | 0 | - |
| APP_STATE_ACTIVE | 1 | - |
| APP_STATE_INACTIVE | 2 | - |




### OIDCAppType {#oidcapptype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_APP_TYPE_WEB | 0 | - |
| OIDC_APP_TYPE_USER_AGENT | 1 | - |
| OIDC_APP_TYPE_NATIVE | 2 | - |




### OIDCAuthMethodType {#oidcauthmethodtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_AUTH_METHOD_TYPE_BASIC | 0 | - |
| OIDC_AUTH_METHOD_TYPE_POST | 1 | - |
| OIDC_AUTH_METHOD_TYPE_NONE | 2 | - |
| OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT | 3 | - |




### OIDCGrantType {#oidcgranttype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_GRANT_TYPE_AUTHORIZATION_CODE | 0 | - |
| OIDC_GRANT_TYPE_IMPLICIT | 1 | - |
| OIDC_GRANT_TYPE_REFRESH_TOKEN | 2 | - |




### OIDCResponseType {#oidcresponsetype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_RESPONSE_TYPE_CODE | 0 | - |
| OIDC_RESPONSE_TYPE_ID_TOKEN | 1 | - |
| OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN | 2 | - |




### OIDCTokenType {#oidctokentype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_TOKEN_TYPE_BEARER | 0 | - |
| OIDC_TOKEN_TYPE_JWT | 1 | - |




### OIDCVersion {#oidcversion}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_VERSION_1_0 | 0 | - |





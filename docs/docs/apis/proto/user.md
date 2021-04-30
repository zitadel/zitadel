---
title: zitadel/user.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### AuthFactor



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  AuthFactorState | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.otp |  AuthFactorOTP | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.u2f |  AuthFactorU2F | - |  |




### AuthFactorOTP





### AuthFactorU2F



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| name |  string | - |  |




### DisplayNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| display_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | - |  |
| is_email_verified |  bool | - |  |




### EmailQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email_address |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### FirstNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### Human



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| profile |  Profile | - |  |
| email |  Email | - |  |
| phone |  Phone | - |  |




### LastNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| last_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### Machine



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - |  |
| description |  string | - |  |




### Membership



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| roles | repeated string | - |  |
| display_name |  string | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.iam |  bool | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.org_id |  string | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.project_id |  string | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.project_grant_id |  string | - |  |




### MembershipIAMQuery
this query is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| iam |  bool | - |  |




### MembershipOrgQuery
this query is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.max_len: 200<br />  |




### MembershipProjectGrantQuery
this query is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_grant_id |  string | - | string.max_len: 200<br />  |




### MembershipProjectQuery
this query is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.max_len: 200<br />  |




### MembershipQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.org_query |  MembershipOrgQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_query |  MembershipProjectQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_grant_query |  MembershipProjectGrantQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.iam_query |  MembershipIAMQuery | - |  |




### NickNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| nick_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### Phone



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| phone |  string | - |  |
| is_phone_verified |  bool | - |  |




### Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - |  |
| last_name |  string | - |  |
| nick_name |  string | - |  |
| display_name |  string | - |  |
| preferred_language |  string | - |  |
| gender |  Gender | - |  |




### SearchQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_name_query |  UserNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.first_name_query |  FirstNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.last_name_query |  LastNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.nick_name_query |  NickNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.display_name_query |  DisplayNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.email_query |  EmailQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.state_query |  StateQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.type_query |  TypeQuery | - |  |




### Session



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| session_id |  string | - |  |
| agent_id |  string | - |  |
| auth_state |  SessionState | - |  |
| user_id |  string | - |  |
| user_name |  string | - |  |
| login_name |  string | - |  |
| display_name |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### StateQuery
UserStateQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  UserState | - | enum.defined_only: true<br />  |




### TypeQuery
UserTypeQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  Type | - | enum.defined_only: true<br />  |




### User



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  UserState | - |  |
| user_name |  string | - |  |
| login_names | repeated string | - |  |
| preferred_login_name |  string | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.human |  Human | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.machine |  Machine | - |  |




### UserGrant



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| role_keys | repeated string | - |  |
| state |  UserGrantState | - |  |
| user_id |  string | - |  |
| user_name |  string | - |  |
| first_name |  string | - |  |
| last_name |  string | - |  |
| email |  string | - | string.email: true<br />  |
| display_name |  string | - | string.max_len: 200<br />  |
| org_id |  string | - |  |
| org_name |  string | - |  |
| org_domain |  string | - |  |
| project_id |  string | - |  |
| project_name |  string | - |  |
| project_grant_id |  string | - |  |




### UserGrantDisplayNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| display_name |  string | - |  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantEmailQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantFirstNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantLastNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| last_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantOrgDomainQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_domain |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantOrgNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantProjectGrantIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_grant_id |  string | - | string.max_len: 200<br />  |




### UserGrantProjectIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.max_len: 200<br />  |




### UserGrantProjectNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_id_query |  UserGrantProjectIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_id_query |  UserGrantUserIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.with_granted_query |  UserGrantWithGrantedQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.role_key_query |  UserGrantRoleKeyQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_grant_id_query |  UserGrantProjectGrantIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_name_query |  UserGrantUserNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.first_name_query |  UserGrantFirstNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.last_name_query |  UserGrantLastNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.email_query |  UserGrantEmailQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.org_name_query |  UserGrantOrgNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.org_domain_query |  UserGrantOrgDomainQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_name_query |  UserGrantProjectNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.display_name_query |  UserGrantDisplayNameQuery | - |  |




### UserGrantRoleKeyQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| role_key |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantUserIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.max_len: 200<br />  |




### UserGrantUserNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### UserGrantWithGrantedQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| with_granted |  bool | - |  |




### UserNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### WebAuthNKey



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| public_key |  bytes | - |  |




### WebAuthNToken



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| state |  AuthFactorState | - |  |
| name |  string | - |  |




### WebAuthNVerification



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| public_key_credential |  bytes | - | bytes.min_len: 55<br />  |
| token_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |






## Enums


### AuthFactorState {#authfactorstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| AUTH_FACTOR_STATE_UNSPECIFIED | 0 | - |
| AUTH_FACTOR_STATE_NOT_READY | 1 | - |
| AUTH_FACTOR_STATE_READY | 2 | - |
| AUTH_FACTOR_STATE_REMOVED | 3 | - |




### Gender {#gender}


| Name | Number | Description |
| ---- | ------ | ----------- |
| GENDER_UNSPECIFIED | 0 | - |
| GENDER_FEMALE | 1 | - |
| GENDER_MALE | 2 | - |
| GENDER_DIVERSE | 3 | - |




### SessionState {#sessionstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SESSION_STATE_UNSPECIFIED | 0 | - |
| SESSION_STATE_ACTIVE | 1 | - |
| SESSION_STATE_TERMINATED | 2 | - |




### Type {#type}


| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 | - |
| TYPE_HUMAN | 1 | - |
| TYPE_MACHINE | 2 | - |




### UserFieldName {#userfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_FIELD_NAME_UNSPECIFIED | 0 | - |
| USER_FIELD_NAME_USER_NAME | 1 | - |
| USER_FIELD_NAME_FIRST_NAME | 2 | - |
| USER_FIELD_NAME_LAST_NAME | 3 | - |
| USER_FIELD_NAME_NICK_NAME | 4 | - |
| USER_FIELD_NAME_DISPLAY_NAME | 5 | - |
| USER_FIELD_NAME_EMAIL | 6 | - |
| USER_FIELD_NAME_STATE | 7 | - |
| USER_FIELD_NAME_TYPE | 8 | - |




### UserGrantState {#usergrantstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_GRANT_STATE_UNSPECIFIED | 0 | - |
| USER_GRANT_STATE_ACTIVE | 1 | - |
| USER_GRANT_STATE_INACTIVE | 2 | - |




### UserState {#userstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_STATE_UNSPECIFIED | 0 | - |
| USER_STATE_ACTIVE | 1 | - |
| USER_STATE_INACTIVE | 2 | - |
| USER_STATE_DELETED | 3 | - |
| USER_STATE_LOCKED | 4 | - |
| USER_STATE_SUSPEND | 5 | - |
| USER_STATE_INITIAL | 6 | - |





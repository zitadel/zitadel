---
title: zitadel/project.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### AllProjectGrantQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_name_query |  GrantProjectNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.role_key_query |  GrantRoleKeyQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_id_query |  ProjectIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.granted_org_id_query |  GrantedOrgIDQuery | - |  |




### GrantProjectNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### GrantRoleKeyQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| role_key |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### GrantedOrgIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| granted_org_id |  string | - | string.max_len: 200<br />  |




### GrantedProject



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| grant_id |  string | - |  |
| granted_org_id |  string | - |  |
| granted_org_name |  string | - |  |
| granted_role_keys | repeated string | - |  |
| state |  ProjectGrantState | - |  |
| project_id |  string | - |  |
| project_name |  string | - |  |
| project_owner_id |  string | - |  |
| project_owner_name |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### Project



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| name |  string | - |  |
| state |  ProjectState | - |  |
| project_role_assertion |  bool | describes if roles of user should be added in token |  |
| project_role_check |  bool | ZITADEL checks if the user has at least one on this project |  |
| has_project_check |  bool | ZITADEL checks if the org of the user has permission to this project |  |
| private_labeling_setting |  PrivateLabelingSetting | Defines from where the private labeling should be triggered |  |




### ProjectGrantQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_name_query |  GrantProjectNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.role_key_query |  GrantRoleKeyQuery | - |  |




### ProjectIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.max_len: 200<br />  |




### ProjectNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### ProjectQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.name_query |  ProjectNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_resource_owner_query |  ProjectResourceOwnerQuery | - |  |




### ProjectResourceOwnerQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| resource_owner |  string | - | string.max_len: 200<br />  |




### Role



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| display_name |  string | - |  |
| group |  string | - |  |




### RoleDisplayNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| display_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### RoleKeyQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### RoleQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.key_query |  RoleKeyQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.display_name_query |  RoleDisplayNameQuery | - |  |






## Enums


### PrivateLabelingSetting {#privatelabelingsetting}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PRIVATE_LABELING_SETTING_UNSPECIFIED | 0 | - |
| PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY | 1 | - |
| PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY | 2 | - |




### ProjectGrantState {#projectgrantstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PROJECT_GRANT_STATE_UNSPECIFIED | 0 | - |
| PROJECT_GRANT_STATE_ACTIVE | 1 | - |
| PROJECT_GRANT_STATE_INACTIVE | 2 | - |




### ProjectState {#projectstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PROJECT_STATE_UNSPECIFIED | 0 | - |
| PROJECT_STATE_ACTIVE | 1 | - |
| PROJECT_STATE_INACTIVE | 2 | - |





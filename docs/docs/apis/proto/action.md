---
title: zitadel/action.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Action



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  ActionState | - |  |
| name |  string | - |  |
| script |  string | - |  |
| timeout |  google.protobuf.Duration | - |  |
| allowed_to_fail |  bool | - |  |




### ActionIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.max_len: 200<br />  |




### ActionNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### ActionStateQuery
ActionStateQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  ActionState | - | enum.defined_only: true<br />  |




### Flow



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  FlowType | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  FlowState | - |  |
| triggers | map Flow.TriggersEntry | - |  |




### Flow.TriggersEntry



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  int32 | - |  |
| value |  TriggerAction | - |  |




### FlowStateQuery
FlowStateQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  FlowState | - | enum.defined_only: true<br />  |




### FlowTypeQuery
FlowTypeQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  FlowType | - | enum.defined_only: true<br />  |




### TriggerAction



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| actions | map TriggerAction.ActionsEntry | - |  |




### TriggerAction.ActionsEntry



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  int32 | - |  |
| value |  Action | - |  |






## Enums


### ActionFieldName {#actionfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACTION_FIELD_NAME_UNSPECIFIED | 0 | - |
| ACTION_FIELD_NAME_NAME | 1 | - |
| ACTION_FIELD_NAME_ID | 2 | - |
| ACTION_FIELD_NAME_STATE | 3 | - |




### ActionState {#actionstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACTION_STATE_UNSPECIFIED | 0 | - |
| ACTION_STATE_INACTIVE | 1 | - |
| ACTION_STATE_ACTIVE | 2 | - |




### FlowFieldName {#flowfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| FLOW_FIELD_NAME_UNSPECIFIED | 0 | - |
| FLOW_FIELD_NAME_TYPE | 1 | - |
| FLOW_FIELD_NAME_STATE | 2 | - |




### FlowState {#flowstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| FLOW_STATE_UNSPECIFIED | 0 | - |
| FLOW_STATE_INACTIVE | 1 | - |
| FLOW_STATE_ACTIVE | 2 | - |




### FlowType {#flowtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| FLOW_TYPE_UNSPECIFIED | 0 | - |
| FLOW_TYPE_REGISTER | 1 | - |
| FLOW_TYPE_LOGIN | 2 | - |




### TriggerType {#triggertype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| TRIGGER_TYPE_UNSPECIFIED | 0 | - |
| TRIGGER_TYPE_PRE_REGISTER | 1 | - |
| TRIGGER_TYPE_POST_REGISTER | 2 | - |
| TRIGGER_TYPE_POST_LOGIN | 3 | - |





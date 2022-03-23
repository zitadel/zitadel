---
title: zitadel/instance.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### DomainsQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domains | repeated string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.ListQueryMethod | - | enum.defined_only: true<br />  |




### IdQuery
IdQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.max_len: 200<br />  |




### Instance



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  State | - |  |
| generated_domain |  string | - |  |
| custom_domains | repeated string | - |  |
| name |  string | - |  |
| version |  string | - |  |




### Query



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.id_query |  IdQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.domains_query |  DomainsQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.state_query |  StateQuery | - |  |




### StateQuery
StateQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  State | - | enum.defined_only: true<br />  |






## Enums


### FieldName {#fieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| FIELD_NAME_UNSPECIFIED | 0 | - |
| FIELD_NAME_ID | 1 | - |
| FIELD_NAME_GENERATED_DOMAIN | 2 | - |
| FIELD_NAME_NAME | 3 | - |
| FIELD_NAME_CREATION_DATE | 4 | - |




### State {#state}


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNSPECIFIED | 0 | - |
| STATE_CREATING | 1 | - |
| STATE_RUNNING | 2 | - |
| STATE_STOPPING | 3 | - |
| STATE_STOPPED | 4 | - |





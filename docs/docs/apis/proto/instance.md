---
title: zitadel/instance.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Domain



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| domain |  string | - |  |
| primary |  bool | - |  |
| generated |  bool | - |  |




### DomainGeneratedQuery
DomainGeneratedQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| generated |  bool | - |  |




### DomainPrimaryQuery
DomainPrimaryQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| primary |  bool | - |  |




### DomainQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### DomainSearchQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.domain_query |  DomainQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.generated_query |  DomainGeneratedQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.primary_query |  DomainPrimaryQuery | - |  |




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
| name |  string | - |  |




### Query



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.id_query |  IdQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.state_query |  StateQuery | - |  |




### StateQuery
StateQuery is always equals


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| state |  State | - | enum.defined_only: true<br />  |






## Enums


### DomainFieldName {#domainfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOMAIN_FIELD_NAME_UNSPECIFIED | 0 | - |
| DOMAIN_FIELD_NAME_DOMAIN | 1 | - |
| DOMAIN_FIELD_NAME_PRIMARY | 2 | - |
| DOMAIN_FIELD_NAME_GENERATED | 3 | - |
| DOMAIN_FIELD_NAME_CREATION_DATE | 4 | - |




### FieldName {#fieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| FIELD_NAME_UNSPECIFIED | 0 | - |
| FIELD_NAME_ID | 1 | - |
| FIELD_NAME_NAME | 2 | - |
| FIELD_NAME_CREATION_DATE | 3 | - |




### State {#state}


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNSPECIFIED | 0 | - |
| STATE_CREATING | 1 | - |
| STATE_RUNNING | 2 | - |
| STATE_STOPPING | 3 | - |
| STATE_STOPPED | 4 | - |





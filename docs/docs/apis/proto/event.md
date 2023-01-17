---
title: zitadel/event.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Aggregate



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| type |  zitadel.v1.LocalizedMessage | - |  |
| resource_owner |  string | - |  |




### Editor



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| display_name |  string | - |  |
| service |  string | - |  |




### Event



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| editor |  Editor | - |  |
| aggregate |  Aggregate | - |  |
| sequence |  uint64 | - |  |
| creation_date |  google.protobuf.Timestamp | The timestamp the event occurred |  |
| payload |  google.protobuf.Struct | - |  |
| type |  zitadel.v1.LocalizedMessage | - |  |







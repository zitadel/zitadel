---
title: zitadel/object.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### ListDetails



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| total_result |  uint64 | - |  |
| processed_sequence |  uint64 | - |  |
| view_timestamp |  google.protobuf.Timestamp | - |  |




### ListQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| offset |  uint64 | - |  |
| limit |  uint32 | - |  |
| asc |  bool | - |  |




### ObjectDetails



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sequence |  uint64 | sequence represents the order of events. It's always upcounting

on read: the sequence of the last event reduced by the projection

on manipulation: the timestamp of the event(s) added by the manipulation |  |
| creation_date |  google.protobuf.Timestamp | creation_date is the timestamp where the first operation on the object was made

on read: the timestamp of the first event of the object

on create: the timestamp of the event(s) added by the manipulation |  |
| change_date |  google.protobuf.Timestamp | change_date is the timestamp when the object was changed

on read: the timestamp of the last event reduced by the projection

on manipulation: the |  |
| resource_owner |  string | resource_owner is the organisation an object belongs to |  |






## Enums


### TextQueryMethod {#textquerymethod}


| Name | Number | Description |
| ---- | ------ | ----------- |
| TEXT_QUERY_METHOD_EQUALS | 0 | - |
| TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE | 1 | - |
| TEXT_QUERY_METHOD_STARTS_WITH | 2 | - |
| TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE | 3 | - |
| TEXT_QUERY_METHOD_CONTAINS | 4 | - |
| TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE | 5 | - |
| TEXT_QUERY_METHOD_ENDS_WITH | 6 | - |
| TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE | 7 | - |





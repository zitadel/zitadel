---
title: zitadel/change.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Change



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| change_date |  google.protobuf.Timestamp | - |  |
| event_type |  zitadel.v1.LocalizedMessage | - |  |
| sequence |  uint64 | - |  |
| editor_id |  string | - |  |
| editor_display_name |  string | - |  |
| resource_owner_id |  string | - |  |
| editor_preferred_login_name |  string | - |  |
| editor_avatar_url |  string | - |  |




### ChangeQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sequence |  uint64 | sequence represents the order of events. It's always upcounting |  |
| limit |  uint32 | - |  |
| asc |  bool | - |  |







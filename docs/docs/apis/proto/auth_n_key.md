---
title: zitadel/auth_n_key.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Key



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| type |  KeyType | - |  |
| expiration_date |  google.protobuf.Timestamp | - |  |






## Enums


### KeyType {#keytype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| KEY_TYPE_UNSPECIFIED | 0 | - |
| KEY_TYPE_JSON | 1 | - |





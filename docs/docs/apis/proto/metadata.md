---
title: zitadel/metadata.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### MetaData



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| key |  string | - |  |
| value |  string | - |  |




### MetaDataKeyQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### MetaDataQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.key_query |  MetaDataKeyQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.value_query |  MetaDataValueQuery | - |  |




### MetaDataValueQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| value |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |







---
title: zitadel/metadata.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Metadata



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| key |  string | - |  |
| value |  bytes | - |  |




### MetadataKeyQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### MetadataQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.key_query |  MetadataKeyQuery | - |  |







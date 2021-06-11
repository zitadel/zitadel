---
title: zitadel/member.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### EmailQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### FirstNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### LastNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| last_name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### Member



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| roles | repeated string | - |  |
| preferred_login_name |  string | - |  |
| email |  string | - |  |
| first_name |  string | - |  |
| last_name |  string | - |  |
| display_name |  string | - |  |
| avatar_url |  string | - |  |




### SearchQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.first_name_query |  FirstNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.last_name_query |  LastNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.email_query |  EmailQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_id_query |  UserIDQuery | - |  |




### UserIDQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.max_len: 200<br />  |







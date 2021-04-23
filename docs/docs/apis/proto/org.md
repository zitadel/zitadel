---
title: zitadel/org.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### Domain



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| domain_name |  string | - |  |
| is_verified |  bool | - |  |
| is_primary |  bool | - |  |
| validation_type |  DomainValidationType | - |  |




### DomainNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### DomainSearchQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.domain_name_query |  DomainNameQuery | - |  |




### Org



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| state |  OrgState | - |  |
| name |  string | - |  |
| primary_domain |  string | - |  |




### OrgDomainQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### OrgNameQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| method |  zitadel.v1.TextQueryMethod | - | enum.defined_only: true<br />  |




### OrgQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.name_query |  OrgNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.domain_query |  OrgDomainQuery | - |  |






## Enums


### DomainValidationType {#domainvalidationtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOMAIN_VALIDATION_TYPE_UNSPECIFIED | 0 | - |
| DOMAIN_VALIDATION_TYPE_HTTP | 1 | - |
| DOMAIN_VALIDATION_TYPE_DNS | 2 | - |




### OrgFieldName {#orgfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| ORG_FIELD_NAME_UNSPECIFIED | 0 | - |
| ORG_FIELD_NAME_NAME | 1 | - |




### OrgState {#orgstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| ORG_STATE_UNSPECIFIED | 0 | - |
| ORG_STATE_ACTIVE | 1 | - |
| ORG_STATE_INACTIVE | 2 | - |





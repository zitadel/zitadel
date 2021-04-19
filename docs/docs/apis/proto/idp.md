---
title: zitadel/idp.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)






## Enums


### IDPFieldName {#idpfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_FIELD_NAME_UNSPECIFIED | 0 | - |
| IDP_FIELD_NAME_NAME | 1 | - |




### IDPOwnerType {#idpownertype}
the owner of the identity provider.

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_OWNER_TYPE_UNSPECIFIED | 0 | - |
| IDP_OWNER_TYPE_SYSTEM | 1 | system is managed by the ZITADEL administrators |
| IDP_OWNER_TYPE_ORG | 2 | org is managed by de organisation administrators |




### IDPState {#idpstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_STATE_UNSPECIFIED | 0 | - |
| IDP_STATE_ACTIVE | 1 | - |
| IDP_STATE_INACTIVE | 2 | - |




### IDPStylingType {#idpstylingtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| STYLING_TYPE_UNSPECIFIED | 0 | - |
| STYLING_TYPE_GOOGLE | 1 | - |




### IDPType {#idptype}
authorization framework of the identity provider

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_TYPE_UNSPECIFIED | 0 | - |
| IDP_TYPE_OIDC | 1 | PLANNED: IDP_TYPE_SAML |




### OIDCMappingField {#oidcmappingfield}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_MAPPING_FIELD_UNSPECIFIED | 0 | - |
| OIDC_MAPPING_FIELD_PREFERRED_USERNAME | 1 | - |
| OIDC_MAPPING_FIELD_EMAIL | 2 | - |





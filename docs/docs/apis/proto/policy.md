---
title: zitadel/policy.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### LabelPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| primary_color |  string | - |  |
| secondary_color |  string | - |  |
| is_default |  bool | - |  |
| hide_login_name_suffix |  bool | - |  |




### LoginPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| allow_username_password |  bool | - |  |
| allow_register |  bool | - |  |
| allow_external_idp |  bool | - |  |
| force_mfa |  bool | - |  |
| passwordless_type |  PasswordlessType | - |  |
| is_default |  bool | - |  |




### OrgIAMPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| user_login_must_be_domain |  bool | - |  |
| is_default |  bool | - |  |




### PasswordAgePolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| max_age_days |  uint64 | - |  |
| expire_warn_days |  uint64 | - |  |
| is_default |  bool | - |  |




### PasswordComplexityPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| min_length |  uint64 | - |  |
| has_uppercase |  bool | - |  |
| has_lowercase |  bool | - |  |
| has_number |  bool | - |  |
| has_symbol |  bool | - |  |
| is_default |  bool | - |  |




### PasswordLockoutPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| max_attempts |  uint64 | - |  |
| show_lockout_failure |  bool | - |  |
| is_default |  bool | - |  |






## Enums


### MultiFactorType {#multifactortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| MULTI_FACTOR_TYPE_UNSPECIFIED | 0 | - |
| MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION | 1 | - |




### PasswordlessType {#passwordlesstype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PASSWORDLESS_TYPE_NOT_ALLOWED | 0 | - |
| PASSWORDLESS_TYPE_ALLOWED | 1 | PLANNED: PASSWORDLESS_TYPE_WITH_CERT |




### SecondFactorType {#secondfactortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SECOND_FACTOR_TYPE_UNSPECIFIED | 0 | - |
| SECOND_FACTOR_TYPE_OTP | 1 | - |
| SECOND_FACTOR_TYPE_U2F | 2 | - |





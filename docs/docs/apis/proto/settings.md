---
title: zitadel/settings.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### OIDCSettings



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| access_token_lifetime |  google.protobuf.Duration | - |  |
| id_token_lifetime |  google.protobuf.Duration | - |  |
| refresh_token_idle_expiration |  google.protobuf.Duration | - |  |
| refresh_token_expiration |  google.protobuf.Duration | - |  |




### SMSProvider



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| id |  string | - |  |
| state |  SMSProviderConfigState | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.twilio |  TwilioConfig | - |  |




### SMTPConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| sender_address |  string | - |  |
| sender_name |  string | - |  |
| tls |  bool | - |  |
| host |  string | - |  |
| user |  string | - |  |




### SecretGenerator



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| generator_type |  SecretGeneratorType | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| length |  uint32 | - |  |
| expiry |  google.protobuf.Duration | - |  |
| include_lower_letters |  bool | - |  |
| include_upper_letters |  bool | - |  |
| include_digits |  bool | - |  |
| include_symbols |  bool | - |  |




### SecretGeneratorQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.type_query |  SecretGeneratorTypeQuery | - |  |




### SecretGeneratorTypeQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| generator_type |  SecretGeneratorType | - |  |




### TwilioConfig



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sid |  string | - |  |
| sender_number |  string | - |  |






## Enums


### SMSProviderConfigState {#smsproviderconfigstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED | 0 | - |
| SMS_PROVIDER_CONFIG_ACTIVE | 1 | - |
| SMS_PROVIDER_CONFIG_INACTIVE | 2 | - |




### SecretGeneratorType {#secretgeneratortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SECRET_GENERATOR_TYPE_UNSPECIFIED | 0 | - |
| SECRET_GENERATOR_TYPE_INIT_CODE | 1 | - |
| SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE | 2 | - |
| SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE | 3 | - |
| SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE | 4 | - |
| SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE | 5 | - |
| SECRET_GENERATOR_TYPE_APP_SECRET | 6 | - |





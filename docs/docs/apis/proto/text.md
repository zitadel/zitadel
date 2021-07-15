---
title: zitadel/text.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### EmailVerificationDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |
| login_button_text |  string | - | string.max_len: 100<br />  |




### EmailVerificationScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| code_label |  string | - | string.max_len: 200<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| resend_button_text |  string | - | string.max_len: 100<br />  |




### ExternalUserNotFoundScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| link_button_text |  string | - | string.max_len: 100<br />  |
| auto_register_button_text |  string | - | string.max_len: 100<br />  |




### FooterText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| tos |  string | - | string.max_len: 200<br />  |
| tos_link |  string | - | string.max_len: 500<br />  |
| privacy_policy |  string | - | string.max_len: 200<br />  |
| privacy_policy_link |  string | - | string.max_len: 500<br />  |
| help |  string | - | string.max_len: 200<br />  |
| help_link |  string | - | string.max_len: 500<br />  |




### InitMFADoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### InitMFAOTPScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| description_otp |  string | - | string.max_len: 500<br />  |
| secret_label |  string | - | string.max_len: 200<br />  |
| code_label |  string | - | string.max_len: 200<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |




### InitMFAPromptScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| otp_option |  string | - | string.max_len: 200<br />  |
| u2f_option |  string | - | string.max_len: 200<br />  |
| skip_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### InitMFAU2FScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| token_name_label |  string | - | string.max_len: 200<br />  |
| not_supported |  string | - | string.max_len: 500<br />  |
| register_token_button_text |  string | - | string.max_len: 100<br />  |
| error_retry |  string | - | string.max_len: 500<br />  |




### InitPasswordDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |




### InitPasswordScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| code_label |  string | - | string.max_len: 200<br />  |
| new_password_label |  string | - | string.max_len: 200<br />  |
| new_password_confirm_label |  string | - | string.max_len: 200<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| resend_button_text |  string | - | string.max_len: 100<br />  |




### InitializeUserDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### InitializeUserScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| code_label |  string | - | string.max_len: 200<br />  |
| new_password_label |  string | - | string.max_len: 200<br />  |
| new_password_confirm_label |  string | - | string.max_len: 200<br />  |
| resend_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### LinkingUserDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### LoginCustomText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| select_account_text |  SelectAccountScreenText | - |  |
| login_text |  LoginScreenText | - |  |
| password_text |  PasswordScreenText | - |  |
| username_change_text |  UsernameChangeScreenText | - |  |
| username_change_done_text |  UsernameChangeDoneScreenText | - |  |
| init_password_text |  InitPasswordScreenText | - |  |
| init_password_done_text |  InitPasswordDoneScreenText | - |  |
| email_verification_text |  EmailVerificationScreenText | - |  |
| email_verification_done_text |  EmailVerificationDoneScreenText | - |  |
| initialize_user_text |  InitializeUserScreenText | - |  |
| initialize_done_text |  InitializeUserDoneScreenText | - |  |
| init_mfa_prompt_text |  InitMFAPromptScreenText | - |  |
| init_mfa_otp_text |  InitMFAOTPScreenText | - |  |
| init_mfa_u2f_text |  InitMFAU2FScreenText | - |  |
| init_mfa_done_text |  InitMFADoneScreenText | - |  |
| mfa_providers_text |  MFAProvidersText | - |  |
| verify_mfa_otp_text |  VerifyMFAOTPScreenText | - |  |
| verify_mfa_u2f_text |  VerifyMFAU2FScreenText | - |  |
| passwordless_text |  PasswordlessScreenText | - |  |
| password_change_text |  PasswordChangeScreenText | - |  |
| password_change_done_text |  PasswordChangeDoneScreenText | - |  |
| password_reset_done_text |  PasswordResetDoneScreenText | - |  |
| registration_option_text |  RegistrationOptionScreenText | - |  |
| registration_user_text |  RegistrationUserScreenText | - |  |
| registration_org_text |  RegistrationOrgScreenText | - |  |
| linking_user_done_text |  LinkingUserDoneScreenText | - |  |
| external_user_not_found_text |  ExternalUserNotFoundScreenText | - |  |
| success_login_text |  SuccessLoginScreenText | - |  |
| logout_text |  LogoutDoneScreenText | - |  |
| footer_text |  FooterText | - |  |




### LoginScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| title_linking_process |  string | - | string.max_len: 200<br />  |
| description_linking_process |  string | - | string.max_len: 500<br />  |
| user_must_be_member_of_org |  string | - | string.max_len: 500<br />  |
| login_name_label |  string | - | string.max_len: 200<br />  |
| register_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| external_user_description |  string | - | string.max_len: 500<br />  |
| user_name_placeholder |  string | - | string.max_len: 200<br />  |
| login_name_placeholder |  string | - | string.max_len: 200<br />  |




### LogoutDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| login_button_text |  string | - | string.max_len: 200<br />  |




### MFAProvidersText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| choose_other |  string | - | string.max_len: 500<br />  |
| otp |  string | - | string.max_len: 200<br />  |
| u2f |  string | - | string.max_len: 200<br />  |




### MessageCustomText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| title |  string | - |  |
| pre_header |  string | - |  |
| subject |  string | - |  |
| greeting |  string | - |  |
| text |  string | - |  |
| button_text |  string | - |  |
| footer_text |  string | - |  |




### PasswordChangeDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### PasswordChangeScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| old_password_label |  string | - | string.max_len: 200<br />  |
| new_password_label |  string | - | string.max_len: 200<br />  |
| new_password_confirm_label |  string | - | string.max_len: 200<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### PasswordResetDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### PasswordScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| password_label |  string | - | string.max_len: 200<br />  |
| reset_link_text |  string | - | string.max_len: 100<br />  |
| back_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |
| min_length |  string | - | string.max_len: 100<br />  |
| has_uppercase |  string | - | string.max_len: 100<br />  |
| has_lowercase |  string | - | string.max_len: 100<br />  |
| has_number |  string | - | string.max_len: 100<br />  |
| has_symbol |  string | - | string.max_len: 100<br />  |
| confirmation |  string | - | string.max_len: 100<br />  |




### PasswordlessScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| login_with_pw_button_text |  string | - | string.max_len: 100<br />  |
| validate_token_button_text |  string | - | string.max_len: 200<br />  |
| not_supported |  string | - | string.max_len: 500<br />  |
| error_retry |  string | - | string.max_len: 500<br />  |




### RegistrationOptionScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| user_name_button_text |  string | - | string.max_len: 200<br />  |
| external_login_description |  string | - | string.max_len: 500<br />  |




### RegistrationOrgScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| orgname_label |  string | - | string.max_len: 200<br />  |
| firstname_label |  string | - | string.max_len: 200<br />  |
| lastname_label |  string | - | string.max_len: 200<br />  |
| username_label |  string | - | string.max_len: 200<br />  |
| email_label |  string | - | string.max_len: 200<br />  |
| password_label |  string | - | string.max_len: 200<br />  |
| password_confirm_label |  string | - | string.max_len: 200<br />  |
| tos_and_privacy_label |  string | - | string.max_len: 200<br />  |
| tos_confirm |  string | - | string.max_len: 200<br />  |
| tos_link |  string | - | string.max_len: 200<br />  |
| tos_link_text |  string | - | string.max_len: 200<br />  |
| privacy_confirm |  string | - | string.max_len: 200<br />  |
| privacy_link |  string | - | string.max_len: 200<br />  |
| privacy_link_text |  string | - | string.max_len: 200<br />  |
| external_login_description |  string | - | string.max_len: 500<br />  |
| save_button_text |  string | - | string.max_len: 200<br />  |




### RegistrationUserScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| description_org_register |  string | - | string.max_len: 500<br />  |
| firstname_label |  string | - | string.max_len: 200<br />  |
| lastname_label |  string | - | string.max_len: 200<br />  |
| email_label |  string | - | string.max_len: 200<br />  |
| username_label |  string | - | string.max_len: 200<br />  |
| language_label |  string | - | string.max_len: 200<br />  |
| gender_label |  string | - | string.max_len: 200<br />  |
| password_label |  string | - | string.max_len: 200<br />  |
| password_confirm_label |  string | - | string.max_len: 200<br />  |
| tos_and_privacy_label |  string | - | string.max_len: 200<br />  |
| tos_confirm |  string | - | string.max_len: 200<br />  |
| tos_link |  string | - | string.max_len: 200<br />  |
| tos_link_text |  string | - | string.max_len: 200<br />  |
| privacy_confirm |  string | - | string.max_len: 200<br />  |
| privacy_link |  string | - | string.max_len: 200<br />  |
| privacy_link_text |  string | - | string.max_len: 200<br />  |
| external_login_description |  string | - | string.max_len: 500<br />  |
| next_button_text |  string | - | string.max_len: 200<br />  |
| back_button_text |  string | - | string.max_len: 200<br />  |




### SelectAccountScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| title_linking_process |  string | - | string.max_len: 200<br />  |
| description_linking_process |  string | - | string.max_len: 500<br />  |
| other_user |  string | - | string.max_len: 500<br />  |
| session_state_active |  string | - | string.max_len: 100<br />  |
| session_state_inactive |  string | - | string.max_len: 100<br />  |
| user_must_be_member_of_org |  string | - | string.max_len: 500<br />  |




### SuccessLoginScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| auto_redirect_description |  string | Text to describe that auto redirect should happen after successful login | string.max_len: 500<br />  |
| redirected_description |  string | Text to describe that the window can be closed after redirect | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 200<br />  |




### UsernameChangeDoneScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### UsernameChangeScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| username_label |  string | - | string.max_len: 200<br />  |
| cancel_button_text |  string | - | string.max_len: 100<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### VerifyMFAOTPScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| code_label |  string | - | string.max_len: 200<br />  |
| next_button_text |  string | - | string.max_len: 100<br />  |




### VerifyMFAU2FScreenText



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| title |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| validate_token_text |  string | - | string.max_len: 500<br />  |
| not_supported |  string | - | string.max_len: 500<br />  |
| error_retry |  string | - | string.max_len: 500<br />  |







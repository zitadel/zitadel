import * as _zitadel_proto_zitadel_object_v2_object_pb_js from '@zitadel/proto/zitadel/object/v2/object_pb.js';
import * as _zitadel_proto_zitadel_idp_v2_idp_service_pb_js from '@zitadel/proto/zitadel/idp/v2/idp_service_pb.js';
import * as _zitadel_proto_zitadel_feature_v2_user_pb_js from '@zitadel/proto/zitadel/feature/v2/user_pb.js';
import * as _zitadel_proto_zitadel_feature_v2_organization_pb_js from '@zitadel/proto/zitadel/feature/v2/organization_pb.js';
import * as _zitadel_proto_zitadel_feature_v2_instance_pb_js from '@zitadel/proto/zitadel/feature/v2/instance_pb.js';
import * as _zitadel_proto_zitadel_feature_v2_system_pb_js from '@zitadel/proto/zitadel/feature/v2/system_pb.js';
import * as _zitadel_proto_zitadel_org_v2_org_service_pb_js from '@zitadel/proto/zitadel/org/v2/org_service_pb.js';
import * as _zitadel_proto_zitadel_saml_v2_saml_service_pb_js from '@zitadel/proto/zitadel/saml/v2/saml_service_pb.js';
import * as _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js from '@zitadel/proto/zitadel/oidc/v2/oidc_service_pb.js';
import * as _zitadel_proto_zitadel_session_v2_session_service_pb_js from '@zitadel/proto/zitadel/session/v2/session_service_pb.js';
import * as _zitadel_proto_zitadel_settings_v2_settings_service_pb_js from '@zitadel/proto/zitadel/settings/v2/settings_service_pb.js';
import * as _bufbuild_protobuf_codegenv1 from '@bufbuild/protobuf/codegenv1';
import * as _zitadel_proto_zitadel_user_v2_user_service_pb_js from '@zitadel/proto/zitadel/user/v2/user_service_pb.js';
import * as _connectrpc_connect from '@connectrpc/connect';

declare const createUserServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    createUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.CreateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.CreateUserResponseSchema;
    };
    addHumanUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddHumanUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddHumanUserResponseSchema;
    };
    getUserByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.GetUserByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.GetUserByIDResponseSchema;
    };
    listUsers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListUsersRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListUsersResponseSchema;
    };
    setEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetEmailResponseSchema;
    };
    resendEmailCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ResendEmailCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ResendEmailCodeResponseSchema;
    };
    sendEmailCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SendEmailCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SendEmailCodeResponseSchema;
    };
    verifyEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyEmailResponseSchema;
    };
    setPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetPhoneResponseSchema;
    };
    removePhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemovePhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemovePhoneResponseSchema;
    };
    resendPhoneCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ResendPhoneCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ResendPhoneCodeResponseSchema;
    };
    verifyPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyPhoneResponseSchema;
    };
    updateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.UpdateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.UpdateUserResponseSchema;
    };
    updateHumanUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.UpdateHumanUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.UpdateHumanUserResponseSchema;
    };
    deactivateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.DeactivateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.DeactivateUserResponseSchema;
    };
    reactivateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ReactivateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ReactivateUserResponseSchema;
    };
    lockUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.LockUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.LockUserResponseSchema;
    };
    unlockUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.UnlockUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.UnlockUserResponseSchema;
    };
    deleteUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.DeleteUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.DeleteUserResponseSchema;
    };
    registerPasskey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RegisterPasskeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RegisterPasskeyResponseSchema;
    };
    verifyPasskeyRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyPasskeyRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyPasskeyRegistrationResponseSchema;
    };
    createPasskeyRegistrationLink: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.CreatePasskeyRegistrationLinkRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.CreatePasskeyRegistrationLinkResponseSchema;
    };
    listPasskeys: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListPasskeysRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListPasskeysResponseSchema;
    };
    removePasskey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemovePasskeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemovePasskeyResponseSchema;
    };
    registerU2F: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RegisterU2FRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RegisterU2FResponseSchema;
    };
    verifyU2FRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyU2FRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyU2FRegistrationResponseSchema;
    };
    removeU2F: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveU2FRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveU2FResponseSchema;
    };
    registerTOTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RegisterTOTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RegisterTOTPResponseSchema;
    };
    verifyTOTPRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyTOTPRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyTOTPRegistrationResponseSchema;
    };
    removeTOTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveTOTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveTOTPResponseSchema;
    };
    addOTPSMS: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddOTPSMSRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddOTPSMSResponseSchema;
    };
    removeOTPSMS: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveOTPSMSRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveOTPSMSResponseSchema;
    };
    addOTPEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddOTPEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddOTPEmailResponseSchema;
    };
    removeOTPEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveOTPEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveOTPEmailResponseSchema;
    };
    startIdentityProviderIntent: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.StartIdentityProviderIntentRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.StartIdentityProviderIntentResponseSchema;
    };
    retrieveIdentityProviderIntent: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RetrieveIdentityProviderIntentRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RetrieveIdentityProviderIntentResponseSchema;
    };
    addIDPLink: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddIDPLinkRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddIDPLinkResponseSchema;
    };
    listIDPLinks: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListIDPLinksRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListIDPLinksResponseSchema;
    };
    removeIDPLink: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveIDPLinkRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveIDPLinkResponseSchema;
    };
    passwordReset: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.PasswordResetRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.PasswordResetResponseSchema;
    };
    setPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetPasswordResponseSchema;
    };
    addSecret: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddSecretRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddSecretResponseSchema;
    };
    removeSecret: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveSecretRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveSecretResponseSchema;
    };
    addKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddKeyResponseSchema;
    };
    removeKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemoveKeyResponseSchema;
    };
    listKeys: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListKeysRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListKeysResponseSchema;
    };
    addPersonalAccessToken: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddPersonalAccessTokenRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.AddPersonalAccessTokenResponseSchema;
    };
    removePersonalAccessToken: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemovePersonalAccessTokenRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.RemovePersonalAccessTokenResponseSchema;
    };
    listPersonalAccessTokens: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListPersonalAccessTokensRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListPersonalAccessTokensResponseSchema;
    };
    listAuthenticationMethodTypes: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListAuthenticationMethodTypesRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListAuthenticationMethodTypesResponseSchema;
    };
    listAuthenticationFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListAuthenticationFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListAuthenticationFactorsResponseSchema;
    };
    createInviteCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.CreateInviteCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.CreateInviteCodeResponseSchema;
    };
    resendInviteCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ResendInviteCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ResendInviteCodeResponseSchema;
    };
    verifyInviteCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyInviteCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.VerifyInviteCodeResponseSchema;
    };
    humanMFAInitSkipped: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.HumanMFAInitSkippedRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.HumanMFAInitSkippedResponseSchema;
    };
    setUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.SetUserMetadataResponseSchema;
    };
    listUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.ListUserMetadataResponseSchema;
    };
    deleteUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.DeleteUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_user_v2_user_service_pb_js.DeleteUserMetadataResponseSchema;
    };
}>>;
declare const createSettingsServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    getGeneralSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetGeneralSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetGeneralSettingsResponseSchema;
    };
    getLoginSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetLoginSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetLoginSettingsResponseSchema;
    };
    getActiveIdentityProviders: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetActiveIdentityProvidersRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetActiveIdentityProvidersResponseSchema;
    };
    getPasswordComplexitySettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetPasswordComplexitySettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetPasswordComplexitySettingsResponseSchema;
    };
    getPasswordExpirySettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetPasswordExpirySettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetPasswordExpirySettingsResponseSchema;
    };
    getBrandingSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetBrandingSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetBrandingSettingsResponseSchema;
    };
    getDomainSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetDomainSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetDomainSettingsResponseSchema;
    };
    getLegalAndSupportSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetLegalAndSupportSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetLegalAndSupportSettingsResponseSchema;
    };
    getLockoutSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetLockoutSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetLockoutSettingsResponseSchema;
    };
    getSecuritySettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetSecuritySettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetSecuritySettingsResponseSchema;
    };
    setSecuritySettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.SetSecuritySettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.SetSecuritySettingsResponseSchema;
    };
    getHostedLoginTranslation: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetHostedLoginTranslationRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.GetHostedLoginTranslationResponseSchema;
    };
    setHostedLoginTranslation: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.SetHostedLoginTranslationRequestSchema;
        output: typeof _zitadel_proto_zitadel_settings_v2_settings_service_pb_js.SetHostedLoginTranslationResponseSchema;
    };
}>>;
declare const createSessionServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    listSessions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.ListSessionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.ListSessionsResponseSchema;
    };
    getSession: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.GetSessionRequestSchema;
        output: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.GetSessionResponseSchema;
    };
    createSession: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.CreateSessionRequestSchema;
        output: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.CreateSessionResponseSchema;
    };
    setSession: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.SetSessionRequestSchema;
        output: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.SetSessionResponseSchema;
    };
    deleteSession: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.DeleteSessionRequestSchema;
        output: typeof _zitadel_proto_zitadel_session_v2_session_service_pb_js.DeleteSessionResponseSchema;
    };
}>>;
declare const createOIDCServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    getAuthRequest: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.GetAuthRequestRequestSchema;
        output: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.GetAuthRequestResponseSchema;
    };
    createCallback: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.CreateCallbackRequestSchema;
        output: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.CreateCallbackResponseSchema;
    };
    getDeviceAuthorizationRequest: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.GetDeviceAuthorizationRequestRequestSchema;
        output: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.GetDeviceAuthorizationRequestResponseSchema;
    };
    authorizeOrDenyDeviceAuthorization: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.AuthorizeOrDenyDeviceAuthorizationRequestSchema;
        output: typeof _zitadel_proto_zitadel_oidc_v2_oidc_service_pb_js.AuthorizeOrDenyDeviceAuthorizationResponseSchema;
    };
}>>;
declare const createSAMLServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    getSAMLRequest: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_saml_v2_saml_service_pb_js.GetSAMLRequestRequestSchema;
        output: typeof _zitadel_proto_zitadel_saml_v2_saml_service_pb_js.GetSAMLRequestResponseSchema;
    };
    createResponse: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_saml_v2_saml_service_pb_js.CreateResponseRequestSchema;
        output: typeof _zitadel_proto_zitadel_saml_v2_saml_service_pb_js.CreateResponseResponseSchema;
    };
}>>;
declare const createOrganizationServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    addOrganization: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_org_v2_org_service_pb_js.AddOrganizationRequestSchema;
        output: typeof _zitadel_proto_zitadel_org_v2_org_service_pb_js.AddOrganizationResponseSchema;
    };
    listOrganizations: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_org_v2_org_service_pb_js.ListOrganizationsRequestSchema;
        output: typeof _zitadel_proto_zitadel_org_v2_org_service_pb_js.ListOrganizationsResponseSchema;
    };
}>>;
declare const createFeatureServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    setSystemFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_system_pb_js.SetSystemFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_system_pb_js.SetSystemFeaturesResponseSchema;
    };
    resetSystemFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_system_pb_js.ResetSystemFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_system_pb_js.ResetSystemFeaturesResponseSchema;
    };
    getSystemFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_system_pb_js.GetSystemFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_system_pb_js.GetSystemFeaturesResponseSchema;
    };
    setInstanceFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_instance_pb_js.SetInstanceFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_instance_pb_js.SetInstanceFeaturesResponseSchema;
    };
    resetInstanceFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_instance_pb_js.ResetInstanceFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_instance_pb_js.ResetInstanceFeaturesResponseSchema;
    };
    getInstanceFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_instance_pb_js.GetInstanceFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_instance_pb_js.GetInstanceFeaturesResponseSchema;
    };
    setOrganizationFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_organization_pb_js.SetOrganizationFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_organization_pb_js.SetOrganizationFeaturesResponseSchema;
    };
    resetOrganizationFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_organization_pb_js.ResetOrganizationFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_organization_pb_js.ResetOrganizationFeaturesResponseSchema;
    };
    getOrganizationFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_organization_pb_js.GetOrganizationFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_organization_pb_js.GetOrganizationFeaturesResponseSchema;
    };
    setUserFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_user_pb_js.SetUserFeatureRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_user_pb_js.SetUserFeaturesResponseSchema;
    };
    resetUserFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_user_pb_js.ResetUserFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_user_pb_js.ResetUserFeaturesResponseSchema;
    };
    getUserFeatures: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_feature_v2_user_pb_js.GetUserFeaturesRequestSchema;
        output: typeof _zitadel_proto_zitadel_feature_v2_user_pb_js.GetUserFeaturesResponseSchema;
    };
}>>;
declare const createIdpServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    getIDPByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_idp_v2_idp_service_pb_js.GetIDPByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_idp_v2_idp_service_pb_js.GetIDPByIDResponseSchema;
    };
}>>;
declare function makeReqCtx(orgId: string | undefined): _zitadel_proto_zitadel_object_v2_object_pb_js.RequestContext;

export { createFeatureServiceClient, createIdpServiceClient, createOIDCServiceClient, createOrganizationServiceClient, createSAMLServiceClient, createSessionServiceClient, createSettingsServiceClient, createUserServiceClient, makeReqCtx };

import * as _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js from '@zitadel/proto/zitadel/resources/user/v3alpha/user_service_pb.js';
import * as _bufbuild_protobuf_codegenv1 from '@bufbuild/protobuf/codegenv1';
import * as _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js from '@zitadel/proto/zitadel/resources/userschema/v3alpha/user_schema_service_pb.js';
import * as _connectrpc_connect from '@connectrpc/connect';

declare const createUserSchemaServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    searchUserSchemas: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.SearchUserSchemasRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.SearchUserSchemasResponseSchema;
    };
    getUserSchema: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.GetUserSchemaRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.GetUserSchemaResponseSchema;
    };
    createUserSchema: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.CreateUserSchemaRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.CreateUserSchemaResponseSchema;
    };
    patchUserSchema: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.PatchUserSchemaRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.PatchUserSchemaResponseSchema;
    };
    deactivateUserSchema: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.DeactivateUserSchemaRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.DeactivateUserSchemaResponseSchema;
    };
    reactivateUserSchema: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.ReactivateUserSchemaRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.ReactivateUserSchemaResponseSchema;
    };
    deleteUserSchema: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.DeleteUserSchemaRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_userschema_v3alpha_user_schema_service_pb_js.DeleteUserSchemaResponseSchema;
    };
}>>;
declare const createUserServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    searchUsers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SearchUsersRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SearchUsersResponseSchema;
    };
    getUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.GetUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.GetUserResponseSchema;
    };
    createUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.CreateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.CreateUserResponseSchema;
    };
    patchUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.PatchUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.PatchUserResponseSchema;
    };
    deactivateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.DeactivateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.DeactivateUserResponseSchema;
    };
    activateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.ActivateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.ActivateUserResponseSchema;
    };
    lockUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.LockUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.LockUserResponseSchema;
    };
    unlockUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.UnlockUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.UnlockUserResponseSchema;
    };
    deleteUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.DeleteUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.DeleteUserResponseSchema;
    };
    setContactEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SetContactEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SetContactEmailResponseSchema;
    };
    verifyContactEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyContactEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyContactEmailResponseSchema;
    };
    resendContactEmailCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.ResendContactEmailCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.ResendContactEmailCodeResponseSchema;
    };
    setContactPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SetContactPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SetContactPhoneResponseSchema;
    };
    verifyContactPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyContactPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyContactPhoneResponseSchema;
    };
    resendContactPhoneCode: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.ResendContactPhoneCodeRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.ResendContactPhoneCodeResponseSchema;
    };
    addUsername: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddUsernameRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddUsernameResponseSchema;
    };
    removeUsername: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveUsernameRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveUsernameResponseSchema;
    };
    setPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SetPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.SetPasswordResponseSchema;
    };
    requestPasswordReset: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RequestPasswordResetRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RequestPasswordResetResponseSchema;
    };
    startWebAuthNRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.StartWebAuthNRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.StartWebAuthNRegistrationResponseSchema;
    };
    verifyWebAuthNRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyWebAuthNRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyWebAuthNRegistrationResponseSchema;
    };
    createWebAuthNRegistrationLink: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.CreateWebAuthNRegistrationLinkRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.CreateWebAuthNRegistrationLinkResponseSchema;
    };
    removeWebAuthNAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveWebAuthNAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveWebAuthNAuthenticatorResponseSchema;
    };
    startTOTPRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.StartTOTPRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.StartTOTPRegistrationResponseSchema;
    };
    verifyTOTPRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyTOTPRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyTOTPRegistrationResponseSchema;
    };
    removeTOTPAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveTOTPAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveTOTPAuthenticatorResponseSchema;
    };
    addOTPSMSAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddOTPSMSAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddOTPSMSAuthenticatorResponseSchema;
    };
    verifyOTPSMSRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyOTPSMSRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyOTPSMSRegistrationResponseSchema;
    };
    removeOTPSMSAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveOTPSMSAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveOTPSMSAuthenticatorResponseSchema;
    };
    addOTPEmailAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddOTPEmailAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddOTPEmailAuthenticatorResponseSchema;
    };
    verifyOTPEmailRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyOTPEmailRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.VerifyOTPEmailRegistrationResponseSchema;
    };
    removeOTPEmailAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveOTPEmailAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveOTPEmailAuthenticatorResponseSchema;
    };
    startIdentityProviderIntent: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.StartIdentityProviderIntentRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.StartIdentityProviderIntentResponseSchema;
    };
    getIdentityProviderIntent: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.GetIdentityProviderIntentRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.GetIdentityProviderIntentResponseSchema;
    };
    addIDPAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddIDPAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.AddIDPAuthenticatorResponseSchema;
    };
    removeIDPAuthenticator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveIDPAuthenticatorRequestSchema;
        output: typeof _zitadel_proto_zitadel_resources_user_v3alpha_user_service_pb_js.RemoveIDPAuthenticatorResponseSchema;
    };
}>>;

export { createUserSchemaServiceClient, createUserServiceClient };

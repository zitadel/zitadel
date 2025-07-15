import * as _zitadel_proto_zitadel_system_pb_js from '@zitadel/proto/zitadel/system_pb.js';
import * as _zitadel_proto_zitadel_management_pb_js from '@zitadel/proto/zitadel/management_pb.js';
import * as _zitadel_proto_zitadel_auth_pb_js from '@zitadel/proto/zitadel/auth_pb.js';
import * as _bufbuild_protobuf_codegenv1 from '@bufbuild/protobuf/codegenv1';
import * as _zitadel_proto_zitadel_admin_pb_js from '@zitadel/proto/zitadel/admin_pb.js';
import * as _connectrpc_connect from '@connectrpc/connect';

declare const createAdminServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    healthz: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.HealthzRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.HealthzResponseSchema;
    };
    getSupportedLanguages: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetSupportedLanguagesRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetSupportedLanguagesResponseSchema;
    };
    getAllowedLanguages: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetAllowedLanguagesRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetAllowedLanguagesResponseSchema;
    };
    setDefaultLanguage: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultLanguageRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultLanguageResponseSchema;
    };
    getDefaultLanguage: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultLanguageRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultLanguageResponseSchema;
    };
    getMyInstance: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetMyInstanceRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetMyInstanceResponseSchema;
    };
    listInstanceDomains: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListInstanceDomainsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListInstanceDomainsResponseSchema;
    };
    listInstanceTrustedDomains: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListInstanceTrustedDomainsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListInstanceTrustedDomainsResponseSchema;
    };
    addInstanceTrustedDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddInstanceTrustedDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddInstanceTrustedDomainResponseSchema;
    };
    removeInstanceTrustedDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveInstanceTrustedDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveInstanceTrustedDomainResponseSchema;
    };
    listSecretGenerators: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListSecretGeneratorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListSecretGeneratorsResponseSchema;
    };
    getSecretGenerator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetSecretGeneratorRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetSecretGeneratorResponseSchema;
    };
    updateSecretGenerator: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSecretGeneratorRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSecretGeneratorResponseSchema;
    };
    getSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetSMTPConfigResponseSchema;
    };
    getSMTPConfigById: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetSMTPConfigByIdRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetSMTPConfigByIdResponseSchema;
    };
    addSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddSMTPConfigResponseSchema;
    };
    updateSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMTPConfigResponseSchema;
    };
    updateSMTPConfigPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMTPConfigPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMTPConfigPasswordResponseSchema;
    };
    activateSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateSMTPConfigResponseSchema;
    };
    deactivateSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateSMTPConfigResponseSchema;
    };
    removeSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveSMTPConfigResponseSchema;
    };
    testSMTPConfigById: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.TestSMTPConfigByIdRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.TestSMTPConfigByIdResponseSchema;
    };
    testSMTPConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.TestSMTPConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.TestSMTPConfigResponseSchema;
    };
    listSMTPConfigs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListSMTPConfigsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListSMTPConfigsResponseSchema;
    };
    listEmailProviders: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListEmailProvidersRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListEmailProvidersResponseSchema;
    };
    getEmailProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetEmailProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetEmailProviderResponseSchema;
    };
    getEmailProviderById: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetEmailProviderByIdRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetEmailProviderByIdResponseSchema;
    };
    addEmailProviderSMTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddEmailProviderSMTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddEmailProviderSMTPResponseSchema;
    };
    updateEmailProviderSMTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateEmailProviderSMTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateEmailProviderSMTPResponseSchema;
    };
    addEmailProviderHTTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddEmailProviderHTTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddEmailProviderHTTPResponseSchema;
    };
    updateEmailProviderHTTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateEmailProviderHTTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateEmailProviderHTTPResponseSchema;
    };
    updateEmailProviderSMTPPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateEmailProviderSMTPPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateEmailProviderSMTPPasswordResponseSchema;
    };
    activateEmailProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateEmailProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateEmailProviderResponseSchema;
    };
    deactivateEmailProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateEmailProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateEmailProviderResponseSchema;
    };
    removeEmailProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveEmailProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveEmailProviderResponseSchema;
    };
    testEmailProviderSMTPById: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.TestEmailProviderSMTPByIdRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.TestEmailProviderSMTPByIdResponseSchema;
    };
    testEmailProviderSMTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.TestEmailProviderSMTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.TestEmailProviderSMTPResponseSchema;
    };
    listSMSProviders: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListSMSProvidersRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListSMSProvidersResponseSchema;
    };
    getSMSProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetSMSProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetSMSProviderResponseSchema;
    };
    addSMSProviderTwilio: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddSMSProviderTwilioRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddSMSProviderTwilioResponseSchema;
    };
    updateSMSProviderTwilio: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMSProviderTwilioRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMSProviderTwilioResponseSchema;
    };
    updateSMSProviderTwilioToken: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMSProviderTwilioTokenRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMSProviderTwilioTokenResponseSchema;
    };
    addSMSProviderHTTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddSMSProviderHTTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddSMSProviderHTTPResponseSchema;
    };
    updateSMSProviderHTTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMSProviderHTTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSMSProviderHTTPResponseSchema;
    };
    activateSMSProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateSMSProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateSMSProviderResponseSchema;
    };
    deactivateSMSProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateSMSProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateSMSProviderResponseSchema;
    };
    removeSMSProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveSMSProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveSMSProviderResponseSchema;
    };
    getOIDCSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetOIDCSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetOIDCSettingsResponseSchema;
    };
    addOIDCSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddOIDCSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddOIDCSettingsResponseSchema;
    };
    updateOIDCSettings: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateOIDCSettingsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateOIDCSettingsResponseSchema;
    };
    getFileSystemNotificationProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetFileSystemNotificationProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetFileSystemNotificationProviderResponseSchema;
    };
    getLogNotificationProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetLogNotificationProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetLogNotificationProviderResponseSchema;
    };
    getSecurityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetSecurityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetSecurityPolicyResponseSchema;
    };
    setSecurityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetSecurityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetSecurityPolicyResponseSchema;
    };
    getOrgByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetOrgByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetOrgByIDResponseSchema;
    };
    isOrgUnique: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.IsOrgUniqueRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.IsOrgUniqueResponseSchema;
    };
    setDefaultOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultOrgResponseSchema;
    };
    getDefaultOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultOrgResponseSchema;
    };
    listOrgs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListOrgsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListOrgsResponseSchema;
    };
    setUpOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetUpOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetUpOrgResponseSchema;
    };
    removeOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveOrgResponseSchema;
    };
    getIDPByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetIDPByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetIDPByIDResponseSchema;
    };
    listIDPs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListIDPsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListIDPsResponseSchema;
    };
    addOIDCIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddOIDCIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddOIDCIDPResponseSchema;
    };
    addJWTIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddJWTIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddJWTIDPResponseSchema;
    };
    updateIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIDPResponseSchema;
    };
    deactivateIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.DeactivateIDPResponseSchema;
    };
    reactivateIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ReactivateIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ReactivateIDPResponseSchema;
    };
    removeIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveIDPResponseSchema;
    };
    updateIDPOIDCConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIDPOIDCConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIDPOIDCConfigResponseSchema;
    };
    updateIDPJWTConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIDPJWTConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIDPJWTConfigResponseSchema;
    };
    listProviders: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListProvidersRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListProvidersResponseSchema;
    };
    getProviderByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetProviderByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetProviderByIDResponseSchema;
    };
    addGenericOAuthProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGenericOAuthProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGenericOAuthProviderResponseSchema;
    };
    updateGenericOAuthProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGenericOAuthProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGenericOAuthProviderResponseSchema;
    };
    addGenericOIDCProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGenericOIDCProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGenericOIDCProviderResponseSchema;
    };
    updateGenericOIDCProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGenericOIDCProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGenericOIDCProviderResponseSchema;
    };
    migrateGenericOIDCProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.MigrateGenericOIDCProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.MigrateGenericOIDCProviderResponseSchema;
    };
    addJWTProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddJWTProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddJWTProviderResponseSchema;
    };
    updateJWTProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateJWTProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateJWTProviderResponseSchema;
    };
    addAzureADProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddAzureADProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddAzureADProviderResponseSchema;
    };
    updateAzureADProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateAzureADProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateAzureADProviderResponseSchema;
    };
    addGitHubProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitHubProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitHubProviderResponseSchema;
    };
    updateGitHubProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitHubProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitHubProviderResponseSchema;
    };
    addGitHubEnterpriseServerProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitHubEnterpriseServerProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitHubEnterpriseServerProviderResponseSchema;
    };
    updateGitHubEnterpriseServerProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitHubEnterpriseServerProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitHubEnterpriseServerProviderResponseSchema;
    };
    addGitLabProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitLabProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitLabProviderResponseSchema;
    };
    updateGitLabProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitLabProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitLabProviderResponseSchema;
    };
    addGitLabSelfHostedProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitLabSelfHostedProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGitLabSelfHostedProviderResponseSchema;
    };
    updateGitLabSelfHostedProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitLabSelfHostedProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGitLabSelfHostedProviderResponseSchema;
    };
    addGoogleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddGoogleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddGoogleProviderResponseSchema;
    };
    updateGoogleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGoogleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateGoogleProviderResponseSchema;
    };
    addLDAPProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddLDAPProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddLDAPProviderResponseSchema;
    };
    updateLDAPProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLDAPProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLDAPProviderResponseSchema;
    };
    addAppleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddAppleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddAppleProviderResponseSchema;
    };
    updateAppleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateAppleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateAppleProviderResponseSchema;
    };
    addSAMLProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddSAMLProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddSAMLProviderResponseSchema;
    };
    updateSAMLProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSAMLProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateSAMLProviderResponseSchema;
    };
    regenerateSAMLProviderCertificate: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RegenerateSAMLProviderCertificateRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RegenerateSAMLProviderCertificateResponseSchema;
    };
    deleteProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.DeleteProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.DeleteProviderResponseSchema;
    };
    getOrgIAMPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetOrgIAMPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetOrgIAMPolicyResponseSchema;
    };
    updateOrgIAMPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateOrgIAMPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateOrgIAMPolicyResponseSchema;
    };
    getCustomOrgIAMPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomOrgIAMPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomOrgIAMPolicyResponseSchema;
    };
    addCustomOrgIAMPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddCustomOrgIAMPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddCustomOrgIAMPolicyResponseSchema;
    };
    updateCustomOrgIAMPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateCustomOrgIAMPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateCustomOrgIAMPolicyResponseSchema;
    };
    resetCustomOrgIAMPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomOrgIAMPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomOrgIAMPolicyToDefaultResponseSchema;
    };
    getDomainPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDomainPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDomainPolicyResponseSchema;
    };
    updateDomainPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateDomainPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateDomainPolicyResponseSchema;
    };
    getCustomDomainPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomDomainPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomDomainPolicyResponseSchema;
    };
    addCustomDomainPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddCustomDomainPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddCustomDomainPolicyResponseSchema;
    };
    updateCustomDomainPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateCustomDomainPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateCustomDomainPolicyResponseSchema;
    };
    resetCustomDomainPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomDomainPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomDomainPolicyToDefaultResponseSchema;
    };
    getLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetLabelPolicyResponseSchema;
    };
    getPreviewLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetPreviewLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetPreviewLabelPolicyResponseSchema;
    };
    updateLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLabelPolicyResponseSchema;
    };
    activateLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateLabelPolicyResponseSchema;
    };
    removeLabelPolicyLogo: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyLogoRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyLogoResponseSchema;
    };
    removeLabelPolicyLogoDark: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyLogoDarkRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyLogoDarkResponseSchema;
    };
    removeLabelPolicyIcon: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyIconRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyIconResponseSchema;
    };
    removeLabelPolicyIconDark: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyIconDarkRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyIconDarkResponseSchema;
    };
    removeLabelPolicyFont: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyFontRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveLabelPolicyFontResponseSchema;
    };
    getLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetLoginPolicyResponseSchema;
    };
    updateLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLoginPolicyResponseSchema;
    };
    listLoginPolicyIDPs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListLoginPolicyIDPsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListLoginPolicyIDPsResponseSchema;
    };
    addIDPToLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddIDPToLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddIDPToLoginPolicyResponseSchema;
    };
    removeIDPFromLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveIDPFromLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveIDPFromLoginPolicyResponseSchema;
    };
    listLoginPolicySecondFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListLoginPolicySecondFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListLoginPolicySecondFactorsResponseSchema;
    };
    addSecondFactorToLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddSecondFactorToLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddSecondFactorToLoginPolicyResponseSchema;
    };
    removeSecondFactorFromLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveSecondFactorFromLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveSecondFactorFromLoginPolicyResponseSchema;
    };
    listLoginPolicyMultiFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListLoginPolicyMultiFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListLoginPolicyMultiFactorsResponseSchema;
    };
    addMultiFactorToLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddMultiFactorToLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddMultiFactorToLoginPolicyResponseSchema;
    };
    removeMultiFactorFromLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveMultiFactorFromLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveMultiFactorFromLoginPolicyResponseSchema;
    };
    getPasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetPasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetPasswordComplexityPolicyResponseSchema;
    };
    updatePasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdatePasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdatePasswordComplexityPolicyResponseSchema;
    };
    getPasswordAgePolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetPasswordAgePolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetPasswordAgePolicyResponseSchema;
    };
    updatePasswordAgePolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdatePasswordAgePolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdatePasswordAgePolicyResponseSchema;
    };
    getLockoutPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetLockoutPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetLockoutPolicyResponseSchema;
    };
    updateLockoutPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLockoutPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateLockoutPolicyResponseSchema;
    };
    getPrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetPrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetPrivacyPolicyResponseSchema;
    };
    updatePrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdatePrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdatePrivacyPolicyResponseSchema;
    };
    addNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddNotificationPolicyResponseSchema;
    };
    getNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetNotificationPolicyResponseSchema;
    };
    updateNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateNotificationPolicyResponseSchema;
    };
    getDefaultInitMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultInitMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultInitMessageTextResponseSchema;
    };
    getCustomInitMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomInitMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomInitMessageTextResponseSchema;
    };
    setDefaultInitMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultInitMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultInitMessageTextResponseSchema;
    };
    resetCustomInitMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomInitMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomInitMessageTextToDefaultResponseSchema;
    };
    getDefaultPasswordResetMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultPasswordResetMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultPasswordResetMessageTextResponseSchema;
    };
    getCustomPasswordResetMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomPasswordResetMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomPasswordResetMessageTextResponseSchema;
    };
    setDefaultPasswordResetMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultPasswordResetMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultPasswordResetMessageTextResponseSchema;
    };
    resetCustomPasswordResetMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomPasswordResetMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomPasswordResetMessageTextToDefaultResponseSchema;
    };
    getDefaultVerifyEmailMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifyEmailMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifyEmailMessageTextResponseSchema;
    };
    getCustomVerifyEmailMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifyEmailMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifyEmailMessageTextResponseSchema;
    };
    setDefaultVerifyEmailMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifyEmailMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifyEmailMessageTextResponseSchema;
    };
    resetCustomVerifyEmailMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifyEmailMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifyEmailMessageTextToDefaultResponseSchema;
    };
    getDefaultVerifyPhoneMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifyPhoneMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifyPhoneMessageTextResponseSchema;
    };
    getCustomVerifyPhoneMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifyPhoneMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifyPhoneMessageTextResponseSchema;
    };
    setDefaultVerifyPhoneMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifyPhoneMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifyPhoneMessageTextResponseSchema;
    };
    resetCustomVerifyPhoneMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifyPhoneMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifyPhoneMessageTextToDefaultResponseSchema;
    };
    getDefaultVerifySMSOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifySMSOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifySMSOTPMessageTextResponseSchema;
    };
    getCustomVerifySMSOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifySMSOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifySMSOTPMessageTextResponseSchema;
    };
    setDefaultVerifySMSOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifySMSOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifySMSOTPMessageTextResponseSchema;
    };
    resetCustomVerifySMSOTPMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifySMSOTPMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifySMSOTPMessageTextToDefaultResponseSchema;
    };
    getDefaultVerifyEmailOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifyEmailOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultVerifyEmailOTPMessageTextResponseSchema;
    };
    getCustomVerifyEmailOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifyEmailOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomVerifyEmailOTPMessageTextResponseSchema;
    };
    setDefaultVerifyEmailOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifyEmailOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultVerifyEmailOTPMessageTextResponseSchema;
    };
    resetCustomVerifyEmailOTPMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifyEmailOTPMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomVerifyEmailOTPMessageTextToDefaultResponseSchema;
    };
    getDefaultDomainClaimedMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultDomainClaimedMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultDomainClaimedMessageTextResponseSchema;
    };
    getCustomDomainClaimedMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomDomainClaimedMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomDomainClaimedMessageTextResponseSchema;
    };
    setDefaultDomainClaimedMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultDomainClaimedMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultDomainClaimedMessageTextResponseSchema;
    };
    resetCustomDomainClaimedMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomDomainClaimedMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomDomainClaimedMessageTextToDefaultResponseSchema;
    };
    getDefaultPasswordlessRegistrationMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultPasswordlessRegistrationMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultPasswordlessRegistrationMessageTextResponseSchema;
    };
    getCustomPasswordlessRegistrationMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomPasswordlessRegistrationMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomPasswordlessRegistrationMessageTextResponseSchema;
    };
    setDefaultPasswordlessRegistrationMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultPasswordlessRegistrationMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultPasswordlessRegistrationMessageTextResponseSchema;
    };
    resetCustomPasswordlessRegistrationMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomPasswordlessRegistrationMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomPasswordlessRegistrationMessageTextToDefaultResponseSchema;
    };
    getDefaultPasswordChangeMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultPasswordChangeMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultPasswordChangeMessageTextResponseSchema;
    };
    getCustomPasswordChangeMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomPasswordChangeMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomPasswordChangeMessageTextResponseSchema;
    };
    setDefaultPasswordChangeMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultPasswordChangeMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultPasswordChangeMessageTextResponseSchema;
    };
    resetCustomPasswordChangeMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomPasswordChangeMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomPasswordChangeMessageTextToDefaultResponseSchema;
    };
    getDefaultInviteUserMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultInviteUserMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultInviteUserMessageTextResponseSchema;
    };
    getCustomInviteUserMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomInviteUserMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomInviteUserMessageTextResponseSchema;
    };
    setDefaultInviteUserMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultInviteUserMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetDefaultInviteUserMessageTextResponseSchema;
    };
    resetCustomInviteUserMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomInviteUserMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomInviteUserMessageTextToDefaultResponseSchema;
    };
    getDefaultLoginTexts: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultLoginTextsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetDefaultLoginTextsResponseSchema;
    };
    getCustomLoginTexts: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomLoginTextsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetCustomLoginTextsResponseSchema;
    };
    setCustomLoginText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetCustomLoginTextsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetCustomLoginTextsResponseSchema;
    };
    resetCustomLoginTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomLoginTextsToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ResetCustomLoginTextsToDefaultResponseSchema;
    };
    listIAMMemberRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListIAMMemberRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListIAMMemberRolesResponseSchema;
    };
    listIAMMembers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListIAMMembersRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListIAMMembersResponseSchema;
    };
    addIAMMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.AddIAMMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.AddIAMMemberResponseSchema;
    };
    updateIAMMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIAMMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.UpdateIAMMemberResponseSchema;
    };
    removeIAMMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveIAMMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveIAMMemberResponseSchema;
    };
    listViews: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListViewsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListViewsResponseSchema;
    };
    listFailedEvents: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListFailedEventsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListFailedEventsResponseSchema;
    };
    removeFailedEvent: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveFailedEventRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.RemoveFailedEventResponseSchema;
    };
    importData: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ImportDataRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ImportDataResponseSchema;
    };
    exportData: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ExportDataRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ExportDataResponseSchema;
    };
    listEventTypes: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListEventTypesRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListEventTypesResponseSchema;
    };
    listEvents: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListEventsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListEventsResponseSchema;
    };
    listAggregateTypes: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListAggregateTypesRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListAggregateTypesResponseSchema;
    };
    activateFeatureLoginDefaultOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateFeatureLoginDefaultOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ActivateFeatureLoginDefaultOrgResponseSchema;
    };
    listMilestones: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.ListMilestonesRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.ListMilestonesResponseSchema;
    };
    setRestrictions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.SetRestrictionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.SetRestrictionsResponseSchema;
    };
    getRestrictions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_admin_pb_js.GetRestrictionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_admin_pb_js.GetRestrictionsResponseSchema;
    };
}>>;
declare const createAuthServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    healthz: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.HealthzRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.HealthzResponseSchema;
    };
    getSupportedLanguages: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetSupportedLanguagesRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetSupportedLanguagesResponseSchema;
    };
    getMyUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyUserResponseSchema;
    };
    removeMyUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyUserResponseSchema;
    };
    listMyUserChanges: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyUserChangesRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyUserChangesResponseSchema;
    };
    listMyUserSessions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyUserSessionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyUserSessionsResponseSchema;
    };
    listMyMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyMetadataResponseSchema;
    };
    getMyMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyMetadataResponseSchema;
    };
    listMyRefreshTokens: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyRefreshTokensRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyRefreshTokensResponseSchema;
    };
    revokeMyRefreshToken: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RevokeMyRefreshTokenRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RevokeMyRefreshTokenResponseSchema;
    };
    revokeAllMyRefreshTokens: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RevokeAllMyRefreshTokensRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RevokeAllMyRefreshTokensResponseSchema;
    };
    updateMyUserName: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.UpdateMyUserNameRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.UpdateMyUserNameResponseSchema;
    };
    getMyPasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyPasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyPasswordComplexityPolicyResponseSchema;
    };
    updateMyPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.UpdateMyPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.UpdateMyPasswordResponseSchema;
    };
    getMyProfile: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyProfileRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyProfileResponseSchema;
    };
    updateMyProfile: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.UpdateMyProfileRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.UpdateMyProfileResponseSchema;
    };
    getMyEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyEmailResponseSchema;
    };
    setMyEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.SetMyEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.SetMyEmailResponseSchema;
    };
    verifyMyEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyEmailResponseSchema;
    };
    resendMyEmailVerification: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ResendMyEmailVerificationRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ResendMyEmailVerificationResponseSchema;
    };
    getMyPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyPhoneResponseSchema;
    };
    setMyPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.SetMyPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.SetMyPhoneResponseSchema;
    };
    verifyMyPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyPhoneResponseSchema;
    };
    resendMyPhoneVerification: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ResendMyPhoneVerificationRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ResendMyPhoneVerificationResponseSchema;
    };
    removeMyPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyPhoneResponseSchema;
    };
    removeMyAvatar: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAvatarRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAvatarResponseSchema;
    };
    listMyLinkedIDPs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyLinkedIDPsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyLinkedIDPsResponseSchema;
    };
    removeMyLinkedIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyLinkedIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyLinkedIDPResponseSchema;
    };
    listMyAuthFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyAuthFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyAuthFactorsResponseSchema;
    };
    addMyAuthFactorOTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorOTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorOTPResponseSchema;
    };
    verifyMyAuthFactorOTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyAuthFactorOTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyAuthFactorOTPResponseSchema;
    };
    removeMyAuthFactorOTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorOTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorOTPResponseSchema;
    };
    addMyAuthFactorOTPSMS: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorOTPSMSRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorOTPSMSResponseSchema;
    };
    removeMyAuthFactorOTPSMS: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorOTPSMSRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorOTPSMSResponseSchema;
    };
    addMyAuthFactorOTPEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorOTPEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorOTPEmailResponseSchema;
    };
    removeMyAuthFactorOTPEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorOTPEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorOTPEmailResponseSchema;
    };
    addMyAuthFactorU2F: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorU2FRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyAuthFactorU2FResponseSchema;
    };
    verifyMyAuthFactorU2F: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyAuthFactorU2FRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyAuthFactorU2FResponseSchema;
    };
    removeMyAuthFactorU2F: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorU2FRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyAuthFactorU2FResponseSchema;
    };
    listMyPasswordless: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyPasswordlessRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyPasswordlessResponseSchema;
    };
    addMyPasswordless: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyPasswordlessRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyPasswordlessResponseSchema;
    };
    addMyPasswordlessLink: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyPasswordlessLinkRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.AddMyPasswordlessLinkResponseSchema;
    };
    sendMyPasswordlessLink: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.SendMyPasswordlessLinkRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.SendMyPasswordlessLinkResponseSchema;
    };
    verifyMyPasswordless: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyPasswordlessRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.VerifyMyPasswordlessResponseSchema;
    };
    removeMyPasswordless: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyPasswordlessRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.RemoveMyPasswordlessResponseSchema;
    };
    listMyUserGrants: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyUserGrantsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyUserGrantsResponseSchema;
    };
    listMyProjectOrgs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyProjectOrgsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyProjectOrgsResponseSchema;
    };
    listMyZitadelPermissions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyZitadelPermissionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyZitadelPermissionsResponseSchema;
    };
    listMyProjectPermissions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyProjectPermissionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyProjectPermissionsResponseSchema;
    };
    listMyMemberships: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyMembershipsRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.ListMyMembershipsResponseSchema;
    };
    getMyLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyLabelPolicyResponseSchema;
    };
    getMyPrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyPrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyPrivacyPolicyResponseSchema;
    };
    getMyLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_auth_pb_js.GetMyLoginPolicyResponseSchema;
    };
}>>;
declare const createManagementServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    healthz: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.HealthzRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.HealthzResponseSchema;
    };
    getOIDCInformation: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetOIDCInformationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetOIDCInformationResponseSchema;
    };
    getIAM: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetIAMRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetIAMResponseSchema;
    };
    getSupportedLanguages: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetSupportedLanguagesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetSupportedLanguagesResponseSchema;
    };
    getUserByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetUserByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetUserByIDResponseSchema;
    };
    getUserByLoginNameGlobal: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetUserByLoginNameGlobalRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetUserByLoginNameGlobalResponseSchema;
    };
    listUsers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListUsersRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListUsersResponseSchema;
    };
    listUserChanges: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListUserChangesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListUserChangesResponseSchema;
    };
    isUserUnique: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.IsUserUniqueRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.IsUserUniqueResponseSchema;
    };
    addHumanUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddHumanUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddHumanUserResponseSchema;
    };
    importHumanUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ImportHumanUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ImportHumanUserResponseSchema;
    };
    addMachineUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddMachineUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddMachineUserResponseSchema;
    };
    deactivateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateUserResponseSchema;
    };
    reactivateUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateUserResponseSchema;
    };
    lockUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.LockUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.LockUserResponseSchema;
    };
    unlockUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UnlockUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UnlockUserResponseSchema;
    };
    removeUser: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveUserRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveUserResponseSchema;
    };
    updateUserName: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateUserNameRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateUserNameResponseSchema;
    };
    setUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetUserMetadataResponseSchema;
    };
    bulkSetUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.BulkSetUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.BulkSetUserMetadataResponseSchema;
    };
    listUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListUserMetadataResponseSchema;
    };
    getUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetUserMetadataResponseSchema;
    };
    removeUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveUserMetadataResponseSchema;
    };
    bulkRemoveUserMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.BulkRemoveUserMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.BulkRemoveUserMetadataResponseSchema;
    };
    getHumanProfile: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetHumanProfileRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetHumanProfileResponseSchema;
    };
    updateHumanProfile: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateHumanProfileRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateHumanProfileResponseSchema;
    };
    getHumanEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetHumanEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetHumanEmailResponseSchema;
    };
    updateHumanEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateHumanEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateHumanEmailResponseSchema;
    };
    resendHumanInitialization: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResendHumanInitializationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResendHumanInitializationResponseSchema;
    };
    resendHumanEmailVerification: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResendHumanEmailVerificationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResendHumanEmailVerificationResponseSchema;
    };
    getHumanPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetHumanPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetHumanPhoneResponseSchema;
    };
    updateHumanPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateHumanPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateHumanPhoneResponseSchema;
    };
    removeHumanPhone: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanPhoneRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanPhoneResponseSchema;
    };
    resendHumanPhoneVerification: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResendHumanPhoneVerificationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResendHumanPhoneVerificationResponseSchema;
    };
    removeHumanAvatar: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAvatarRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAvatarResponseSchema;
    };
    setHumanInitialPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetHumanInitialPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetHumanInitialPasswordResponseSchema;
    };
    setHumanPassword: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetHumanPasswordRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetHumanPasswordResponseSchema;
    };
    sendHumanResetPasswordNotification: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SendHumanResetPasswordNotificationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SendHumanResetPasswordNotificationResponseSchema;
    };
    listHumanAuthFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListHumanAuthFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListHumanAuthFactorsResponseSchema;
    };
    removeHumanAuthFactorOTP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorOTPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorOTPResponseSchema;
    };
    removeHumanAuthFactorU2F: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorU2FRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorU2FResponseSchema;
    };
    removeHumanAuthFactorOTPSMS: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorOTPSMSRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorOTPSMSResponseSchema;
    };
    removeHumanAuthFactorOTPEmail: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorOTPEmailRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanAuthFactorOTPEmailResponseSchema;
    };
    listHumanPasswordless: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListHumanPasswordlessRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListHumanPasswordlessResponseSchema;
    };
    addPasswordlessRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddPasswordlessRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddPasswordlessRegistrationResponseSchema;
    };
    sendPasswordlessRegistration: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SendPasswordlessRegistrationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SendPasswordlessRegistrationResponseSchema;
    };
    removeHumanPasswordless: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanPasswordlessRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanPasswordlessResponseSchema;
    };
    updateMachine: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateMachineRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateMachineResponseSchema;
    };
    generateMachineSecret: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GenerateMachineSecretRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GenerateMachineSecretResponseSchema;
    };
    removeMachineSecret: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveMachineSecretRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveMachineSecretResponseSchema;
    };
    getMachineKeyByIDs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetMachineKeyByIDsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetMachineKeyByIDsResponseSchema;
    };
    listMachineKeys: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListMachineKeysRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListMachineKeysResponseSchema;
    };
    addMachineKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddMachineKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddMachineKeyResponseSchema;
    };
    removeMachineKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveMachineKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveMachineKeyResponseSchema;
    };
    getPersonalAccessTokenByIDs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetPersonalAccessTokenByIDsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetPersonalAccessTokenByIDsResponseSchema;
    };
    listPersonalAccessTokens: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListPersonalAccessTokensRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListPersonalAccessTokensResponseSchema;
    };
    addPersonalAccessToken: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddPersonalAccessTokenRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddPersonalAccessTokenResponseSchema;
    };
    removePersonalAccessToken: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemovePersonalAccessTokenRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemovePersonalAccessTokenResponseSchema;
    };
    listHumanLinkedIDPs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListHumanLinkedIDPsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListHumanLinkedIDPsResponseSchema;
    };
    removeHumanLinkedIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanLinkedIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveHumanLinkedIDPResponseSchema;
    };
    listUserMemberships: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListUserMembershipsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListUserMembershipsResponseSchema;
    };
    getMyOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetMyOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetMyOrgResponseSchema;
    };
    getOrgByDomainGlobal: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgByDomainGlobalRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgByDomainGlobalResponseSchema;
    };
    listOrgChanges: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgChangesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgChangesResponseSchema;
    };
    addOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgResponseSchema;
    };
    updateOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgResponseSchema;
    };
    deactivateOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateOrgResponseSchema;
    };
    reactivateOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateOrgResponseSchema;
    };
    removeOrg: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgResponseSchema;
    };
    setOrgMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetOrgMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetOrgMetadataResponseSchema;
    };
    bulkSetOrgMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.BulkSetOrgMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.BulkSetOrgMetadataResponseSchema;
    };
    listOrgMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgMetadataResponseSchema;
    };
    getOrgMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgMetadataResponseSchema;
    };
    removeOrgMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgMetadataResponseSchema;
    };
    bulkRemoveOrgMetadata: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.BulkRemoveOrgMetadataRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.BulkRemoveOrgMetadataResponseSchema;
    };
    addOrgDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgDomainResponseSchema;
    };
    listOrgDomains: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgDomainsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgDomainsResponseSchema;
    };
    removeOrgDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgDomainResponseSchema;
    };
    generateOrgDomainValidation: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GenerateOrgDomainValidationRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GenerateOrgDomainValidationResponseSchema;
    };
    validateOrgDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ValidateOrgDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ValidateOrgDomainResponseSchema;
    };
    setPrimaryOrgDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetPrimaryOrgDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetPrimaryOrgDomainResponseSchema;
    };
    listOrgMemberRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgMemberRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgMemberRolesResponseSchema;
    };
    listOrgMembers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgMembersRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgMembersResponseSchema;
    };
    addOrgMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgMemberResponseSchema;
    };
    updateOrgMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgMemberResponseSchema;
    };
    removeOrgMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgMemberResponseSchema;
    };
    getProjectByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetProjectByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetProjectByIDResponseSchema;
    };
    getGrantedProjectByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetGrantedProjectByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetGrantedProjectByIDResponseSchema;
    };
    listProjects: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectsResponseSchema;
    };
    listGrantedProjects: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListGrantedProjectsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListGrantedProjectsResponseSchema;
    };
    listGrantedProjectRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListGrantedProjectRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListGrantedProjectRolesResponseSchema;
    };
    listProjectChanges: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectChangesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectChangesResponseSchema;
    };
    addProject: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectResponseSchema;
    };
    updateProject: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectResponseSchema;
    };
    deactivateProject: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateProjectRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateProjectResponseSchema;
    };
    reactivateProject: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateProjectRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateProjectResponseSchema;
    };
    removeProject: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectResponseSchema;
    };
    listProjectRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectRolesResponseSchema;
    };
    addProjectRole: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectRoleRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectRoleResponseSchema;
    };
    bulkAddProjectRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.BulkAddProjectRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.BulkAddProjectRolesResponseSchema;
    };
    updateProjectRole: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectRoleRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectRoleResponseSchema;
    };
    removeProjectRole: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectRoleRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectRoleResponseSchema;
    };
    listProjectMemberRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectMemberRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectMemberRolesResponseSchema;
    };
    listProjectMembers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectMembersRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectMembersResponseSchema;
    };
    addProjectMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectMemberResponseSchema;
    };
    updateProjectMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectMemberResponseSchema;
    };
    removeProjectMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectMemberResponseSchema;
    };
    getAppByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetAppByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetAppByIDResponseSchema;
    };
    listApps: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListAppsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListAppsResponseSchema;
    };
    listAppChanges: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListAppChangesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListAppChangesResponseSchema;
    };
    addOIDCApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddOIDCAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddOIDCAppResponseSchema;
    };
    addSAMLApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddSAMLAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddSAMLAppResponseSchema;
    };
    addAPIApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddAPIAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddAPIAppResponseSchema;
    };
    updateApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAppResponseSchema;
    };
    updateOIDCAppConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOIDCAppConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOIDCAppConfigResponseSchema;
    };
    updateSAMLAppConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateSAMLAppConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateSAMLAppConfigResponseSchema;
    };
    updateAPIAppConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAPIAppConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAPIAppConfigResponseSchema;
    };
    deactivateApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateAppResponseSchema;
    };
    reactivateApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateAppResponseSchema;
    };
    removeApp: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveAppRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveAppResponseSchema;
    };
    regenerateOIDCClientSecret: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RegenerateOIDCClientSecretRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RegenerateOIDCClientSecretResponseSchema;
    };
    regenerateAPIClientSecret: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RegenerateAPIClientSecretRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RegenerateAPIClientSecretResponseSchema;
    };
    getAppKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetAppKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetAppKeyResponseSchema;
    };
    listAppKeys: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListAppKeysRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListAppKeysResponseSchema;
    };
    addAppKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddAppKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddAppKeyResponseSchema;
    };
    removeAppKey: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveAppKeyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveAppKeyResponseSchema;
    };
    listProjectGrantChanges: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantChangesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantChangesResponseSchema;
    };
    getProjectGrantByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetProjectGrantByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetProjectGrantByIDResponseSchema;
    };
    listProjectGrants: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantsResponseSchema;
    };
    listAllProjectGrants: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListAllProjectGrantsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListAllProjectGrantsResponseSchema;
    };
    addProjectGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectGrantResponseSchema;
    };
    updateProjectGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectGrantResponseSchema;
    };
    deactivateProjectGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateProjectGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateProjectGrantResponseSchema;
    };
    reactivateProjectGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateProjectGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateProjectGrantResponseSchema;
    };
    removeProjectGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectGrantResponseSchema;
    };
    listProjectGrantMemberRoles: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantMemberRolesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantMemberRolesResponseSchema;
    };
    listProjectGrantMembers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantMembersRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProjectGrantMembersResponseSchema;
    };
    addProjectGrantMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectGrantMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddProjectGrantMemberResponseSchema;
    };
    updateProjectGrantMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectGrantMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateProjectGrantMemberResponseSchema;
    };
    removeProjectGrantMember: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectGrantMemberRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveProjectGrantMemberResponseSchema;
    };
    getUserGrantByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetUserGrantByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetUserGrantByIDResponseSchema;
    };
    listUserGrants: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListUserGrantResponseSchema;
    };
    addUserGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddUserGrantResponseSchema;
    };
    updateUserGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateUserGrantResponseSchema;
    };
    deactivateUserGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateUserGrantResponseSchema;
    };
    reactivateUserGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateUserGrantResponseSchema;
    };
    removeUserGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveUserGrantResponseSchema;
    };
    bulkRemoveUserGrant: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.BulkRemoveUserGrantRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.BulkRemoveUserGrantResponseSchema;
    };
    getOrgIAMPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgIAMPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgIAMPolicyResponseSchema;
    };
    getDomainPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDomainPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDomainPolicyResponseSchema;
    };
    getLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetLoginPolicyResponseSchema;
    };
    getDefaultLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLoginPolicyResponseSchema;
    };
    addCustomLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomLoginPolicyResponseSchema;
    };
    updateCustomLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomLoginPolicyResponseSchema;
    };
    resetLoginPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetLoginPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetLoginPolicyToDefaultResponseSchema;
    };
    listLoginPolicyIDPs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListLoginPolicyIDPsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListLoginPolicyIDPsResponseSchema;
    };
    addIDPToLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddIDPToLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddIDPToLoginPolicyResponseSchema;
    };
    removeIDPFromLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveIDPFromLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveIDPFromLoginPolicyResponseSchema;
    };
    listLoginPolicySecondFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListLoginPolicySecondFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListLoginPolicySecondFactorsResponseSchema;
    };
    addSecondFactorToLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddSecondFactorToLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddSecondFactorToLoginPolicyResponseSchema;
    };
    removeSecondFactorFromLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveSecondFactorFromLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveSecondFactorFromLoginPolicyResponseSchema;
    };
    listLoginPolicyMultiFactors: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListLoginPolicyMultiFactorsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListLoginPolicyMultiFactorsResponseSchema;
    };
    addMultiFactorToLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddMultiFactorToLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddMultiFactorToLoginPolicyResponseSchema;
    };
    removeMultiFactorFromLoginPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveMultiFactorFromLoginPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveMultiFactorFromLoginPolicyResponseSchema;
    };
    getPasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetPasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetPasswordComplexityPolicyResponseSchema;
    };
    getDefaultPasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordComplexityPolicyResponseSchema;
    };
    addCustomPasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomPasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomPasswordComplexityPolicyResponseSchema;
    };
    updateCustomPasswordComplexityPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomPasswordComplexityPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomPasswordComplexityPolicyResponseSchema;
    };
    resetPasswordComplexityPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetPasswordComplexityPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetPasswordComplexityPolicyToDefaultResponseSchema;
    };
    getPasswordAgePolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetPasswordAgePolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetPasswordAgePolicyResponseSchema;
    };
    getDefaultPasswordAgePolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordAgePolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordAgePolicyResponseSchema;
    };
    addCustomPasswordAgePolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomPasswordAgePolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomPasswordAgePolicyResponseSchema;
    };
    updateCustomPasswordAgePolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomPasswordAgePolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomPasswordAgePolicyResponseSchema;
    };
    resetPasswordAgePolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetPasswordAgePolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetPasswordAgePolicyToDefaultResponseSchema;
    };
    getLockoutPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetLockoutPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetLockoutPolicyResponseSchema;
    };
    getDefaultLockoutPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLockoutPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLockoutPolicyResponseSchema;
    };
    addCustomLockoutPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomLockoutPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomLockoutPolicyResponseSchema;
    };
    updateCustomLockoutPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomLockoutPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomLockoutPolicyResponseSchema;
    };
    resetLockoutPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetLockoutPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetLockoutPolicyToDefaultResponseSchema;
    };
    getPrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetPrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetPrivacyPolicyResponseSchema;
    };
    getDefaultPrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPrivacyPolicyResponseSchema;
    };
    addCustomPrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomPrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomPrivacyPolicyResponseSchema;
    };
    updateCustomPrivacyPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomPrivacyPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomPrivacyPolicyResponseSchema;
    };
    resetPrivacyPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetPrivacyPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetPrivacyPolicyToDefaultResponseSchema;
    };
    getNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetNotificationPolicyResponseSchema;
    };
    getDefaultNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultNotificationPolicyResponseSchema;
    };
    addCustomNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomNotificationPolicyResponseSchema;
    };
    updateCustomNotificationPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomNotificationPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomNotificationPolicyResponseSchema;
    };
    resetNotificationPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetNotificationPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetNotificationPolicyToDefaultResponseSchema;
    };
    getLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetLabelPolicyResponseSchema;
    };
    getPreviewLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetPreviewLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetPreviewLabelPolicyResponseSchema;
    };
    getDefaultLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLabelPolicyResponseSchema;
    };
    addCustomLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddCustomLabelPolicyResponseSchema;
    };
    updateCustomLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateCustomLabelPolicyResponseSchema;
    };
    activateCustomLabelPolicy: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ActivateCustomLabelPolicyRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ActivateCustomLabelPolicyResponseSchema;
    };
    removeCustomLabelPolicyLogo: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyLogoRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyLogoResponseSchema;
    };
    removeCustomLabelPolicyLogoDark: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyLogoDarkRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyLogoDarkResponseSchema;
    };
    removeCustomLabelPolicyIcon: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyIconRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyIconResponseSchema;
    };
    removeCustomLabelPolicyIconDark: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyIconDarkRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyIconDarkResponseSchema;
    };
    removeCustomLabelPolicyFont: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyFontRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveCustomLabelPolicyFontResponseSchema;
    };
    resetLabelPolicyToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetLabelPolicyToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetLabelPolicyToDefaultResponseSchema;
    };
    getCustomInitMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomInitMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomInitMessageTextResponseSchema;
    };
    getDefaultInitMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultInitMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultInitMessageTextResponseSchema;
    };
    setCustomInitMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomInitMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomInitMessageTextResponseSchema;
    };
    resetCustomInitMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomInitMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomInitMessageTextToDefaultResponseSchema;
    };
    getCustomPasswordResetMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomPasswordResetMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomPasswordResetMessageTextResponseSchema;
    };
    getDefaultPasswordResetMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordResetMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordResetMessageTextResponseSchema;
    };
    setCustomPasswordResetMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomPasswordResetMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomPasswordResetMessageTextResponseSchema;
    };
    resetCustomPasswordResetMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomPasswordResetMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomPasswordResetMessageTextToDefaultResponseSchema;
    };
    getCustomVerifyEmailMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifyEmailMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifyEmailMessageTextResponseSchema;
    };
    getDefaultVerifyEmailMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifyEmailMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifyEmailMessageTextResponseSchema;
    };
    setCustomVerifyEmailMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifyEmailMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifyEmailMessageTextResponseSchema;
    };
    resetCustomVerifyEmailMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifyEmailMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifyEmailMessageTextToDefaultResponseSchema;
    };
    getCustomVerifyPhoneMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifyPhoneMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifyPhoneMessageTextResponseSchema;
    };
    getDefaultVerifyPhoneMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifyPhoneMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifyPhoneMessageTextResponseSchema;
    };
    setCustomVerifyPhoneMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifyPhoneMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifyPhoneMessageTextResponseSchema;
    };
    resetCustomVerifyPhoneMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifyPhoneMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifyPhoneMessageTextToDefaultResponseSchema;
    };
    getCustomVerifySMSOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifySMSOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifySMSOTPMessageTextResponseSchema;
    };
    getDefaultVerifySMSOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifySMSOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifySMSOTPMessageTextResponseSchema;
    };
    setCustomVerifySMSOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifySMSOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifySMSOTPMessageTextResponseSchema;
    };
    resetCustomVerifySMSOTPMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifySMSOTPMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifySMSOTPMessageTextToDefaultResponseSchema;
    };
    getCustomVerifyEmailOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifyEmailOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomVerifyEmailOTPMessageTextResponseSchema;
    };
    getDefaultVerifyEmailOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifyEmailOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultVerifyEmailOTPMessageTextResponseSchema;
    };
    setCustomVerifyEmailOTPMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifyEmailOTPMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomVerifyEmailOTPMessageTextResponseSchema;
    };
    resetCustomVerifyEmailOTPMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifyEmailOTPMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomVerifyEmailOTPMessageTextToDefaultResponseSchema;
    };
    getCustomDomainClaimedMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomDomainClaimedMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomDomainClaimedMessageTextResponseSchema;
    };
    getDefaultDomainClaimedMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultDomainClaimedMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultDomainClaimedMessageTextResponseSchema;
    };
    setCustomDomainClaimedMessageCustomText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomDomainClaimedMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomDomainClaimedMessageTextResponseSchema;
    };
    resetCustomDomainClaimedMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomDomainClaimedMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomDomainClaimedMessageTextToDefaultResponseSchema;
    };
    getCustomPasswordlessRegistrationMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomPasswordlessRegistrationMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomPasswordlessRegistrationMessageTextResponseSchema;
    };
    getDefaultPasswordlessRegistrationMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordlessRegistrationMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordlessRegistrationMessageTextResponseSchema;
    };
    setCustomPasswordlessRegistrationMessageCustomText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomPasswordlessRegistrationMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomPasswordlessRegistrationMessageTextResponseSchema;
    };
    resetCustomPasswordlessRegistrationMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomPasswordlessRegistrationMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomPasswordlessRegistrationMessageTextToDefaultResponseSchema;
    };
    getCustomPasswordChangeMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomPasswordChangeMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomPasswordChangeMessageTextResponseSchema;
    };
    getDefaultPasswordChangeMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordChangeMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultPasswordChangeMessageTextResponseSchema;
    };
    setCustomPasswordChangeMessageCustomText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomPasswordChangeMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomPasswordChangeMessageTextResponseSchema;
    };
    resetCustomPasswordChangeMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomPasswordChangeMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomPasswordChangeMessageTextToDefaultResponseSchema;
    };
    getCustomInviteUserMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomInviteUserMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomInviteUserMessageTextResponseSchema;
    };
    getDefaultInviteUserMessageText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultInviteUserMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultInviteUserMessageTextResponseSchema;
    };
    setCustomInviteUserMessageCustomText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomInviteUserMessageTextRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomInviteUserMessageTextResponseSchema;
    };
    resetCustomInviteUserMessageTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomInviteUserMessageTextToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomInviteUserMessageTextToDefaultResponseSchema;
    };
    getCustomLoginTexts: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomLoginTextsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetCustomLoginTextsResponseSchema;
    };
    getDefaultLoginTexts: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLoginTextsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetDefaultLoginTextsResponseSchema;
    };
    setCustomLoginText: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomLoginTextsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetCustomLoginTextsResponseSchema;
    };
    resetCustomLoginTextToDefault: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomLoginTextsToDefaultRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ResetCustomLoginTextsToDefaultResponseSchema;
    };
    getOrgIDPByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgIDPByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetOrgIDPByIDResponseSchema;
    };
    listOrgIDPs: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgIDPsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListOrgIDPsResponseSchema;
    };
    addOrgOIDCIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgOIDCIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgOIDCIDPResponseSchema;
    };
    addOrgJWTIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgJWTIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddOrgJWTIDPResponseSchema;
    };
    deactivateOrgIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateOrgIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateOrgIDPResponseSchema;
    };
    reactivateOrgIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateOrgIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateOrgIDPResponseSchema;
    };
    removeOrgIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RemoveOrgIDPResponseSchema;
    };
    updateOrgIDP: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgIDPRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgIDPResponseSchema;
    };
    updateOrgIDPOIDCConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgIDPOIDCConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgIDPOIDCConfigResponseSchema;
    };
    updateOrgIDPJWTConfig: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgIDPJWTConfigRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateOrgIDPJWTConfigResponseSchema;
    };
    listProviders: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListProvidersRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListProvidersResponseSchema;
    };
    getProviderByID: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetProviderByIDRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetProviderByIDResponseSchema;
    };
    addGenericOAuthProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGenericOAuthProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGenericOAuthProviderResponseSchema;
    };
    updateGenericOAuthProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGenericOAuthProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGenericOAuthProviderResponseSchema;
    };
    addGenericOIDCProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGenericOIDCProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGenericOIDCProviderResponseSchema;
    };
    updateGenericOIDCProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGenericOIDCProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGenericOIDCProviderResponseSchema;
    };
    migrateGenericOIDCProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.MigrateGenericOIDCProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.MigrateGenericOIDCProviderResponseSchema;
    };
    addJWTProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddJWTProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddJWTProviderResponseSchema;
    };
    updateJWTProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateJWTProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateJWTProviderResponseSchema;
    };
    addAzureADProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddAzureADProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddAzureADProviderResponseSchema;
    };
    updateAzureADProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAzureADProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAzureADProviderResponseSchema;
    };
    addGitHubProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGitHubProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGitHubProviderResponseSchema;
    };
    updateGitHubProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitHubProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitHubProviderResponseSchema;
    };
    addGitHubEnterpriseServerProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGitHubEnterpriseServerProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGitHubEnterpriseServerProviderResponseSchema;
    };
    updateGitHubEnterpriseServerProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitHubEnterpriseServerProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitHubEnterpriseServerProviderResponseSchema;
    };
    addGitLabProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGitLabProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGitLabProviderResponseSchema;
    };
    updateGitLabProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitLabProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitLabProviderResponseSchema;
    };
    addGitLabSelfHostedProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGitLabSelfHostedProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGitLabSelfHostedProviderResponseSchema;
    };
    updateGitLabSelfHostedProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitLabSelfHostedProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGitLabSelfHostedProviderResponseSchema;
    };
    addGoogleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddGoogleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddGoogleProviderResponseSchema;
    };
    updateGoogleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGoogleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateGoogleProviderResponseSchema;
    };
    addLDAPProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddLDAPProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddLDAPProviderResponseSchema;
    };
    updateLDAPProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateLDAPProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateLDAPProviderResponseSchema;
    };
    addAppleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddAppleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddAppleProviderResponseSchema;
    };
    updateAppleProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAppleProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateAppleProviderResponseSchema;
    };
    addSAMLProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.AddSAMLProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.AddSAMLProviderResponseSchema;
    };
    updateSAMLProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateSAMLProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateSAMLProviderResponseSchema;
    };
    regenerateSAMLProviderCertificate: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.RegenerateSAMLProviderCertificateRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.RegenerateSAMLProviderCertificateResponseSchema;
    };
    deleteProvider: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeleteProviderRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeleteProviderResponseSchema;
    };
    listActions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListActionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListActionsResponseSchema;
    };
    getAction: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetActionRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetActionResponseSchema;
    };
    createAction: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.CreateActionRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.CreateActionResponseSchema;
    };
    updateAction: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.UpdateActionRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.UpdateActionResponseSchema;
    };
    deactivateAction: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateActionRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeactivateActionResponseSchema;
    };
    reactivateAction: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateActionRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ReactivateActionResponseSchema;
    };
    deleteAction: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.DeleteActionRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.DeleteActionResponseSchema;
    };
    listFlowTypes: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListFlowTypesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListFlowTypesResponseSchema;
    };
    listFlowTriggerTypes: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ListFlowTriggerTypesRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ListFlowTriggerTypesResponseSchema;
    };
    getFlow: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.GetFlowRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.GetFlowResponseSchema;
    };
    clearFlow: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.ClearFlowRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.ClearFlowResponseSchema;
    };
    setTriggerActions: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_management_pb_js.SetTriggerActionsRequestSchema;
        output: typeof _zitadel_proto_zitadel_management_pb_js.SetTriggerActionsResponseSchema;
    };
}>>;
declare const createSystemServiceClient: (transport: _connectrpc_connect.Transport) => _connectrpc_connect.Client<_bufbuild_protobuf_codegenv1.GenService<{
    healthz: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.HealthzRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.HealthzResponseSchema;
    };
    listInstances: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ListInstancesRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ListInstancesResponseSchema;
    };
    getInstance: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.GetInstanceRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.GetInstanceResponseSchema;
    };
    addInstance: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.AddInstanceRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.AddInstanceResponseSchema;
    };
    updateInstance: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.UpdateInstanceRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.UpdateInstanceResponseSchema;
    };
    createInstance: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.CreateInstanceRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.CreateInstanceResponseSchema;
    };
    removeInstance: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.RemoveInstanceRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.RemoveInstanceResponseSchema;
    };
    listIAMMembers: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ListIAMMembersRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ListIAMMembersResponseSchema;
    };
    existsDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ExistsDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ExistsDomainResponseSchema;
    };
    listDomains: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ListDomainsRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ListDomainsResponseSchema;
    };
    addDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.AddDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.AddDomainResponseSchema;
    };
    removeDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.RemoveDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.RemoveDomainResponseSchema;
    };
    setPrimaryDomain: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.SetPrimaryDomainRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.SetPrimaryDomainResponseSchema;
    };
    listViews: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ListViewsRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ListViewsResponseSchema;
    };
    clearView: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ClearViewRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ClearViewResponseSchema;
    };
    listFailedEvents: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ListFailedEventsRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ListFailedEventsResponseSchema;
    };
    removeFailedEvent: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.RemoveFailedEventRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.RemoveFailedEventResponseSchema;
    };
    addQuota: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.AddQuotaRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.AddQuotaResponseSchema;
    };
    setQuota: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.SetQuotaRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.SetQuotaResponseSchema;
    };
    removeQuota: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.RemoveQuotaRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.RemoveQuotaResponseSchema;
    };
    setInstanceFeature: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.SetInstanceFeatureRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.SetInstanceFeatureResponseSchema;
    };
    setLimits: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.SetLimitsRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.SetLimitsResponseSchema;
    };
    bulkSetLimits: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.BulkSetLimitsRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.BulkSetLimitsResponseSchema;
    };
    resetLimits: {
        methodKind: "unary";
        input: typeof _zitadel_proto_zitadel_system_pb_js.ResetLimitsRequestSchema;
        output: typeof _zitadel_proto_zitadel_system_pb_js.ResetLimitsResponseSchema;
    };
}>>;

export { createAdminServiceClient, createAuthServiceClient, createManagementServiceClient, createSystemServiceClient };

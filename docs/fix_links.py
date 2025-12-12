import os

# Mapping of old paths to new paths
mapping = {
    # User Service V2
    "apis/resources/user_service_v2/user-service-create-user.api.mdx": "/docs/references/api/user/zitadel.user.v2.UserService.CreateUser",
    "apis/resources/user_service_v2/user-service-update-user.api.mdx": "/docs/references/api/user/zitadel.user.v2.UserService.UpdateUser",
    "apis/resources/user_service_v2/user-service-create-invite-code.api.mdx": "/docs/references/api/user/zitadel.user.v2.UserService.CreateInviteCode",
    "apis/resources/user_service_v2/user-service-delete-user-metadata.api.mdx": "/docs/references/api/user/zitadel.user.v2.UserService.DeleteUserMetadata",
    "apis/resources/user_service_v2/user-service-add-idp-link.api.mdx": "/docs/references/api/user/zitadel.user.v2.UserService.AddIDPLink",
    "apis/resources/user_service_v2/user-service-remove-idp-link.api.mdx": "/docs/references/api/user/zitadel.user.v2.UserService.RemoveIDPLink",

    # Org Service V2 Beta
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-update-organization.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.UpdateOrganization",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-deactivate-organization.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.DeactivateOrganization",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-activate-organization.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.ActivateOrganization",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-delete-organization.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.DeleteOrganization",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-set-organization-metadata.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.SetOrganizationMetadata",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-list-organization-metadata.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.ListOrganizationMetadata",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-delete-organization-metadata.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.DeleteOrganizationMetadata",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-add-organization-domain.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.AddOrganizationDomain",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-list-organization-domains.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.ListOrganizationDomains",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-delete-organization-domain.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.DeleteOrganizationDomain",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-generate-organization-domain-validation.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.GenerateOrganizationDomainValidation",
    "apis/resources/org_service_v2beta/zitadel-org-v-2-beta-organization-service-verify-organization-domain.api.mdx": "/docs/references/api/org/zitadel.org.v2beta.OrganizationService.VerifyOrganizationDomain",

    # Application Service V2 (replacing app/v2beta)
    "apis/resources/application_service_v2/application-service-create-application.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.CreateApplication",
    "apis/resources/application_service_v2/zitadel-app-v-2-application-service-update-application.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.UpdateApplication",
    "apis/resources/application_service_v2/application-service-get-application.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.GetApplication",
    "apis/resources/application_service_v2/application-service-delete-application.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.DeleteApplication",
    "apis/resources/application_service_v2/application-service-deactivate-application.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.DeactivateApplication",
    "apis/resources/application_service_v2/application-service-reactivate-application.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.ReactivateApplication",
    "apis/resources/application_service_v2/application-service-generate-client-secret.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.GenerateClientSecret",
    "apis/resources/application_service_v2/application-service-list-applications.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.ListApplications",
    "apis/resources/application_service_v2/application-service-create-application-key.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.CreateApplicationKey",
    "apis/resources/application_service_v2/application-service-delete-application-key.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.DeleteApplicationKey",
    "apis/resources/application_service_v2/application-service-get-application-key.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.GetApplicationKey",
    "apis/resources/application_service_v2/application-service-list-application-keys.api.mdx": "/docs/references/api/application/zitadel.application.v2.ApplicationService.ListApplicationKeys",

    # User Schema Service
    "/docs/references/resources/user_schema_service_v3/user-schema-service-create-user-schema": "/docs/references/api/user_schema/zitadel.resources.userschema.v3alpha.ZITADELUserSchemas.CreateUserSchema",
    "/docs/references/resources/user_schema_service_v3/user-schema-service-update-user-schema": "/docs/references/api/user_schema/zitadel.resources.userschema.v3alpha.ZITADELUserSchemas.PatchUserSchema",
    "/apis/resources/user_schema_service_v3/user-schema-service-list-user-schemas": "/docs/references/api/user_schema/zitadel.resources.userschema.v3alpha.ZITADELUserSchemas.SearchUserSchemas",
    "/apis/resources/user_schema_service_v3/user-schema-service-get-user-schema-by-id": "/docs/references/api/user_schema/zitadel.resources.userschema.v3alpha.ZITADELUserSchemas.GetUserSchema",
    
    # Admin Service
    "/docs/references/resources/admin/admin-service-add-idp-to-login-policy": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.AddIdPToLoginPolicy",
    "/docs/references/resources/admin/admin-service-update-login-policy": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.UpdateLoginPolicy",
    "/docs/references/resources/admin/identity-providers": "/docs/references/api-v1/admin",
    "/docs/references/resources/admin/admin-service-migrate-generic-oidc-provider": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.MigrateGenericOIDCProvider",
    "/docs/references/resources/admin/admin-service-add-saml-provider": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.AddSAMLProvider",
    "/docs/references/resources/admin/admin-service-update-saml-provider": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.UpdateSAMLProvider",
    "/docs/references/resources/admin/admin-service-add-instance-trusted-domain": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.AddInstanceTrustedDomain",
    "/docs/references/resources/admin/admin-service-set-restrictions": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.SetRestrictions",
    "/docs/references/resources/admin/admin-service-set-up-org": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.SetUpOrg",
    "/docs/references/resources/admin/events": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.ListEvents",
    "/apis/resources/admin/admin-service-add-sms-provider-http": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.AddSMSProviderHTTP",
    "/apis/resources/admin/admin-service-activate-sms-provider": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.ActivateSMSProvider",
    "/apis/resources/admin/sms-provider": "/docs/references/api-v1/admin",
    "/apis/resources/admin/admin-service-add-email-provider-http": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.AddEmailProviderHTTP",
    "/apis/resources/admin/admin-service-activate-email-provider": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.ActivateEmailProvider",
    "/apis/resources/admin/admin-service-list-email-providers": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.ListEmailProviders",
    "/apis/resources/admin/feature-restrictions": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.GetRestrictions",
    "/docs/references/resources/admin/admin-service-import-data": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.ImportData",
    "/docs/references/resources/admin/admin-service-activate-feature-login-default-org": "/docs/references/api-v1/admin/zitadel.admin.v1.AdminService.ActivateFeatureLoginDefaultOrg",

    # Management Service
    "/docs/references/resources/mgmt/members": "/docs/references/api-v1/management",
    "/docs/references/resources/mgmt/user-grants": "/docs/references/api-v1/management",
    "/docs/references/resources/mgmt/management-service-list-user-grants": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListUserGrants",
    "/docs/references/resources/mgmt/management-service-get-user-grant-by-id": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.GetUserGrantByID",
    "/docs/references/resources/mgmt/management-service-add-idp-to-login-policy": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.AddIdPToLoginPolicy",
    "/docs/references/resources/mgmt/management-service-update-custom-login-policy": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.UpdateCustomLoginPolicy",
    "/docs/references/resources/mgmt/management-service-migrate-generic-oidc-provider": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.MigrateGenericOIDCProvider",
    "/docs/references/resources/mgmt/management-service-add-saml-provider": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.AddSAMLProvider",
    "/docs/references/resources/mgmt/management-service-update-saml-provider": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.UpdateSAMLProvider",
    "/docs/references/resources/mgmt/management-service-set-org-metadata": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.SetOrgMetadata",
    "/docs/references/resources/mgmt/management-service-bulk-set-org-metadata": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.BulkSetOrgMetadata",
    "/docs/references/resources/mgmt/management-service-set-user-metadata": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.SetUserMetadata",
    "/docs/references/resources/mgmt/management-service-bulk-set-user-metadata": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.BulkSetUserMetadata",
    "/docs/references/resources/mgmt/user-machine": "/docs/references/api-v1/management",
    "/docs/references/resources/mgmt/management-service-update-project": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.UpdateProject",
    "/docs/references/resources/mgmt/management-service-add-machine-user": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.AddMachineUser",
    "/docs/references/resources/mgmt/identity-providers": "/docs/references/api-v1/management",
    "/docs/references/resources/mgmt/applications": "/docs/references/api-v1/management",
    "/docs/references/resources/mgmt/management-service-list-org-member-roles": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListOrgMemberRoles",
    "/docs/references/resources/mgmt/users": "/docs/references/api-v1/management",
    "/docs/references/resources/mgmt/management-service-list-user-metadata": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListUserMetadata",
    "/docs/references/resources/mgmt/management-service-get-user-metadata": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.GetUserMetadata",
    "/docs/references/resources/mgmt/management-service-import-human-user": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ImportHumanUser",
    "/apis/resources/mgmt/management-service-list-user-changes": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListUserChanges",
    "/apis/resources/mgmt/management-service-list-app-changes": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListAppChanges",
    "/apis/resources/mgmt/management-service-list-org-changes": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListOrgChanges",
    "/apis/resources/mgmt/management-service-list-project-changes": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListProjectChanges",
    "/apis/resources/mgmt/management-service-list-project-grant-changes": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.ListProjectGrantChanges",
    "/apis/resources/mgmt/management-service-add-project-grant": "/docs/references/api-v1/management/zitadel.management.v1.ManagementService.AddProjectGrant",

    # Auth Service
    "/docs/references/resources/auth/user-authorizations-grants": "/docs/references/api-v1/auth",
    "/docs/references/resources/auth/auth-service-list-my-project-permissions": "/docs/references/api-v1/auth/zitadel.auth.v1.AuthService.ListMyProjectPermissions",
    "/docs/references/resources/auth/auth-service-list-my-zitadel-permissions": "/docs/references/api-v1/auth/zitadel.auth.v1.AuthService.ListMyZitadelPermissions",
    "/docs/references/resources/auth/auth-service-list-my-user-grants": "/docs/references/api-v1/auth/zitadel.auth.v1.AuthService.ListMyUserGrants",
    "/docs/references/resources/auth/auth-service-list-my-metadata": "/docs/references/api-v1/auth/zitadel.auth.v1.AuthService.ListMyMetadata",
    "/apis/resources/auth/auth-service-list-my-user-changes": "/docs/references/api-v1/auth/zitadel.auth.v1.AuthService.ListMyUserChanges",

    # System Service
    "/docs/references/resources/system/system-service": "/docs/references/api-v1/system",
    "/apis/resources/system/system-service-add-domain": "/docs/references/api-v1/system/zitadel.system.v1.SystemService.AddDomain",
    "/apis/resources/system/limits": "/docs/references/api-v1/system/zitadel.system.v1.SystemService.SetLimits",
    "/apis/resources/system/quotas": "/docs/references/api-v1/system/zitadel.system.v1.SystemService.SetQuota",

    # Generic Service Links (Must be last to avoid partial matches)
    "/apis/resources/admin": "/docs/references/api-v1/admin",
    "/apis/resources/mgmt": "/docs/references/api-v1/management",
    "/apis/resources/system": "/docs/references/api-v1/system",
    "/apis/resources/auth": "/docs/references/api-v1/auth",
}

# Files to update
files_to_update = [
    "content/docs/references/introduction.mdx",
    "content/docs/build-and-integrate/retrieve-user-roles.md",
    "content/docs/build-and-integrate/identity-providers/_activate.mdx",
    "content/docs/build-and-integrate/identity-providers/_custom_login_policy.mdx",
    "content/docs/build-and-integrate/identity-providers/introduction.md",
    "content/docs/build-and-integrate/identity-providers/migrate.mdx",
    "content/docs/build-and-integrate/identity-providers/okta_saml.mdx",
    "content/docs/build-and-integrate/login-ui/login-app.mdx",
    "content/docs/build-and-integrate/onboarding/b2b.mdx",
    "content/docs/build-and-integrate/service-users/authenticate-service-users.md",
    "content/docs/build-and-integrate/service-users/client-credentials.md",
    "content/docs/build-and-integrate/solution-scenarios/b2c.mdx",
    "content/docs/build-and-integrate/solution-scenarios/restrict-console.mdx",
    "content/docs/build-and-integrate/zitadel-apis/access-zitadel-apis.md",
    "content/docs/build-and-integrate/zitadel-apis/access-zitadel-system-api.md",
    "content/docs/build-and-integrate/zitadel-apis/event-api.md",
    "content/docs/concepts/features/audit-trail.md",
    "content/docs/concepts/features/external-user-grant.md",
    "content/docs/concepts/features/identity-brokering.md",
    "content/docs/concepts/structure/applications.md",
    "content/docs/concepts/structure/instance.mdx",
    "content/docs/concepts/structure/managers.mdx",
    "content/docs/concepts/structure/users.md",
    "content/docs/manage-and-govern/customize/notification-providers.mdx",
    "content/docs/manage-and-govern/customize/restrictions.md",
    "content/docs/manage-and-govern/customize/user-metadata.md",
    "content/docs/manage-and-govern/customize/user-schema.md",
    "content/docs/manage-and-govern/migrate/users.md",
    "content/docs/manage-and-govern/user/reg-create-user.md",
    "content/docs/operate-and-self-host/manage/cache.md",
    "content/docs/operate-and-self-host/manage/custom-domain.md",
    "content/docs/operate-and-self-host/manage/usage_control.md",
    "content/docs/manage-and-govern/support/advisory/a10003.md",
    "content/docs/manage-and-govern/support/advisory/a10014.md"
]

for file_path in files_to_update:
    if not os.path.exists(file_path):
        print(f"File not found: {file_path}")
        continue

    with open(file_path, "r") as f:
        content = f.read()

    updated_content = content
    for old, new in mapping.items():
        updated_content = updated_content.replace(old, new)

    if updated_content != content:
        with open(file_path, "w") as f:
            f.write(updated_content)
        print(f"Updated {file_path}")
    else:
        print(f"No changes in {file_path}")

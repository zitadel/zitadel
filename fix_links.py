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
}

# Files to update
files_to_update = [
    "proto/zitadel/management.proto",
    "proto/zitadel/system.proto",
    "proto/zitadel/admin.proto",
    "proto/zitadel/user/v2/user_service.proto",
    "proto/zitadel/app/v2beta/app_service.proto"
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

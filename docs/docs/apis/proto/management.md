---
title: zitadel/management.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)


## ManagementService {#zitadelmanagementv1managementservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)





    GET: /healthz


### GetOIDCInformation

> **rpc** GetOIDCInformation([GetOIDCInformationRequest](#getoidcinformationrequest))
[GetOIDCInformationResponse](#getoidcinformationresponse)





    GET: /zitadel/docs


### GetIAM

> **rpc** GetIAM([GetIAMRequest](#getiamrequest))
[GetIAMResponse](#getiamresponse)

Returns some needed settings of the IAM (Global Organisation ID, Zitadel Project ID)



    GET: /iam


### GetUserByID

> **rpc** GetUserByID([GetUserByIDRequest](#getuserbyidrequest))
[GetUserByIDResponse](#getuserbyidresponse)

Returns the requested full blown user (human or machine)



    GET: /users/{id}


### GetUserByLoginNameGlobal

> **rpc** GetUserByLoginNameGlobal([GetUserByLoginNameGlobalRequest](#getuserbyloginnameglobalrequest))
[GetUserByLoginNameGlobalResponse](#getuserbyloginnameglobalresponse)

Searches a user over all organisations
the login name has to match exactly



    GET: /global/users/_by_login_name


### ListUsers

> **rpc** ListUsers([ListUsersRequest](#listusersrequest))
[ListUsersResponse](#listusersresponse)

Return the users matching the query
Limit should always be set, there is a default limit set by the service



    POST: /users/_search


### ListUserChanges

> **rpc** ListUserChanges([ListUserChangesRequest](#listuserchangesrequest))
[ListUserChangesResponse](#listuserchangesresponse)

Returns the history of the user (each event)
Limit should always be set, there is a default limit set by the service



    POST: /users/{user_id}/changes/_search


### IsUserUnique

> **rpc** IsUserUnique([IsUserUniqueRequest](#isuseruniquerequest))
[IsUserUniqueResponse](#isuseruniqueresponse)

Returns if a user with the searched email or username is unique



    GET: /users/_is_unique


### AddHumanUser

> **rpc** AddHumanUser([AddHumanUserRequest](#addhumanuserrequest))
[AddHumanUserResponse](#addhumanuserresponse)

Create a user of the type human
A email will be sent to the user if email is not verified or no password is set
If a password is given, the user has to change on the next login



    POST: /users/human


### ImportHumanUser

> **rpc** ImportHumanUser([ImportHumanUserRequest](#importhumanuserrequest))
[ImportHumanUserResponse](#importhumanuserresponse)

Create a user of the type human
A email will be sent to the user if email is not verified or no password is set
If a password is given, the user doesn't have to change on the next login



    POST: /users/human/_import


### AddMachineUser

> **rpc** AddMachineUser([AddMachineUserRequest](#addmachineuserrequest))
[AddMachineUserResponse](#addmachineuserresponse)

Create a user of the type machine



    POST: /users/machine


### DeactivateUser

> **rpc** DeactivateUser([DeactivateUserRequest](#deactivateuserrequest))
[DeactivateUserResponse](#deactivateuserresponse)

Changes the user state to deactivated
The user will not be able to login
returns an error if user state is already deactivated



    POST: /users/{id}/_deactivate


### ReactivateUser

> **rpc** ReactivateUser([ReactivateUserRequest](#reactivateuserrequest))
[ReactivateUserResponse](#reactivateuserresponse)

Changes the user state to active
returns an error if user state is not deactivated



    POST: /users/{id}/_reactivate


### LockUser

> **rpc** LockUser([LockUserRequest](#lockuserrequest))
[LockUserResponse](#lockuserresponse)

Changes the user state to deactivated
The user will not be able to login
returns an error if user state is already locked



    POST: /users/{id}/_lock


### UnlockUser

> **rpc** UnlockUser([UnlockUserRequest](#unlockuserrequest))
[UnlockUserResponse](#unlockuserresponse)

Changes the user state to active
returns an error if user state is not locked



    POST: /users/{id}/_unlock


### RemoveUser

> **rpc** RemoveUser([RemoveUserRequest](#removeuserrequest))
[RemoveUserResponse](#removeuserresponse)

Changes the user state to deleted



    DELETE: /users/{id}


### UpdateUserName

> **rpc** UpdateUserName([UpdateUserNameRequest](#updateusernamerequest))
[UpdateUserNameResponse](#updateusernameresponse)

Changes the username



    GET: /users/{user_id}/username


### GetHumanProfile

> **rpc** GetHumanProfile([GetHumanProfileRequest](#gethumanprofilerequest))
[GetHumanProfileResponse](#gethumanprofileresponse)

Returns the profile of the human



    GET: /users/{user_id}/profile


### UpdateHumanProfile

> **rpc** UpdateHumanProfile([UpdateHumanProfileRequest](#updatehumanprofilerequest))
[UpdateHumanProfileResponse](#updatehumanprofileresponse)

Changes the profile of the human



    PUT: /users/{user_id}/profile


### GetHumanEmail

> **rpc** GetHumanEmail([GetHumanEmailRequest](#gethumanemailrequest))
[GetHumanEmailResponse](#gethumanemailresponse)

GetHumanEmail returns the email and verified state of the human



    GET: /users/{user_id}/email


### UpdateHumanEmail

> **rpc** UpdateHumanEmail([UpdateHumanEmailRequest](#updatehumanemailrequest))
[UpdateHumanEmailResponse](#updatehumanemailresponse)

Changes the email of the human
If state is not verified, the user will get a verification email



    PUT: /users/{user_id}/email


### ResendHumanInitialization

> **rpc** ResendHumanInitialization([ResendHumanInitializationRequest](#resendhumaninitializationrequest))
[ResendHumanInitializationResponse](#resendhumaninitializationresponse)

Resends an email to the given email address to finish the initialization process of the user
Changes the email address of the user if it is provided



    POST: /users/{user_id}/_resend_initialization


### ResendHumanEmailVerification

> **rpc** ResendHumanEmailVerification([ResendHumanEmailVerificationRequest](#resendhumanemailverificationrequest))
[ResendHumanEmailVerificationResponse](#resendhumanemailverificationresponse)

Resends an email to the given email address to finish the email verification process of the user



    POST: /users/{user_id}/email/_resend_verification


### GetHumanPhone

> **rpc** GetHumanPhone([GetHumanPhoneRequest](#gethumanphonerequest))
[GetHumanPhoneResponse](#gethumanphoneresponse)

Returns the phone and verified state of the human phone



    GET: /users/{user_id}/phone


### UpdateHumanPhone

> **rpc** UpdateHumanPhone([UpdateHumanPhoneRequest](#updatehumanphonerequest))
[UpdateHumanPhoneResponse](#updatehumanphoneresponse)

Changes the phone number
If verified is not set, the user will get an sms to verify the number



    PUT: /users/{user_id}/phone


### RemoveHumanPhone

> **rpc** RemoveHumanPhone([RemoveHumanPhoneRequest](#removehumanphonerequest))
[RemoveHumanPhoneResponse](#removehumanphoneresponse)

Removes the phone number of the human



    DELETE: /users/{user_id}/phone


### ResendHumanPhoneVerification

> **rpc** ResendHumanPhoneVerification([ResendHumanPhoneVerificationRequest](#resendhumanphoneverificationrequest))
[ResendHumanPhoneVerificationResponse](#resendhumanphoneverificationresponse)

An sms will be sent to the given phone number to finish the phone verification process of the user



    POST: /users/{user_id}/phone/_resend_verification


### SetHumanInitialPassword

> **rpc** SetHumanInitialPassword([SetHumanInitialPasswordRequest](#sethumaninitialpasswordrequest))
[SetHumanInitialPasswordResponse](#sethumaninitialpasswordresponse)

A Manager is only allowed to set an initial password, on the next login the user has to change his password



    POST: /users/{user_id}/password/_initialize


### SendHumanResetPasswordNotification

> **rpc** SendHumanResetPasswordNotification([SendHumanResetPasswordNotificationRequest](#sendhumanresetpasswordnotificationrequest))
[SendHumanResetPasswordNotificationResponse](#sendhumanresetpasswordnotificationresponse)

An email will be sent to the given address to reset the password of the user



    POST: /users/{user_id}/password/_reset


### ListHumanAuthFactors

> **rpc** ListHumanAuthFactors([ListHumanAuthFactorsRequest](#listhumanauthfactorsrequest))
[ListHumanAuthFactorsResponse](#listhumanauthfactorsresponse)

Returns a list of all factors (second and multi) which are configured on the user



    POST: /users/{user_id}/auth_factors/_search


### RemoveHumanAuthFactorOTP

> **rpc** RemoveHumanAuthFactorOTP([RemoveHumanAuthFactorOTPRequest](#removehumanauthfactorotprequest))
[RemoveHumanAuthFactorOTPResponse](#removehumanauthfactorotpresponse)

The otp second factor will be removed from the user
Because only one otp can be configured per user, the configured one will be removed



    DELETE: /users/{user_id}/auth_factors/otp


### RemoveHumanAuthFactorU2F

> **rpc** RemoveHumanAuthFactorU2F([RemoveHumanAuthFactorU2FRequest](#removehumanauthfactoru2frequest))
[RemoveHumanAuthFactorU2FResponse](#removehumanauthfactoru2fresponse)

The u2f (universial second factor) will be removed from the user



    DELETE: /users/{user_id}/auth_factors/u2f/{token_id}


### ListHumanPasswordless

> **rpc** ListHumanPasswordless([ListHumanPasswordlessRequest](#listhumanpasswordlessrequest))
[ListHumanPasswordlessResponse](#listhumanpasswordlessresponse)

Returns all configured passwordless authentications



    POST: /users/{user_id}/passwordless/_search


### RemoveHumanPasswordless

> **rpc** RemoveHumanPasswordless([RemoveHumanPasswordlessRequest](#removehumanpasswordlessrequest))
[RemoveHumanPasswordlessResponse](#removehumanpasswordlessresponse)

Removed a configured passwordless authentication



    DELETE: /users/{user_id}/passwordless/{token_id}


### UpdateMachine

> **rpc** UpdateMachine([UpdateMachineRequest](#updatemachinerequest))
[UpdateMachineResponse](#updatemachineresponse)

Changes a machine user



    PUT: /users/{user_id}/machine


### GetMachineKeyByIDs

> **rpc** GetMachineKeyByIDs([GetMachineKeyByIDsRequest](#getmachinekeybyidsrequest))
[GetMachineKeyByIDsResponse](#getmachinekeybyidsresponse)

Returns a machine key of a (machine) user



    GET: /users/{user_id}/keys/{key_id}


### ListMachineKeys

> **rpc** ListMachineKeys([ListMachineKeysRequest](#listmachinekeysrequest))
[ListMachineKeysResponse](#listmachinekeysresponse)

Returns all machine keys of a (machine) user which match the query
Limit should always be set, there is a default limit set by the service



    POST: /users/{user_id}/keys/_search


### AddMachineKey

> **rpc** AddMachineKey([AddMachineKeyRequest](#addmachinekeyrequest))
[AddMachineKeyResponse](#addmachinekeyresponse)

Generates a new machine key, details should be stored after return



    POST: /users/{user_id}/keys


### RemoveMachineKey

> **rpc** RemoveMachineKey([RemoveMachineKeyRequest](#removemachinekeyrequest))
[RemoveMachineKeyResponse](#removemachinekeyresponse)

Removed a machine key



    DELETE: /users/{user_id}/keys/{key_id}


### ListHumanLinkedIDPs

> **rpc** ListHumanLinkedIDPs([ListHumanLinkedIDPsRequest](#listhumanlinkedidpsrequest))
[ListHumanLinkedIDPsResponse](#listhumanlinkedidpsresponse)

Lists all identity providers (social logins) which a human has configured (e.g Google, Microsoft, AD, etc..)
Limit should always be set, there is a default limit set by the service



    POST: /users/{user_id}/idps/_search


### RemoveHumanLinkedIDP

> **rpc** RemoveHumanLinkedIDP([RemoveHumanLinkedIDPRequest](#removehumanlinkedidprequest))
[RemoveHumanLinkedIDPResponse](#removehumanlinkedidpresponse)

Removed a configured identity provider (social login) of a human



    DELETE: /users/{user_id}/idps/{idp_id}/{linked_user_id}


### ListUserMemberships

> **rpc** ListUserMemberships([ListUserMembershipsRequest](#listusermembershipsrequest))
[ListUserMembershipsResponse](#listusermembershipsresponse)

Show all the permissions a user has iin ZITADEL (ZITADEL Manager)
Limit should always be set, there is a default limit set by the service



    POST: /users/{user_id}/memberships/_search


### GetMyOrg

> **rpc** GetMyOrg([GetMyOrgRequest](#getmyorgrequest))
[GetMyOrgResponse](#getmyorgresponse)

Returns the org given in the header



    GET: /orgs/me


### GetOrgByDomainGlobal

> **rpc** GetOrgByDomainGlobal([GetOrgByDomainGlobalRequest](#getorgbydomainglobalrequest))
[GetOrgByDomainGlobalResponse](#getorgbydomainglobalresponse)

Search a org over all organisations
Domain must match exactly



    GET: /global/orgs/_by_domain


### ListOrgChanges

> **rpc** ListOrgChanges([ListOrgChangesRequest](#listorgchangesrequest))
[ListOrgChangesResponse](#listorgchangesresponse)

Returns the history of my organisation (each event)
Limit should always be set, there is a default limit set by the service



    POST: /orgs/me/changes/_search


### AddOrg

> **rpc** AddOrg([AddOrgRequest](#addorgrequest))
[AddOrgResponse](#addorgresponse)

Creates a new organisation



    POST: /orgs


### DeactivateOrg

> **rpc** DeactivateOrg([DeactivateOrgRequest](#deactivateorgrequest))
[DeactivateOrgResponse](#deactivateorgresponse)

Sets the state of my organisation to deactivated
Users of this organisation will not be able login



    POST: /orgs/me/_deactivate


### ReactivateOrg

> **rpc** ReactivateOrg([ReactivateOrgRequest](#reactivateorgrequest))
[ReactivateOrgResponse](#reactivateorgresponse)

Sets the state of my organisation to active



    POST: /orgs/me/_reactivate


### ListOrgDomains

> **rpc** ListOrgDomains([ListOrgDomainsRequest](#listorgdomainsrequest))
[ListOrgDomainsResponse](#listorgdomainsresponse)

Returns all registered domains of my organisation
Limit should always be set, there is a default limit set by the service



    POST: /orgs/me/domains/_search


### AddOrgDomain

> **rpc** AddOrgDomain([AddOrgDomainRequest](#addorgdomainrequest))
[AddOrgDomainResponse](#addorgdomainresponse)

Adds a new domain to my organisation



    POST: /orgs/me/domains


### RemoveOrgDomain

> **rpc** RemoveOrgDomain([RemoveOrgDomainRequest](#removeorgdomainrequest))
[RemoveOrgDomainResponse](#removeorgdomainresponse)

Removed the domain from my organisation



    DELETE: /orgs/me/domains/{domain}


### GenerateOrgDomainValidation

> **rpc** GenerateOrgDomainValidation([GenerateOrgDomainValidationRequest](#generateorgdomainvalidationrequest))
[GenerateOrgDomainValidationResponse](#generateorgdomainvalidationresponse)

Generates a new file to validate you domain



    POST: /orgs/me/domains/{domain}/validation/_generate


### ValidateOrgDomain

> **rpc** ValidateOrgDomain([ValidateOrgDomainRequest](#validateorgdomainrequest))
[ValidateOrgDomainResponse](#validateorgdomainresponse)

Validates your domain with the choosen method
Validated domains must be unique



    POST: /orgs/me/domains/{domain}/validation/_validate


### SetPrimaryOrgDomain

> **rpc** SetPrimaryOrgDomain([SetPrimaryOrgDomainRequest](#setprimaryorgdomainrequest))
[SetPrimaryOrgDomainResponse](#setprimaryorgdomainresponse)

Sets the domain as primary
Primary domain is shown as suffix on the preferred username on the users of the organisation



    POST: /orgs/me/domains/{domain}/_set_primary


### ListOrgMemberRoles

> **rpc** ListOrgMemberRoles([ListOrgMemberRolesRequest](#listorgmemberrolesrequest))
[ListOrgMemberRolesResponse](#listorgmemberrolesresponse)

Returns all ZITADEL roles which are for organisation managers



    POST: /orgs/members/roles/_search


### ListOrgMembers

> **rpc** ListOrgMembers([ListOrgMembersRequest](#listorgmembersrequest))
[ListOrgMembersResponse](#listorgmembersresponse)

Returns all ZITADEL managers of this organisation (Project and Project Grant managers not included)
Limit should always be set, there is a default limit set by the service



    POST: /orgs/me/members/_search


### AddOrgMember

> **rpc** AddOrgMember([AddOrgMemberRequest](#addorgmemberrequest))
[AddOrgMemberResponse](#addorgmemberresponse)

Adds a new organisation manager, which is allowed to administrate ZITADEL



    POST: /orgs/me/members


### UpdateOrgMember

> **rpc** UpdateOrgMember([UpdateOrgMemberRequest](#updateorgmemberrequest))
[UpdateOrgMemberResponse](#updateorgmemberresponse)

Changes the organisation manager



    PUT: /orgs/me/members/{user_id}


### RemoveOrgMember

> **rpc** RemoveOrgMember([RemoveOrgMemberRequest](#removeorgmemberrequest))
[RemoveOrgMemberResponse](#removeorgmemberresponse)

Removes an organisation manager



    DELETE: /orgs/me/members/{user_id}


### GetProjectByID

> **rpc** GetProjectByID([GetProjectByIDRequest](#getprojectbyidrequest))
[GetProjectByIDResponse](#getprojectbyidresponse)

Returns a project from my organisation (no granted projects)



    GET: /projects/{id}


### GetGrantedProjectByID

> **rpc** GetGrantedProjectByID([GetGrantedProjectByIDRequest](#getgrantedprojectbyidrequest))
[GetGrantedProjectByIDResponse](#getgrantedprojectbyidresponse)

returns a project my organisation got granted from another organisation



    GET: /granted_projects/{project_id}/grants/{grant_id}


### ListProjects

> **rpc** ListProjects([ListProjectsRequest](#listprojectsrequest))
[ListProjectsResponse](#listprojectsresponse)

Returns all projects my organisation is the owner (no granted projects)
Limit should always be set, there is a default limit set by the service



    POST: /projects/_search


### ListGrantedProjects

> **rpc** ListGrantedProjects([ListGrantedProjectsRequest](#listgrantedprojectsrequest))
[ListGrantedProjectsResponse](#listgrantedprojectsresponse)

returns all projects my organisation got granted from another organisation
Limit should always be set, there is a default limit set by the service



    POST: /granted_projects/_search


### ListGrantedProjectRoles

> **rpc** ListGrantedProjectRoles([ListGrantedProjectRolesRequest](#listgrantedprojectrolesrequest))
[ListGrantedProjectRolesResponse](#listgrantedprojectrolesresponse)

returns all roles of a project grant
Limit should always be set, there is a default limit set by the service



    GET: /granted_projects/{project_id}/grants/{grant_id}/roles/_search


### ListProjectChanges

> **rpc** ListProjectChanges([ListProjectChangesRequest](#listprojectchangesrequest))
[ListProjectChangesResponse](#listprojectchangesresponse)

Returns the history of the project (each event)
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/changes/_search


### AddProject

> **rpc** AddProject([AddProjectRequest](#addprojectrequest))
[AddProjectResponse](#addprojectresponse)

Adds an new project to the organisation



    POST: /projects


### UpdateProject

> **rpc** UpdateProject([UpdateProjectRequest](#updateprojectrequest))
[UpdateProjectResponse](#updateprojectresponse)

Changes a project



    PUT: /projects/{id}


### DeactivateProject

> **rpc** DeactivateProject([DeactivateProjectRequest](#deactivateprojectrequest))
[DeactivateProjectResponse](#deactivateprojectresponse)

Sets the state of a project to deactivated
Returns an error if project is already deactivated



    POST: /projects/{id}/_deactivate


### ReactivateProject

> **rpc** ReactivateProject([ReactivateProjectRequest](#reactivateprojectrequest))
[ReactivateProjectResponse](#reactivateprojectresponse)

Sets the state of a project to active
Returns an error if project is not deactivated



    POST: /projects/{id}/_reactivate


### RemoveProject

> **rpc** RemoveProject([RemoveProjectRequest](#removeprojectrequest))
[RemoveProjectResponse](#removeprojectresponse)

Removes a project
All project grants, applications and user grants for this project will be removed



    DELETE: /projects/{id}


### ListProjectRoles

> **rpc** ListProjectRoles([ListProjectRolesRequest](#listprojectrolesrequest))
[ListProjectRolesResponse](#listprojectrolesresponse)

Returns all roles of a project matching the search query
If no limit is requested, default limit will be set, if the limit is higher then the default an error will be returned



    POST: /projects/{project_id}/roles/_search


### AddProjectRole

> **rpc** AddProjectRole([AddProjectRoleRequest](#addprojectrolerequest))
[AddProjectRoleResponse](#addprojectroleresponse)

Adds a role to a project, key must be unique in the project



    POST: /projects/{project_id}/roles


### BulkAddProjectRoles

> **rpc** BulkAddProjectRoles([BulkAddProjectRolesRequest](#bulkaddprojectrolesrequest))
[BulkAddProjectRolesResponse](#bulkaddprojectrolesresponse)

add a list of project roles in one request



    POST: /projects/{project_id}/roles/_bulk


### UpdateProjectRole

> **rpc** UpdateProjectRole([UpdateProjectRoleRequest](#updateprojectrolerequest))
[UpdateProjectRoleResponse](#updateprojectroleresponse)

Changes a project role, key is not editable
If a key should change, remove the role and create a new



    PUT: /projects/{project_id}/roles/{role_key}


### RemoveProjectRole

> **rpc** RemoveProjectRole([RemoveProjectRoleRequest](#removeprojectrolerequest))
[RemoveProjectRoleResponse](#removeprojectroleresponse)

Removes role from UserGrants, ProjectGrants and from Project



    DELETE: /projects/{project_id}/roles/{role_key}


### ListProjectMemberRoles

> **rpc** ListProjectMemberRoles([ListProjectMemberRolesRequest](#listprojectmemberrolesrequest))
[ListProjectMemberRolesResponse](#listprojectmemberrolesresponse)

Returns all ZITADEL roles which are for project managers



    POST: /projects/members/roles/_search


### ListProjectMembers

> **rpc** ListProjectMembers([ListProjectMembersRequest](#listprojectmembersrequest))
[ListProjectMembersResponse](#listprojectmembersresponse)

Returns all ZITADEL managers of a projects
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/members/_search


### AddProjectMember

> **rpc** AddProjectMember([AddProjectMemberRequest](#addprojectmemberrequest))
[AddProjectMemberResponse](#addprojectmemberresponse)

Adds a new project manager, which is allowed to administrate in ZITADEL



    POST: /projects/{project_id}/members


### UpdateProjectMember

> **rpc** UpdateProjectMember([UpdateProjectMemberRequest](#updateprojectmemberrequest))
[UpdateProjectMemberResponse](#updateprojectmemberresponse)

Change project manager, which is allowed to administrate in ZITADEL



    PUT: /projects/{project_id}/members/{user_id}


### RemoveProjectMember

> **rpc** RemoveProjectMember([RemoveProjectMemberRequest](#removeprojectmemberrequest))
[RemoveProjectMemberResponse](#removeprojectmemberresponse)

Remove project manager, which is allowed to administrate in ZITADEL



    DELETE: /projects/{project_id}/members/{user_id}


### GetAppByID

> **rpc** GetAppByID([GetAppByIDRequest](#getappbyidrequest))
[GetAppByIDResponse](#getappbyidresponse)

Returns an application (oidc or api)



    GET: /projects/{project_id}/apps/{app_id}


### ListApps

> **rpc** ListApps([ListAppsRequest](#listappsrequest))
[ListAppsResponse](#listappsresponse)

Returns all applications of a project matching the query
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/apps/_search


### ListAppChanges

> **rpc** ListAppChanges([ListAppChangesRequest](#listappchangesrequest))
[ListAppChangesResponse](#listappchangesresponse)

Returns the history of the application (each event)
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/apps/{app_id}/changes/_search


### AddOIDCApp

> **rpc** AddOIDCApp([AddOIDCAppRequest](#addoidcapprequest))
[AddOIDCAppResponse](#addoidcappresponse)

Adds a new oidc client
Returns a client id
Returns a new generated secret if needed (Depending on the configuration)



    POST: /projects/{project_id}/apps/oidc


### AddAPIApp

> **rpc** AddAPIApp([AddAPIAppRequest](#addapiapprequest))
[AddAPIAppResponse](#addapiappresponse)

Adds a new api application
Returns a client id
Returns a new generated secret if needed (Depending on the configuration)



    POST: /projects/{project_id}/apps/api


### UpdateApp

> **rpc** UpdateApp([UpdateAppRequest](#updateapprequest))
[UpdateAppResponse](#updateappresponse)

Changes application



    PUT: /projects/{project_id}/apps/{app_id}


### UpdateOIDCAppConfig

> **rpc** UpdateOIDCAppConfig([UpdateOIDCAppConfigRequest](#updateoidcappconfigrequest))
[UpdateOIDCAppConfigResponse](#updateoidcappconfigresponse)

Changes the configuration of the oidc client



    PUT: /projects/{project_id}/apps/{app_id}/oidc_config


### UpdateAPIAppConfig

> **rpc** UpdateAPIAppConfig([UpdateAPIAppConfigRequest](#updateapiappconfigrequest))
[UpdateAPIAppConfigResponse](#updateapiappconfigresponse)

Changes the configuration of the api application



    PUT: /projects/{project_id}/apps/{app_id}/api_config


### DeactivateApp

> **rpc** DeactivateApp([DeactivateAppRequest](#deactivateapprequest))
[DeactivateAppResponse](#deactivateappresponse)

Set the state to deactivated
Its not possible to request tokens for deactivated apps
Returns an error if already deactivated



    POST: /projects/{project_id}/apps/{app_id}/_deactivate


### ReactivateApp

> **rpc** ReactivateApp([ReactivateAppRequest](#reactivateapprequest))
[ReactivateAppResponse](#reactivateappresponse)

Set the state to active
Returns an error if not deactivated



    POST: /projects/{project_id}/apps/{app_id}/_reactivate


### RemoveApp

> **rpc** RemoveApp([RemoveAppRequest](#removeapprequest))
[RemoveAppResponse](#removeappresponse)

Removed the application



    DELETE: /projects/{project_id}/apps/{app_id}


### RegenerateOIDCClientSecret

> **rpc** RegenerateOIDCClientSecret([RegenerateOIDCClientSecretRequest](#regenerateoidcclientsecretrequest))
[RegenerateOIDCClientSecretResponse](#regenerateoidcclientsecretresponse)

Generates a new client secret for the oidc client, make sure to save the response



    POST: /projects/{project_id}/apps/{app_id}/oidc_config/_generate_client_secret


### RegenerateAPIClientSecret

> **rpc** RegenerateAPIClientSecret([RegenerateAPIClientSecretRequest](#regenerateapiclientsecretrequest))
[RegenerateAPIClientSecretResponse](#regenerateapiclientsecretresponse)

Generates a new client secret for the api application, make sure to save the response



    POST: /projects/{project_id}/apps/{app_id}/api_config/_generate_client_secret


### GetAppKey

> **rpc** GetAppKey([GetAppKeyRequest](#getappkeyrequest))
[GetAppKeyResponse](#getappkeyresponse)

Returns an application key



    GET: /projects/{project_id}/apps/{app_id}/keys/{key_id}


### ListAppKeys

> **rpc** ListAppKeys([ListAppKeysRequest](#listappkeysrequest))
[ListAppKeysResponse](#listappkeysresponse)

Returns all application keys matching the result
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/apps/{app_id}/keys/_search


### AddAppKey

> **rpc** AddAppKey([AddAppKeyRequest](#addappkeyrequest))
[AddAppKeyResponse](#addappkeyresponse)

Creates a new app key
Will return key details in result, make sure to save it



    POST: /projects/{project_id}/apps/{app_id}/keys


### RemoveAppKey

> **rpc** RemoveAppKey([RemoveAppKeyRequest](#removeappkeyrequest))
[RemoveAppKeyResponse](#removeappkeyresponse)

Removes an app key



    DELETE: /projects/{project_id}/apps/{app_id}/keys/{key_id}


### GetProjectGrantByID

> **rpc** GetProjectGrantByID([GetProjectGrantByIDRequest](#getprojectgrantbyidrequest))
[GetProjectGrantByIDResponse](#getprojectgrantbyidresponse)

Returns a project grant (ProjectGrant = Grant another organisation for my project)



    GET: /projects/{project_id}/grants/{grant_id}


### ListProjectGrants

> **rpc** ListProjectGrants([ListProjectGrantsRequest](#listprojectgrantsrequest))
[ListProjectGrantsResponse](#listprojectgrantsresponse)

Returns all project grants matching the query, (ProjectGrant = Grant another organisation for my project)
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/grants/_search


### AddProjectGrant

> **rpc** AddProjectGrant([AddProjectGrantRequest](#addprojectgrantrequest))
[AddProjectGrantResponse](#addprojectgrantresponse)

Add a new project grant (ProjectGrant = Grant another organisation for my project)
Project Grant will be listed in granted project of the other organisation



    POST: /projects/{project_id}/grants


### UpdateProjectGrant

> **rpc** UpdateProjectGrant([UpdateProjectGrantRequest](#updateprojectgrantrequest))
[UpdateProjectGrantResponse](#updateprojectgrantresponse)

Change project grant (ProjectGrant = Grant another organisation for my project)
Project Grant will be listed in granted project of the other organisation



    PUT: /projects/{project_id}/grants/{grant_id}


### DeactivateProjectGrant

> **rpc** DeactivateProjectGrant([DeactivateProjectGrantRequest](#deactivateprojectgrantrequest))
[DeactivateProjectGrantResponse](#deactivateprojectgrantresponse)

Set state of project grant to deactivated (ProjectGrant = Grant another organisation for my project)
Returns error if project not active



    POST: /projects/{project_id}/grants/{grant_id}/_deactivate


### ReactivateProjectGrant

> **rpc** ReactivateProjectGrant([ReactivateProjectGrantRequest](#reactivateprojectgrantrequest))
[ReactivateProjectGrantResponse](#reactivateprojectgrantresponse)

Set state of project grant to active (ProjectGrant = Grant another organisation for my project)
Returns error if project not deactivated



    POST: /projects/{project_id}/grants/{grant_id}/_reactivate


### RemoveProjectGrant

> **rpc** RemoveProjectGrant([RemoveProjectGrantRequest](#removeprojectgrantrequest))
[RemoveProjectGrantResponse](#removeprojectgrantresponse)

Removes project grant and all user grants for this project grant



    DELETE: /projects/{project_id}/grants/{grant_id}


### ListProjectGrantMemberRoles

> **rpc** ListProjectGrantMemberRoles([ListProjectGrantMemberRolesRequest](#listprojectgrantmemberrolesrequest))
[ListProjectGrantMemberRolesResponse](#listprojectgrantmemberrolesresponse)

Returns all ZITADEL roles which are for project grant managers



    POST: /projects/grants/members/roles/_search


### ListProjectGrantMembers

> **rpc** ListProjectGrantMembers([ListProjectGrantMembersRequest](#listprojectgrantmembersrequest))
[ListProjectGrantMembersResponse](#listprojectgrantmembersresponse)

Returns all ZITADEL managers of this project grant
Limit should always be set, there is a default limit set by the service



    POST: /projects/{project_id}/grants/{grant_id}/members/_search


### AddProjectGrantMember

> **rpc** AddProjectGrantMember([AddProjectGrantMemberRequest](#addprojectgrantmemberrequest))
[AddProjectGrantMemberResponse](#addprojectgrantmemberresponse)

Adds a new project grant manager, which is allowed to administrate in ZITADEL



    POST: /projects/{project_id}/grants/{grant_id}/members


### UpdateProjectGrantMember

> **rpc** UpdateProjectGrantMember([UpdateProjectGrantMemberRequest](#updateprojectgrantmemberrequest))
[UpdateProjectGrantMemberResponse](#updateprojectgrantmemberresponse)

Changes project grant manager, which is allowed to administrate in ZITADEL



    PUT: /projects/{project_id}/grants/{grant_id}/members/{user_id}


### RemoveProjectGrantMember

> **rpc** RemoveProjectGrantMember([RemoveProjectGrantMemberRequest](#removeprojectgrantmemberrequest))
[RemoveProjectGrantMemberResponse](#removeprojectgrantmemberresponse)

Removed project grant manager



    DELETE: /projects/{project_id}/grants/{grant_id}/members/{user_id}


### GetUserGrantByID

> **rpc** GetUserGrantByID([GetUserGrantByIDRequest](#getusergrantbyidrequest))
[GetUserGrantByIDResponse](#getusergrantbyidresponse)

Returns a user grant (authorization of a user for a project)



    GET: /users/{user_id}/grants/{grant_id}


### ListUserGrants

> **rpc** ListUserGrants([ListUserGrantRequest](#listusergrantrequest))
[ListUserGrantResponse](#listusergrantresponse)

Returns al user grant matching the query (authorizations of user for projects)
Limit should always be set, there is a default limit set by the service



    POST: /users/grants/_search


### AddUserGrant

> **rpc** AddUserGrant([AddUserGrantRequest](#addusergrantrequest))
[AddUserGrantResponse](#addusergrantresponse)

Creates a new user grant (authorization of a user for a project with specified roles)



    POST: /users/{user_id}/grants


### UpdateUserGrant

> **rpc** UpdateUserGrant([UpdateUserGrantRequest](#updateusergrantrequest))
[UpdateUserGrantResponse](#updateusergrantresponse)

Changes a user grant (authorization of a user for a project with specified roles)



    PUT: /users/{user_id}/grants/{grant_id}


### DeactivateUserGrant

> **rpc** DeactivateUserGrant([DeactivateUserGrantRequest](#deactivateusergrantrequest))
[DeactivateUserGrantResponse](#deactivateusergrantresponse)

Sets the state of a user grant to deactivated
User will not be able to use the granted project anymore
Returns an error if user grant is already deactivated



    POST: /users/{user_id}/grants/{grant_id}/_deactivate


### ReactivateUserGrant

> **rpc** ReactivateUserGrant([ReactivateUserGrantRequest](#reactivateusergrantrequest))
[ReactivateUserGrantResponse](#reactivateusergrantresponse)

Sets the state of a user grant to active
Returns an error if user grant is not deactivated



    POST: /users/{user_id}/grants/{grant_id}/_reactivate


### RemoveUserGrant

> **rpc** RemoveUserGrant([RemoveUserGrantRequest](#removeusergrantrequest))
[RemoveUserGrantResponse](#removeusergrantresponse)

Removes a user grant



    DELETE: /users/{user_id}/grants/{grant_id}


### BulkRemoveUserGrant

> **rpc** BulkRemoveUserGrant([BulkRemoveUserGrantRequest](#bulkremoveusergrantrequest))
[BulkRemoveUserGrantResponse](#bulkremoveusergrantresponse)

remove a list of user grants in one request



    DELETE: /user_grants/_bulk


### GetFeatures

> **rpc** GetFeatures([GetFeaturesRequest](#getfeaturesrequest))
[GetFeaturesResponse](#getfeaturesresponse)





    GET: /features


### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)

Returns the org iam policy (this policy is managed by the iam administrator)



    GET: /policies/orgiam


### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)

Returns the login policy of the organisation
With this policy the login gui can be configured



    GET: /policies/login


### GetDefaultLoginPolicy

> **rpc** GetDefaultLoginPolicy([GetDefaultLoginPolicyRequest](#getdefaultloginpolicyrequest))
[GetDefaultLoginPolicyResponse](#getdefaultloginpolicyresponse)

Returns the default login policy configured in the IAM



    GET: /policies/default/login


### AddCustomLoginPolicy

> **rpc** AddCustomLoginPolicy([AddCustomLoginPolicyRequest](#addcustomloginpolicyrequest))
[AddCustomLoginPolicyResponse](#addcustomloginpolicyresponse)

Add a custom login policy for the organisation
With this policy the login gui can be configured



    POST: /policies/login


### UpdateCustomLoginPolicy

> **rpc** UpdateCustomLoginPolicy([UpdateCustomLoginPolicyRequest](#updatecustomloginpolicyrequest))
[UpdateCustomLoginPolicyResponse](#updatecustomloginpolicyresponse)

Change the custom login policy for the organisation
With this policy the login gui can be configured



    PUT: /policies/login


### ResetLoginPolicyToDefault

> **rpc** ResetLoginPolicyToDefault([ResetLoginPolicyToDefaultRequest](#resetloginpolicytodefaultrequest))
[ResetLoginPolicyToDefaultResponse](#resetloginpolicytodefaultresponse)

Removes the custom login policy of the organisation
The default policy of the IAM will trigger after



    DELETE: /policies/login


### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)

Lists all possible identity providers configured on the organisation
Limit should always be set, there is a default limit set by the service



    POST: /policies/login/idps/_search


### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)

Add a (preconfigured) identity provider to the custom login policy



    POST: /policies/login/idps


### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)

Remove a identity provider from the custom login policy



    DELETE: /policies/login/idps/{idp_id}


### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)

Returns all configured second factors of the custom login policy



    POST: /policies/login/second_factors/_search


### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)

Adds a new second factor to the custom login policy



    POST: /policies/login/second_factors


### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)

Remove a second factor from the custom login policy



    DELETE: /policies/login/second_factors/{type}


### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)

Returns all configured multi factors of the custom login policy



    POST: /policies/login/auth_factors/_search


### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)

Adds a new multi factor to the custom login policy



    POST: /policies/login/multi_factors


### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)

Remove a multi factor from the custom login policy



    DELETE: /policies/login/multi_factors/{type}


### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)

Returns the password complexity policy of the organisation
With this policy the password strength can be configured



    GET: /policies/password/complexity


### GetDefaultPasswordComplexityPolicy

> **rpc** GetDefaultPasswordComplexityPolicy([GetDefaultPasswordComplexityPolicyRequest](#getdefaultpasswordcomplexitypolicyrequest))
[GetDefaultPasswordComplexityPolicyResponse](#getdefaultpasswordcomplexitypolicyresponse)

Returns the default password complexity policy of the IAM
With this policy the password strength can be configured



    GET: /policies/default/password/complexity


### AddCustomPasswordComplexityPolicy

> **rpc** AddCustomPasswordComplexityPolicy([AddCustomPasswordComplexityPolicyRequest](#addcustompasswordcomplexitypolicyrequest))
[AddCustomPasswordComplexityPolicyResponse](#addcustompasswordcomplexitypolicyresponse)

Add a custom password complexity policy for the organisation
With this policy the password strength can be configured



    POST: /policies/password/complexity


### UpdateCustomPasswordComplexityPolicy

> **rpc** UpdateCustomPasswordComplexityPolicy([UpdateCustomPasswordComplexityPolicyRequest](#updatecustompasswordcomplexitypolicyrequest))
[UpdateCustomPasswordComplexityPolicyResponse](#updatecustompasswordcomplexitypolicyresponse)

Update the custom password complexity policy for the organisation
With this policy the password strength can be configured



    PUT: /policies/password/complexity


### ResetPasswordComplexityPolicyToDefault

> **rpc** ResetPasswordComplexityPolicyToDefault([ResetPasswordComplexityPolicyToDefaultRequest](#resetpasswordcomplexitypolicytodefaultrequest))
[ResetPasswordComplexityPolicyToDefaultResponse](#resetpasswordcomplexitypolicytodefaultresponse)

Removes the custom password complexity policy of the organisation
The default policy of the IAM will trigger after



    DELETE: /policies/password/complexity


### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)

The password age policy is not used at the moment



    GET: /policies/password/age


### GetDefaultPasswordAgePolicy

> **rpc** GetDefaultPasswordAgePolicy([GetDefaultPasswordAgePolicyRequest](#getdefaultpasswordagepolicyrequest))
[GetDefaultPasswordAgePolicyResponse](#getdefaultpasswordagepolicyresponse)

The password age policy is not used at the moment



    GET: /policies/default/password/age


### AddCustomPasswordAgePolicy

> **rpc** AddCustomPasswordAgePolicy([AddCustomPasswordAgePolicyRequest](#addcustompasswordagepolicyrequest))
[AddCustomPasswordAgePolicyResponse](#addcustompasswordagepolicyresponse)

The password age policy is not used at the moment



    POST: /policies/password/age


### UpdateCustomPasswordAgePolicy

> **rpc** UpdateCustomPasswordAgePolicy([UpdateCustomPasswordAgePolicyRequest](#updatecustompasswordagepolicyrequest))
[UpdateCustomPasswordAgePolicyResponse](#updatecustompasswordagepolicyresponse)

The password age policy is not used at the moment



    PUT: /policies/password/age


### ResetPasswordAgePolicyToDefault

> **rpc** ResetPasswordAgePolicyToDefault([ResetPasswordAgePolicyToDefaultRequest](#resetpasswordagepolicytodefaultrequest))
[ResetPasswordAgePolicyToDefaultResponse](#resetpasswordagepolicytodefaultresponse)

The password age policy is not used at the moment



    DELETE: /policies/password/age


### GetPasswordLockoutPolicy

> **rpc** GetPasswordLockoutPolicy([GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest))
[GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)

The password lockout policy is not used at the moment



    GET: /policies/password/lockout


### GetDefaultPasswordLockoutPolicy

> **rpc** GetDefaultPasswordLockoutPolicy([GetDefaultPasswordLockoutPolicyRequest](#getdefaultpasswordlockoutpolicyrequest))
[GetDefaultPasswordLockoutPolicyResponse](#getdefaultpasswordlockoutpolicyresponse)

The password lockout policy is not used at the moment



    GET: /policies/default/password/lockout


### AddCustomPasswordLockoutPolicy

> **rpc** AddCustomPasswordLockoutPolicy([AddCustomPasswordLockoutPolicyRequest](#addcustompasswordlockoutpolicyrequest))
[AddCustomPasswordLockoutPolicyResponse](#addcustompasswordlockoutpolicyresponse)

The password lockout policy is not used at the moment



    POST: /policies/password/lockout


### UpdateCustomPasswordLockoutPolicy

> **rpc** UpdateCustomPasswordLockoutPolicy([UpdateCustomPasswordLockoutPolicyRequest](#updatecustompasswordlockoutpolicyrequest))
[UpdateCustomPasswordLockoutPolicyResponse](#updatecustompasswordlockoutpolicyresponse)

The password lockout policy is not used at the moment



    PUT: /policies/password/lockout


### ResetPasswordLockoutPolicyToDefault

> **rpc** ResetPasswordLockoutPolicyToDefault([ResetPasswordLockoutPolicyToDefaultRequest](#resetpasswordlockoutpolicytodefaultrequest))
[ResetPasswordLockoutPolicyToDefaultResponse](#resetpasswordlockoutpolicytodefaultresponse)

The password lockout policy is not used at the moment



    DELETE: /policies/password/lockout


### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the label policy of the organisation
With this policy the private labeling can be configured (colors, etc.)



    GET: /policies/label


### GetDefaultLabelPolicy

> **rpc** GetDefaultLabelPolicy([GetDefaultLabelPolicyRequest](#getdefaultlabelpolicyrequest))
[GetDefaultLabelPolicyResponse](#getdefaultlabelpolicyresponse)

Returns the default label policy of the IAM
With this policy the private labeling can be configured (colors, etc.)



    GET: /policies/default/label


### AddCustomLabelPolicy

> **rpc** AddCustomLabelPolicy([AddCustomLabelPolicyRequest](#addcustomlabelpolicyrequest))
[AddCustomLabelPolicyResponse](#addcustomlabelpolicyresponse)

Add a custom label policy for the organisation
With this policy the private labeling can be configured (colors, etc.)



    POST: /policies/label


### UpdateCustomLabelPolicy

> **rpc** UpdateCustomLabelPolicy([UpdateCustomLabelPolicyRequest](#updatecustomlabelpolicyrequest))
[UpdateCustomLabelPolicyResponse](#updatecustomlabelpolicyresponse)

Changes the custom label policy for the organisation
With this policy the private labeling can be configured (colors, etc.)



    PUT: /policies/label


### ResetLabelPolicyToDefault

> **rpc** ResetLabelPolicyToDefault([ResetLabelPolicyToDefaultRequest](#resetlabelpolicytodefaultrequest))
[ResetLabelPolicyToDefaultResponse](#resetlabelpolicytodefaultresponse)

Removes the custom label policy of the organisation
The default policy of the IAM will trigger after



    DELETE: /policies/label


### GetOrgIDPByID

> **rpc** GetOrgIDPByID([GetOrgIDPByIDRequest](#getorgidpbyidrequest))
[GetOrgIDPByIDResponse](#getorgidpbyidresponse)

Returns a identity provider configuration of the organisation



    GET: /idps/{id}


### ListOrgIDPs

> **rpc** ListOrgIDPs([ListOrgIDPsRequest](#listorgidpsrequest))
[ListOrgIDPsResponse](#listorgidpsresponse)

Returns all identity provider configuration in the organisation, which match the query
Limit should always be set, there is a default limit set by the service



    POST: /idps/_search


### AddOrgOIDCIDP

> **rpc** AddOrgOIDCIDP([AddOrgOIDCIDPRequest](#addorgoidcidprequest))
[AddOrgOIDCIDPResponse](#addorgoidcidpresponse)

Add a new identity provider configuration in the organisation
Provider must be OIDC compliant



    POST: /idps/oidc


### DeactivateOrgIDP

> **rpc** DeactivateOrgIDP([DeactivateOrgIDPRequest](#deactivateorgidprequest))
[DeactivateOrgIDPResponse](#deactivateorgidpresponse)

Deactivate identity provider configuration
Users will not be able to use this provider for login (e.g Google, Microsoft, AD, etc)
Returns error if already deactivated



    POST: /idps/{idp_id}/_deactivate


### ReactivateOrgIDP

> **rpc** ReactivateOrgIDP([ReactivateOrgIDPRequest](#reactivateorgidprequest))
[ReactivateOrgIDPResponse](#reactivateorgidpresponse)

Activate identity provider configuration
Returns error if not deactivated



    POST: /idps/{idp_id}/_reactivate


### RemoveOrgIDP

> **rpc** RemoveOrgIDP([RemoveOrgIDPRequest](#removeorgidprequest))
[RemoveOrgIDPResponse](#removeorgidpresponse)

Removes identity provider configuration
Will remove all linked providers of this configuration on the users



    DELETE: /idps/{idp_id}


### UpdateOrgIDP

> **rpc** UpdateOrgIDP([UpdateOrgIDPRequest](#updateorgidprequest))
[UpdateOrgIDPResponse](#updateorgidpresponse)

Change identity provider configuration of the organisation



    PUT: /idps/{idp_id}


### UpdateOrgIDPOIDCConfig

> **rpc** UpdateOrgIDPOIDCConfig([UpdateOrgIDPOIDCConfigRequest](#updateorgidpoidcconfigrequest))
[UpdateOrgIDPOIDCConfigResponse](#updateorgidpoidcconfigresponse)

Change OIDC identity provider configuration of the organisation



    PUT: /idps/{idp_id}/oidc_config







## Messages


### AddAPIAppRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| name |  string | - |
| auth_method_type |  zitadel.app.v1.APIAuthMethodType | - |



### AddAPIAppResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| app_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |
| client_id |  string | - |
| client_secret |  string | - |



### AddAppKeyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |
| type |  zitadel.authn.v1.KeyType | - |
| expiration_date |  google.protobuf.Timestamp | - |



### AddAppKeyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |
| key_details |  bytes | - |



### AddCustomLabelPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| primary_color |  string | - |
| secondary_color |  string | - |
| hide_login_name_suffix |  bool | - |



### AddCustomLabelPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddCustomLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| allow_username_password |  bool | - |
| allow_register |  bool | - |
| allow_external_idp |  bool | - |
| force_mfa |  bool | - |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - |



### AddCustomLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddCustomPasswordAgePolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| max_age_days |  uint32 | - |
| expire_warn_days |  uint32 | - |



### AddCustomPasswordAgePolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddCustomPasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| min_length |  uint64 | - |
| has_uppercase |  bool | - |
| has_lowercase |  bool | - |
| has_number |  bool | - |
| has_symbol |  bool | - |



### AddCustomPasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddCustomPasswordLockoutPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| max_attempts |  uint32 | - |
| show_lockout_failure |  bool | - |



### AddCustomPasswordLockoutPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddHumanUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name |  string | - |
| profile |  AddHumanUserRequest.Profile | - |
| email |  AddHumanUserRequest.Email | - |
| phone |  AddHumanUserRequest.Phone | - |
| initial_password |  string | - |



### AddHumanUserRequest.Email


| Field | Type | Description |
| ----- | ---- | ----------- |
| email |  string | TODO: check if no value is allowed |
| is_email_verified |  bool | - |



### AddHumanUserRequest.Phone


| Field | Type | Description |
| ----- | ---- | ----------- |
| phone |  string | has to be a global number |
| is_phone_verified |  bool | - |



### AddHumanUserRequest.Profile


| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name |  string | - |
| last_name |  string | - |
| nick_name |  string | - |
| display_name |  string | - |
| preferred_language |  string | - |
| gender |  zitadel.user.v1.Gender | - |



### AddHumanUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddIDPToLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |
| ownerType |  zitadel.idp.v1.IDPOwnerType | - |



### AddIDPToLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddMachineKeyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| type |  zitadel.authn.v1.KeyType | - |
| expiration_date |  google.protobuf.Timestamp | - |



### AddMachineKeyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| key_id |  string | - |
| key_details |  bytes | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddMachineUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name |  string | - |
| name |  string | - |
| description |  string | - |



### AddMachineUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddMultiFactorToLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - |



### AddMultiFactorToLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddOIDCAppRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| name |  string | - |
| redirect_uris | repeated string | - |
| response_types | repeated zitadel.app.v1.OIDCResponseType | - |
| grant_types | repeated zitadel.app.v1.OIDCGrantType | - |
| app_type |  zitadel.app.v1.OIDCAppType | - |
| auth_method_type |  zitadel.app.v1.OIDCAuthMethodType | - |
| post_logout_redirect_uris | repeated string | - |
| version |  zitadel.app.v1.OIDCVersion | - |
| dev_mode |  bool | - |
| access_token_type |  zitadel.app.v1.OIDCTokenType | - |
| access_token_role_assertion |  bool | - |
| id_token_role_assertion |  bool | - |
| id_token_userinfo_assertion |  bool | - |
| clock_skew |  google.protobuf.Duration | - |



### AddOIDCAppResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| app_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |
| client_id |  string | - |
| client_secret |  string | - |
| none_compliant |  bool | - |
| compliance_problems | repeated zitadel.v1.LocalizedMessage | - |



### AddOrgDomainRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| domain |  string | - |



### AddOrgDomainResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddOrgMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| roles | repeated string | - |



### AddOrgMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddOrgOIDCIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| name |  string | - |
| styling_type |  zitadel.idp.v1.IDPStylingType | - |
| client_id |  string | - |
| client_secret |  string | - |
| issuer |  string | - |
| scopes | repeated string | - |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - |



### AddOrgOIDCIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| idp_id |  string | - |



### AddOrgRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| name |  string | - |



### AddOrgResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddProjectGrantMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |
| user_id |  string | - |
| roles | repeated string | - |



### AddProjectGrantMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddProjectGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| granted_org_id |  string | - |
| role_keys | repeated string | - |



### AddProjectGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| grant_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddProjectMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| user_id |  string | - |
| roles | repeated string | - |



### AddProjectMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddProjectRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| name |  string | - |
| project_role_assertion |  bool | - |
| project_role_check |  bool | - |



### AddProjectResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddProjectRoleRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| role_key |  string | - |
| display_name |  string | - |
| group |  string | - |



### AddProjectRoleResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddSecondFactorToLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - |



### AddSecondFactorToLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| project_id |  string | - |
| project_grant_id |  string | - |
| role_keys | repeated string | - |



### AddUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_grant_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### BulkAddProjectRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| roles | repeated BulkAddProjectRolesRequest.Role | - |



### BulkAddProjectRolesRequest.Role


| Field | Type | Description |
| ----- | ---- | ----------- |
| key |  string | - |
| display_name |  string | - |
| group |  string | - |



### BulkAddProjectRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### BulkRemoveUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| grant_id | repeated string | - |



### BulkRemoveUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### DeactivateAppRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |



### DeactivateAppResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### DeactivateOrgIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### DeactivateOrgIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### DeactivateOrgRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### DeactivateOrgResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### DeactivateProjectGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |



### DeactivateProjectGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### DeactivateProjectRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### DeactivateProjectResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### DeactivateUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| grant_id |  string | - |



### DeactivateUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### DeactivateUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### DeactivateUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### GenerateOrgDomainValidationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| domain |  string | - |
| type |  zitadel.org.v1.DomainValidationType | - |



### GenerateOrgDomainValidationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| token |  string | - |
| url |  string | - |



### GetAppByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |



### GetAppByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| app |  zitadel.app.v1.App | - |



### GetAppKeyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |
| key_id |  string | - |



### GetAppKeyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| key |  zitadel.authn.v1.Key | - |



### GetDefaultLabelPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetDefaultLabelPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |



### GetDefaultLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetDefaultLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |



### GetDefaultPasswordAgePolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetDefaultPasswordAgePolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |



### GetDefaultPasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetDefaultPasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |



### GetDefaultPasswordLockoutPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetDefaultPasswordLockoutPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordLockoutPolicy | - |



### GetFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| features |  zitadel.features.v1.Features | - |



### GetGrantedProjectByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |



### GetGrantedProjectByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| granted_project |  zitadel.project.v1.GrantedProject | - |



### GetHumanEmailRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### GetHumanEmailResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| email |  zitadel.user.v1.Email | - |



### GetHumanPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### GetHumanPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| phone |  zitadel.user.v1.Phone | - |



### GetHumanProfileRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### GetHumanProfileResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| profile |  zitadel.user.v1.Profile | - |



### GetIAMRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetIAMResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| global_org_id |  string | - |
| iam_project_id |  string | - |



### GetLabelPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetLabelPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |
| is_default |  bool | - |



### GetLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |
| is_default |  bool | - |



### GetMachineKeyByIDsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| key_id |  string | - |



### GetMachineKeyByIDsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| key |  zitadel.authn.v1.Key | - |



### GetMyOrgRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetMyOrgResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| org |  zitadel.org.v1.Org | - |



### GetOIDCInformationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetOIDCInformationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| issuer |  string | - |
| discovery_endpoint |  string | - |



### GetOrgByDomainGlobalRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| domain |  string | - |



### GetOrgByDomainGlobalResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| org |  zitadel.org.v1.Org | - |



### GetOrgIAMPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetOrgIAMPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |



### GetOrgIDPByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### GetOrgIDPByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp |  zitadel.idp.v1.IDP | - |



### GetPasswordAgePolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetPasswordAgePolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |
| is_default |  bool | - |



### GetPasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetPasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |
| is_default |  bool | - |



### GetPasswordLockoutPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetPasswordLockoutPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordLockoutPolicy | - |
| is_default |  bool | - |



### GetProjectByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### GetProjectByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| project |  zitadel.project.v1.Project | - |



### GetProjectGrantByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |



### GetProjectGrantByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_grant |  zitadel.project.v1.GrantedProject | - |



### GetUserByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### GetUserByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user |  zitadel.user.v1.User | - |



### GetUserByLoginNameGlobalRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| login_name |  string | - |



### GetUserByLoginNameGlobalResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user |  zitadel.user.v1.User | - |



### GetUserGrantByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| grant_id |  string | - |



### GetUserGrantByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_grant |  zitadel.user.v1.UserGrant | - |



### HealthzRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### HealthzResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### IDPQuery


| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query |  zitadel.idp.v1.IDPIDQuery | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query |  zitadel.idp.v1.IDPNameQuery | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.owner_type_query |  zitadel.idp.v1.IDPOwnerTypeQuery | - |



### ImportHumanUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name |  string | - |
| profile |  ImportHumanUserRequest.Profile | - |
| email |  ImportHumanUserRequest.Email | - |
| phone |  ImportHumanUserRequest.Phone | - |
| password |  string | - |
| password_change_required |  bool | - |



### ImportHumanUserRequest.Email


| Field | Type | Description |
| ----- | ---- | ----------- |
| email |  string | TODO: check if no value is allowed |
| is_email_verified |  bool | - |



### ImportHumanUserRequest.Phone


| Field | Type | Description |
| ----- | ---- | ----------- |
| phone |  string | has to be a global number |
| is_phone_verified |  bool | - |



### ImportHumanUserRequest.Profile


| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name |  string | - |
| last_name |  string | - |
| nick_name |  string | - |
| display_name |  string | - |
| preferred_language |  string | - |
| gender |  zitadel.user.v1.Gender | - |



### ImportHumanUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### IsUserUniqueRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name |  string | - |
| email |  string | - |



### IsUserUniqueResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| is_unique |  bool | - |



### ListAppChangesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |
| project_id |  string | - |
| app_id |  string | - |



### ListAppChangesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.change.v1.Change | - |



### ListAppKeysRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| app_id |  string | - |
| project_id |  string | - |



### ListAppKeysResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.authn.v1.Key | - |



### ListAppsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.app.v1.AppQuery | criterias the client is looking for |



### ListAppsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.app.v1.App | - |



### ListGrantedProjectRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.project.v1.RoleQuery | criterias the client is looking for |



### ListGrantedProjectRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.project.v1.Role | - |



### ListGrantedProjectsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.project.v1.ProjectQuery | criterias the client is looking for |



### ListGrantedProjectsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.project.v1.GrantedProject | - |



### ListHumanAuthFactorsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### ListHumanAuthFactorsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated zitadel.user.v1.AuthFactor | - |



### ListHumanLinkedIDPsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |



### ListHumanLinkedIDPsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.idp.v1.IDPUserLink | - |



### ListHumanPasswordlessRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### ListHumanPasswordlessResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated zitadel.user.v1.WebAuthNToken | - |



### ListLoginPolicyIDPsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | - |



### ListLoginPolicyIDPsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.idp.v1.IDPLoginPolicyLink | - |



### ListLoginPolicyMultiFactorsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListLoginPolicyMultiFactorsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.policy.v1.MultiFactorType | - |



### ListLoginPolicySecondFactorsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListLoginPolicySecondFactorsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.policy.v1.SecondFactorType | - |



### ListMachineKeysRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |



### ListMachineKeysResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.authn.v1.Key | - |



### ListOrgChangesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |



### ListOrgChangesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.change.v1.Change | - |



### ListOrgDomainsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.org.v1.DomainSearchQuery | criterias the client is looking for |



### ListOrgDomainsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.org.v1.Domain | - |



### ListOrgIDPsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| sorting_column |  zitadel.idp.v1.IDPFieldName | the field the result is sorted |
| queries | repeated IDPQuery | criterias the client is looking for |



### ListOrgIDPsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| sorting_column |  zitadel.idp.v1.IDPFieldName | - |
| result | repeated zitadel.idp.v1.IDP | - |



### ListOrgMemberRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListOrgMemberRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated string | - |



### ListOrgMembersRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |



### ListOrgMembersResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | list limitations and ordering |
| result | repeated zitadel.member.v1.Member | criterias the client is looking for |



### ListProjectChangesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |
| project_id |  string | - |



### ListProjectChangesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.change.v1.Change | - |



### ListProjectGrantMemberRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | - |
| result | repeated string | - |



### ListProjectGrantMemberRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated string | - |



### ListProjectGrantMembersRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |



### ListProjectGrantMembersResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.member.v1.Member | - |



### ListProjectGrantsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.project.v1.ProjectGrantQuery | criterias the client is looking for |



### ListProjectGrantsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.project.v1.GrantedProject | - |



### ListProjectMemberRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListProjectMemberRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated string | - |



### ListProjectMembersRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |



### ListProjectMembersResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.member.v1.Member | - |



### ListProjectRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.project.v1.RoleQuery | criterias the client is looking for |



### ListProjectRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.project.v1.Role | - |



### ListProjectsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.project.v1.ProjectQuery | criterias the client is looking for |



### ListProjectsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.project.v1.Project | - |



### ListUserChangesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |
| user_id |  string | - |



### ListUserChangesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.change.v1.Change | - |



### ListUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.user.v1.UserGrantQuery | criterias the client is looking for |



### ListUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.user.v1.UserGrant | - |



### ListUserMembershipsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | list limitations and ordering |
| query |  zitadel.v1.ListQuery | the field the result is sorted |
| queries | repeated zitadel.user.v1.MembershipQuery | criterias the client is looking for |



### ListUserMembershipsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.user.v1.Membership | - |



### ListUsersRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| sorting_column |  zitadel.user.v1.UserFieldName | the field the result is sorted |
| queries | repeated zitadel.user.v1.SearchQuery | criterias the client is looking for |



### ListUsersResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| sorting_column |  zitadel.user.v1.UserFieldName | - |
| result | repeated zitadel.user.v1.User | - |



### LockUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### LockUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateAppRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |



### ReactivateAppResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateOrgIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### ReactivateOrgIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateOrgRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ReactivateOrgResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateProjectGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |



### ReactivateProjectGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateProjectRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### ReactivateProjectResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| grant_id |  string | - |



### ReactivateUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ReactivateUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### ReactivateUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RegenerateAPIClientSecretRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |



### RegenerateAPIClientSecretResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| client_secret |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### RegenerateOIDCClientSecretRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |



### RegenerateOIDCClientSecretResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| client_secret |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveAppKeyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |
| key_id |  string | - |



### RemoveAppKeyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveAppRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |



### RemoveAppResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveHumanAuthFactorOTPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### RemoveHumanAuthFactorOTPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveHumanAuthFactorU2FRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| token_id |  string | - |



### RemoveHumanAuthFactorU2FResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveHumanLinkedIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| idp_id |  string | - |
| linked_user_id |  string | - |



### RemoveHumanLinkedIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveHumanPasswordlessRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| token_id |  string | - |



### RemoveHumanPasswordlessResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveHumanPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### RemoveHumanPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveIDPFromLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### RemoveIDPFromLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMachineKeyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| key_id |  string | - |



### RemoveMachineKeyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMultiFactorFromLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - |



### RemoveMultiFactorFromLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveOrgDomainRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| domain |  string | - |



### RemoveOrgDomainResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveOrgIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### RemoveOrgIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### RemoveOrgMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### RemoveOrgMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveProjectGrantMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |
| user_id |  string | - |



### RemoveProjectGrantMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveProjectGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |



### RemoveProjectGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveProjectMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| user_id |  string | - |



### RemoveProjectMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveProjectRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### RemoveProjectResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveProjectRoleRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| role_key |  string | - |



### RemoveProjectRoleResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveSecondFactorFromLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - |



### RemoveSecondFactorFromLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| grant_id |  string | - |



### RemoveUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### RemoveUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResendHumanEmailVerificationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### ResendHumanEmailVerificationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResendHumanInitializationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| email |  string | - |



### ResendHumanInitializationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResendHumanPhoneVerificationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### ResendHumanPhoneVerificationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetLabelPolicyToDefaultRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResetLabelPolicyToDefaultResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetLoginPolicyToDefaultRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResetLoginPolicyToDefaultResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetPasswordAgePolicyToDefaultRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResetPasswordAgePolicyToDefaultResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetPasswordComplexityPolicyToDefaultRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResetPasswordComplexityPolicyToDefaultResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetPasswordLockoutPolicyToDefaultRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResetPasswordLockoutPolicyToDefaultResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SendHumanResetPasswordNotificationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| type |  SendHumanResetPasswordNotificationRequest.Type | - |



### SendHumanResetPasswordNotificationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetHumanInitialPasswordRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| password |  string | - |



### SetHumanInitialPasswordResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetPrimaryOrgDomainRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| domain |  string | - |



### SetPrimaryOrgDomainResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UnlockUserRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### UnlockUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateAPIAppConfigRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |
| auth_method_type |  zitadel.app.v1.APIAuthMethodType | - |



### UpdateAPIAppConfigResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateAppRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |
| name |  string | - |



### UpdateAppResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateCustomLabelPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| primary_color |  string | - |
| secondary_color |  string | - |
| hide_login_name_suffix |  bool | - |



### UpdateCustomLabelPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateCustomLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| allow_username_password |  bool | - |
| allow_register |  bool | - |
| allow_external_idp |  bool | - |
| force_mfa |  bool | - |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - |



### UpdateCustomLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateCustomPasswordAgePolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| max_age_days |  uint32 | - |
| expire_warn_days |  uint32 | - |



### UpdateCustomPasswordAgePolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateCustomPasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| min_length |  uint64 | - |
| has_uppercase |  bool | - |
| has_lowercase |  bool | - |
| has_number |  bool | - |
| has_symbol |  bool | - |



### UpdateCustomPasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateCustomPasswordLockoutPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| max_attempts |  uint32 | - |
| show_lockout_failure |  bool | - |



### UpdateCustomPasswordLockoutPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateHumanEmailRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| email |  string | - |
| is_email_verified |  bool | - |



### UpdateHumanEmailResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateHumanPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| phone |  string | - |
| is_phone_verified |  bool | - |



### UpdateHumanPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateHumanProfileRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| first_name |  string | - |
| last_name |  string | - |
| nick_name |  string | - |
| display_name |  string | - |
| preferred_language |  string | - |
| gender |  zitadel.user.v1.Gender | - |



### UpdateHumanProfileResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateMachineRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| description |  string | - |
| name |  string | - |



### UpdateMachineResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateOIDCAppConfigRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| app_id |  string | - |
| redirect_uris | repeated string | - |
| response_types | repeated zitadel.app.v1.OIDCResponseType | - |
| grant_types | repeated zitadel.app.v1.OIDCGrantType | - |
| app_type |  zitadel.app.v1.OIDCAppType | - |
| auth_method_type |  zitadel.app.v1.OIDCAuthMethodType | - |
| post_logout_redirect_uris | repeated string | - |
| dev_mode |  bool | - |
| access_token_type |  zitadel.app.v1.OIDCTokenType | - |
| access_token_role_assertion |  bool | - |
| id_token_role_assertion |  bool | - |
| id_token_userinfo_assertion |  bool | - |
| clock_skew |  google.protobuf.Duration | - |



### UpdateOIDCAppConfigResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateOrgIDPOIDCConfigRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |
| client_id |  string | - |
| client_secret |  string | - |
| issuer |  string | - |
| scopes | repeated string | - |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - |



### UpdateOrgIDPOIDCConfigResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateOrgIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |
| name |  string | - |
| styling_type |  zitadel.idp.v1.IDPStylingType | - |



### UpdateOrgIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateOrgMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| roles | repeated string | - |



### UpdateOrgMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateProjectGrantMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |
| user_id |  string | - |
| roles | repeated string | - |



### UpdateProjectGrantMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateProjectGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| grant_id |  string | - |
| role_keys | repeated string | - |



### UpdateProjectGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateProjectMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| user_id |  string | - |
| roles | repeated string | - |



### UpdateProjectMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateProjectRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |
| name |  string | - |
| project_role_assertion |  bool | - |
| project_role_check |  bool | - |



### UpdateProjectResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateProjectRoleRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id |  string | - |
| role_key |  string | - |
| display_name |  string | - |
| group |  string | - |



### UpdateProjectRoleResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateUserGrantRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| grant_id |  string | - |
| role_keys | repeated string | - |



### UpdateUserGrantResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateUserNameRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| user_name |  string | - |



### UpdateUserNameResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ValidateOrgDomainRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| domain |  string | - |



### ValidateOrgDomainResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |






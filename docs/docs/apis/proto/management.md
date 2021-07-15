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


### GetSupportedLanguages

> **rpc** GetSupportedLanguages([GetSupportedLanguagesRequest](#getsupportedlanguagesrequest))
[GetSupportedLanguagesResponse](#getsupportedlanguagesresponse)

Returns the default languages



    GET: /languages


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


### SetUserMetaData

> **rpc** SetUserMetaData([SetUserMetaDataRequest](#setusermetadatarequest))
[SetUserMetaDataResponse](#setusermetadataresponse)

Sets a user meta data by key



    POST: /users/{id}/metadata/{key}


### BulkSetUserMetaData

> **rpc** BulkSetUserMetaData([BulkSetUserMetaDataRequest](#bulksetusermetadatarequest))
[BulkSetUserMetaDataResponse](#bulksetusermetadataresponse)

Set a list of user meta data



    POST: /users/{id}/metadata/_bulk


### ListUserMetaData

> **rpc** ListUserMetaData([ListUserMetaDataRequest](#listusermetadatarequest))
[ListUserMetaDataResponse](#listusermetadataresponse)

Returns the user meta data



    POST: /users/{id}/metadata/_search


### GetUserMetaData

> **rpc** GetUserMetaData([GetUserMetaDataRequest](#getusermetadatarequest))
[GetUserMetaDataResponse](#getusermetadataresponse)

Returns the user meta data by key



    GET: /users/{id}/metadata/{key}


### RemoveUserMetaData

> **rpc** RemoveUserMetaData([RemoveUserMetaDataRequest](#removeusermetadatarequest))
[RemoveUserMetaDataResponse](#removeusermetadataresponse)

Removes a user meta data by key



    DELETE: /users/{id}/metadata/{key}


### BulkRemoveUserMetaData

> **rpc** BulkRemoveUserMetaData([BulkRemoveUserMetaDataRequest](#bulkremoveusermetadatarequest))
[BulkRemoveUserMetaDataResponse](#bulkremoveusermetadataresponse)

Set a list of user meta data



    DELETE: /users/{id}/metadata/_bulk


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


### RemoveMyAvatar

> **rpc** RemoveMyAvatar([RemoveHumanAvatarRequest](#removehumanavatarrequest))
[RemoveHumanAvatarResponse](#removehumanavatarresponse)

Removes the avatar number of the human



    DELETE: /users/{user_id}/avatar


### SetHumanInitialPassword

> **rpc** SetHumanInitialPassword([SetHumanInitialPasswordRequest](#sethumaninitialpasswordrequest))
[SetHumanInitialPasswordResponse](#sethumaninitialpasswordresponse)

deprecated: use SetHumanPassword



    POST: /users/{user_id}/password/_initialize


### SetHumanPassword

> **rpc** SetHumanPassword([SetHumanPasswordRequest](#sethumanpasswordrequest))
[SetHumanPasswordResponse](#sethumanpasswordresponse)

Set a new password for a user, on default the user has to change the password on the next login
Set no_change_required to true if the user does not have to change the password on the next login



    POST: /users/{user_id}/password


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


### UpdateOrg

> **rpc** UpdateOrg([UpdateOrgRequest](#updateorgrequest))
[UpdateOrgResponse](#updateorgresponse)

Changes my organisation



    PUT: /orgs/me


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


### GetPrivacyPolicy

> **rpc** GetPrivacyPolicy([GetPrivacyPolicyRequest](#getprivacypolicyrequest))
[GetPrivacyPolicyResponse](#getprivacypolicyresponse)

Returns the privacy policy of the organisation
With this policy privacy relevant things can be configured (e.g. tos link)



    GET: /policies/privacy


### GetDefaultPrivacyPolicy

> **rpc** GetDefaultPrivacyPolicy([GetDefaultPrivacyPolicyRequest](#getdefaultprivacypolicyrequest))
[GetDefaultPrivacyPolicyResponse](#getdefaultprivacypolicyresponse)

Returns the default privacy policy of the IAM
With this policy the privacy relevant things can be configured (e.g tos link)



    GET: /policies/default/privacy


### AddCustomPrivacyPolicy

> **rpc** AddCustomPrivacyPolicy([AddCustomPrivacyPolicyRequest](#addcustomprivacypolicyrequest))
[AddCustomPrivacyPolicyResponse](#addcustomprivacypolicyresponse)

Add a custom privacy policy for the organisation
With this policy privacy relevant things can be configured (e.g. tos link)



    POST: /policies/privacy


### UpdateCustomPrivacyPolicy

> **rpc** UpdateCustomPrivacyPolicy([UpdateCustomPrivacyPolicyRequest](#updatecustomprivacypolicyrequest))
[UpdateCustomPrivacyPolicyResponse](#updatecustomprivacypolicyresponse)

Update the privacy complexity policy for the organisation
With this policy privacy relevant things can be configured (e.g. tos link)



    PUT: /policies/privacy


### ResetPrivacyPolicyToDefault

> **rpc** ResetPrivacyPolicyToDefault([ResetPrivacyPolicyToDefaultRequest](#resetprivacypolicytodefaultrequest))
[ResetPrivacyPolicyToDefaultResponse](#resetprivacypolicytodefaultresponse)

Removes the privacy policy of the organisation
The default policy of the IAM will trigger after



    DELETE: /policies/privacy


### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the active label policy of the organisation
With this policy the private labeling can be configured (colors, etc.)



    GET: /policies/label


### GetPreviewLabelPolicy

> **rpc** GetPreviewLabelPolicy([GetPreviewLabelPolicyRequest](#getpreviewlabelpolicyrequest))
[GetPreviewLabelPolicyResponse](#getpreviewlabelpolicyresponse)

Returns the preview label policy of the organisation
With this policy the private labeling can be configured (colors, etc.)



    GET: /policies/label/_preview


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


### ActivateCustomLabelPolicy

> **rpc** ActivateCustomLabelPolicy([ActivateCustomLabelPolicyRequest](#activatecustomlabelpolicyrequest))
[ActivateCustomLabelPolicyResponse](#activatecustomlabelpolicyresponse)

Activates all changes of the label policy



    POST: /policies/label/_activate


### RemoveCustomLabelPolicyLogo

> **rpc** RemoveCustomLabelPolicyLogo([RemoveCustomLabelPolicyLogoRequest](#removecustomlabelpolicylogorequest))
[RemoveCustomLabelPolicyLogoResponse](#removecustomlabelpolicylogoresponse)

Removes the logo of the label policy



    DELETE: /policies/label/logo


### RemoveCustomLabelPolicyLogoDark

> **rpc** RemoveCustomLabelPolicyLogoDark([RemoveCustomLabelPolicyLogoDarkRequest](#removecustomlabelpolicylogodarkrequest))
[RemoveCustomLabelPolicyLogoDarkResponse](#removecustomlabelpolicylogodarkresponse)

Removes the logo dark of the label policy



    DELETE: /policies/label/logo_dark


### RemoveCustomLabelPolicyIcon

> **rpc** RemoveCustomLabelPolicyIcon([RemoveCustomLabelPolicyIconRequest](#removecustomlabelpolicyiconrequest))
[RemoveCustomLabelPolicyIconResponse](#removecustomlabelpolicyiconresponse)

Removes the icon of the label policy



    DELETE: /policies/label/icon


### RemoveCustomLabelPolicyIconDark

> **rpc** RemoveCustomLabelPolicyIconDark([RemoveCustomLabelPolicyIconDarkRequest](#removecustomlabelpolicyicondarkrequest))
[RemoveCustomLabelPolicyIconDarkResponse](#removecustomlabelpolicyicondarkresponse)

Removes the logo dark of the label policy



    DELETE: /policies/label/icon_dark


### RemoveCustomLabelPolicyFont

> **rpc** RemoveCustomLabelPolicyFont([RemoveCustomLabelPolicyFontRequest](#removecustomlabelpolicyfontrequest))
[RemoveCustomLabelPolicyFontResponse](#removecustomlabelpolicyfontresponse)

Removes the font of the label policy



    DELETE: /policies/label/font


### ResetLabelPolicyToDefault

> **rpc** ResetLabelPolicyToDefault([ResetLabelPolicyToDefaultRequest](#resetlabelpolicytodefaultrequest))
[ResetLabelPolicyToDefaultResponse](#resetlabelpolicytodefaultresponse)

Removes the custom label policy of the organisation
The default policy of the IAM will trigger after



    DELETE: /policies/label


### GetCustomInitMessageText

> **rpc** GetCustomInitMessageText([GetCustomInitMessageTextRequest](#getcustominitmessagetextrequest))
[GetCustomInitMessageTextResponse](#getcustominitmessagetextresponse)

Returns the custom text for initial message



    GET: /text/message/init/{language}


### GetDefaultInitMessageText

> **rpc** GetDefaultInitMessageText([GetDefaultInitMessageTextRequest](#getdefaultinitmessagetextrequest))
[GetDefaultInitMessageTextResponse](#getdefaultinitmessagetextresponse)

Returns the default text for initial message



    GET: /text/default/message/init/{language}


### SetCustomInitMessageText

> **rpc** SetCustomInitMessageText([SetCustomInitMessageTextRequest](#setcustominitmessagetextrequest))
[SetCustomInitMessageTextResponse](#setcustominitmessagetextresponse)

Sets the default custom text for initial message
it impacts all organisations without customized initial message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}



    PUT: /text/message/init/{language}


### ResetCustomInitMessageTextToDefault

> **rpc** ResetCustomInitMessageTextToDefault([ResetCustomInitMessageTextToDefaultRequest](#resetcustominitmessagetexttodefaultrequest))
[ResetCustomInitMessageTextToDefaultResponse](#resetcustominitmessagetexttodefaultresponse)

Removes the custom init message text of the organisation
The default text of the IAM will trigger after



    DELETE: /text/message/init/{language}


### GetCustomPasswordResetMessageText

> **rpc** GetCustomPasswordResetMessageText([GetCustomPasswordResetMessageTextRequest](#getcustompasswordresetmessagetextrequest))
[GetCustomPasswordResetMessageTextResponse](#getcustompasswordresetmessagetextresponse)

Returns the custom text for password reset message



    GET: /text/message/passwordreset/{language}


### GetDefaultPasswordResetMessageText

> **rpc** GetDefaultPasswordResetMessageText([GetDefaultPasswordResetMessageTextRequest](#getdefaultpasswordresetmessagetextrequest))
[GetDefaultPasswordResetMessageTextResponse](#getdefaultpasswordresetmessagetextresponse)

Returns the default text for password reset message



    GET: /text/default/message/passwordreset/{language}


### SetCustomPasswordResetMessageText

> **rpc** SetCustomPasswordResetMessageText([SetCustomPasswordResetMessageTextRequest](#setcustompasswordresetmessagetextrequest))
[SetCustomPasswordResetMessageTextResponse](#setcustompasswordresetmessagetextresponse)

Sets the default custom text for password reset message
it impacts all organisations without customized password reset message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}



    PUT: /text/message/passwordreset/{language}


### ResetCustomPasswordResetMessageTextToDefault

> **rpc** ResetCustomPasswordResetMessageTextToDefault([ResetCustomPasswordResetMessageTextToDefaultRequest](#resetcustompasswordresetmessagetexttodefaultrequest))
[ResetCustomPasswordResetMessageTextToDefaultResponse](#resetcustompasswordresetmessagetexttodefaultresponse)

Removes the custom password reset message text of the organisation
The default text of the IAM will trigger after



    DELETE: /text/message/verifyemail/{language}


### GetCustomVerifyEmailMessageText

> **rpc** GetCustomVerifyEmailMessageText([GetCustomVerifyEmailMessageTextRequest](#getcustomverifyemailmessagetextrequest))
[GetCustomVerifyEmailMessageTextResponse](#getcustomverifyemailmessagetextresponse)

Returns the custom text for verify email message



    GET: /text/message/verifyemail/{language}


### GetDefaultVerifyEmailMessageText

> **rpc** GetDefaultVerifyEmailMessageText([GetDefaultVerifyEmailMessageTextRequest](#getdefaultverifyemailmessagetextrequest))
[GetDefaultVerifyEmailMessageTextResponse](#getdefaultverifyemailmessagetextresponse)

Returns the default text for verify email message



    GET: /text/default/message/verifyemail/{language}


### SetCustomVerifyEmailMessageText

> **rpc** SetCustomVerifyEmailMessageText([SetCustomVerifyEmailMessageTextRequest](#setcustomverifyemailmessagetextrequest))
[SetCustomVerifyEmailMessageTextResponse](#setcustomverifyemailmessagetextresponse)

Sets the default custom text for verify email message
it impacts all organisations without customized verify email message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}



    PUT: /text/message/verifyemail/{language}


### ResetCustomVerifyEmailMessageTextToDefault

> **rpc** ResetCustomVerifyEmailMessageTextToDefault([ResetCustomVerifyEmailMessageTextToDefaultRequest](#resetcustomverifyemailmessagetexttodefaultrequest))
[ResetCustomVerifyEmailMessageTextToDefaultResponse](#resetcustomverifyemailmessagetexttodefaultresponse)

Removes the custom verify email message text of the organisation
The default text of the IAM will trigger after



    DELETE: /text/message/verifyemail/{language}


### GetCustomVerifyPhoneMessageText

> **rpc** GetCustomVerifyPhoneMessageText([GetCustomVerifyPhoneMessageTextRequest](#getcustomverifyphonemessagetextrequest))
[GetCustomVerifyPhoneMessageTextResponse](#getcustomverifyphonemessagetextresponse)

Returns the custom text for verify email message



    GET: /text/message/verifyphone/{language}


### GetDefaultVerifyPhoneMessageText

> **rpc** GetDefaultVerifyPhoneMessageText([GetDefaultVerifyPhoneMessageTextRequest](#getdefaultverifyphonemessagetextrequest))
[GetDefaultVerifyPhoneMessageTextResponse](#getdefaultverifyphonemessagetextresponse)

Returns the custom text for verify email message



    GET: /text/default/message/verifyphone/{language}


### SetCustomVerifyPhoneMessageText

> **rpc** SetCustomVerifyPhoneMessageText([SetCustomVerifyPhoneMessageTextRequest](#setcustomverifyphonemessagetextrequest))
[SetCustomVerifyPhoneMessageTextResponse](#setcustomverifyphonemessagetextresponse)

Sets the default custom text for verify email message
it impacts all organisations without customized verify email message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}



    PUT: /text/message/verifyphone/{language}


### ResetCustomVerifyPhoneMessageTextToDefault

> **rpc** ResetCustomVerifyPhoneMessageTextToDefault([ResetCustomVerifyPhoneMessageTextToDefaultRequest](#resetcustomverifyphonemessagetexttodefaultrequest))
[ResetCustomVerifyPhoneMessageTextToDefaultResponse](#resetcustomverifyphonemessagetexttodefaultresponse)

Removes the custom verify phone text of the organisation
The default text of the IAM will trigger after



    DELETE: /text/message/verifyphone/{language}


### GetCustomDomainClaimedMessageText

> **rpc** GetCustomDomainClaimedMessageText([GetCustomDomainClaimedMessageTextRequest](#getcustomdomainclaimedmessagetextrequest))
[GetCustomDomainClaimedMessageTextResponse](#getcustomdomainclaimedmessagetextresponse)

Returns the custom text for domain claimed message



    GET: /text/message/domainclaimed/{language}


### GetDefaultDomainClaimedMessageText

> **rpc** GetDefaultDomainClaimedMessageText([GetDefaultDomainClaimedMessageTextRequest](#getdefaultdomainclaimedmessagetextrequest))
[GetDefaultDomainClaimedMessageTextResponse](#getdefaultdomainclaimedmessagetextresponse)

Returns the custom text for domain claimed message



    GET: /text/default/message/domainclaimed/{language}


### SetCustomDomainClaimedMessageCustomText

> **rpc** SetCustomDomainClaimedMessageCustomText([SetCustomDomainClaimedMessageTextRequest](#setcustomdomainclaimedmessagetextrequest))
[SetCustomDomainClaimedMessageTextResponse](#setcustomdomainclaimedmessagetextresponse)

Sets the default custom text for domain claimed message
it impacts all organisations without customized domain claimed message text
The Following Variables can be used:
{{.Domain}} {{.TempUsername}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}



    PUT: /text/message/domainclaimed/{language}


### ResetCustomDomainClaimedMessageTextToDefault

> **rpc** ResetCustomDomainClaimedMessageTextToDefault([ResetCustomDomainClaimedMessageTextToDefaultRequest](#resetcustomdomainclaimedmessagetexttodefaultrequest))
[ResetCustomDomainClaimedMessageTextToDefaultResponse](#resetcustomdomainclaimedmessagetexttodefaultresponse)

Removes the custom init message text of the organisation
The default text of the IAM will trigger after



    DELETE: /text/message/domainclaimed/{language}


### GetCustomLoginTexts

> **rpc** GetCustomLoginTexts([GetCustomLoginTextsRequest](#getcustomlogintextsrequest))
[GetCustomLoginTextsResponse](#getcustomlogintextsresponse)

Returns the custom texts for login ui



    GET: /text/login/{language}


### GetDefaultLoginTexts

> **rpc** GetDefaultLoginTexts([GetDefaultLoginTextsRequest](#getdefaultlogintextsrequest))
[GetDefaultLoginTextsResponse](#getdefaultlogintextsresponse)

Returns the custom texts for login ui



    GET: /text/default/login/{language}


### SetCustomLoginText

> **rpc** SetCustomLoginText([SetCustomLoginTextsRequest](#setcustomlogintextsrequest))
[SetCustomLoginTextsResponse](#setcustomlogintextsresponse)

Sets the default custom text for login ui
it impacts all organisations without customized login ui texts



    PUT: /text/login/{language}


### ResetCustomLoginTextToDefault

> **rpc** ResetCustomLoginTextToDefault([ResetCustomLoginTextsToDefaultRequest](#resetcustomlogintextstodefaultrequest))
[ResetCustomLoginTextsToDefaultResponse](#resetcustomlogintextstodefaultresponse)

Removes the custom login text of the organisation
The default text of the IAM will trigger after



    DELETE: /text/login/{language}


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


### ActivateCustomLabelPolicyRequest
This is an empty request




### ActivateCustomLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddAPIAppRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| auth_method_type |  zitadel.app.v1.APIAuthMethodType | - | enum.defined_only: true<br />  |




### AddAPIAppResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| app_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| client_id |  string | - |  |
| client_secret |  string | - |  |




### AddAppKeyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| type |  zitadel.authn.v1.KeyType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |
| expiration_date |  google.protobuf.Timestamp | - |  |




### AddAppKeyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| key_details |  bytes | - |  |




### AddCustomLabelPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| primary_color |  string | - | string.max_len: 50<br />  |
| hide_login_name_suffix |  bool | hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set. Details about this scope in https://docs.zitadel.ch/concepts#Reserved_Scopes |  |
| warn_color |  string | - | string.max_len: 50<br />  |
| background_color |  string | - | string.max_len: 50<br />  |
| font_color |  string | - | string.max_len: 50<br />  |
| primary_color_dark |  string | - | string.max_len: 50<br />  |
| background_color_dark |  string | - | string.max_len: 50<br />  |
| warn_color_dark |  string | - | string.max_len: 50<br />  |
| font_color_dark |  string | - | string.max_len: 50<br />  |
| disable_watermark |  bool | - |  |




### AddCustomLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| allow_username_password |  bool | - |  |
| allow_register |  bool | - |  |
| allow_external_idp |  bool | - |  |
| force_mfa |  bool | - |  |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - | enum.defined_only: true<br />  |
| hide_password_reset |  bool | - |  |




### AddCustomLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomPasswordAgePolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_age_days |  uint32 | - |  |
| expire_warn_days |  uint32 | - |  |




### AddCustomPasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomPasswordComplexityPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| min_length |  uint64 | - |  |
| has_uppercase |  bool | - |  |
| has_lowercase |  bool | - |  |
| has_number |  bool | - |  |
| has_symbol |  bool | - |  |




### AddCustomPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomPasswordLockoutPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_attempts |  uint32 | - |  |
| show_lockout_failure |  bool | - |  |




### AddCustomPasswordLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomPrivacyPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| tos_link |  string | - |  |
| privacy_link |  string | - |  |




### AddCustomPrivacyPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddHumanUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| profile |  AddHumanUserRequest.Profile | - | message.required: true<br />  |
| email |  AddHumanUserRequest.Email | - | message.required: true<br />  |
| phone |  AddHumanUserRequest.Phone | - |  |
| initial_password |  string | - |  |




### AddHumanUserRequest.Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | TODO: check if no value is allowed | string.email: true<br />  |
| is_email_verified |  bool | - |  |




### AddHumanUserRequest.Phone



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| phone |  string | has to be a global number | string.min_len: 1<br /> string.max_len: 50<br /> string.prefix: +<br />  |
| is_phone_verified |  bool | - |  |




### AddHumanUserRequest.Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| last_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| nick_name |  string | - | string.max_len: 200<br />  |
| display_name |  string | - | string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |
| gender |  zitadel.user.v1.Gender | - |  |




### AddHumanUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddIDPToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| ownerType |  zitadel.idp.v1.IDPOwnerType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddIDPToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddMachineKeyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br />  |
| type |  zitadel.authn.v1.KeyType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |
| expiration_date |  google.protobuf.Timestamp | - |  |




### AddMachineKeyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key_id |  string | - |  |
| key_details |  bytes | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddMachineUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |




### AddMachineUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddMultiFactorToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddMultiFactorToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddOIDCAppRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| redirect_uris | repeated string | - |  |
| response_types | repeated zitadel.app.v1.OIDCResponseType | - |  |
| grant_types | repeated zitadel.app.v1.OIDCGrantType | - |  |
| app_type |  zitadel.app.v1.OIDCAppType | - | enum.defined_only: true<br />  |
| auth_method_type |  zitadel.app.v1.OIDCAuthMethodType | - | enum.defined_only: true<br />  |
| post_logout_redirect_uris | repeated string | - |  |
| version |  zitadel.app.v1.OIDCVersion | - | enum.defined_only: true<br />  |
| dev_mode |  bool | - |  |
| access_token_type |  zitadel.app.v1.OIDCTokenType | - | enum.defined_only: true<br />  |
| access_token_role_assertion |  bool | - |  |
| id_token_role_assertion |  bool | - |  |
| id_token_userinfo_assertion |  bool | - |  |
| clock_skew |  google.protobuf.Duration | - | duration.lte.seconds: 5<br /> duration.lte.nanos: 0<br /> duration.gte.seconds: 0<br /> duration.gte.nanos: 0<br />  |
| additional_origins | repeated string | - |  |




### AddOIDCAppResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| app_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |
| client_id |  string | - |  |
| client_secret |  string | - |  |
| none_compliant |  bool | - |  |
| compliance_problems | repeated zitadel.v1.LocalizedMessage | - |  |




### AddOrgDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### AddOrgDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddOrgMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | - |  |




### AddOrgMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddOrgOIDCIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| scopes | repeated string | - |  |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |




### AddOrgOIDCIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| idp_id |  string | - |  |




### AddOrgRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### AddOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddProjectGrantMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | - |  |




### AddProjectGrantMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddProjectGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| granted_org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| role_keys | repeated string | - |  |




### AddProjectGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| grant_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddProjectMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | - |  |




### AddProjectMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddProjectRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| project_role_assertion |  bool | - |  |
| project_role_check |  bool | - |  |




### AddProjectResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddProjectRoleRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| role_key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| display_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| group |  string | - | string.max_len: 200<br />  |




### AddProjectRoleResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddSecondFactorToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddSecondFactorToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| project_grant_id |  string | - | string.max_len: 200<br />  |
| role_keys | repeated string | - |  |




### AddUserGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_grant_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### BulkAddProjectRolesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated BulkAddProjectRolesRequest.Role | - |  |




### BulkAddProjectRolesRequest.Role



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| display_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| group |  string | - | string.max_len: 200<br />  |




### BulkAddProjectRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### BulkRemoveUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| grant_id | repeated string | - |  |




### BulkRemoveUserGrantResponse





### BulkRemoveUserMetaDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| keys | repeated string | - |  |




### BulkRemoveUserMetaDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### BulkSetUserMetaDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| meta_data | repeated BulkSetUserMetaDataRequest.MetaData | - |  |




### BulkSetUserMetaDataRequest.MetaData



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| value |  string | - | string.min_len: 1<br />  |




### BulkSetUserMetaDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateAppRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateAppResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateOrgIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateOrgIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateOrgRequest
This is an empty request




### DeactivateOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateProjectGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateProjectGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateProjectRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateProjectResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateUserGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### GenerateOrgDomainValidationRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| type |  zitadel.org.v1.DomainValidationType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### GenerateOrgDomainValidationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| token |  string | - |  |
| url |  string | - |  |




### GetAppByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetAppByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| app |  zitadel.app.v1.App | - |  |




### GetAppKeyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetAppKeyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  zitadel.authn.v1.Key | - |  |




### GetCustomDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomLoginTextsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomLoginTextsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.LoginCustomText | - |  |




### GetCustomPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultLabelPolicyRequest
This is an empty request




### GetDefaultLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |  |




### GetDefaultLoginPolicyRequest





### GetDefaultLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |  |




### GetDefaultLoginTextsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultLoginTextsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.LoginCustomText | - |  |




### GetDefaultPasswordAgePolicyRequest
This is an empty request




### GetDefaultPasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |  |




### GetDefaultPasswordComplexityPolicyRequest
This is an empty request




### GetDefaultPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |  |




### GetDefaultPasswordLockoutPolicyRequest
This is an empty request




### GetDefaultPasswordLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordLockoutPolicy | - |  |




### GetDefaultPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultPrivacyPolicyRequest
This is an empty request




### GetDefaultPrivacyPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PrivacyPolicy | - |  |




### GetDefaultVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetFeaturesRequest





### GetFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| features |  zitadel.features.v1.Features | - |  |




### GetGrantedProjectByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetGrantedProjectByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| granted_project |  zitadel.project.v1.GrantedProject | - |  |




### GetHumanEmailRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetHumanEmailResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| email |  zitadel.user.v1.Email | - |  |




### GetHumanPhoneRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetHumanPhoneResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| phone |  zitadel.user.v1.Phone | - |  |




### GetHumanProfileRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetHumanProfileResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| profile |  zitadel.user.v1.Profile | - |  |




### GetIAMRequest
This is an empty request




### GetIAMResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| global_org_id |  string | - |  |
| iam_project_id |  string | - |  |




### GetLabelPolicyRequest
This is an empty request




### GetLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |  |
| is_default |  bool | - |  |




### GetLoginPolicyRequest





### GetLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |  |
| is_default |  bool | - |  |




### GetMachineKeyByIDsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetMachineKeyByIDsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  zitadel.authn.v1.Key | - |  |




### GetMyOrgRequest
This is an empty request




### GetMyOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  zitadel.org.v1.Org | - |  |




### GetOIDCInformationRequest
This is an empty request




### GetOIDCInformationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| issuer |  string | - |  |
| discovery_endpoint |  string | - |  |




### GetOrgByDomainGlobalRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetOrgByDomainGlobalResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  zitadel.org.v1.Org | - |  |




### GetOrgIAMPolicyRequest





### GetOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |  |




### GetOrgIDPByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetOrgIDPByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp |  zitadel.idp.v1.IDP | - |  |




### GetPasswordAgePolicyRequest
This is an empty request




### GetPasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |  |
| is_default |  bool | - |  |




### GetPasswordComplexityPolicyRequest





### GetPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |  |
| is_default |  bool | - |  |




### GetPasswordLockoutPolicyRequest
This is an empty request




### GetPasswordLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordLockoutPolicy | - |  |
| is_default |  bool | - |  |




### GetPreviewLabelPolicyRequest
This is an empty request




### GetPreviewLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |  |
| is_default |  bool | - |  |




### GetPrivacyPolicyRequest
This is an empty request




### GetPrivacyPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PrivacyPolicy | - |  |




### GetProjectByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetProjectByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project |  zitadel.project.v1.Project | - |  |




### GetProjectGrantByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetProjectGrantByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_grant |  zitadel.project.v1.GrantedProject | - |  |




### GetSupportedLanguagesRequest
This is an empty request




### GetSupportedLanguagesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| languages | repeated string | - |  |




### GetUserByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetUserByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user |  zitadel.user.v1.User | - |  |




### GetUserByLoginNameGlobalRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| login_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetUserByLoginNameGlobalResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user |  zitadel.user.v1.User | - |  |




### GetUserGrantByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetUserGrantByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_grant |  zitadel.user.v1.UserGrant | - |  |




### GetUserMetaDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetUserMetaDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| meta_data |  zitadel.metadata.v1.MetaData | - |  |




### HealthzRequest
This is an empty request




### HealthzResponse
This is an empty response




### IDPQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query |  zitadel.idp.v1.IDPIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query |  zitadel.idp.v1.IDPNameQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.owner_type_query |  zitadel.idp.v1.IDPOwnerTypeQuery | - |  |




### ImportHumanUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| profile |  ImportHumanUserRequest.Profile | - | message.required: true<br />  |
| email |  ImportHumanUserRequest.Email | - | message.required: true<br />  |
| phone |  ImportHumanUserRequest.Phone | - |  |
| password |  string | - |  |
| password_change_required |  bool | - |  |




### ImportHumanUserRequest.Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | TODO: check if no value is allowed | string.email: true<br />  |
| is_email_verified |  bool | - |  |




### ImportHumanUserRequest.Phone



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| phone |  string | has to be a global number | string.min_len: 1<br /> string.max_len: 50<br /> string.prefix: +<br />  |
| is_phone_verified |  bool | - |  |




### ImportHumanUserRequest.Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| last_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| nick_name |  string | - | string.max_len: 200<br />  |
| display_name |  string | - | string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |
| gender |  zitadel.user.v1.Gender | - |  |




### ImportHumanUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### IsUserUniqueRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.max_len: 200<br />  |
| email |  string | - | string.max_len: 200<br />  |




### IsUserUniqueResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| is_unique |  bool | - |  |




### ListAppChangesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |  |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ListAppChangesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.change.v1.Change | - |  |




### ListAppKeysRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ListAppKeysResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.authn.v1.Key | - |  |




### ListAppsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.app.v1.AppQuery | criterias the client is looking for |  |




### ListAppsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.app.v1.App | - |  |




### ListGrantedProjectRolesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.project.v1.RoleQuery | criterias the client is looking for |  |




### ListGrantedProjectRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.project.v1.Role | - |  |




### ListGrantedProjectsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.project.v1.ProjectQuery | criterias the client is looking for |  |




### ListGrantedProjectsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.project.v1.GrantedProject | - |  |




### ListHumanAuthFactorsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ListHumanAuthFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.AuthFactor | - |  |




### ListHumanLinkedIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |




### ListHumanLinkedIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.IDPUserLink | - |  |




### ListHumanPasswordlessRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ListHumanPasswordlessResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.WebAuthNToken | - |  |




### ListLoginPolicyIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | - |  |




### ListLoginPolicyIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.IDPLoginPolicyLink | - |  |




### ListLoginPolicyMultiFactorsRequest





### ListLoginPolicyMultiFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.policy.v1.MultiFactorType | - |  |




### ListLoginPolicySecondFactorsRequest





### ListLoginPolicySecondFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.policy.v1.SecondFactorType | - |  |




### ListMachineKeysRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |




### ListMachineKeysResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.authn.v1.Key | - |  |




### ListOrgChangesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |  |




### ListOrgChangesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.change.v1.Change | - |  |




### ListOrgDomainsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.org.v1.DomainSearchQuery | criterias the client is looking for |  |




### ListOrgDomainsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.org.v1.Domain | - |  |




### ListOrgIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.idp.v1.IDPFieldName | the field the result is sorted |  |
| queries | repeated IDPQuery | criterias the client is looking for |  |




### ListOrgIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.idp.v1.IDPFieldName | - |  |
| result | repeated zitadel.idp.v1.IDP | - |  |




### ListOrgMemberRolesRequest
This is an empty request




### ListOrgMemberRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |




### ListOrgMembersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |  |




### ListOrgMembersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | list limitations and ordering |  |
| result | repeated zitadel.member.v1.Member | criterias the client is looking for |  |




### ListProjectChangesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |  |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ListProjectChangesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.change.v1.Change | - |  |




### ListProjectGrantMemberRolesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | - |  |
| result | repeated string | - |  |




### ListProjectGrantMemberRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated string | - |  |




### ListProjectGrantMembersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |  |




### ListProjectGrantMembersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.member.v1.Member | - |  |




### ListProjectGrantsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.project.v1.ProjectGrantQuery | criterias the client is looking for |  |




### ListProjectGrantsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.project.v1.GrantedProject | - |  |




### ListProjectMemberRolesRequest
This is an empty request




### ListProjectMemberRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated string | - |  |




### ListProjectMembersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |  |




### ListProjectMembersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.member.v1.Member | - |  |




### ListProjectRolesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.project.v1.RoleQuery | criterias the client is looking for |  |




### ListProjectRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.project.v1.Role | - |  |




### ListProjectsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.project.v1.ProjectQuery | criterias the client is looking for |  |




### ListProjectsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.project.v1.Project | - |  |




### ListUserChangesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | list limitations and ordering |  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ListUserChangesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.change.v1.Change | - |  |




### ListUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.user.v1.UserGrantQuery | criterias the client is looking for |  |




### ListUserGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.user.v1.UserGrant | - |  |




### ListUserMembershipsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | list limitations and ordering | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | the field the result is sorted |  |
| queries | repeated zitadel.user.v1.MembershipQuery | criterias the client is looking for |  |




### ListUserMembershipsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.user.v1.Membership | - |  |




### ListUserMetaDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| query |  zitadel.v1.ListQuery | - |  |
| queries | repeated zitadel.metadata.v1.MetaDataQuery | - |  |




### ListUserMetaDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.metadata.v1.MetaData | - |  |




### ListUsersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.user.v1.UserFieldName | the field the result is sorted |  |
| queries | repeated zitadel.user.v1.SearchQuery | criterias the client is looking for |  |




### ListUsersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.user.v1.UserFieldName | - |  |
| result | repeated zitadel.user.v1.User | - |  |




### LockUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### LockUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateAppRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateAppResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateOrgIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateOrgIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateOrgRequest
This is an empty request




### ReactivateOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateProjectGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateProjectGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateProjectRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateProjectResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateUserGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ReactivateUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RegenerateAPIClientSecretRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RegenerateAPIClientSecretResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| client_secret |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### RegenerateOIDCClientSecretRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RegenerateOIDCClientSecretResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| client_secret |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveAppKeyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveAppKeyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveAppRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveAppResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveCustomLabelPolicyFontRequest
This is an empty request




### RemoveCustomLabelPolicyFontResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveCustomLabelPolicyIconDarkRequest
This is an empty request




### RemoveCustomLabelPolicyIconDarkResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveCustomLabelPolicyIconRequest
This is an empty request




### RemoveCustomLabelPolicyIconResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveCustomLabelPolicyLogoDarkRequest
This is an empty request




### RemoveCustomLabelPolicyLogoDarkResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveCustomLabelPolicyLogoRequest
This is an empty request




### RemoveCustomLabelPolicyLogoResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveHumanAuthFactorOTPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveHumanAuthFactorOTPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveHumanAuthFactorU2FRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| token_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveHumanAuthFactorU2FResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveHumanAvatarRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveHumanAvatarResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveHumanLinkedIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| linked_user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveHumanLinkedIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveHumanPasswordlessRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| token_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveHumanPasswordlessResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveHumanPhoneRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveHumanPhoneResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveIDPFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIDPFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveMachineKeyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveMachineKeyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveMultiFactorFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### RemoveMultiFactorFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveOrgDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveOrgDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveOrgIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveOrgIDPResponse
This is an empty response




### RemoveOrgMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveOrgMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveProjectGrantMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveProjectGrantMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveProjectGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveProjectGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveProjectMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveProjectMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveProjectRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveProjectResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveProjectRoleRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| role_key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveProjectRoleResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveSecondFactorFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### RemoveSecondFactorFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveUserGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveUserMetaDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveUserMetaDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResendHumanEmailVerificationRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResendHumanEmailVerificationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResendHumanInitializationRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| email |  string | - | string.email: true<br />  |




### ResendHumanInitializationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResendHumanPhoneVerificationRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResendHumanPhoneVerificationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomDomainClaimedMessageTextToDefaultRequest
This is an empty request


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomDomainClaimedMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomInitMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomInitMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomLoginTextsToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomLoginTextsToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomPasswordResetMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomPasswordResetMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomVerifyEmailMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomVerifyEmailMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomVerifyPhoneMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomVerifyPhoneMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetLabelPolicyToDefaultRequest
This is an empty request




### ResetLabelPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetLoginPolicyToDefaultRequest





### ResetLoginPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetPasswordAgePolicyToDefaultRequest
This is an empty request




### ResetPasswordAgePolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetPasswordComplexityPolicyToDefaultRequest
This is an empty request




### ResetPasswordComplexityPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetPasswordLockoutPolicyToDefaultRequest
This is an empty request




### ResetPasswordLockoutPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetPrivacyPolicyToDefaultRequest
This is an empty request




### ResetPrivacyPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SendHumanResetPasswordNotificationRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| type |  SendHumanResetPasswordNotificationRequest.Type | - | enum.defined_only: true<br />  |




### SendHumanResetPasswordNotificationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetCustomDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetCustomInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomLoginTextsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| select_account_text |  zitadel.text.v1.SelectAccountScreenText | - |  |
| login_text |  zitadel.text.v1.LoginScreenText | - |  |
| password_text |  zitadel.text.v1.PasswordScreenText | - |  |
| username_change_text |  zitadel.text.v1.UsernameChangeScreenText | - |  |
| username_change_done_text |  zitadel.text.v1.UsernameChangeDoneScreenText | - |  |
| init_password_text |  zitadel.text.v1.InitPasswordScreenText | - |  |
| init_password_done_text |  zitadel.text.v1.InitPasswordDoneScreenText | - |  |
| email_verification_text |  zitadel.text.v1.EmailVerificationScreenText | - |  |
| email_verification_done_text |  zitadel.text.v1.EmailVerificationDoneScreenText | - |  |
| initialize_user_text |  zitadel.text.v1.InitializeUserScreenText | - |  |
| initialize_done_text |  zitadel.text.v1.InitializeUserDoneScreenText | - |  |
| init_mfa_prompt_text |  zitadel.text.v1.InitMFAPromptScreenText | - |  |
| init_mfa_otp_text |  zitadel.text.v1.InitMFAOTPScreenText | - |  |
| init_mfa_u2f_text |  zitadel.text.v1.InitMFAU2FScreenText | - |  |
| init_mfa_done_text |  zitadel.text.v1.InitMFADoneScreenText | - |  |
| mfa_providers_text |  zitadel.text.v1.MFAProvidersText | - |  |
| verify_mfa_otp_text |  zitadel.text.v1.VerifyMFAOTPScreenText | - |  |
| verify_mfa_u2f_text |  zitadel.text.v1.VerifyMFAU2FScreenText | - |  |
| passwordless_text |  zitadel.text.v1.PasswordlessScreenText | - |  |
| password_change_text |  zitadel.text.v1.PasswordChangeScreenText | - |  |
| password_change_done_text |  zitadel.text.v1.PasswordChangeDoneScreenText | - |  |
| password_reset_done_text |  zitadel.text.v1.PasswordResetDoneScreenText | - |  |
| registration_option_text |  zitadel.text.v1.RegistrationOptionScreenText | - |  |
| registration_user_text |  zitadel.text.v1.RegistrationUserScreenText | - |  |
| registration_org_text |  zitadel.text.v1.RegistrationOrgScreenText | - |  |
| linking_user_done_text |  zitadel.text.v1.LinkingUserDoneScreenText | - |  |
| external_user_not_found_text |  zitadel.text.v1.ExternalUserNotFoundScreenText | - |  |
| success_login_text |  zitadel.text.v1.SuccessLoginScreenText | - |  |
| logout_text |  zitadel.text.v1.LogoutDoneScreenText | - |  |
| footer_text |  zitadel.text.v1.FooterText | - |  |




### SetCustomLoginTextsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetCustomPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetCustomVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetCustomVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetHumanInitialPasswordRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br />  |
| password |  string | - | string.min_len: 1<br /> string.max_len: 72<br />  |




### SetHumanInitialPasswordResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetHumanPasswordRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br />  |
| password |  string | - | string.min_len: 1<br /> string.max_len: 72<br />  |
| no_change_required |  bool | - |  |




### SetHumanPasswordResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetPrimaryOrgDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### SetPrimaryOrgDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetUserMetaDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| value |  string | - | string.min_len: 1<br />  |




### SetUserMetaDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| details |  zitadel.v1.ObjectDetails | - |  |




### UnlockUserRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UnlockUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateAPIAppConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| auth_method_type |  zitadel.app.v1.APIAuthMethodType | - | enum.defined_only: true<br />  |




### UpdateAPIAppConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateAppRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateAppResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomLabelPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| primary_color |  string | - | string.max_len: 50<br />  |
| hide_login_name_suffix |  bool | - |  |
| warn_color |  string | - | string.max_len: 50<br />  |
| background_color |  string | - | string.max_len: 50<br />  |
| font_color |  string | - | string.max_len: 50<br />  |
| primary_color_dark |  string | - | string.max_len: 50<br />  |
| background_color_dark |  string | - | string.max_len: 50<br />  |
| warn_color_dark |  string | - | string.max_len: 50<br />  |
| font_color_dark |  string | - | string.max_len: 50<br />  |
| disable_watermark |  bool | - |  |




### UpdateCustomLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| allow_username_password |  bool | - |  |
| allow_register |  bool | - |  |
| allow_external_idp |  bool | - |  |
| force_mfa |  bool | - |  |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - | enum.defined_only: true<br />  |
| hide_password_reset |  bool | - |  |




### UpdateCustomLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomPasswordAgePolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_age_days |  uint32 | - |  |
| expire_warn_days |  uint32 | - |  |




### UpdateCustomPasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomPasswordComplexityPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| min_length |  uint64 | - |  |
| has_uppercase |  bool | - |  |
| has_lowercase |  bool | - |  |
| has_number |  bool | - |  |
| has_symbol |  bool | - |  |




### UpdateCustomPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomPasswordLockoutPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_attempts |  uint32 | - |  |
| show_lockout_failure |  bool | - |  |




### UpdateCustomPasswordLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomPrivacyPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| tos_link |  string | - |  |
| privacy_link |  string | - |  |




### UpdateCustomPrivacyPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateHumanEmailRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| email |  string | - | string.email: true<br />  |
| is_email_verified |  bool | - |  |




### UpdateHumanEmailResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateHumanPhoneRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| phone |  string | - | string.min_len: 1<br /> string.max_len: 50<br /> string.prefix: +<br />  |
| is_phone_verified |  bool | - |  |




### UpdateHumanPhoneResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateHumanProfileRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| first_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| last_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| nick_name |  string | - | string.max_len: 200<br />  |
| display_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |
| gender |  zitadel.user.v1.Gender | - |  |




### UpdateHumanProfileResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateMachineRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| description |  string | - | string.max_len: 500<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateMachineResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOIDCAppConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| app_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| redirect_uris | repeated string | - |  |
| response_types | repeated zitadel.app.v1.OIDCResponseType | - |  |
| grant_types | repeated zitadel.app.v1.OIDCGrantType | - |  |
| app_type |  zitadel.app.v1.OIDCAppType | - | enum.defined_only: true<br />  |
| auth_method_type |  zitadel.app.v1.OIDCAuthMethodType | - | enum.defined_only: true<br />  |
| post_logout_redirect_uris | repeated string | - |  |
| dev_mode |  bool | - |  |
| access_token_type |  zitadel.app.v1.OIDCTokenType | - | enum.defined_only: true<br />  |
| access_token_role_assertion |  bool | - |  |
| id_token_role_assertion |  bool | - |  |
| id_token_userinfo_assertion |  bool | - |  |
| clock_skew |  google.protobuf.Duration | - | duration.lte.seconds: 5<br /> duration.lte.nanos: 0<br /> duration.gte.seconds: 0<br /> duration.gte.nanos: 0<br />  |
| additional_origins | repeated string | - |  |




### UpdateOIDCAppConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOrgIDPOIDCConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| scopes | repeated string | - |  |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |




### UpdateOrgIDPOIDCConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOrgIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |




### UpdateOrgIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOrgMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | - |  |




### UpdateOrgMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOrgRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateProjectGrantMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | - |  |




### UpdateProjectGrantMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateProjectGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| role_keys | repeated string | - |  |




### UpdateProjectGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateProjectMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | - |  |




### UpdateProjectMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateProjectRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| project_role_assertion |  bool | - |  |
| project_role_check |  bool | - |  |




### UpdateProjectResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateProjectRoleRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| role_key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| display_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| group |  string | - | string.max_len: 200<br />  |




### UpdateProjectRoleResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateUserGrantRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| grant_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| role_keys | repeated string | - |  |




### UpdateUserGrantResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateUserNameRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateUserNameResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ValidateOrgDomainRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ValidateOrgDomainResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |







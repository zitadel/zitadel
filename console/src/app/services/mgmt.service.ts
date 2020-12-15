import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject } from 'rxjs';

import { MultiFactorsResult } from '../proto/generated/admin_pb';
import {
    AddMachineKeyRequest,
    AddMachineKeyResponse,
    AddOrgDomainRequest,
    AddOrgMemberRequest,
    Application,
    ApplicationID,
    ApplicationSearchQuery,
    ApplicationSearchRequest,
    ApplicationSearchResponse,
    ApplicationUpdate,
    ApplicationView,
    ChangeOrgMemberRequest,
    ChangeRequest,
    Changes,
    CreateHumanRequest,
    CreateMachineRequest,
    CreateUserRequest,
    Domain,
    ExternalIDPRemoveRequest,
    ExternalIDPSearchRequest,
    ExternalIDPSearchResponse,
    Gender,
    GrantedProjectSearchRequest,
    Iam,
    Idp,
    IdpID,
    IdpProviderAdd,
    IdpProviderID,
    IdpProviderSearchRequest,
    IdpProviderSearchResponse,
    IdpProviderType,
    IdpSearchQuery,
    IdpSearchRequest,
    IdpSearchResponse,
    IdpUpdate,
    IdpView,
    InitialMailRequest,
    LoginName,
    LoginPolicy,
    LoginPolicyRequest,
    LoginPolicyView,
    MachineKeyIDRequest,
    MachineKeySearchRequest,
    MachineKeySearchResponse,
    MachineKeyType,
    MachineResponse,
    MultiFactor,
    NotificationType,
    OIDCApplicationCreate,
    OIDCConfig,
    OIDCConfigUpdate,
    OidcIdpConfig,
    OidcIdpConfigCreate,
    OidcIdpConfigUpdate,
    Org,
    OrgCreateRequest,
    OrgDomain,
    OrgDomainSearchQuery,
    OrgDomainSearchRequest,
    OrgDomainSearchResponse,
    OrgDomainValidationRequest,
    OrgDomainValidationResponse,
    OrgDomainValidationType,
    OrgIamPolicyView,
    OrgMember,
    OrgMemberRoles,
    OrgMemberSearchRequest,
    OrgMemberSearchResponse,
    OrgView,
    PasswordAgePolicy,
    PasswordAgePolicyRequest,
    PasswordAgePolicyView,
    PasswordComplexityPolicy,
    PasswordComplexityPolicyRequest,
    PasswordComplexityPolicyView,
    PasswordLockoutPolicy,
    PasswordLockoutPolicyRequest,
    PasswordRequest,
    PrimaryOrgDomainRequest,
    Project,
    ProjectCreateRequest,
    ProjectGrant,
    ProjectGrantCreate,
    ProjectGrantID,
    ProjectGrantMember,
    ProjectGrantMemberAdd,
    ProjectGrantMemberChange,
    ProjectGrantMemberRemove,
    ProjectGrantMemberRoles,
    ProjectGrantMemberSearchQuery,
    ProjectGrantMemberSearchRequest,
    ProjectGrantSearchRequest,
    ProjectGrantSearchResponse,
    ProjectGrantUpdate,
    ProjectGrantView,
    ProjectID,
    ProjectMember,
    ProjectMemberAdd,
    ProjectMemberChange,
    ProjectMemberRemove,
    ProjectMemberRoles,
    ProjectMemberSearchQuery,
    ProjectMemberSearchRequest,
    ProjectMemberSearchResponse,
    ProjectRole,
    ProjectRoleAdd,
    ProjectRoleAddBulk,
    ProjectRoleChange,
    ProjectRoleRemove,
    ProjectRoleSearchQuery,
    ProjectRoleSearchRequest,
    ProjectRoleSearchResponse,
    ProjectSearchQuery,
    ProjectSearchRequest,
    ProjectSearchResponse,
    ProjectUpdateRequest,
    ProjectView,
    RemoveOrgDomainRequest,
    RemoveOrgMemberRequest,
    SecondFactor,
    SecondFactorsResult,
    SetPasswordNotificationRequest,
    UpdateMachineRequest,
    UpdateUserAddressRequest,
    UpdateUserEmailRequest,
    UpdateUserPhoneRequest,
    UpdateUserProfileRequest,
    UserAddress,
    UserEmail,
    UserGrant,
    UserGrantCreate,
    UserGrantID,
    UserGrantRemoveBulk,
    UserGrantSearchQuery,
    UserGrantSearchRequest,
    UserGrantSearchResponse,
    UserGrantUpdate,
    UserGrantView,
    UserID,
    UserMembershipSearchQuery,
    UserMembershipSearchRequest,
    UserMembershipSearchResponse,
    UserMultiFactors,
    UserPhone,
    UserProfile,
    UserResponse,
    UserSearchQuery,
    UserSearchRequest,
    UserSearchResponse,
    UserView,
    ValidateOrgDomainRequest,
    WebAuthNTokenID,
    ZitadelDocs,
} from '../proto/generated/management_pb';
import { GrpcService } from './grpc.service';

export type ResponseMapper<TResp, TMappedResp> = (resp: TResp) => TMappedResp;

@Injectable({
    providedIn: 'root',
})
export class ManagementService {
    public ownedProjectsCount: BehaviorSubject<number> = new BehaviorSubject(0);
    public grantedProjectsCount: BehaviorSubject<number> = new BehaviorSubject(0);

    constructor(private readonly grpcService: GrpcService) { }

    public SearchIdps(
        limit?: number,
        offset?: number,
        queryList?: IdpSearchQuery[],
    ): Promise<IdpSearchResponse> {
        const req = new IdpSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchIdps(req);
    }

    public GetLoginPolicyMultiFactors(): Promise<MultiFactorsResult> {
        const req = new Empty();
        return this.grpcService.mgmt.getLoginPolicyMultiFactors(req);
    }

    public AddMultiFactorToLoginPolicy(req: MultiFactor): Promise<MultiFactor> {
        return this.grpcService.mgmt.addMultiFactorToLoginPolicy(req);
    }

    public RemoveMultiFactorFromLoginPolicy(req: MultiFactor): Promise<Empty> {
        return this.grpcService.mgmt.removeMultiFactorFromLoginPolicy(req);
    }

    public GetLoginPolicySecondFactors(): Promise<SecondFactorsResult> {
        const req = new Empty();
        return this.grpcService.mgmt.getLoginPolicySecondFactors(req);
    }

    public AddSecondFactorToLoginPolicy(req: SecondFactor): Promise<SecondFactor> {
        return this.grpcService.mgmt.addSecondFactorToLoginPolicy(req);
    }

    public RemoveSecondFactorFromLoginPolicy(req: SecondFactor): Promise<Empty> {
        return this.grpcService.mgmt.removeSecondFactorFromLoginPolicy(req);
    }

    public GetLoginPolicy(): Promise<LoginPolicyView> {
        const req = new Empty();
        return this.grpcService.mgmt.getLoginPolicy(req);
    }

    public UpdateLoginPolicy(req: LoginPolicyRequest): Promise<LoginPolicy> {
        return this.grpcService.mgmt.updateLoginPolicy(req);
    }

    public CreateLoginPolicy(req: LoginPolicyRequest): Promise<LoginPolicy> {
        return this.grpcService.mgmt.createLoginPolicy(req);
    }

    public RemoveLoginPolicy(): Promise<Empty> {
        return this.grpcService.mgmt.removeLoginPolicy(new Empty());
    }

    public addIdpProviderToLoginPolicy(configId: string, idpType: IdpProviderType): Promise<IdpProviderID> {
        const req = new IdpProviderAdd();
        req.setIdpProviderType(idpType);
        req.setIdpConfigId(configId);
        return this.grpcService.mgmt.addIdpProviderToLoginPolicy(req);
    }

    public RemoveIdpProviderFromLoginPolicy(configId: string): Promise<Empty> {
        const req = new IdpProviderID();
        req.setIdpConfigId(configId);
        return this.grpcService.mgmt.removeIdpProviderFromLoginPolicy(req);
    }

    public GetLoginPolicyIdpProviders(limit?: number, offset?: number): Promise<IdpProviderSearchResponse> {
        const req = new IdpProviderSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }
        return this.grpcService.mgmt.getLoginPolicyIdpProviders(req);
    }

    public IdpByID(
        id: string,
    ): Promise<IdpView> {
        const req = new IdpID();
        req.setId(id);
        return this.grpcService.mgmt.idpByID(req);
    }

    public UpdateIdp(
        req: IdpUpdate,
    ): Promise<Idp> {
        return this.grpcService.mgmt.updateIdpConfig(req);
    }

    public CreateOidcIdp(
        req: OidcIdpConfigCreate,
    ): Promise<Idp> {
        return this.grpcService.mgmt.createOidcIdp(req);
    }

    public UpdateOidcIdpConfig(
        req: OidcIdpConfigUpdate,
    ): Promise<OidcIdpConfig> {
        return this.grpcService.mgmt.updateOidcIdpConfig(req);
    }

    public RemoveIdpConfig(
        id: string,
    ): Promise<Empty> {
        const req = new IdpID;
        req.setId(id);
        return this.grpcService.mgmt.removeIdpConfig(req);
    }

    public DeactivateIdpConfig(
        id: string,
    ): Promise<Empty> {
        const req = new IdpID;
        req.setId(id);
        return this.grpcService.mgmt.deactivateIdpConfig(req);
    }

    public ReactivateIdpConfig(
        id: string,
    ): Promise<Empty> {
        const req = new IdpID;
        req.setId(id);
        return this.grpcService.mgmt.reactivateIdpConfig(req);
    }

    public CreateUserHuman(username: string, user: CreateHumanRequest): Promise<UserResponse> {
        const req = new CreateUserRequest();

        req.setUserName(username);
        req.setHuman(user);

        return this.grpcService.mgmt.createUser(req);
    }

    public CreateUserMachine(username: string, user: CreateMachineRequest): Promise<UserResponse> {
        const req = new CreateUserRequest();

        req.setUserName(username);
        req.setMachine(user);

        return this.grpcService.mgmt.createUser(req);
    }

    public UpdateUserMachine(
        id: string,
        description?: string,
    ): Promise<MachineResponse> {
        const req = new UpdateMachineRequest();
        req.setId(id);
        if (description) {
            req.setDescription(description);
        }
        return this.grpcService.mgmt.updateUserMachine(req);
    }

    public AddMachineKey(
        userId: string,
        type: MachineKeyType,
        date?: Timestamp,
    ): Promise<AddMachineKeyResponse> {
        const req = new AddMachineKeyRequest();
        req.setType(type);
        req.setUserId(userId);
        if (date) {
            req.setExpirationDate(date);
        }
        return this.grpcService.mgmt.addMachineKey(req);
    }

    public DeleteMachineKey(
        keyId: string,
        userId: string,
    ): Promise<Empty> {
        const req = new MachineKeyIDRequest();
        req.setKeyId(keyId);
        req.setUserId(userId);

        return this.grpcService.mgmt.deleteMachineKey(req);
    }

    public SearchMachineKeys(
        userId: string,
        limit: number,
        offset: number,
        asc?: boolean,
    ): Promise<MachineKeySearchResponse> {
        const req = new MachineKeySearchRequest();
        req.setUserId(userId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (asc) {
            req.setAsc(asc);
        }
        return this.grpcService.mgmt.searchMachineKeys(req);
    }

    public RemoveExternalIDP(
        externalUserId: string,
        idpConfigId: string,
        userId: string,
    ): Promise<Empty> {
        const req = new ExternalIDPRemoveRequest();
        req.setUserId(userId);
        req.setExternalUserId(externalUserId);
        req.setIdpConfigId(idpConfigId);
        return this.grpcService.mgmt.removeExternalIDP(req);
    }

    public SearchUserExternalIDPs(
        limit: number,
        offset: number,
        userId: string,
    ): Promise<ExternalIDPSearchResponse> {
        const req = new ExternalIDPSearchRequest();
        req.setUserId(userId);
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.mgmt.searchUserExternalIDPs(req);
    }

    public GetIam(): Promise<Iam> {
        const req = new Empty();
        return this.grpcService.mgmt.getIam(req);
    }

    public GetDefaultPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        const req = new Empty();
        return this.grpcService.mgmt.getDefaultPasswordComplexityPolicy(req);
    }

    public GetMyOrg(): Promise<OrgView> {
        const req = new Empty();
        return this.grpcService.mgmt.getMyOrg(req);
    }

    public AddMyOrgDomain(domain: string): Promise<OrgDomain> {
        const req: AddOrgDomainRequest = new AddOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.addMyOrgDomain(req);
    }

    public RemoveMyOrgDomain(domain: string): Promise<Empty> {
        const req: RemoveOrgDomainRequest = new AddOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.removeMyOrgDomain(req);
    }

    public SearchMyOrgDomains(queryList?: OrgDomainSearchQuery[]):
        Promise<OrgDomainSearchResponse> {
        const req: OrgDomainSearchRequest = new OrgDomainSearchRequest();
        if (queryList) {
            req.setQueriesList(queryList);
        }

        return this.grpcService.mgmt.searchMyOrgDomains(req);
    }

    public setMyPrimaryOrgDomain(domain: string): Promise<Empty> {
        const req: PrimaryOrgDomainRequest = new PrimaryOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.setMyPrimaryOrgDomain(req);
    }

    public GenerateMyOrgDomainValidation(domain: string, type: OrgDomainValidationType):
        Promise<OrgDomainValidationResponse> {
        const req: OrgDomainValidationRequest = new OrgDomainValidationRequest();
        req.setDomain(domain);
        req.setType(type);

        return this.grpcService.mgmt.generateMyOrgDomainValidation(req);
    }

    public ValidateMyOrgDomain(domain: string):
        Promise<Empty> {
        const req: ValidateOrgDomainRequest = new ValidateOrgDomainRequest();
        req.setDomain(domain);

        return this.grpcService.mgmt.validateMyOrgDomain(req);
    }

    public SearchMyOrgMembers(limit: number, offset: number): Promise<OrgMemberSearchResponse> {
        const req = new OrgMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.mgmt.searchMyOrgMembers(req);
    }

    public getOrgByDomainGlobal(domain: string): Promise<Org> {
        const req = new Domain();
        req.setDomain(domain);
        return this.grpcService.mgmt.getOrgByDomainGlobal(req);
    }

    public CreateOrg(name: string): Promise<Org> {
        const req = new OrgCreateRequest();
        req.setName(name);
        return this.grpcService.mgmt.createOrg(req);
    }

    public AddMyOrgMember(userId: string, rolesList: string[]): Promise<Empty> {
        const req = new AddOrgMemberRequest();
        req.setUserId(userId);
        if (rolesList) {
            req.setRolesList(rolesList);
        }
        return this.grpcService.mgmt.addMyOrgMember(req);
    }

    public ChangeMyOrgMember(userId: string, rolesList: string[]): Promise<OrgMember> {
        const req = new ChangeOrgMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeMyOrgMember(req);
    }


    public RemoveMyOrgMember(userId: string): Promise<Empty> {
        const req = new RemoveOrgMemberRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.removeMyOrgMember(req);
    }

    public DeactivateMyOrg(): Promise<Org> {
        const req = new Empty();
        return this.grpcService.mgmt.deactivateMyOrg(req);
    }

    public ReactivateMyOrg(): Promise<Org> {
        const req = new Empty();
        return this.grpcService.mgmt.reactivateMyOrg(req);
    }

    public CreateProjectGrant(
        orgId: string,
        projectId: string,
        roleKeysList: string[],
    ): Promise<ProjectGrant> {
        const req = new ProjectGrantCreate();
        req.setProjectId(projectId);
        req.setGrantedOrgId(orgId);
        req.setRoleKeysList(roleKeysList);
        return this.grpcService.mgmt.createProjectGrant(req);
    }

    public GetOrgMemberRoles(): Promise<OrgMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getOrgMemberRoles(req);
    }

    // Policy

    public GetMyOrgIamPolicy(): Promise<OrgIamPolicyView> {
        const req = new Empty();
        return this.grpcService.mgmt.getMyOrgIamPolicy(req);
    }

    public GetPasswordAgePolicy(): Promise<PasswordAgePolicyView> {
        const req = new Empty();

        return this.grpcService.mgmt.getPasswordAgePolicy(req);
    }

    public CreatePasswordAgePolicy(
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<PasswordAgePolicy> {
        const req = new PasswordAgePolicyRequest();
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);

        return this.grpcService.mgmt.createPasswordAgePolicy(req);
    }

    public RemovePasswordAgePolicy(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.mgmt.removePasswordAgePolicy(req);
    }

    public UpdatePasswordAgePolicy(
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<PasswordAgePolicy> {
        const req = new PasswordAgePolicyRequest();
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);
        return this.grpcService.mgmt.updatePasswordAgePolicy(req);
    }

    public GetPasswordComplexityPolicy(): Promise<PasswordComplexityPolicyView> {
        const req = new Empty();
        return this.grpcService.mgmt.getPasswordComplexityPolicy(req);
    }

    public CreatePasswordComplexityPolicy(
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<PasswordComplexityPolicy> {
        const req = new PasswordComplexityPolicyRequest();
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.mgmt.createPasswordComplexityPolicy(req);
    }

    public removePasswordComplexityPolicy(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.mgmt.removePasswordComplexityPolicy(req);
    }

    public UpdatePasswordComplexityPolicy(
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<PasswordComplexityPolicy> {
        const req = new PasswordComplexityPolicy();
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.mgmt.updatePasswordComplexityPolicy(req);
    }

    public GetPasswordLockoutPolicy(): Promise<PasswordLockoutPolicy> {
        const req = new Empty();

        return this.grpcService.mgmt.getPasswordLockoutPolicy(req);
    }

    public CreatePasswordLockoutPolicy(
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<PasswordLockoutPolicy> {
        const req = new PasswordLockoutPolicyRequest();
        req.setMaxAttempts(maxAttempts);
        req.setShowLockoutFailure(showLockoutFailures);

        return this.grpcService.mgmt.createPasswordLockoutPolicy(req);
    }

    public RemovePasswordLockoutPolicy(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.mgmt.removePasswordLockoutPolicy(req);
    }

    public UpdatePasswordLockoutPolicy(
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<PasswordLockoutPolicy> {
        const req = new PasswordLockoutPolicy();
        req.setMaxAttempts(maxAttempts);
        req.setShowLockoutFailure(showLockoutFailures);
        return this.grpcService.mgmt.updatePasswordLockoutPolicy(req);
    }

    public getLocalizedComplexityPolicyPatternErrorString(policy: PasswordComplexityPolicy.AsObject): string {
        if (policy.hasNumber && policy.hasSymbol) {
            return 'POLICY.PWD_COMPLEXITY.SYMBOLANDNUMBERERROR';
        } else if (policy.hasNumber) {
            return 'POLICY.PWD_COMPLEXITY.NUMBERERROR';
        } else if (policy.hasSymbol) {
            return 'POLICY.PWD_COMPLEXITY.SYMBOLERROR';
        } else {
            return 'POLICY.PWD_COMPLEXITY.PATTERNERROR';
        }
    }

    public GetUserByID(id: string): Promise<UserView> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserByID(req);
    }

    public DeleteUser(id: string): Promise<Empty> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.deleteUser(req);
    }

    public SearchProjectMembers(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ProjectMemberSearchQuery[],
    ): Promise<ProjectMemberSearchResponse> {
        const req = new ProjectMemberSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchProjectMembers(req);
    }

    public SearchUserMemberships(userId: string,
        limit: number, offset: number, queryList?: UserMembershipSearchQuery[]): Promise<UserMembershipSearchResponse> {
        const req = new UserMembershipSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        req.setUserId(userId);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUserMemberships(req);
    }

    public GetUserProfile(id: string): Promise<UserProfile> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserProfile(req);
    }

    public getUserMfas(id: string): Promise<UserMultiFactors> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserMfas(req);
    }

    public removeMfaOTP(id: string): Promise<Empty> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.removeMfaOTP(req);
    }

    public RemoveMfaU2F(userid: string, id: string): Promise<Empty> {
        const req = new WebAuthNTokenID();
        req.setId(id);
        req.setUserId(userid);
        return this.grpcService.mgmt.removeMfaU2F(req);
    }

    public SaveUserProfile(
        id: string,
        firstName?: string,
        lastName?: string,
        nickName?: string,
        preferredLanguage?: string,
        gender?: Gender,
    ): Promise<UserProfile> {
        const req = new UpdateUserProfileRequest();
        req.setId(id);
        if (firstName) {
            req.setFirstName(firstName);
        }
        if (lastName) {
            req.setLastName(lastName);
        }
        if (nickName) {
            req.setNickName(nickName);
        }
        if (gender) {
            req.setGender(gender);
        }
        if (preferredLanguage) {
            req.setPreferredLanguage(preferredLanguage);
        }
        return this.grpcService.mgmt.updateUserProfile(req);
    }

    public GetUserEmail(id: string): Promise<UserEmail> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserEmail(req);
    }

    public SaveUserEmail(id: string, email: string): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setId(id);
        req.setEmail(email);
        return this.grpcService.mgmt.changeUserEmail(req);
    }

    public GetUserPhone(id: string): Promise<UserPhone> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserPhone(req);
    }

    public SaveUserPhone(id: string, phone: string): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setId(id);
        req.setPhone(phone);
        return this.grpcService.mgmt.changeUserPhone(req);
    }

    public RemoveUserPhone(id: string): Promise<Empty> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.removeUserPhone(req);
    }

    public DeactivateUser(id: string): Promise<UserResponse> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.deactivateUser(req);
    }

    public CreateUserGrant(
        userId: string,
        roleNamesList: string[],
        projectId?: string,
        grantId?: string,
    ): Promise<UserGrant> {
        const req = new UserGrantCreate();
        if (projectId) { req.setProjectId(projectId); }
        if (grantId) { req.setGrantId(grantId); }
        req.setUserId(userId);
        req.setRoleKeysList(roleNamesList);

        return this.grpcService.mgmt.createUserGrant(req);
    }

    public ReactivateUser(id: string): Promise<UserResponse> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.reactivateUser(req);
    }

    public AddRole(id: string, key: string, displayName: string, group: string): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(id);
        req.setKey(key);
        if (displayName) {
            req.setDisplayName(displayName);
        }
        req.setGroup(group);
        return this.grpcService.mgmt.addProjectRole(req);
    }

    public GetUserAddress(id: string): Promise<UserAddress> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserAddress(req);
    }

    public ResendEmailVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.resendEmailVerificationMail(req);
    }

    public ResendInitialMail(userId: string, newemail: string): Promise<Empty> {
        const req = new InitialMailRequest();
        if (newemail) {
            req.setEmail(newemail);
        }
        req.setId(userId);

        return this.grpcService.mgmt.resendInitialMail(req);
    }

    public ResendPhoneVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.resendPhoneVerificationCode(req);
    }

    public SetInitialPassword(id: string, password: string): Promise<any> {
        const req = new PasswordRequest();
        req.setId(id);
        req.setPassword(password);
        return this.grpcService.mgmt.setInitialPassword(req);
    }

    public SendSetPasswordNotification(id: string, type: NotificationType): Promise<any> {
        const req = new SetPasswordNotificationRequest();
        req.setId(id);
        req.setType(type);
        return this.grpcService.mgmt.sendSetPasswordNotification(req);
    }

    public SaveUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setId(address.id);
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return this.grpcService.mgmt.updateUserAddress(req);
    }

    public SearchUsers(limit: number, offset: number, queryList?: UserSearchQuery[]): Promise<UserSearchResponse> {
        const req = new UserSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUsers(req);
    }

    public GetUserByLoginNameGlobal(loginname: string): Promise<UserView> {
        const req = new LoginName();
        req.setLoginName(loginname);
        return this.grpcService.mgmt.getUserByLoginNameGlobal(req);
    }

    // USER GRANTS

    public SearchUserGrants(
        limit?: number,
        offset?: number,
        queryList?: UserGrantSearchQuery[],
    ): Promise<UserGrantSearchResponse> {
        const req = new UserGrantSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUserGrants(req);
    }


    public UserGrantByID(
        id: string,
        userId: string,
    ): Promise<UserGrantView> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return this.grpcService.mgmt.userGrantByID(req);
    }

    public UpdateUserGrant(
        id: string,
        userId: string,
        roleKeysList: string[],
    ): Promise<UserGrant> {
        const req = new UserGrantUpdate();
        req.setId(id);
        req.setRoleKeysList(roleKeysList);
        req.setUserId(userId);

        return this.grpcService.mgmt.updateUserGrant(req);
    }

    public RemoveUserGrant(
        id: string,
        userId: string,
    ): Promise<Empty> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return this.grpcService.mgmt.removeUserGrant(req);
    }

    public BulkRemoveUserGrant(
        idsList: string[],
    ): Promise<Empty> {
        const req = new UserGrantRemoveBulk();
        req.setIdsList(idsList);

        return this.grpcService.mgmt.bulkRemoveUserGrant(req);
    }

    //

    public ApplicationChanges(id: string, secId: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setSecId(secId);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.applicationChanges(req);
    }

    public OrgChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.orgChanges(req);
    }

    public ProjectChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.projectChanges(req);
    }

    public UserChanges(id: string, limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return this.grpcService.mgmt.userChanges(req);
    }

    // project

    public SearchProjects(
        limit?: number, offset?: number, queryList?: ProjectSearchQuery[]): Promise<ProjectSearchResponse> {
        const req = new ProjectSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }

        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchProjects(req).then(value => {
            const count = value.toObject().resultList.length;
            if (count >= 0) {
                this.ownedProjectsCount.next(count);
            }

            return value;
        });
    }

    public SearchGrantedProjects(
        limit: number, offset: number, queryList?: ProjectSearchQuery[]): Promise<ProjectGrantSearchResponse> {
        const req = new GrantedProjectSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchGrantedProjects(req).then(value => {
            this.grantedProjectsCount.next(value.toObject().resultList.length);
            return value;
        });
    }

    public GetZitadelDocs(): Promise<ZitadelDocs> {
        const req = new Empty();
        return this.grpcService.mgmt.getZitadelDocs(req);
    }

    public GetProjectById(projectId: string): Promise<ProjectView> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.projectByID(req);
    }

    public GetGrantedProjectByID(projectId: string, id: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.getGrantedProjectByID(req);
    }

    public CreateProject(project: ProjectCreateRequest.AsObject): Promise<Project> {
        const req = new ProjectCreateRequest();
        req.setName(project.name);
        return this.grpcService.mgmt.createProject(req).then(value => {
            const current = this.ownedProjectsCount.getValue();
            this.ownedProjectsCount.next(current + 1);
            return value;
        });
    }

    public UpdateProject(id: string, projectView: ProjectView.AsObject): Promise<Project> {
        const req = new ProjectUpdateRequest();
        req.setId(id);
        req.setName(projectView.name);
        req.setProjectRoleAssertion(projectView.projectRoleAssertion);
        req.setProjectRoleCheck(projectView.projectRoleCheck);
        return this.grpcService.mgmt.updateProject(req);
    }

    public UpdateProjectGrant(id: string, projectId: string, rolesList: string[]): Promise<ProjectGrant> {
        const req = new ProjectGrantUpdate();
        req.setRoleKeysList(rolesList);
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.updateProjectGrant(req);
    }

    public RemoveProjectGrant(id: string, projectId: string): Promise<Empty> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.removeProjectGrant(req);
    }

    public DeactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.deactivateProject(req);
    }

    public ReactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.reactivateProject(req);
    }

    public SearchProjectGrants(projectId: string, limit: number, offset: number): Promise<ProjectGrantSearchResponse> {
        const req = new ProjectGrantSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.mgmt.searchProjectGrants(req);
    }

    public GetProjectGrantMemberRoles(): Promise<ProjectGrantMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getProjectGrantMemberRoles(req);
    }

    public AddProjectMember(id: string, userId: string, rolesList: string[]): Promise<Empty> {
        const req = new ProjectMemberAdd();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.addProjectMember(req);
    }

    public ChangeProjectMember(id: string, userId: string, rolesList: string[]): Promise<ProjectMember> {
        const req = new ProjectMemberChange();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeProjectMember(req);
    }

    public AddProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
        rolesList: string[],
    ): Promise<Empty> {
        const req = new ProjectGrantMemberAdd();
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.addProjectGrantMember(req);
    }

    public ChangeProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
        rolesList: string[],
    ): Promise<ProjectGrantMember> {
        const req = new ProjectGrantMemberChange();
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeProjectGrantMember(req);
    }

    public SearchProjectGrantMembers(
        projectId: string,
        grantId: string,
        limit: number,
        offset: number,
        queryList?: ProjectGrantMemberSearchQuery[],
    ): Promise<ProjectMemberSearchResponse> {
        const req = new ProjectGrantMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        return this.grpcService.mgmt.searchProjectGrantMembers(req);
    }

    public RemoveProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
    ): Promise<Empty> {
        const req = new ProjectGrantMemberRemove();
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.removeProjectGrantMember(req);
    }

    public ReactivateApplication(projectId: string, appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);

        return this.grpcService.mgmt.reactivateApplication(req);
    }

    public DeactivateApplication(projectId: string, appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);

        return this.grpcService.mgmt.deactivateApplication(req);
    }

    public RegenerateOIDCClientSecret(id: string, projectId: string): Promise<any> {
        const req = new ApplicationID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.regenerateOIDCClientSecret(req);
    }

    public SearchProjectRoles(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ProjectRoleSearchQuery[],
    ): Promise<ProjectRoleSearchResponse> {
        const req = new ProjectRoleSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchProjectRoles(req);
    }

    public AddProjectRole(role: ProjectRoleAdd.AsObject): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(role.id);
        if (role.displayName) {
            req.setDisplayName(role.displayName);
        }
        req.setKey(role.key);
        req.setGroup(role.group);
        return this.grpcService.mgmt.addProjectRole(req);
    }

    public BulkAddProjectRole(
        id: string,
        rolesList: ProjectRoleAdd[],
    ): Promise<Empty> {
        const req = new ProjectRoleAddBulk();
        req.setId(id);
        req.setProjectRolesList(rolesList);
        return this.grpcService.mgmt.bulkAddProjectRole(req);
    }

    public RemoveProjectRole(projectId: string, key: string): Promise<Empty> {
        const req = new ProjectRoleRemove();
        req.setId(projectId);
        req.setKey(key);
        return this.grpcService.mgmt.removeProjectRole(req);
    }


    public ChangeProjectRole(projectId: string, key: string, displayName: string, group: string):
        Promise<ProjectRole> {
        const req = new ProjectRoleChange();
        req.setId(projectId);
        req.setKey(key);
        req.setGroup(group);
        req.setDisplayName(displayName);
        return this.grpcService.mgmt.changeProjectRole(req);
    }


    public RemoveProjectMember(id: string, userId: string): Promise<Empty> {
        const req = new ProjectMemberRemove();
        req.setId(id);
        req.setUserId(userId);
        return this.grpcService.mgmt.removeProjectMember(req);
    }

    public SearchApplications(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ApplicationSearchQuery[]): Promise<ApplicationSearchResponse> {
        const req = new ApplicationSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchApplications(req);
    }

    public GetApplicationById(projectId: string, applicationId: string): Promise<ApplicationView> {
        const req = new ApplicationID();
        req.setProjectId(projectId);
        req.setId(applicationId);
        return this.grpcService.mgmt.applicationByID(req);
    }

    public GetProjectMemberRoles(): Promise<ProjectMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getProjectMemberRoles(req);
    }

    public ProjectGrantByID(id: string, projectId: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.projectGrantByID(req);
    }

    public RemoveProject(id: string): Promise<Empty> {
        const req = new ProjectID();
        req.setId(id);
        return this.grpcService.mgmt.removeProject(req).then(value => {
            const current = this.ownedProjectsCount.getValue();
            this.ownedProjectsCount.next(current > 0 ? current - 1 : 0);
            return value;
        });
    }


    public DeactivateProjectGrant(id: string, projectId: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.deactivateProjectGrant(req);
    }

    public ReactivateProjectGrant(id: string, projectId: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.reactivateProjectGrant(req);
    }

    public CreateOIDCApp(app: OIDCApplicationCreate.AsObject): Promise<Application> {
        const req = new OIDCApplicationCreate();
        req.setProjectId(app.projectId);
        req.setName(app.name);
        req.setRedirectUrisList(app.redirectUrisList);
        req.setResponseTypesList(app.responseTypesList);
        req.setGrantTypesList(app.grantTypesList);
        req.setApplicationType(app.applicationType);
        req.setAuthMethodType(app.authMethodType);
        req.setPostLogoutRedirectUrisList(app.postLogoutRedirectUrisList);

        return this.grpcService.mgmt.createOIDCApplication(req);
    }

    public UpdateApplication(projectId: string, appId: string, name: string): Promise<Application> {
        const req = new ApplicationUpdate();
        req.setId(appId);
        req.setName(name);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.updateApplication(req);
    }

    public UpdateOIDCAppConfig(projectId: string,
        appId: string, oidcConfig: OIDCConfig.AsObject): Promise<OIDCConfig> {
        const req = new OIDCConfigUpdate();
        req.setProjectId(projectId);
        req.setApplicationId(appId);
        req.setRedirectUrisList(oidcConfig.redirectUrisList);
        req.setResponseTypesList(oidcConfig.responseTypesList);
        req.setAuthMethodType(oidcConfig.authMethodType);
        req.setPostLogoutRedirectUrisList(oidcConfig.postLogoutRedirectUrisList);
        req.setGrantTypesList(oidcConfig.grantTypesList);
        req.setApplicationType(oidcConfig.applicationType);
        req.setDevMode(oidcConfig.devMode);
        req.setAccessTokenType(oidcConfig.accessTokenType);
        req.setAccessTokenRoleAssertion(oidcConfig.accessTokenRoleAssertion);
        req.setIdTokenRoleAssertion(oidcConfig.idTokenRoleAssertion);
        return this.grpcService.mgmt.updateApplicationOIDCConfig(req);
    }
}

import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';

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
    Gender,
    GrantedProjectSearchRequest,
    Iam,
    LoginName,
    MachineKeyIDRequest,
    MachineKeySearchRequest,
    MachineKeySearchResponse,
    MachineKeyType,
    MachineResponse,
    MultiFactors,
    NotificationType,
    OIDCApplicationCreate,
    OIDCConfig,
    OIDCConfigUpdate,
    Org,
    OrgCreateRequest,
    OrgDomain,
    OrgDomainSearchQuery,
    OrgDomainSearchRequest,
    OrgDomainSearchResponse,
    OrgDomainValidationRequest,
    OrgDomainValidationResponse,
    OrgDomainValidationType,
    OrgIamPolicy,
    OrgMember,
    OrgMemberRoles,
    OrgMemberSearchRequest,
    OrgMemberSearchResponse,
    OrgView,
    PasswordAgePolicy,
    PasswordAgePolicyCreate,
    PasswordAgePolicyID,
    PasswordAgePolicyUpdate,
    PasswordComplexityPolicy,
    PasswordComplexityPolicyCreate,
    PasswordComplexityPolicyID,
    PasswordComplexityPolicyUpdate,
    PasswordLockoutPolicy,
    PasswordLockoutPolicyCreate,
    PasswordLockoutPolicyID,
    PasswordLockoutPolicyUpdate,
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
    UserPhone,
    UserProfile,
    UserResponse,
    UserSearchQuery,
    UserSearchRequest,
    UserSearchResponse,
    UserView,
    ValidateOrgDomainRequest,
    ZitadelDocs,
} from '../proto/generated/management_pb';
import { GrpcService } from './grpc.service';

export type ResponseMapper<TResp, TMappedResp> = (resp: TResp) => TMappedResp;

@Injectable({
    providedIn: 'root',
})
export class ManagementService {
    constructor(private readonly grpcService: GrpcService) { }

    public async CreateUserHuman(username: string, user: CreateHumanRequest): Promise<UserResponse> {
        const req = new CreateUserRequest();

        req.setUserName(username);
        req.setHuman(user);

        return this.grpcService.mgmt.createUser(req);
    }

    public async CreateUserMachine(username: string, user: CreateMachineRequest): Promise<UserResponse> {
        const req = new CreateUserRequest();

        req.setUserName(username);
        req.setMachine(user);

        return this.grpcService.mgmt.createUser(req);
    }

    public async UpdateUserMachine(
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

    public async AddMachineKey(
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

    public async DeleteMachineKey(
        keyId: string,
        userId: string,
    ): Promise<Empty> {
        const req = new MachineKeyIDRequest();
        req.setKeyId(keyId);
        req.setUserId(userId);

        return this.grpcService.mgmt.deleteMachineKey(req);
    }

    public async SearchMachineKeys(
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

    public async GetIam(): Promise<Iam> {
        const req = new Empty();
        return this.grpcService.mgmt.getIam(req);
    }

    public async GetDefaultPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        const req = new Empty();
        return this.grpcService.mgmt.getDefaultPasswordComplexityPolicy(req);
    }

    public async GetMyOrg(): Promise<OrgView> {
        const req = new Empty();
        return this.grpcService.mgmt.getMyOrg(req);
    }

    public async AddMyOrgDomain(domain: string): Promise<OrgDomain> {
        const req: AddOrgDomainRequest = new AddOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.addMyOrgDomain(req);
    }

    public async RemoveMyOrgDomain(domain: string): Promise<Empty> {
        const req: RemoveOrgDomainRequest = new AddOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.removeMyOrgDomain(req);
    }

    public async SearchMyOrgDomains(offset: number, limit: number, queryList?: OrgDomainSearchQuery[]):
        Promise<OrgDomainSearchResponse> {
        const req: OrgDomainSearchRequest = new OrgDomainSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }

        return this.grpcService.mgmt.searchMyOrgDomains(req);
    }

    public async setMyPrimaryOrgDomain(domain: string): Promise<Empty> {
        const req: PrimaryOrgDomainRequest = new PrimaryOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.setMyPrimaryOrgDomain(req);
    }

    public async GenerateMyOrgDomainValidation(domain: string, type: OrgDomainValidationType):
        Promise<OrgDomainValidationResponse> {
        const req: OrgDomainValidationRequest = new OrgDomainValidationRequest();
        req.setDomain(domain);
        req.setType(type);

        return this.grpcService.mgmt.generateMyOrgDomainValidation(req);
    }

    public async ValidateMyOrgDomain(domain: string):
        Promise<Empty> {
        const req: ValidateOrgDomainRequest = new ValidateOrgDomainRequest();
        req.setDomain(domain);

        return this.grpcService.mgmt.validateMyOrgDomain(req);
    }

    public async SearchMyOrgMembers(limit: number, offset: number): Promise<OrgMemberSearchResponse> {
        const req = new OrgMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.mgmt.searchMyOrgMembers(req);
    }

    public async getOrgByDomainGlobal(domain: string): Promise<Org> {
        const req = new Domain();
        req.setDomain(domain);
        return this.grpcService.mgmt.getOrgByDomainGlobal(req);
    }

    public async CreateOrg(name: string): Promise<Org> {
        const req = new OrgCreateRequest();
        req.setName(name);
        return this.grpcService.mgmt.createOrg(req);
    }

    public async AddMyOrgMember(userId: string, rolesList: string[]): Promise<Empty> {
        const req = new AddOrgMemberRequest();
        req.setUserId(userId);
        if (rolesList) {
            req.setRolesList(rolesList);
        }
        return this.grpcService.mgmt.addMyOrgMember(req);
    }

    public async ChangeMyOrgMember(userId: string, rolesList: string[]): Promise<OrgMember> {
        const req = new ChangeOrgMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeMyOrgMember(req);
    }


    public async RemoveMyOrgMember(userId: string): Promise<Empty> {
        const req = new RemoveOrgMemberRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.removeMyOrgMember(req);
    }

    public async DeactivateMyOrg(): Promise<Org> {
        const req = new Empty();
        return this.grpcService.mgmt.deactivateMyOrg(req);
    }

    public async ReactivateMyOrg(): Promise<Org> {
        const req = new Empty();
        return this.grpcService.mgmt.reactivateMyOrg(req);
    }

    public async CreateProjectGrant(
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

    public async GetOrgMemberRoles(): Promise<OrgMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getOrgMemberRoles(req);
    }

    // Policy

    public async GetMyOrgIamPolicy(): Promise<OrgIamPolicy> {
        const req = new Empty();
        return this.grpcService.mgmt.getMyOrgIamPolicy(req);
    }

    public async GetPasswordAgePolicy(): Promise<PasswordAgePolicy> {
        const req = new Empty();

        return this.grpcService.mgmt.getPasswordAgePolicy(req);
    }

    public async CreatePasswordAgePolicy(
        description: string,
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<PasswordAgePolicy> {
        const req = new PasswordAgePolicyCreate();
        req.setDescription(description);
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);

        return this.grpcService.mgmt.createPasswordAgePolicy(req);
    }

    public async DeletePasswordAgePolicy(id: string): Promise<Empty> {
        const req = new PasswordAgePolicyID();
        req.setId(id);
        return this.grpcService.mgmt.deletePasswordAgePolicy(req);
    }

    public async UpdatePasswordAgePolicy(
        description: string,
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<PasswordAgePolicy> {
        const req = new PasswordAgePolicyUpdate();
        req.setDescription(description);
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);
        return this.grpcService.mgmt.updatePasswordAgePolicy(req);
    }

    public async GetPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        const req = new Empty();
        return this.grpcService.mgmt.getPasswordComplexityPolicy(req);
    }

    public async CreatePasswordComplexityPolicy(
        description: string,
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<PasswordComplexityPolicy> {
        const req = new PasswordComplexityPolicyCreate();
        req.setDescription(description);
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.mgmt.createPasswordComplexityPolicy(req);
    }

    public async DeletePasswordComplexityPolicy(id: string): Promise<Empty> {
        const req = new PasswordComplexityPolicyID();
        req.setId(id);
        return this.grpcService.mgmt.deletePasswordComplexityPolicy(req);
    }

    public async UpdatePasswordComplexityPolicy(
        description: string,
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<PasswordComplexityPolicy> {
        const req = new PasswordComplexityPolicyUpdate();
        req.setDescription(description);
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.mgmt.updatePasswordComplexityPolicy(req);
    }

    public async GetPasswordLockoutPolicy(): Promise<PasswordLockoutPolicy> {
        const req = new Empty();

        return this.grpcService.mgmt.getPasswordLockoutPolicy(req);
    }

    public async CreatePasswordLockoutPolicy(
        description: string,
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<PasswordLockoutPolicy> {
        const req = new PasswordLockoutPolicyCreate();
        req.setDescription(description);
        req.setMaxAttempts(maxAttempts);
        req.setShowLockOutFailures(showLockoutFailures);

        return this.grpcService.mgmt.createPasswordLockoutPolicy(req);
    }

    public async DeletePasswordLockoutPolicy(id: string): Promise<Empty> {
        const req = new PasswordLockoutPolicyID();
        req.setId(id);

        return this.grpcService.mgmt.deletePasswordLockoutPolicy(req);
    }

    public async UpdatePasswordLockoutPolicy(
        description: string,
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<PasswordLockoutPolicy> {
        const req = new PasswordLockoutPolicyUpdate();
        req.setDescription(description);
        req.setMaxAttempts(maxAttempts);
        req.setShowLockOutFailures(showLockoutFailures);
        return this.grpcService.mgmt.updatePasswordLockoutPolicy(req);
    }

    public getLocalizedComplexityPolicyPatternErrorString(policy: PasswordComplexityPolicy.AsObject): string {
        if (policy.hasNumber && policy.hasSymbol) {
            return 'ORG.POLICY.PWD_COMPLEXITY.SYMBOLANDNUMBERERROR';
        } else if (policy.hasNumber) {
            return 'ORG.POLICY.PWD_COMPLEXITY.NUMBERERROR';
        } else if (policy.hasSymbol) {
            return 'ORG.POLICY.PWD_COMPLEXITY.SYMBOLERROR';
        } else {
            return 'ORG.POLICY.PWD_COMPLEXITY.PATTERNERROR';
        }
    }

    public async GetUserByID(id: string): Promise<UserView> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserByID(req);
    }

    public async SearchProjectMembers(
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

    public async SearchUserMemberships(userId: string,
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

    public async GetUserProfile(id: string): Promise<UserProfile> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserProfile(req);
    }

    public async getUserMfas(id: string): Promise<MultiFactors> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserMfas(req);
    }

    public async SaveUserProfile(
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

    public async GetUserEmail(id: string): Promise<UserEmail> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserEmail(req);
    }

    public async SaveUserEmail(id: string, email: string): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setId(id);
        req.setEmail(email);
        return this.grpcService.mgmt.changeUserEmail(req);
    }

    public async GetUserPhone(id: string): Promise<UserPhone> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserPhone(req);
    }

    public async SaveUserPhone(id: string, phone: string): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setId(id);
        req.setPhone(phone);
        return this.grpcService.mgmt.changeUserPhone(req);
    }

    public async RemoveUserPhone(id: string): Promise<Empty> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.removeUserPhone(req);
    }

    public async DeactivateUser(id: string): Promise<UserResponse> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.deactivateUser(req);
    }

    public async CreateUserGrant(
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

    public async ReactivateUser(id: string): Promise<UserResponse> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.reactivateUser(req);
    }

    public async AddRole(id: string, key: string, displayName: string, group: string): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(id);
        req.setKey(key);
        if (displayName) {
            req.setDisplayName(displayName);
        }
        req.setGroup(group);
        return this.grpcService.mgmt.addProjectRole(req);
    }

    public async GetUserAddress(id: string): Promise<UserAddress> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserAddress(req);
    }

    public async ResendEmailVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.resendEmailVerificationMail(req);
    }

    public async ResendPhoneVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.resendPhoneVerificationCode(req);
    }

    public async SetInitialPassword(id: string, password: string): Promise<any> {
        const req = new PasswordRequest();
        req.setId(id);
        req.setPassword(password);
        return this.grpcService.mgmt.setInitialPassword(req);
    }

    public async SendSetPasswordNotification(id: string, type: NotificationType): Promise<any> {
        const req = new SetPasswordNotificationRequest();
        req.setId(id);
        req.setType(type);
        return this.grpcService.mgmt.sendSetPasswordNotification(req);
    }

    public async SaveUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setId(address.id);
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return this.grpcService.mgmt.updateUserAddress(req);
    }

    public async SearchUsers(limit: number, offset: number, queryList?: UserSearchQuery[]): Promise<UserSearchResponse> {
        const req = new UserSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUsers(req);
    }

    public async GetUserByLoginNameGlobal(loginname: string): Promise<UserView> {
        const req = new LoginName();
        req.setLoginName(loginname);
        return this.grpcService.mgmt.getUserByLoginNameGlobal(req);
    }

    // USER GRANTS

    public async SearchUserGrants(
        limit: number,
        offset: number,
        queryList?: UserGrantSearchQuery[],
    ): Promise<UserGrantSearchResponse> {
        const req = new UserGrantSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUserGrants(req);
    }


    public async UserGrantByID(
        id: string,
        userId: string,
    ): Promise<UserGrantView> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return this.grpcService.mgmt.userGrantByID(req);
    }

    public async UpdateUserGrant(
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

    public async RemoveUserGrant(
        id: string,
        userId: string,
    ): Promise<Empty> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return this.grpcService.mgmt.removeUserGrant(req);
    }

    public async BulkRemoveUserGrant(
        idsList: string[],
    ): Promise<Empty> {
        const req = new UserGrantRemoveBulk();
        req.setIdsList(idsList);

        return this.grpcService.mgmt.bulkRemoveUserGrant(req);
    }

    //

    public async ApplicationChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.applicationChanges(req);
    }

    public async OrgChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.orgChanges(req);
    }

    public async ProjectChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.projectChanges(req);
    }

    public async UserChanges(id: string, limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return this.grpcService.mgmt.userChanges(req);
    }

    // project

    public async SearchProjects(
        limit: number, offset: number, queryList?: ProjectSearchQuery[]): Promise<ProjectSearchResponse> {
        const req = new ProjectSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchProjects(req);
    }

    public async SearchGrantedProjects(
        limit: number, offset: number, queryList?: ProjectSearchQuery[]): Promise<ProjectGrantSearchResponse> {
        const req = new GrantedProjectSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchGrantedProjects(req);
    }


    public async GetZitadelDocs(): Promise<ZitadelDocs> {
        const req = new Empty();
        return this.grpcService.mgmt.getZitadelDocs(req);
    }

    public async GetProjectById(projectId: string): Promise<ProjectView> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.projectByID(req);
    }

    public async GetGrantedProjectByID(projectId: string, id: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.getGrantedProjectByID(req);
    }

    public async CreateProject(project: ProjectCreateRequest.AsObject): Promise<Project> {
        const req = new ProjectCreateRequest();
        req.setName(project.name);
        return this.grpcService.mgmt.createProject(req);
    }

    public async UpdateProject(id: string, name: string): Promise<Project> {
        const req = new ProjectUpdateRequest();
        req.setName(name);
        req.setId(id);
        return this.grpcService.mgmt.updateProject(req);
    }

    public async UpdateProjectGrant(id: string, projectId: string, rolesList: string[]): Promise<ProjectGrant> {
        const req = new ProjectGrantUpdate();
        req.setRoleKeysList(rolesList);
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.updateProjectGrant(req);
    }

    public async RemoveProjectGrant(id: string, projectId: string): Promise<Empty> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.removeProjectGrant(req);
    }

    public async DeactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.deactivateProject(req);
    }

    public async ReactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.reactivateProject(req);
    }

    public async SearchProjectGrants(projectId: string, limit: number, offset: number): Promise<ProjectGrantSearchResponse> {
        const req = new ProjectGrantSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.mgmt.searchProjectGrants(req);
    }

    public async GetProjectGrantMemberRoles(): Promise<ProjectGrantMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getProjectGrantMemberRoles(req);
    }

    public async AddProjectMember(id: string, userId: string, rolesList: string[]): Promise<Empty> {
        const req = new ProjectMemberAdd();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.addProjectMember(req);
    }

    public async ChangeProjectMember(id: string, userId: string, rolesList: string[]): Promise<ProjectMember> {
        const req = new ProjectMemberChange();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeProjectMember(req);
    }

    public async AddProjectGrantMember(
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

    public async ChangeProjectGrantMember(
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

    public async SearchProjectGrantMembers(
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

    public async RemoveProjectGrantMember(
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

    public async ReactivateApplication(projectId: string, appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);

        return this.grpcService.mgmt.reactivateApplication(req);
    }

    public async DeactivateApplication(projectId: string, appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);

        return this.grpcService.mgmt.deactivateApplication(req);
    }

    public async RegenerateOIDCClientSecret(id: string, projectId: string): Promise<any> {
        const req = new ApplicationID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.regenerateOIDCClientSecret(req);
    }

    public async SearchProjectRoles(
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

    public async AddProjectRole(role: ProjectRoleAdd.AsObject): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(role.id);
        if (role.displayName) {
            req.setDisplayName(role.displayName);
        }
        req.setKey(role.key);
        req.setGroup(role.group);
        return this.grpcService.mgmt.addProjectRole(req);
    }

    public async BulkAddProjectRole(
        id: string,
        rolesList: ProjectRoleAdd[],
    ): Promise<Empty> {
        const req = new ProjectRoleAddBulk();
        req.setId(id);
        req.setProjectRolesList(rolesList);
        return this.grpcService.mgmt.bulkAddProjectRole(req);
    }

    public async RemoveProjectRole(projectId: string, key: string): Promise<Empty> {
        const req = new ProjectRoleRemove();
        req.setId(projectId);
        req.setKey(key);
        return this.grpcService.mgmt.removeProjectRole(req);
    }


    public async ChangeProjectRole(projectId: string, key: string, displayName: string, group: string):
        Promise<ProjectRole> {
        const req = new ProjectRoleChange();
        req.setId(projectId);
        req.setKey(key);
        req.setGroup(group);
        req.setDisplayName(displayName);
        return this.grpcService.mgmt.changeProjectRole(req);
    }


    public async RemoveProjectMember(id: string, userId: string): Promise<Empty> {
        const req = new ProjectMemberRemove();
        req.setId(id);
        req.setUserId(userId);
        return this.grpcService.mgmt.removeProjectMember(req);
    }

    public async SearchApplications(
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

    public async GetApplicationById(projectId: string, applicationId: string): Promise<ApplicationView> {
        const req = new ApplicationID();
        req.setProjectId(projectId);
        req.setId(applicationId);
        return this.grpcService.mgmt.applicationByID(req);
    }

    public async GetProjectMemberRoles(): Promise<ProjectMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getProjectMemberRoles(req);
    }

    public async ProjectGrantByID(id: string, projectId: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.projectGrantByID(req);
    }

    public async RemoveProject(id: string): Promise<Empty> {
        const req = new ProjectID();
        req.setId(id);
        return this.grpcService.mgmt.removeProject(req);
    }


    public async DeactivateProjectGrant(id: string, projectId: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.deactivateProjectGrant(req);
    }

    public async ReactivateProjectGrant(id: string, projectId: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.reactivateProjectGrant(req);
    }

    public async CreateOIDCApp(app: OIDCApplicationCreate.AsObject): Promise<Application> {
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

    public async UpdateApplication(projectId: string, appId: string, name: string): Promise<Application> {
        const req = new ApplicationUpdate();
        req.setId(appId);
        req.setName(name);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.updateApplication(req);
    }

    public async UpdateOIDCAppConfig(projectId: string,
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
        return this.grpcService.mgmt.updateApplicationOIDCConfig(req);
    }
}

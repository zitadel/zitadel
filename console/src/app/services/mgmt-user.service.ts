import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';

import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import {
    ChangeRequest,
    Changes,
    CreateUserRequest,
    Email,
    Gender,
    MultiFactors,
    NotificationType,
    PasswordRequest,
    ProjectGrantMemberSearchQuery,
    ProjectGrantMemberSearchRequest,
    ProjectGrantMemberSearchResponse,
    ProjectGrantUserGrantID,
    ProjectGrantUserGrantSearchRequest,
    ProjectGrantUserGrantUpdate,
    ProjectRoleAdd,
    SetPasswordNotificationRequest,
    UpdateUserAddressRequest,
    UpdateUserEmailRequest,
    UpdateUserPhoneRequest,
    UpdateUserProfileRequest,
    User,
    UserAddress,
    UserEmail,
    UserGrant,
    UserGrantCreate,
    UserGrantID,
    UserGrantSearchQuery,
    UserGrantSearchRequest,
    UserGrantSearchResponse,
    UserGrantUpdate,
    UserGrantView,
    UserID,
    UserPhone,
    UserProfile,
    UserSearchQuery,
    UserSearchRequest,
    UserSearchResponse,
    UserView,
} from '../proto/generated/management_pb';
import { GrpcBackendService } from './grpc-backend.service';
import { GrpcService, RequestFactory, ResponseMapper } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class MgmtUserService {
    constructor(private grpcService: GrpcService, private grpcBackendService: GrpcBackendService) { }

    public async request<TReq, TResp, TMappedResp>(
        requestFn: RequestFactory<ManagementServicePromiseClient, TReq, TResp>,
        request: TReq,
        responseMapper: ResponseMapper<TResp, TMappedResp>,
        metadata?: Metadata,
    ): Promise<TMappedResp> {
        const mappedRequestFn = requestFn(this.grpcService.mgmt).bind(this.grpcService.mgmt);
        const response = await this.grpcBackendService.runRequest(
            mappedRequestFn,
            request,
            metadata,
        );
        return responseMapper(response);
    }

    public async CreateUser(user: CreateUserRequest.AsObject): Promise<User> {
        const req = new CreateUserRequest();
        req.setEmail(user.email);
        req.setUserName(user.userName);
        req.setFirstName(user.firstName);
        req.setLastName(user.lastName);
        req.setNickName(user.nickName);
        req.setPassword(user.password);
        req.setPreferredLanguage(user.preferredLanguage);
        req.setGender(user.gender);
        req.setPhone(user.phone);
        req.setStreetAddress(user.streetAddress);
        req.setPostalCode(user.postalCode);
        req.setLocality(user.locality);
        req.setRegion(user.region);
        req.setCountry(user.country);
        return await this.request(
            c => c.createUser,
            req,
            f => f,
        );
    }

    public async GetUserByID(id: string): Promise<UserView> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.getUserByID,
            req,
            f => f,
        );
    }

    public async GetUserProfile(id: string): Promise<UserProfile> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.getUserProfile,
            req,
            f => f,
        );
    }

    public async getUserMfas(id: string): Promise<MultiFactors> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.getUserMfas,
            req,
            f => f,
        );
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
        return await this.request(
            c => c.updateUserProfile,
            req,
            f => f,
        );
    }

    public async GetUserEmail(id: string): Promise<UserEmail> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.getUserEmail,
            req,
            f => f,
        );
    }

    public async SaveUserEmail(id: string, email: string): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setId(id);
        req.setEmail(email);
        return await this.request(
            c => c.changeUserEmail,
            req,
            f => f,
        );
    }

    public async GetUserPhone(id: string): Promise<UserPhone> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.getUserPhone,
            req,
            f => f,
        );
    }

    public async SaveUserPhone(id: string, phone: string): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setId(id);
        req.setPhone(phone);
        return await this.request(
            c => c.changeUserPhone,
            req,
            f => f,
        );
    }

    public async DeactivateUser(id: string): Promise<UserPhone> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.deactivateUser,
            req,
            f => f,
        );
    }

    public async CreateUserGrant(
        projectId: string,
        userId: string,
        roleNamesList: string[],
    ): Promise<UserGrant> {
        const req = new UserGrantCreate();
        req.setProjectId(projectId);
        req.setUserId(userId);
        req.setRoleKeysList(roleNamesList);

        return await this.request(
            c => c.createUserGrant,
            req,
            f => f,
        );
    }

    public async ReactivateUser(id: string): Promise<UserPhone> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.reactivateUser,
            req,
            f => f,
        );
    }

    public async AddRole(id: string, key: string, displayName: string, group: string): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(id);
        req.setKey(key);
        if (displayName) {
            req.setDisplayName(displayName);
        }
        req.setGroup(group);
        return await this.request(
            c => c.addProjectRole,
            req,
            f => f,
        );
    }

    public async GetUserAddress(id: string): Promise<UserAddress> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.getUserAddress,
            req,
            f => f,
        );
    }

    public async ResendEmailVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.resendEmailVerificationMail,
            req,
            f => f,
        );
    }

    public async ResendPhoneVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.resendPhoneVerificationCode,
            req,
            f => f,
        );
    }

    public async SetInitialPassword(id: string, password: string): Promise<any> {
        const req = new PasswordRequest();
        req.setId(id);
        req.setPassword(password);
        return await this.request(
            c => c.setInitialPassword,
            req,
            f => f,
        );
    }

    public async SendSetPasswordNotification(id: string, type: NotificationType): Promise<any> {
        const req = new SetPasswordNotificationRequest();
        req.setId(id);
        req.setType(type);
        return await this.request(
            c => c.sendSetPasswordNotification,
            req,
            f => f,
        );
    }

    public async SaveUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setId(address.id);
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return await this.request(
            c => c.updateUserAddress,
            req,
            f => f,
        );
    }

    public async SearchProjectGrantMembers(
        limit: number,
        offset: number,
        query?: ProjectGrantMemberSearchQuery[],
    ): Promise<ProjectGrantMemberSearchResponse> {
        const req = new ProjectGrantMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (query) {
            req.setQueriesList(query);
        }
        return await this.request(
            c => c.searchProjectGrantMembers,
            req,
            f => f,
        );
    }

    public async SearchUsers(limit: number, offset: number, queryList?: UserSearchQuery[]): Promise<UserSearchResponse> {
        const req = new UserSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchUsers,
            req,
            f => f,
        );
    }

    public async GetUserByEmailGlobal(email: string): Promise<User> {
        const req = new Email();
        req.setEmail(email);
        return await this.request(
            c => c.getUserByEmailGlobal,
            req,
            f => f,
        );
    }

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
        return await this.request(
            c => c.searchUserGrants,
            req,
            f => f,
        );
    }


    public async searchProjectGrantUserGrants(
        limit: number,
        offset: number,
        queryList?: UserGrantSearchQuery[],
    ): Promise<UserGrantSearchResponse> {
        const req = new ProjectGrantUserGrantSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchProjectGrantUserGrants,
            req,
            f => f,
        );
    }

    public async projectGrantUserGrantByID(
        id: string,
        userId: string,
        projectGrantId: string,
    ): Promise<UserGrant> {
        const req = new ProjectGrantUserGrantID();
        req.setId(id);
        req.setUserId(userId);
        req.setProjectGrantId(projectGrantId);

        return await this.request(
            c => c.projectGrantUserGrantByID,
            req,
            f => f,
        );
    }

    public async UserGrantByID(
        id: string,
        userId: string,
    ): Promise<UserGrantView> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return await this.request(
            c => c.userGrantByID,
            req,
            f => f,
        );
    }

    public async UpdateUserGrant(
        id: string,
        roleKeysList: string[],
        userId: string,
    ): Promise<UserGrant> {
        const req = new UserGrantUpdate();
        req.setId(id);
        req.setRoleKeysList(roleKeysList);
        req.setUserId(userId);

        return await this.request(
            c => c.updateUserGrant,
            req,
            f => f,
        );
    }

    public async updateProjectGrantUserGrant(
        id: string,
        roleKeysList: string[],
        userId: string,
        projectGrantId: string,
    ): Promise<UserGrant> {
        const req = new ProjectGrantUserGrantUpdate();
        req.setId(id);
        req.setRoleKeysList(roleKeysList);
        req.setUserId(userId);
        req.setProjectGrantId(projectGrantId);

        return await this.request(
            c => c.updateProjectGrantUserGrant,
            req,
            f => f,
        );
    }

    public async ApplicationChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return await this.request(
            c => c.applicationChanges,
            req,
            f => f,
        );
    }

    public async OrgChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return await this.request(
            c => c.orgChanges,
            req,
            f => f,
        );
    }

    public async ProjectChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return await this.request(
            c => c.projectChanges,
            req,
            f => f,
        );
    }

    public async UserChanges(id: string, limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return await this.request(
            c => c.userChanges,
            req,
            f => f,
        );
    }
}

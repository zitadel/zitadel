import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';
import { from, Observable, of } from 'rxjs';
import { catchError, switchMap } from 'rxjs/operators';

import { AuthServicePromiseClient } from '../proto/generated/auth_grpc_web_pb';
import {
    Changes,
    ChangesRequest,
    Gender,
    MfaOtpResponse,
    MultiFactors,
    MyProjectOrgSearchQuery,
    MyProjectOrgSearchRequest,
    MyProjectOrgSearchResponse,
    PasswordChange,
    PasswordComplexityPolicy,
    UpdateUserAddressRequest,
    UpdateUserEmailRequest,
    UpdateUserPhoneRequest,
    UpdateUserProfileRequest,
    UserAddress,
    UserEmail,
    UserPhone,
    UserProfile,
    UserProfileView,
    UserSessionViews,
    UserView,
    VerifyMfaOtp,
    VerifyUserPhoneRequest,
} from '../proto/generated/auth_pb';
import { GrpcBackendService } from './grpc-backend.service';
import { GrpcService, RequestFactory, ResponseMapper } from './grpc.service';


@Injectable({
    providedIn: 'root',
})
export class AuthUserService {
    private _roleCache: string[] = [];

    constructor(private readonly grpcClient: GrpcService,
        private grpcBackendService: GrpcBackendService,
    ) { }

    public async request<TReq, TResp, TMappedResp>(
        requestFn: RequestFactory<AuthServicePromiseClient, TReq, TResp>,
        request: TReq,
        responseMapper: ResponseMapper<TResp, TMappedResp>,
        metadata?: Metadata,
    ): Promise<TMappedResp> {
        const mappedRequestFn = requestFn(this.grpcClient.auth).bind(this.grpcClient.auth);
        const response = await this.grpcBackendService.runRequest(
            mappedRequestFn,
            request,
            metadata,
        );
        return responseMapper(response);
    }

    public async GetMyUserProfile(): Promise<UserProfileView> {
        return await this.request(
            c => c.getMyUserProfile,
            new Empty(),
            f => f,
        );
    }

    public async GetMyPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        return await this.request(
            c => c.getMyPasswordComplexityPolicy,
            new Empty(),
            f => f,
        );
    }


    public async GetMyUser(): Promise<UserView> {
        return await this.request(
            c => c.getMyUser,
            new Empty(),
            f => f,
        );
    }

    public async GetMyMfas(): Promise<MultiFactors> {
        return await this.request(
            c => c.getMyMfas,
            new Empty(),
            f => f,
        );
    }

    public async SearchMyProjectOrgs(
        limit: number,
        offset: number,
        queryList?: MyProjectOrgSearchQuery[],
    ): Promise<MyProjectOrgSearchResponse> {
        const req: MyProjectOrgSearchRequest = new MyProjectOrgSearchRequest();
        req.setOffset(offset);
        req.setLimit(limit);
        if (queryList) {
            req.setQueriesList(queryList);
        }

        return await this.request(
            c => c.searchMyProjectOrgs,
            req,
            f => f,
        );
    }

    public async SaveMyUserProfile(
        firstName?: string,
        lastName?: string,
        nickName?: string,
        preferredLanguage?: string,
        gender?: Gender,
    ): Promise<UserProfile> {
        const req = new UpdateUserProfileRequest();
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
            c => c.updateMyUserProfile,
            req,
            f => f,
        );
    }

    public async getMyUserSessions(): Promise<UserSessionViews> {
        return await this.request(
            c => c.getMyUserSessions,
            new Empty(),
            f => f,
        );
    }

    public async GetMyUserEmail(): Promise<UserEmail> {
        return await this.request(
            c => c.getMyUserEmail,
            new Empty(),
            f => f,
        );
    }

    public async SaveMyUserEmail(email: string): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setEmail(email);
        return await this.request(
            c => c.changeMyUserEmail,
            req,
            f => f,
        );
    }

    public async RemoveMyUserPhone(): Promise<Empty> {
        return await this.request(
            c => c.removeMyUserPhone,
            new Empty(),
            f => f,
        );
    }

    private async getMyzitadelPermissions(): Promise<any> {
        return await this.request(
            c => c.getMyZitadelPermissions,
            new Empty(),
            f => f,
        );
    }

    public GetMyzitadelPermissions(): Observable<any> {
        return from(this.getMyzitadelPermissions());
    }

    public hasRoles(userRoles: string[], requestedRoles: string[], each: boolean = false): boolean {
        return each ?
            requestedRoles.every(role => userRoles.includes(role)) :
            requestedRoles.findIndex(role => {
                return userRoles.findIndex(i => i.includes(role)) > -1;
                // return userRoles.includes(role);
            }) > -1;
    }

    public async GetMyUserPhone(): Promise<UserPhone> {
        // return this.grpcClient.auth.getMyUserPhone(new Empty());
        return await this.request(
            c => c.getMyUserPhone,
            new Empty(),
            (f: UserPhone) => f,
        );
    }

    public async SaveMyUserPhone(phone: string): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setPhone(phone);
        return await this.request(
            c => c.changeMyUserPhone,
            req,
            f => f,
        );
    }

    public async GetMyUserAddress(): Promise<UserAddress> {
        return await this.request(
            c => c.getMyUserAddress,
            new Empty(),
            f => f,
        );
    }

    public async ResendEmailVerification(): Promise<Empty> {
        const req = new Empty();
        return await this.request(
            c => c.resendMyEmailVerificationMail,
            req,
            f => f,
        );
    }

    public async ResendPhoneVerification(): Promise<Empty> {
        const req = new Empty();
        return await this.request(
            c => c.resendMyPhoneVerificationCode,
            req,
            f => f,
        );
    }

    public async ChangeMyPassword(oldPassword: string, newPassword: string): Promise<Empty> {
        const req = new PasswordChange();
        req.setOldPassword(oldPassword);
        req.setNewPassword(newPassword);
        return await this.request(
            c => c.changeMyPassword,
            req,
            f => f,
        );
    }

    public async AddMfaOTP(): Promise<MfaOtpResponse> {
        return await this.request(
            c => c.addMfaOTP,
            new Empty(),
            f => f,
        );
    }

    public async RemoveMfaOTP(): Promise<Empty> {
        return await this.request(
            c => c.removeMfaOTP,
            new Empty(),
            f => f,
        );
    }

    public async VerifyMfaOTP(code: string): Promise<Empty> {
        const req = new VerifyMfaOtp();
        req.setCode(code);
        return await this.request(
            c => c.verifyMfaOTP,
            req,
            f => f,
        );
    }

    public async VerifyMyUserPhone(code: string): Promise<Empty> {
        const req = new VerifyUserPhoneRequest();
        req.setCode(code);
        return await this.request(
            c => c.verifyMyUserPhone,
            req,
            f => f,
        );
    }

    public async SaveMyUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return await this.request(
            c => c.updateMyUserAddress,
            req,
            f => f,
        );
    }

    public async GetMyUserChanges(limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangesRequest();
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return await this.request(
            c => c.getMyUserChanges,
            req,
            f => f,
        );
    }

    public isAllowed(roles: string[], each: boolean = false): Observable<boolean> {
        if (roles && roles.length > 0) {
            if (this._roleCache.length > 0) {
                return of(this.hasRoles(this._roleCache, roles));
            }

            return this.GetMyzitadelPermissions().pipe(
                switchMap(response => {
                    let userRoles = [];
                    if (response.toObject().permissionsList) {
                        userRoles = response.toObject().permissionsList;
                    } else {
                        userRoles = ['user.resourceowner'];
                    }
                    this._roleCache = userRoles;
                    return of(this.hasRoles(userRoles, roles, each));
                }),
                catchError((err) => {
                    return of(false);
                }),
            );
        } else {
            return of(false);
        }
    }

}

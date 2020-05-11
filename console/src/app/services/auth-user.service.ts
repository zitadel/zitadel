import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';
import { from, Observable, of } from 'rxjs';
import { switchMap } from 'rxjs/operators';

import { AuthServicePromiseClient } from '../proto/generated/auth_grpc_web_pb';
import {
    GrantSearchQuery,
    GrantSearchRequest,
    GrantSearchResponse,
    MfaOtpResponse,
    MultiFactors,
    MyProjectOrgSearchQuery,
    MyProjectOrgSearchRequest,
    MyProjectOrgSearchResponse,
    PasswordChange,
    PasswordRequest,
    UpdateUserAddressRequest,
    UpdateUserEmailRequest,
    UpdateUserPhoneRequest,
    UpdateUserProfileRequest,
    UserAddress,
    UserEmail,
    UserID,
    UserPhone,
    UserProfile,
    UserSessionViews,
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

    public async GetMyUserProfile(): Promise<UserProfile> {
        return await this.request(
            c => c.getMyUserProfile,
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

    public async SaveMyUserProfile(profile: UserProfile.AsObject): Promise<UserProfile> {
        const req = new UpdateUserProfileRequest();
        req.setFirstName(profile.firstName);
        req.setLastName(profile.lastName);
        req.setNickName(profile.nickName);
        req.setDisplayName(profile.displayName);
        req.setPreferredLanguage(profile.preferredLanguage);
        req.setGender(profile.gender);
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

    public async SaveMyUserEmail(email: UserEmail.AsObject): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setEmail(email.email);
        return await this.request(
            c => c.changeMyUserEmail,
            req,
            f => f,
        );
    }

    private async getMyCitadelPermissions(): Promise<any> {
        return await this.request(
            c => c.getMyCitadelPermissions,
            new Empty(),
            f => f,
        );
    }

    public GetMyCitadelPermissions(): Observable<any> {
        return from(this.getMyCitadelPermissions());
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

    public async SaveMyUserPhone(phone: UserPhone.AsObject): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setPhone(phone.phone);
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

    public async ResendEmailVerification(id: string): Promise<Empty> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.resendMyEmailVerificationMail,
            req,
            f => f,
        );
    }

    public async ResendPhoneVerification(id: string): Promise<Empty> {
        const req = new UserID();
        req.setId(id);
        return await this.request(
            c => c.resendMyPhoneVerificationCode,
            req,
            f => f,
        );
    }

    public async SetMyPassword(id: string): Promise<Empty> {
        const req = new PasswordRequest();
        req.setPassword(id);

        return await this.request(
            c => c.setMyPassword,
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

    public async VerifyMfaOTP(code: string): Promise<MfaOtpResponse> {
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

    public async SearchGrant(
        limit: number,
        offset: number,
        queryList?: GrantSearchQuery[],
    ): Promise<GrantSearchResponse> {
        const req = new GrantSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchGrant,
            req,
            f => f,
        );
    }

    public isAllowed(roles: string[], each: boolean = false): Observable<boolean> {
        if (roles && roles.length > 0) {
            if (this._roleCache.length > 0) {
                return of(this.hasRoles(this._roleCache, roles));
            }

            return this.GetMyCitadelPermissions().pipe(
                switchMap(response => {
                    const userRoles = response.toObject().permissionsList;
                    this._roleCache = userRoles;
                    return of(this.hasRoles(userRoles, roles, each));
                }),
            );
        } else {
            return of(false);
        }
    }

}

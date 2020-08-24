import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

import {
    Changes,
    ChangesRequest,
    Gender,
    MfaOtpResponse,
    MultiFactors,
    MyPermissions,
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
import { GrpcService } from './grpc.service';


@Injectable({
    providedIn: 'root',
})
export class AuthService {
    constructor(private readonly grpcService: GrpcService) { }

    public async GetMyUserProfile(): Promise<UserProfileView> {
        return this.grpcService.auth.getMyUserProfile(new Empty());
    }

    public async GetMyPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        return this.grpcService.auth.getMyPasswordComplexityPolicy(
            new Empty(),
        );
    }

    public async GetMyUser(): Promise<UserView> {
        return this.grpcService.auth.getMyUser(
            new Empty(),
        );
    }

    public async GetMyMfas(): Promise<MultiFactors> {
        return this.grpcService.auth.getMyMfas(
            new Empty(),
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

        return this.grpcService.auth.searchMyProjectOrgs(req);
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
        return this.grpcService.auth.updateMyUserProfile(req);
    }

    public async getMyUserSessions(): Promise<UserSessionViews> {
        return this.grpcService.auth.getMyUserSessions(
            new Empty(),
        );
    }

    public async GetMyUserEmail(): Promise<UserEmail> {
        return this.grpcService.auth.getMyUserEmail(
            new Empty(),
        );
    }

    public async SaveMyUserEmail(email: string): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setEmail(email);
        return this.grpcService.auth.changeMyUserEmail(req);
    }

    public async RemoveMyUserPhone(): Promise<Empty> {
        return this.grpcService.auth.removeMyUserPhone(
            new Empty(),
        );
    }

    public async GetMyzitadelPermissions(): Promise<MyPermissions> {
        return this.grpcService.auth.getMyZitadelPermissions(
            new Empty(),
        );
    }

    public async GetMyUserPhone(): Promise<UserPhone> {
        // return this.grpcClient.auth.getMyUserPhone(new Empty());
        return this.grpcService.auth.getMyUserPhone(
            new Empty(),
        );
    }

    public async SaveMyUserPhone(phone: string): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setPhone(phone);
        return this.grpcService.auth.changeMyUserPhone(req);
    }

    public async GetMyUserAddress(): Promise<UserAddress> {
        return this.grpcService.auth.getMyUserAddress(
            new Empty(),
        );
    }

    public async ResendEmailVerification(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.auth.resendMyEmailVerificationMail(req);
    }

    public async ResendPhoneVerification(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.auth.resendMyPhoneVerificationCode(req);
    }

    public async ChangeMyPassword(oldPassword: string, newPassword: string): Promise<Empty> {
        const req = new PasswordChange();
        req.setOldPassword(oldPassword);
        req.setNewPassword(newPassword);
        return this.grpcService.auth.changeMyPassword(req);
    }

    public async AddMfaOTP(): Promise<MfaOtpResponse> {
        return this.grpcService.auth.addMfaOTP(
            new Empty(),
        );
    }

    public async RemoveMfaOTP(): Promise<Empty> {
        return this.grpcService.auth.removeMfaOTP(
            new Empty(),
        );
    }

    public async VerifyMfaOTP(code: string): Promise<Empty> {
        const req = new VerifyMfaOtp();
        req.setCode(code);
        return this.grpcService.auth.verifyMfaOTP(req);
    }

    public async VerifyMyUserPhone(code: string): Promise<Empty> {
        const req = new VerifyUserPhoneRequest();
        req.setCode(code);
        return this.grpcService.auth.verifyMyUserPhone(req);
    }

    public async SaveMyUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return this.grpcService.auth.updateMyUserAddress(req);
    }

    public async GetMyUserChanges(limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangesRequest();
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return this.grpcService.auth.getMyUserChanges(req);
    }
}

import { Injectable } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { BehaviorSubject, from, merge, Observable, of, Subject } from 'rxjs';
import { catchError, filter, finalize, first, map, mergeMap, switchMap, take, timeout } from 'rxjs/operators';

import {
    GetMyEmailRequest,
    GetMyEmailResponse,
    GetMyPasswordComplexityPolicyRequest,
    GetMyPasswordComplexityPolicyResponse,
    GetMyProfileRequest,
    GetMyProfileResponse,
    GetMyUserRequest,
    GetMyUserResponse,
    ListMyMultiFactorsRequest,
    ListMyMultiFactorsResponse,
    ListMyProjectOrgsRequest,
    ListMyProjectOrgsResponse,
    ListMyUserGrantsRequest,
    ListMyUserSessionsRequest,
    ListMyUserSessionsResponse,
    SetMyEmailResponse,
    UpdateMyProfileRequest,
    SetMyEmailRequest,
    UpdateMyProfileResponse,
    ListMyUserGrantsResponse,
    ResendMyEmailVerificationRequest,
    ListMyZitadelPermissionsResponse,
    ResendMyEmailVerificationResponse,
    RemoveMyPhoneResponse,
    RemoveMyPhoneRequest,
    ListMyZitadelPermissionsRequest,
    GetMyPhoneResponse,
    GetMyPhoneRequest,
    SetMyPhoneRequest,
    SetMyPhoneResponse,
    GetMyAddressResponse,
    GetMyAddressRequest,
    RemoveMyMultiFactorOTPRequest,
    ResendMyPhoneVerificationRequest,
    ResendMyPhoneVerificationResponse,
    UpdateMyPasswordRequest,
    UpdateMyPasswordResponse,
    RemoveMyLinkedIDPRequest,
    VerifyMyMultiFactorOTPRequest,
    RemoveMyLinkedIDPResponse,
    ListMyLinkedIDPsResponse,
    ListMyLinkedIDPsRequest,
    AddMyMultiFactorU2FResponse,
    AddMyMultiFactorU2FRequest,
    RemoveMyMultiFactorU2FRequest,
    RemoveMyMultiFactorU2FResponse,
    VerifyMyMultiFactorU2FResponse,
    VerifyMyMultiFactorU2FRequest,
    ListMyPasswordlessResponse,
    ListMyPasswordlessRequest,
    AddMyMultiFactorOTPResponse,
    AddMyMultiFactorOTPRequest,
    AddMyPasswordlessRequest,
    AddMyPasswordlessResponse,
    RemoveMyPasswordlessRequest,
    VerifyMyPasswordlessRequest,
    VerifyMyPhoneRequest,
    VerifyMyPhoneResponse,
    UpdateMyAddressRequest,
    UpdateMyAddressResponse,
    ListMyUserChangesRequest,
    ListMyUserChangesResponse
} from '../proto/generated/zitadel/auth_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { Org, OrgQuery } from '../proto/generated/zitadel/org_pb';
import { Profile, Gender, WebAuthNVerification, Address } from '../proto/generated/zitadel/user_pb';
import { GrpcService } from './grpc.service';
import { StorageKey, StorageService } from './storage.service';


@Injectable({
    providedIn: 'root',
})
export class GrpcAuthService {
    private _activeOrgChanged: Subject<Org.AsObject> = new Subject();
    public user!: Observable<Profile.AsObject | undefined>;
    private zitadelPermissions: BehaviorSubject<string[]> = new BehaviorSubject(['user.resourceowner']);
    public readonly fetchedZitadelPermissions: BehaviorSubject<boolean> = new BehaviorSubject(false as boolean);

    private cachedOrgs: Org.AsObject[] = [];

    constructor(
        private readonly grpcService: GrpcService,
        private oauthService: OAuthService,
        private storage: StorageService,
    ) {
        this.user = merge(
            of(this.oauthService.getAccessToken()).pipe(
                filter(token => token ? true : false),
            ),
            this.oauthService.events.pipe(
                filter(e => e.type === 'token_received'),
                timeout(this.oauthService.waitForTokenInMsec || 0),
                catchError(_ => of(null)), // timeout is not an error
                map(_ => this.oauthService.getAccessToken()),
            ),
        ).pipe(
            take(1),
            mergeMap(() => {
                return from(this.getMyProfile().then(resp => {
                    const user = resp.toObject();
                    const profile = user.profile;
                    if (profile) {
                        return profile;
                    } else {
                        return undefined;
                    }
                }));
            }),
            finalize(() => {
                this.loadPermissions();
            }),
        );

        this.activeOrgChanged.subscribe(() => {
            this.loadPermissions();
        });
    }

    public async getActiveOrg(id?: string): Promise<Org.AsObject> {
        if (id) {
            const org = this.storage.getItem<Org.AsObject>(StorageKey.organization);
            if (org && this.cachedOrgs.find(tmp => tmp.id === org.id)) {
                return org;
            }
            return Promise.reject(new Error('no cached org'));
        } else {
            let orgs = this.cachedOrgs;
            if (orgs.length === 0) {
                orgs = (await this.listMyProjectOrgs(10, 0)).toObject().resultList;
                this.cachedOrgs = orgs;
            }

            const org = this.storage.getItem<Org.AsObject>(StorageKey.organization);
            if (org && orgs.find(tmp => tmp.id === org.id)) {
                return org;
            }

            if (orgs.length === 0) {
                return Promise.reject(new Error('No organizations found!'));
            }
            const orgToSet = orgs.find(element => element.id !== '0' && element.name !== '');

            if (orgToSet) {
                this.setActiveOrg(orgToSet);
                return Promise.resolve(orgToSet);
            }
            return Promise.resolve(orgs[0]);
        }
    }

    public get activeOrgChanged(): Observable<Org.AsObject> {
        return this._activeOrgChanged;
    }

    public setActiveOrg(org: Org.AsObject): void {
        this.storage.setItem(StorageKey.organization, org);
        this._activeOrgChanged.next(org);
    }

    private loadPermissions(): void {
        merge([
            // this.authenticationChanged,
            this.activeOrgChanged.pipe(map(org => !!org)),
        ]).pipe(
            first(),
            switchMap(() => from(this.listMyZitadelPermissions())),
            map(rolesResp => rolesResp.toObject().resultList),
            catchError(_ => {
                return of([]);
            }),
            finalize(() => {
                this.fetchedZitadelPermissions.next(true);
            }),
        ).subscribe(roles => {
            this.zitadelPermissions.next(roles);
        });
    }

    /**
     * returns true if user has one of the provided roles
     * @param roles roles of the user
     */
    public isAllowed(roles: string[] | RegExp[]): Observable<boolean> {
        if (roles && roles.length > 0) {
            return this.zitadelPermissions.pipe(switchMap(zroles => of(this.hasRoles(zroles, roles))));
        } else {
            return of(false);
        }
    }

    /**
     * returns true if user has one of the provided roles
     * @param userRoles roles of the user
     * @param requestedRoles required roles for accessing the respective component
     */
    public hasRoles(userRoles: string[], requestedRoles: string[] | RegExp[]): boolean {
        return requestedRoles.findIndex((regexp: any) => {
            return userRoles.findIndex(role => {
                return new RegExp(regexp).test(role);
            }) > -1;
        }) > -1;
    }

    public getMyProfile(): Promise<GetMyProfileResponse> {
        return this.grpcService.auth.getMyProfile(new GetMyProfileRequest());
    }

    public getMyPasswordComplexityPolicy(): Promise<GetMyPasswordComplexityPolicyResponse> {
        return this.grpcService.auth.getMyPasswordComplexityPolicy(
            new GetMyPasswordComplexityPolicyRequest(),
        );
    }

    public getMyUser(): Promise<GetMyUserResponse> {
        return this.grpcService.auth.getMyUser(
            new GetMyUserRequest(),
        );
    }

    public listMyMultiFactors(): Promise<ListMyMultiFactorsResponse> {
        return this.grpcService.auth.listMyMultiFactors(
            new ListMyMultiFactorsRequest(),
        );
    }

    public listMyProjectOrgs(
        limit: number,
        offset: number,
        queryList?: OrgQuery[],
    ): Promise<ListMyProjectOrgsResponse> {
        const req = new ListMyProjectOrgsRequest();
        const metadata = new ListQuery();
        if (offset) {
            metadata.setOffset(offset);
        }
        if (limit) {
            metadata.setLimit(limit);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }

        return this.grpcService.auth.listMyProjectOrgs(req);
    }

    public updateMyProfile(
        firstName?: string,
        lastName?: string,
        nickName?: string,
        preferredLanguage?: string,
        gender?: Gender,
    ): Promise<UpdateMyProfileResponse> {
        const req = new UpdateMyProfileRequest();
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
        return this.grpcService.auth.updateMyProfile(req);
    }

    public get zitadelPermissionsChanged(): Observable<string[]> {
        return this.zitadelPermissions;
    }

    public listMyUserSessions(): Promise<ListMyUserSessionsResponse> {
        return this.grpcService.auth.listMyUserSessions(
            new ListMyUserSessionsRequest(),
        );
    }

    public listMyUserGrants(limit?: number, offset?: number, queryList?: ListQuery[]): Promise<ListMyUserGrantsResponse> {
        const req = new ListMyUserGrantsRequest();
        const query = new ListQuery();
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        req.setQuery(query);
        return this.grpcService.auth.listMyUserGrants(req);
    }

    public getMyEmail(): Promise<GetMyEmailResponse> {
        return this.grpcService.auth.getMyEmail(
            new GetMyEmailRequest(),
        );
    }

    public setMyEmail(email: string): Promise<SetMyEmailResponse> {
        const req = new SetMyEmailRequest();
        req.setEmail(email);
        return this.grpcService.auth.setMyEmail(req);
    }

    public resendMyEmailVerification(): Promise<ResendMyEmailVerificationResponse> {
        return this.grpcService.auth.resendMyEmailVerification(
            new ResendMyEmailVerificationRequest(),
        );
    }

    public removeMyPhone(): Promise<RemoveMyPhoneResponse> {
        return this.grpcService.auth.removeMyPhone(
            new RemoveMyPhoneRequest(),
        );
    }

    public listMyZitadelPermissions(): Promise<ListMyZitadelPermissionsResponse> {
        return this.grpcService.auth.listMyZitadelPermissions(
            new ListMyZitadelPermissionsRequest(),
        );
    }

    public getMyPhone(): Promise<GetMyPhoneResponse> {
        return this.grpcService.auth.getMyPhone(
            new GetMyPhoneRequest(),
        );
    }

    public setMyPhone(phone: string): Promise<SetMyPhoneResponse> {
        const req = new SetMyPhoneRequest();
        req.setPhone(phone);
        return this.grpcService.auth.setMyPhone(req);
    }

    public getMyAddress(): Promise<GetMyAddressResponse> {
        return this.grpcService.auth.getMyAddress(
            new GetMyAddressRequest(),
        );
    }

    public resendMyPhoneVerification(): Promise<ResendMyPhoneVerificationResponse> {
        const req = new ResendMyPhoneVerificationRequest();
        return this.grpcService.auth.resendMyPhoneVerification(req);
    }

    public updateMyPassword(oldPassword: string, newPassword: string): Promise<UpdateMyPasswordResponse> {
        const req = new UpdateMyPasswordRequest();
        req.setOldPassword(oldPassword);
        req.setNewPassword(newPassword);
        return this.grpcService.auth.updateMyPassword(req);
    }

    public removeMyLinkedIDP(
        externalUserId: string,
        idpId: string,
    ): Promise<RemoveMyLinkedIDPResponse> {
        const req = new RemoveMyLinkedIDPRequest();
        req.setLinkedUserId(externalUserId);
        req.setIdpId(idpId);
        return this.grpcService.auth.removeMyLinkedIDP(req);
    }

    public listMyLinkedIDPs(
        limit: number,
        offset: number,
    ): Promise<ListMyLinkedIDPsResponse> {
        const req = new ListMyLinkedIDPsRequest();
        const metadata = new ListQuery();
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        req.setMetaData(metadata);
        return this.grpcService.auth.listMyLinkedIDPs(req);
    }

    public addMyMultiFactorOTP(): Promise<AddMyMultiFactorOTPResponse> {
        return this.grpcService.auth.addMyMultiFactorOTP(
            new AddMyMultiFactorOTPRequest(),
        );
    }

    public addMyMultiFactorU2F(): Promise<AddMyMultiFactorU2FResponse> {
        return this.grpcService.auth.addMyMultiFactorU2F(
            new AddMyMultiFactorU2FRequest(),
        );
    }

    public removeMyMultiFactorU2F(tokenId: string): Promise<RemoveMyMultiFactorU2FResponse> {
        const req = new RemoveMyMultiFactorU2FRequest();
        req.setTokenId(tokenId);
        return this.grpcService.auth.removeMyMultiFactorU2F(req);
    }

    public verifyMyMultiFactorU2F(credential: string, tokenname: string): Promise<VerifyMyMultiFactorU2FResponse> {
        const req = new VerifyMyMultiFactorU2FRequest();
        const verification = new WebAuthNVerification();
        verification.setPublicKeyCredential(credential);
        verification.setTokenName(tokenname);
        req.setVerification(verification);

        return this.grpcService.auth.verifyMyMultiFactorU2F(req);
    }

    public listMyPasswordless(): Promise<ListMyPasswordlessResponse> {
        return this.grpcService.auth.listMyPasswordless(
            new ListMyPasswordlessRequest(),
        );
    }

    public addMyPasswordless(): Promise<AddMyPasswordlessResponse> {
        return this.grpcService.auth.addMyPasswordless(
            new AddMyPasswordlessRequest(),
        );
    }

    public removeMyPasswordless(tokenId: string): Promise<Empty> {
        const req = new RemoveMyPasswordlessRequest();
        req.setTokenId(tokenId);
        return this.grpcService.auth.removeMyPasswordless(req);
    }

    public verifyMyPasswordless(credential: string, tokenname: string): Promise<Empty> {
        const req = new VerifyMyPasswordlessRequest();
        const verification = new WebAuthNVerification();
        verification.setTokenName(tokenname);
        verification.setPublicKeyCredential(credential);
        req.setVerification(verification);

        return this.grpcService.auth.verifyMyPasswordless(
            req,
        );
    }

    public removeMyMultiFactorOTP(): Promise<Empty> {
        return this.grpcService.auth.removeMyMultiFactorOTP(
            new RemoveMyMultiFactorOTPRequest(),
        );
    }

    public verifyMyMultiFactorOTP(code: string): Promise<Empty> {
        const req = new VerifyMyMultiFactorOTPRequest();
        req.setCode(code);
        return this.grpcService.auth.verifyMyMultiFactorOTP(req);
    }

    public verifyMyPhone(code: string): Promise<VerifyMyPhoneResponse> {
        const req = new VerifyMyPhoneRequest();
        req.setCode(code);
        return this.grpcService.auth.verifyMyPhone(req);
    }

    public updateMyUserAddress(address: Address.AsObject): Promise<UpdateMyAddressResponse> {
        const req = new UpdateMyAddressRequest();
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return this.grpcService.auth.updateMyUserAddress(req);
    }

    public listMyUserChanges(limit: number, offset: number): Promise<ListMyUserChangesResponse> {
        const req = new ListMyUserChangesRequest();
        const query = new ListQuery();
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        req.setQuery(query);
        return this.grpcService.auth.listMyUserChanges(req);
    }
}

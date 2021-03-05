import { Injectable } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, merge, Observable, of, Subject } from 'rxjs';
import { catchError, filter, finalize, first, map, mergeMap, switchMap, take, timeout } from 'rxjs/operators';

import {
    AddMyMultiFactorOTPRequest,
    AddMyMultiFactorOTPResponse,
    AddMyMultiFactorU2FRequest,
    AddMyMultiFactorU2FResponse,
    AddMyPasswordlessRequest,
    AddMyPasswordlessResponse,
    GetMyEmailRequest,
    GetMyEmailResponse,
    GetMyPasswordComplexityPolicyRequest,
    GetMyPasswordComplexityPolicyResponse,
    GetMyPhoneRequest,
    GetMyPhoneResponse,
    GetMyProfileRequest,
    GetMyProfileResponse,
    GetMyUserRequest,
    GetMyUserResponse,
    ListMyLinkedIDPsRequest,
    ListMyLinkedIDPsResponse,
    ListMyMultiFactorsRequest,
    ListMyMultiFactorsResponse,
    ListMyPasswordlessRequest,
    ListMyPasswordlessResponse,
    ListMyProjectOrgsRequest,
    ListMyProjectOrgsResponse,
    ListMyUserChangesRequest,
    ListMyUserChangesResponse,
    ListMyUserGrantsRequest,
    ListMyUserGrantsResponse,
    ListMyUserSessionsRequest,
    ListMyUserSessionsResponse,
    ListMyZitadelPermissionsRequest,
    ListMyZitadelPermissionsResponse,
    RemoveMyLinkedIDPRequest,
    RemoveMyLinkedIDPResponse,
    RemoveMyMultiFactorOTPRequest,
    RemoveMyMultiFactorOTPResponse,
    RemoveMyMultiFactorU2FRequest,
    RemoveMyMultiFactorU2FResponse,
    RemoveMyPasswordlessRequest,
    RemoveMyPasswordlessResponse,
    RemoveMyPhoneRequest,
    RemoveMyPhoneResponse,
    ResendMyEmailVerificationRequest,
    ResendMyEmailVerificationResponse,
    ResendMyPhoneVerificationRequest,
    ResendMyPhoneVerificationResponse,
    SetMyEmailRequest,
    SetMyEmailResponse,
    SetMyPhoneRequest,
    SetMyPhoneResponse,
    UpdateMyPasswordRequest,
    UpdateMyPasswordResponse,
    UpdateMyProfileRequest,
    UpdateMyProfileResponse,
    VerifyMyMultiFactorOTPRequest,
    VerifyMyMultiFactorOTPResponse,
    VerifyMyMultiFactorU2FRequest,
    VerifyMyMultiFactorU2FResponse,
    VerifyMyPasswordlessRequest,
    VerifyMyPasswordlessResponse,
    VerifyMyPhoneRequest,
    VerifyMyPhoneResponse,
} from '../proto/generated/zitadel/auth_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { Org, OrgQuery } from '../proto/generated/zitadel/org_pb';
import { Gender, Profile, WebAuthNVerification } from '../proto/generated/zitadel/user_pb';
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
                    const profile = resp.profile;
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
                orgs = (await this.listMyProjectOrgs(10, 0)).resultList;
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
            map(rolesResp => rolesResp.resultList),
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

    public getMyProfile(): Promise<GetMyProfileResponse.AsObject> {
        return this.grpcService.auth.getMyProfile(new GetMyProfileRequest(), null).then(resp => resp.toObject());
    }

    public getMyPasswordComplexityPolicy(): Promise<GetMyPasswordComplexityPolicyResponse.AsObject> {
        return this.grpcService.auth.getMyPasswordComplexityPolicy(
            new GetMyPasswordComplexityPolicyRequest(), null
        ).then(resp => resp.toObject());
    }

    public getMyUser(): Promise<GetMyUserResponse.AsObject> {
        return this.grpcService.auth.getMyUser(
            new GetMyUserRequest(), null
        ).then(resp => resp.toObject());
    }

    public listMyMultiFactors(): Promise<ListMyMultiFactorsResponse.AsObject> {
        return this.grpcService.auth.listMyMultiFactors(
            new ListMyMultiFactorsRequest(), null
        ).then(resp => resp.toObject());
    }

    public listMyProjectOrgs(
        limit: number,
        offset: number,
        queryList?: OrgQuery[],
    ): Promise<ListMyProjectOrgsResponse.AsObject> {
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

        return this.grpcService.auth.listMyProjectOrgs(req, null).then(resp => resp.toObject());
    }

    public updateMyProfile(
        firstName?: string,
        lastName?: string,
        nickName?: string,
        preferredLanguage?: string,
        gender?: Gender,
    ): Promise<UpdateMyProfileResponse.AsObject> {
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
        return this.grpcService.auth.updateMyProfile(req, null).then(resp => resp.toObject());
    }

    public get zitadelPermissionsChanged(): Observable<string[]> {
        return this.zitadelPermissions;
    }

    public listMyUserSessions(): Promise<ListMyUserSessionsResponse.AsObject> {
        const req = new ListMyUserSessionsRequest();
        return this.grpcService.auth.listMyUserSessions(req, null).then(resp => resp.toObject());
    }

    public listMyUserGrants(limit?: number, offset?: number, queryList?: ListQuery[]): Promise<ListMyUserGrantsResponse.AsObject> {
        const req = new ListMyUserGrantsRequest();
        const query = new ListQuery();
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        req.setQuery(query);
        return this.grpcService.auth.listMyUserGrants(req, null).then(resp => resp.toObject());
    }

    public getMyEmail(): Promise<GetMyEmailResponse.AsObject> {
        const req = new GetMyEmailRequest();
        return this.grpcService.auth.getMyEmail(req, null).then(resp => resp.toObject());
    }

    public setMyEmail(email: string): Promise<SetMyEmailResponse.AsObject> {
        const req = new SetMyEmailRequest();
        req.setEmail(email);
        return this.grpcService.auth.setMyEmail(req, null).then(resp => resp.toObject());
    }

    public resendMyEmailVerification(): Promise<ResendMyEmailVerificationResponse.AsObject> {
        const req = new ResendMyEmailVerificationRequest();
        return this.grpcService.auth.resendMyEmailVerification(req, null).then(resp => resp.toObject());
    }

    public removeMyPhone(): Promise<RemoveMyPhoneResponse.AsObject> {
        return this.grpcService.auth.removeMyPhone(
            new RemoveMyPhoneRequest(), null
        ).then(resp => resp.toObject());
    }

    public listMyZitadelPermissions(): Promise<ListMyZitadelPermissionsResponse.AsObject> {
        return this.grpcService.auth.listMyZitadelPermissions(
            new ListMyZitadelPermissionsRequest(), null
        ).then(resp => resp.toObject());
    }

    public getMyPhone(): Promise<GetMyPhoneResponse.AsObject> {
        return this.grpcService.auth.getMyPhone(
            new GetMyPhoneRequest(), null
        ).then(resp => resp.toObject());
    }

    public setMyPhone(phone: string): Promise<SetMyPhoneResponse.AsObject> {
        const req = new SetMyPhoneRequest();
        req.setPhone(phone);
        return this.grpcService.auth.setMyPhone(req, null).then(resp => resp.toObject());
    }

    public resendMyPhoneVerification(): Promise<ResendMyPhoneVerificationResponse.AsObject> {
        const req = new ResendMyPhoneVerificationRequest();
        return this.grpcService.auth.resendMyPhoneVerification(req, null).then(resp => resp.toObject());
    }

    public updateMyPassword(oldPassword: string, newPassword: string): Promise<UpdateMyPasswordResponse.AsObject> {
        const req = new UpdateMyPasswordRequest();
        req.setOldPassword(oldPassword);
        req.setNewPassword(newPassword);
        return this.grpcService.auth.updateMyPassword(req, null).then(resp => resp.toObject());
    }

    public removeMyLinkedIDP(
        externalUserId: string,
        idpId: string,
    ): Promise<RemoveMyLinkedIDPResponse.AsObject> {
        const req = new RemoveMyLinkedIDPRequest();
        req.setLinkedUserId(externalUserId);
        req.setIdpId(idpId);
        return this.grpcService.auth.removeMyLinkedIDP(req, null).then(resp => resp.toObject());
    }

    public listMyLinkedIDPs(
        limit: number,
        offset: number,
    ): Promise<ListMyLinkedIDPsResponse.AsObject> {
        const req = new ListMyLinkedIDPsRequest();
        const metadata = new ListQuery();
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        req.setQuery(metadata);
        return this.grpcService.auth.listMyLinkedIDPs(req, null).then(resp => resp.toObject());
    }

    public addMyMultiFactorOTP(): Promise<AddMyMultiFactorOTPResponse.AsObject> {
        return this.grpcService.auth.addMyMultiFactorOTP(
            new AddMyMultiFactorOTPRequest(), null
        ).then(resp => resp.toObject());
    }

    public addMyMultiFactorU2F(): Promise<AddMyMultiFactorU2FResponse.AsObject> {
        return this.grpcService.auth.addMyMultiFactorU2F(
            new AddMyMultiFactorU2FRequest(), null
        ).then(resp => resp.toObject());
    }

    public removeMyMultiFactorU2F(tokenId: string): Promise<RemoveMyMultiFactorU2FResponse.AsObject> {
        const req = new RemoveMyMultiFactorU2FRequest();
        req.setTokenId(tokenId);
        return this.grpcService.auth.removeMyMultiFactorU2F(req, null).then(resp => resp.toObject());
    }

    public verifyMyMultiFactorU2F(credential: string, tokenname: string): Promise<VerifyMyMultiFactorU2FResponse.AsObject> {
        const req = new VerifyMyMultiFactorU2FRequest();
        const verification = new WebAuthNVerification();
        verification.setPublicKeyCredential(credential);
        verification.setTokenName(tokenname);
        req.setVerification(verification);

        return this.grpcService.auth.verifyMyMultiFactorU2F(req, null).then(resp => resp.toObject());
    }

    public listMyPasswordless(): Promise<ListMyPasswordlessResponse.AsObject> {
        return this.grpcService.auth.listMyPasswordless(
            new ListMyPasswordlessRequest(), null
        ).then(resp => resp.toObject());
    }

    public addMyPasswordless(): Promise<AddMyPasswordlessResponse.AsObject> {
        return this.grpcService.auth.addMyPasswordless(
            new AddMyPasswordlessRequest(), null
        ).then(resp => resp.toObject());
    }

    public removeMyPasswordless(tokenId: string): Promise<RemoveMyPasswordlessResponse.AsObject> {
        const req = new RemoveMyPasswordlessRequest();
        req.setTokenId(tokenId);
        return this.grpcService.auth.removeMyPasswordless(req, null).then(resp => resp.toObject());
    }

    public verifyMyPasswordless(credential: string, tokenname: string): Promise<VerifyMyPasswordlessResponse.AsObject> {
        const req = new VerifyMyPasswordlessRequest();
        const verification = new WebAuthNVerification();
        verification.setTokenName(tokenname);
        verification.setPublicKeyCredential(credential);
        req.setVerification(verification);

        return this.grpcService.auth.verifyMyPasswordless(
            req, null
        ).then(resp => resp.toObject());
    }

    public removeMyMultiFactorOTP(): Promise<RemoveMyMultiFactorOTPResponse.AsObject> {
        return this.grpcService.auth.removeMyMultiFactorOTP(
            new RemoveMyMultiFactorOTPRequest(), null
        ).then(resp => resp.toObject());
    }

    public verifyMyMultiFactorOTP(code: string): Promise<VerifyMyMultiFactorOTPResponse.AsObject> {
        const req = new VerifyMyMultiFactorOTPRequest();
        req.setCode(code);
        return this.grpcService.auth.verifyMyMultiFactorOTP(req, null).then(resp => resp.toObject());
    }

    public verifyMyPhone(code: string): Promise<VerifyMyPhoneResponse.AsObject> {
        const req = new VerifyMyPhoneRequest();
        req.setCode(code);
        return this.grpcService.auth.verifyMyPhone(req, null).then(resp => resp.toObject());
    }

    public listMyUserChanges(limit: number, offset: number): Promise<ListMyUserChangesResponse.AsObject> {
        const req = new ListMyUserChangesRequest();
        const query = new ListQuery();
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        req.setQuery(query);
        return this.grpcService.auth.listMyUserChanges(req, null).then(resp => resp.toObject());
    }
}

import { Injectable } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { BehaviorSubject, from, merge, Observable, of, Subject } from 'rxjs';
import { catchError, filter, finalize, first, map, mergeMap, switchMap, take, timeout } from 'rxjs/operators';

import {
    Changes,
    ChangesRequest,
    ExternalIDPRemoveRequest,
    ExternalIDPSearchRequest,
    ExternalIDPSearchResponse,
    Gender,
    MfaOtpResponse,
    MultiFactors,
    MyPermissions,
    MyProjectOrgSearchQuery,
    MyProjectOrgSearchRequest,
    MyProjectOrgSearchResponse,
    Org,
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
    VerifyWebAuthN,
    WebAuthNResponse,
    WebAuthNTokenID,
    WebAuthNTokens,
} from '../proto/generated/auth_pb';
import { GrpcService } from './grpc.service';
import { StorageKey, StorageService } from './storage.service';


@Injectable({
    providedIn: 'root',
})
export class GrpcAuthService {
    private _activeOrgChanged: Subject<Org.AsObject> = new Subject();
    public user!: Observable<UserProfileView.AsObject>;
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
                return from(this.GetMyUserProfile().then(userprofile => userprofile.toObject()));
            }),
            finalize(() => {
                this.loadPermissions();
            }),
        );

        this.activeOrgChanged.subscribe(() => {
            this.loadPermissions();
        });
    }

    public async GetActiveOrg(id?: string): Promise<Org.AsObject> {
        if (id) {
            const org = this.storage.getItem<Org.AsObject>(StorageKey.organization);
            if (org && this.cachedOrgs.find(tmp => tmp.id === org.id)) {
                return org;
            }
            return Promise.reject(new Error('no cached org'));
        } else {
            let orgs = this.cachedOrgs;
            if (orgs.length === 0) {
                orgs = (await this.SearchMyProjectOrgs(10, 0)).toObject().resultList;
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
            switchMap(() => from(this.GetMyzitadelPermissions())),
            map(rolesResp => rolesResp.toObject().permissionsList),
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

    public GetMyUserProfile(): Promise<UserProfileView> {
        return this.grpcService.auth.getMyUserProfile(new Empty());
    }

    public GetMyPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        return this.grpcService.auth.getMyPasswordComplexityPolicy(
            new Empty(),
        );
    }

    public GetMyUser(): Promise<UserView> {
        return this.grpcService.auth.getMyUser(
            new Empty(),
        );
    }

    public GetMyMfas(): Promise<MultiFactors> {
        return this.grpcService.auth.getMyMfas(
            new Empty(),
        );
    }

    public SearchMyProjectOrgs(
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

    public SaveMyUserProfile(
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

    public get zitadelPermissionsChanged(): Observable<string[]> {
        return this.zitadelPermissions;
    }

    public getMyUserSessions(): Promise<UserSessionViews> {
        return this.grpcService.auth.getMyUserSessions(
            new Empty(),
        );
    }

    public GetMyUserEmail(): Promise<UserEmail> {
        return this.grpcService.auth.getMyUserEmail(
            new Empty(),
        );
    }

    public SaveMyUserEmail(email: string): Promise<UserEmail> {
        const req = new UpdateUserEmailRequest();
        req.setEmail(email);
        return this.grpcService.auth.changeMyUserEmail(req);
    }

    public ResendMyEmailVerificationMail(): Promise<Empty> {
        return this.grpcService.auth.resendMyEmailVerificationMail(
            new Empty(),
        );
    }

    public RemoveMyUserPhone(): Promise<Empty> {
        return this.grpcService.auth.removeMyUserPhone(
            new Empty(),
        );
    }

    public GetMyzitadelPermissions(): Promise<MyPermissions> {
        return this.grpcService.auth.getMyZitadelPermissions(
            new Empty(),
        );
    }

    public GetMyUserPhone(): Promise<UserPhone> {
        return this.grpcService.auth.getMyUserPhone(
            new Empty(),
        );
    }

    public SaveMyUserPhone(phone: string): Promise<UserPhone> {
        const req = new UpdateUserPhoneRequest();
        req.setPhone(phone);
        return this.grpcService.auth.changeMyUserPhone(req);
    }

    public GetMyUserAddress(): Promise<UserAddress> {
        return this.grpcService.auth.getMyUserAddress(
            new Empty(),
        );
    }

    public ResendEmailVerification(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.auth.resendMyEmailVerificationMail(req);
    }

    public ResendPhoneVerification(): Promise<Empty> {
        const req = new Empty();
        return this.grpcService.auth.resendMyPhoneVerificationCode(req);
    }

    public ChangeMyPassword(oldPassword: string, newPassword: string): Promise<Empty> {
        const req = new PasswordChange();
        req.setOldPassword(oldPassword);
        req.setNewPassword(newPassword);
        return this.grpcService.auth.changeMyPassword(req);
    }

    public RemoveExternalIDP(
        externalUserId: string,
        idpConfigId: string,
    ): Promise<Empty> {
        const req = new ExternalIDPRemoveRequest();
        req.setExternalUserId(externalUserId);
        req.setIdpConfigId(idpConfigId);
        return this.grpcService.auth.removeMyExternalIDP(req);
    }

    public SearchMyExternalIdps(
        limit: number,
        offset: number,
    ): Promise<ExternalIDPSearchResponse> {
        const req = new ExternalIDPSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.auth.searchMyExternalIDPs(req);
    }

    public AddMfaOTP(): Promise<MfaOtpResponse> {
        return this.grpcService.auth.addMfaOTP(
            new Empty(),
        );
    }

    public AddMyMfaU2F(): Promise<WebAuthNResponse> {
        return this.grpcService.auth.addMyMfaU2F(
            new Empty(),
        );
    }

    public RemoveMyMfaU2F(id: string): Promise<Empty> {
        const req = new WebAuthNTokenID();
        req.setId(id);
        return this.grpcService.auth.removeMyMfaU2F(req);
    }

    public VerifyMyMfaU2F(credential: string, tokenname: string): Promise<Empty> {
        const req = new VerifyWebAuthN();
        req.setPublicKeyCredential(credential);
        req.setTokenName(tokenname);

        return this.grpcService.auth.verifyMyMfaU2F(
            req,
        );
    }

    public GetMyPasswordless(): Promise<WebAuthNTokens> {
        return this.grpcService.auth.getMyPasswordless(
            new Empty(),
        );
    }

    public AddMyPasswordless(): Promise<WebAuthNResponse> {
        return this.grpcService.auth.addMyPasswordless(
            new Empty(),
        );
    }

    public RemoveMyPasswordless(id: string): Promise<Empty> {
        const req = new WebAuthNTokenID();
        req.setId(id);
        return this.grpcService.auth.removeMyPasswordless(req);
    }

    public verifyMyPasswordless(credential: string, tokenname: string): Promise<Empty> {
        const req = new VerifyWebAuthN();
        req.setPublicKeyCredential(credential);
        req.setTokenName(tokenname);

        return this.grpcService.auth.verifyMyPasswordless(
            req,
        );
    }

    public RemoveMfaOTP(): Promise<Empty> {
        return this.grpcService.auth.removeMfaOTP(
            new Empty(),
        );
    }

    public VerifyMfaOTP(code: string): Promise<Empty> {
        const req = new VerifyMfaOtp();
        req.setCode(code);
        return this.grpcService.auth.verifyMfaOTP(req);
    }

    public VerifyMyUserPhone(code: string): Promise<Empty> {
        const req = new VerifyUserPhoneRequest();
        req.setCode(code);
        return this.grpcService.auth.verifyMyUserPhone(req);
    }

    public SaveMyUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return this.grpcService.auth.updateMyUserAddress(req);
    }

    public GetMyUserChanges(limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangesRequest();
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return this.grpcService.auth.getMyUserChanges(req);
    }
}

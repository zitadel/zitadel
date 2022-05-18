import { Injectable } from '@angular/core';
import { SortDirection } from '@angular/material/sort';
import { OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, merge, Observable, of, Subject } from 'rxjs';
import { catchError, filter, finalize, map, mergeMap, switchMap, take, timeout } from 'rxjs/operators';

import {
    AddMyAuthFactorOTPRequest,
    AddMyAuthFactorOTPResponse,
    AddMyAuthFactorU2FRequest,
    AddMyAuthFactorU2FResponse,
    AddMyPasswordlessLinkRequest,
    AddMyPasswordlessLinkResponse,
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
    GetSupportedLanguagesRequest,
    GetSupportedLanguagesResponse,
    ListMyAuthFactorsRequest,
    ListMyAuthFactorsResponse,
    ListMyLinkedIDPsRequest,
    ListMyLinkedIDPsResponse,
    ListMyMembershipsRequest,
    ListMyMembershipsResponse,
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
    RemoveMyAuthFactorOTPRequest,
    RemoveMyAuthFactorOTPResponse,
    RemoveMyAuthFactorU2FRequest,
    RemoveMyAuthFactorU2FResponse,
    RemoveMyAvatarRequest,
    RemoveMyAvatarResponse,
    RemoveMyLinkedIDPRequest,
    RemoveMyLinkedIDPResponse,
    RemoveMyPasswordlessRequest,
    RemoveMyPasswordlessResponse,
    RemoveMyPhoneRequest,
    RemoveMyPhoneResponse,
    RemoveMyUserRequest,
    RemoveMyUserResponse,
    ResendMyEmailVerificationRequest,
    ResendMyEmailVerificationResponse,
    ResendMyPhoneVerificationRequest,
    ResendMyPhoneVerificationResponse,
    SendMyPasswordlessLinkRequest,
    SendMyPasswordlessLinkResponse,
    SetMyEmailRequest,
    SetMyEmailResponse,
    SetMyPhoneRequest,
    SetMyPhoneResponse,
    UpdateMyPasswordRequest,
    UpdateMyPasswordResponse,
    UpdateMyProfileRequest,
    UpdateMyProfileResponse,
    UpdateMyUserNameRequest,
    UpdateMyUserNameResponse,
    VerifyMyAuthFactorOTPRequest,
    VerifyMyAuthFactorOTPResponse,
    VerifyMyAuthFactorU2FRequest,
    VerifyMyAuthFactorU2FResponse,
    VerifyMyPasswordlessRequest,
    VerifyMyPasswordlessResponse,
    VerifyMyPhoneRequest,
    VerifyMyPhoneResponse,
} from '../proto/generated/zitadel/auth_pb';
import { ChangeQuery } from '../proto/generated/zitadel/change_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { Org, OrgFieldName, OrgQuery } from '../proto/generated/zitadel/org_pb';
import { Gender, MembershipQuery, User, WebAuthNVerification } from '../proto/generated/zitadel/user_pb';
import { GrpcService } from './grpc.service';
import { StorageKey, StorageLocation, StorageService } from './storage.service';

@Injectable({
  providedIn: 'root',
})
export class GrpcAuthService {
  private _activeOrgChanged: Subject<Org.AsObject> = new Subject();
  public user!: Observable<User.AsObject | undefined>;
  public userSubject: BehaviorSubject<User.AsObject | undefined> = new BehaviorSubject<User.AsObject | undefined>(undefined);
  private zitadelPermissions: BehaviorSubject<string[]> = new BehaviorSubject(['user.resourceowner']);

  public readonly fetchedZitadelPermissions: BehaviorSubject<boolean> = new BehaviorSubject(false as boolean);

  private cachedOrgs: Org.AsObject[] = [];

  constructor(
    private readonly grpcService: GrpcService,
    private oauthService: OAuthService,
    private storage: StorageService,
  ) {
    this.user = merge(
      of(this.oauthService.getAccessToken()).pipe(filter((token) => (token ? true : false))),
      this.oauthService.events.pipe(
        filter((e) => e.type === 'token_received'),
        timeout(this.oauthService.waitForTokenInMsec || 0),
        catchError((_) => of(null)), // timeout is not an error
        map((_) => this.oauthService.getAccessToken()),
      ),
    ).pipe(
      take(1),
      mergeMap(() => {
        return from(
          this.getMyUser().then((resp) => {
            const user = resp.user;
            if (user) {
              return user;
            } else {
              return undefined;
            }
          }),
        );
      }),
      finalize(() => {
        this.loadPermissions();
      }),
    );

    this.user.subscribe(this.userSubject);

    this.activeOrgChanged.subscribe(() => {
      this.loadPermissions();
    });
  }

  public async getActiveOrg(id?: string): Promise<Org.AsObject> {
    if (id) {
      const find = this.cachedOrgs.find((tmp) => tmp.id === id);
      if (find) {
        this.setActiveOrg(find);
        return Promise.resolve(find);
      } else {
        const orgs = (await this.listMyProjectOrgs(10, 0)).resultList;
        this.cachedOrgs = orgs;

        const toFind = this.cachedOrgs.find((tmp) => tmp.id === id);
        if (toFind) {
          this.setActiveOrg(toFind);
          return Promise.resolve(toFind);
        } else {
          return Promise.reject(new Error('requested organization not found'));
        }
      }
    } else {
      let orgs = this.cachedOrgs;
      if (orgs.length === 0) {
        orgs = (await this.listMyProjectOrgs()).resultList;
        this.cachedOrgs = orgs;
      }

      const org = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.local);
      if (org && orgs.find((tmp) => tmp.id === org.id)) {
        this.storage.setItem(StorageKey.organization, org, StorageLocation.session);
        return org;
      }

      if (orgs.length === 0) {
        return Promise.reject(new Error('No organizations found!'));
      }
      const orgToSet = orgs.find((element) => element.id !== '0' && element.name !== '');

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
    // Set organization in localstorage to get the last used organization in a new tab
    this.storage.setItem(StorageKey.organization, org, StorageLocation.local);
    this.storage.setItem(StorageKey.organization, org, StorageLocation.session);
    this._activeOrgChanged.next(org);
  }

  private loadPermissions(): void {
    from(this.listMyZitadelPermissions())
      .pipe(
        map((rolesResp) => rolesResp.resultList),
        catchError((_) => {
          return of([]);
        }),
        finalize(() => {
          this.fetchedZitadelPermissions.next(true);
        }),
      )
      .subscribe((roles) => {
        this.zitadelPermissions.next(roles);
      });
  }

  /**
   * returns true if user has one of the provided roles
   * @param roles roles of the user
   */
  public isAllowed(roles: string[] | RegExp[]): Observable<boolean> {
    if (roles && roles.length > 0) {
      return this.zitadelPermissions.pipe(switchMap((zroles) => of(this.hasRoles(zroles, roles))));
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
    return (
      requestedRoles.findIndex((regexp: any) => {
        return (
          userRoles.findIndex((role) => {
            return new RegExp(regexp).test(role);
          }) > -1
        );
      }) > -1
    );
  }

  public getMyProfile(): Promise<GetMyProfileResponse.AsObject> {
    return this.grpcService.auth.getMyProfile(new GetMyProfileRequest(), null).then((resp) => resp.toObject());
  }

  public getMyPasswordComplexityPolicy(): Promise<GetMyPasswordComplexityPolicyResponse.AsObject> {
    return this.grpcService.auth
      .getMyPasswordComplexityPolicy(new GetMyPasswordComplexityPolicyRequest(), null)
      .then((resp) => resp.toObject());
  }

  public getMyUser(): Promise<GetMyUserResponse.AsObject> {
    return this.grpcService.auth.getMyUser(new GetMyUserRequest(), null).then((resp) => resp.toObject());
  }

  public listMyMultiFactors(): Promise<ListMyAuthFactorsResponse.AsObject> {
    return this.grpcService.auth.listMyAuthFactors(new ListMyAuthFactorsRequest(), null).then((resp) => resp.toObject());
  }

  public listMyProjectOrgs(
    limit?: number,
    offset?: number,
    queryList?: OrgQuery[],
    sortingColumn?: OrgFieldName,
    sortingDirection?: SortDirection,
  ): Promise<ListMyProjectOrgsResponse.AsObject> {
    const req = new ListMyProjectOrgsRequest();
    const query = new ListQuery();
    if (offset) {
      query.setOffset(offset);
    }
    if (limit) {
      query.setLimit(limit);
    }
    if (queryList) {
      req.setQueriesList(queryList);
    }
    // if (sortingColumn) {
    //     req.setSortingColumn(sortingColumn);
    // }

    req.setQuery(query);

    return this.grpcService.auth.listMyProjectOrgs(req, null).then((resp) => resp.toObject());
  }

  public updateMyProfile(
    firstName?: string,
    lastName?: string,
    nickName?: string,
    displayName?: string,
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
    if (displayName) {
      req.setDisplayName(displayName);
    }
    if (gender) {
      req.setGender(gender);
    }
    if (preferredLanguage) {
      req.setPreferredLanguage(preferredLanguage);
    }
    return this.grpcService.auth.updateMyProfile(req, null).then((resp) => resp.toObject());
  }

  public get zitadelPermissionsChanged(): Observable<string[]> {
    return this.zitadelPermissions;
  }

  public listMyUserSessions(): Promise<ListMyUserSessionsResponse.AsObject> {
    const req = new ListMyUserSessionsRequest();
    return this.grpcService.auth.listMyUserSessions(req, null).then((resp) => resp.toObject());
  }

  public listMyUserGrants(
    limit?: number,
    offset?: number,
    queryList?: ListQuery[],
  ): Promise<ListMyUserGrantsResponse.AsObject> {
    const req = new ListMyUserGrantsRequest();
    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }
    req.setQuery(query);
    return this.grpcService.auth.listMyUserGrants(req, null).then((resp) => resp.toObject());
  }

  public listMyMemberships(
    limit: number,
    offset: number,
    queryList?: MembershipQuery[],
  ): Promise<ListMyMembershipsResponse.AsObject> {
    const req = new ListMyMembershipsRequest();
    const metadata = new ListQuery();
    if (limit) {
      metadata.setLimit(limit);
    }
    if (offset) {
      metadata.setOffset(offset);
    }
    if (queryList) {
      req.setQueriesList(queryList);
    }
    req.setQuery(metadata);
    return this.grpcService.auth.listMyMemberships(req, null).then((resp) => resp.toObject());
  }

  public RemoveMyUser(): Promise<RemoveMyUserResponse.AsObject> {
    const req = new RemoveMyUserRequest();
    return this.grpcService.auth.removeMyUser(req, null).then((resp) => resp.toObject());
  }

  public getMyEmail(): Promise<GetMyEmailResponse.AsObject> {
    const req = new GetMyEmailRequest();
    return this.grpcService.auth.getMyEmail(req, null).then((resp) => resp.toObject());
  }

  public setMyEmail(email: string): Promise<SetMyEmailResponse.AsObject> {
    const req = new SetMyEmailRequest();
    req.setEmail(email);
    return this.grpcService.auth.setMyEmail(req, null).then((resp) => resp.toObject());
  }

  public resendMyEmailVerification(): Promise<ResendMyEmailVerificationResponse.AsObject> {
    const req = new ResendMyEmailVerificationRequest();
    return this.grpcService.auth.resendMyEmailVerification(req, null).then((resp) => resp.toObject());
  }

  public getSupportedLanguages(): Promise<GetSupportedLanguagesResponse.AsObject> {
    const req = new GetSupportedLanguagesRequest();
    return this.grpcService.auth.getSupportedLanguages(req, null).then((resp) => resp.toObject());
  }

  public removeMyPhone(): Promise<RemoveMyPhoneResponse.AsObject> {
    return this.grpcService.auth.removeMyPhone(new RemoveMyPhoneRequest(), null).then((resp) => resp.toObject());
  }

  public updateMyUserName(username: string): Promise<UpdateMyUserNameResponse.AsObject> {
    const req = new UpdateMyUserNameRequest();
    req.setUserName(username);
    return this.grpcService.auth.updateMyUserName(req, null).then((resp) => resp.toObject());
  }

  public listMyZitadelPermissions(): Promise<ListMyZitadelPermissionsResponse.AsObject> {
    return this.grpcService.auth
      .listMyZitadelPermissions(new ListMyZitadelPermissionsRequest(), null)
      .then((resp) => resp.toObject());
  }

  public getMyPhone(): Promise<GetMyPhoneResponse.AsObject> {
    return this.grpcService.auth.getMyPhone(new GetMyPhoneRequest(), null).then((resp) => resp.toObject());
  }

  public setMyPhone(phone: string): Promise<SetMyPhoneResponse.AsObject> {
    const req = new SetMyPhoneRequest();
    req.setPhone(phone);
    return this.grpcService.auth.setMyPhone(req, null).then((resp) => resp.toObject());
  }

  public resendMyPhoneVerification(): Promise<ResendMyPhoneVerificationResponse.AsObject> {
    const req = new ResendMyPhoneVerificationRequest();
    return this.grpcService.auth.resendMyPhoneVerification(req, null).then((resp) => resp.toObject());
  }

  public updateMyPassword(oldPassword: string, newPassword: string): Promise<UpdateMyPasswordResponse.AsObject> {
    const req = new UpdateMyPasswordRequest();
    req.setOldPassword(oldPassword);
    req.setNewPassword(newPassword);
    return this.grpcService.auth.updateMyPassword(req, null).then((resp) => resp.toObject());
  }

  public removeMyLinkedIDP(idpId: string, linkedUserId: string): Promise<RemoveMyLinkedIDPResponse.AsObject> {
    const req = new RemoveMyLinkedIDPRequest();
    req.setLinkedUserId(linkedUserId);
    req.setIdpId(idpId);
    return this.grpcService.auth.removeMyLinkedIDP(req, null).then((resp) => resp.toObject());
  }

  public removeMyAvatar(): Promise<RemoveMyAvatarResponse.AsObject> {
    const req = new RemoveMyAvatarRequest();
    return this.grpcService.auth.removeMyAvatar(req, null).then((resp) => resp.toObject());
  }

  public listMyLinkedIDPs(limit: number, offset: number): Promise<ListMyLinkedIDPsResponse.AsObject> {
    const req = new ListMyLinkedIDPsRequest();
    const metadata = new ListQuery();
    if (limit) {
      metadata.setLimit(limit);
    }
    if (offset) {
      metadata.setOffset(offset);
    }
    req.setQuery(metadata);
    return this.grpcService.auth.listMyLinkedIDPs(req, null).then((resp) => resp.toObject());
  }

  public addMyMultiFactorOTP(): Promise<AddMyAuthFactorOTPResponse.AsObject> {
    return this.grpcService.auth.addMyAuthFactorOTP(new AddMyAuthFactorOTPRequest(), null).then((resp) => resp.toObject());
  }

  public addMyMultiFactorU2F(): Promise<AddMyAuthFactorU2FResponse.AsObject> {
    return this.grpcService.auth.addMyAuthFactorU2F(new AddMyAuthFactorU2FRequest(), null).then((resp) => resp.toObject());
  }

  public removeMyMultiFactorU2F(tokenId: string): Promise<RemoveMyAuthFactorU2FResponse.AsObject> {
    const req = new RemoveMyAuthFactorU2FRequest();
    req.setTokenId(tokenId);
    return this.grpcService.auth.removeMyAuthFactorU2F(req, null).then((resp) => resp.toObject());
  }

  public verifyMyMultiFactorU2F(credential: string, tokenname: string): Promise<VerifyMyAuthFactorU2FResponse.AsObject> {
    const req = new VerifyMyAuthFactorU2FRequest();
    const verification = new WebAuthNVerification();
    verification.setPublicKeyCredential(credential);
    verification.setTokenName(tokenname);
    req.setVerification(verification);

    return this.grpcService.auth.verifyMyAuthFactorU2F(req, null).then((resp) => resp.toObject());
  }

  public listMyPasswordless(): Promise<ListMyPasswordlessResponse.AsObject> {
    return this.grpcService.auth.listMyPasswordless(new ListMyPasswordlessRequest(), null).then((resp) => resp.toObject());
  }

  public addMyPasswordless(): Promise<AddMyPasswordlessResponse.AsObject> {
    return this.grpcService.auth.addMyPasswordless(new AddMyPasswordlessRequest(), null).then((resp) => resp.toObject());
  }

  public removeMyPasswordless(tokenId: string): Promise<RemoveMyPasswordlessResponse.AsObject> {
    const req = new RemoveMyPasswordlessRequest();
    req.setTokenId(tokenId);
    return this.grpcService.auth.removeMyPasswordless(req, null).then((resp) => resp.toObject());
  }

  public verifyMyPasswordless(credential: string, tokenname: string): Promise<VerifyMyPasswordlessResponse.AsObject> {
    const req = new VerifyMyPasswordlessRequest();
    const verification = new WebAuthNVerification();
    verification.setTokenName(tokenname);
    verification.setPublicKeyCredential(credential);
    req.setVerification(verification);

    return this.grpcService.auth.verifyMyPasswordless(req, null).then((resp) => resp.toObject());
  }

  public sendMyPasswordlessLink(): Promise<SendMyPasswordlessLinkResponse.AsObject> {
    const req = new SendMyPasswordlessLinkRequest();
    return this.grpcService.auth.sendMyPasswordlessLink(req, null).then((resp) => resp.toObject());
  }

  public addMyPasswordlessLink(): Promise<AddMyPasswordlessLinkResponse.AsObject> {
    const req = new AddMyPasswordlessLinkRequest();
    return this.grpcService.auth.addMyPasswordlessLink(req, null).then((resp) => resp.toObject());
  }

  public removeMyMultiFactorOTP(): Promise<RemoveMyAuthFactorOTPResponse.AsObject> {
    return this.grpcService.auth
      .removeMyAuthFactorOTP(new RemoveMyAuthFactorOTPRequest(), null)
      .then((resp) => resp.toObject());
  }

  public verifyMyMultiFactorOTP(code: string): Promise<VerifyMyAuthFactorOTPResponse.AsObject> {
    const req = new VerifyMyAuthFactorOTPRequest();
    req.setCode(code);
    return this.grpcService.auth.verifyMyAuthFactorOTP(req, null).then((resp) => resp.toObject());
  }

  public verifyMyPhone(code: string): Promise<VerifyMyPhoneResponse.AsObject> {
    const req = new VerifyMyPhoneRequest();
    req.setCode(code);
    return this.grpcService.auth.verifyMyPhone(req, null).then((resp) => resp.toObject());
  }

  public listMyUserChanges(limit: number, sequence: number): Promise<ListMyUserChangesResponse.AsObject> {
    const req = new ListMyUserChangesRequest();
    const query = new ChangeQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (sequence) {
      query.setSequence(sequence);
    }
    req.setQuery(query);
    return this.grpcService.auth.listMyUserChanges(req, null).then((resp) => resp.toObject());
  }
}

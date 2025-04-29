import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  AddHumanUserRequestSchema,
  AddHumanUserResponse,
  CreateInviteCodeRequestSchema,
  CreateInviteCodeResponse,
  CreatePasskeyRegistrationLinkRequestSchema,
  CreatePasskeyRegistrationLinkResponse,
  DeactivateUserRequestSchema,
  DeactivateUserResponse,
  DeleteUserRequestSchema,
  DeleteUserResponse,
  GetUserByIDResponse,
  ListAuthenticationFactorsRequestSchema,
  ListAuthenticationFactorsResponse,
  ListPasskeysRequestSchema,
  ListPasskeysResponse,
  ListUsersRequestSchema,
  ListUsersResponse,
  LockUserRequestSchema,
  LockUserResponse,
  PasswordResetRequestSchema,
  ReactivateUserRequestSchema,
  ReactivateUserResponse,
  RemoveOTPEmailRequestSchema,
  RemoveOTPEmailResponse,
  RemoveOTPSMSRequestSchema,
  RemoveOTPSMSResponse,
  RemovePasskeyRequestSchema,
  RemovePasskeyResponse,
  RemovePhoneRequestSchema,
  RemovePhoneResponse,
  RemoveTOTPRequestSchema,
  RemoveTOTPResponse,
  RemoveU2FRequestSchema,
  RemoveU2FResponse,
  ResendInviteCodeRequestSchema,
  ResendInviteCodeResponse,
  SetEmailRequestSchema,
  SetEmailResponse,
  SetPasswordRequestSchema,
  SetPasswordResponse,
  SetPhoneRequestSchema,
  SetPhoneResponse,
  UnlockUserRequestSchema,
  UnlockUserResponse,
  UpdateHumanUserRequestSchema,
  UpdateHumanUserResponse,
} from '@zitadel/proto/zitadel/user/v2/user_service_pb';
import type { MessageInitShape } from '@bufbuild/protobuf';
import {
  AccessTokenType,
  Gender,
  HumanProfile,
  HumanProfileSchema,
  HumanUser,
  HumanUserSchema,
  MachineUser,
  MachineUserSchema,
  User as UserV2,
  UserSchema,
  UserState,
} from '@zitadel/proto/zitadel/user/v2/user_pb';
import { create } from '@bufbuild/protobuf';
import { Timestamp as TimestampV2, TimestampSchema } from '@bufbuild/protobuf/wkt';
import { Details, DetailsSchema } from '@zitadel/proto/zitadel/object/v2/object_pb';
import { Human, Machine, Phone, Profile, User } from '../proto/generated/zitadel/user_pb';
import { ObjectDetails } from '../proto/generated/zitadel/object_pb';
import { Timestamp } from '../proto/generated/google/protobuf/timestamp_pb';
import { HumanPhone, HumanPhoneSchema } from '@zitadel/proto/zitadel/user/v2/phone_pb';
import { OAuthService } from 'angular-oauth2-oidc';
import { debounceTime, EMPTY, Observable, of, ReplaySubject, shareReplay, switchAll, switchMap } from 'rxjs';
import { catchError, filter, map, startWith } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private user$$ = new ReplaySubject<Observable<UserV2>>(1);
  public user$ = this.user$$.pipe(
    startWith(this.getUser()),
    // makes sure if many subscribers reset the observable only one wins
    debounceTime(10),
    switchAll(),
    catchError((err) => {
      // reset user observable on error
      this.user$$.next(this.getUser());
      throw err;
    }),
  );

  constructor(
    private readonly grpcService: GrpcService,
    private readonly oauthService: OAuthService,
  ) {}

  private getUserId() {
    return this.oauthService.events.pipe(
      filter((event) => event.type === 'token_received'),
      map(() => this.oauthService.getIdToken()),
      startWith(this.oauthService.getIdToken()),
      filter(Boolean),
      switchMap((token) => {
        // we do this in a try catch so the observable will retry this logic if it fails
        try {
          // split jwt and get base64 encoded payload
          const unparsedPayload = atob(token.split('.')[1]);
          // parse payload
          const payload: unknown = JSON.parse(unparsedPayload);
          // check if sub is in payload and is a string
          if (payload && typeof payload === 'object' && 'sub' in payload && typeof payload.sub === 'string') {
            return of(payload.sub);
          }
          return EMPTY;
        } catch {
          return EMPTY;
        }
      }),
    );
  }

  private getUser() {
    return this.getUserId().pipe(
      switchMap((id) => this.getUserById(id)),
      map((resp) => resp.user),
      filter(Boolean),
      shareReplay({ refCount: true, bufferSize: 1 }),
    );
  }

  public addHumanUser(req: MessageInitShape<typeof AddHumanUserRequestSchema>): Promise<AddHumanUserResponse> {
    return this.grpcService.userNew.addHumanUser(create(AddHumanUserRequestSchema, req));
  }

  public listUsers(req: MessageInitShape<typeof ListUsersRequestSchema>): Promise<ListUsersResponse> {
    return this.grpcService.userNew.listUsers(req);
  }

  public getUserById(userId: string): Promise<GetUserByIDResponse> {
    return this.grpcService.userNew.getUserByID({ userId });
  }

  public deactivateUser(userId: string): Promise<DeactivateUserResponse> {
    return this.grpcService.userNew.deactivateUser(create(DeactivateUserRequestSchema, { userId }));
  }

  public reactivateUser(userId: string): Promise<ReactivateUserResponse> {
    return this.grpcService.userNew.reactivateUser(create(ReactivateUserRequestSchema, { userId }));
  }

  public deleteUser(userId: string): Promise<DeleteUserResponse> {
    return this.grpcService.userNew.deleteUser(create(DeleteUserRequestSchema, { userId }));
  }

  public updateUser(req: MessageInitShape<typeof UpdateHumanUserRequestSchema>): Promise<UpdateHumanUserResponse> {
    return this.grpcService.userNew.updateHumanUser(create(UpdateHumanUserRequestSchema, req));
  }

  public lockUser(userId: string): Promise<LockUserResponse> {
    return this.grpcService.userNew.lockUser(create(LockUserRequestSchema, { userId }));
  }

  public unlockUser(userId: string): Promise<UnlockUserResponse> {
    return this.grpcService.userNew.unlockUser(create(UnlockUserRequestSchema, { userId }));
  }

  public listAuthenticationFactors(
    req: MessageInitShape<typeof ListAuthenticationFactorsRequestSchema>,
  ): Promise<ListAuthenticationFactorsResponse> {
    return this.grpcService.userNew.listAuthenticationFactors(create(ListAuthenticationFactorsRequestSchema, req));
  }

  public listPasskeys(req: MessageInitShape<typeof ListPasskeysRequestSchema>): Promise<ListPasskeysResponse> {
    return this.grpcService.userNew.listPasskeys(create(ListPasskeysRequestSchema, req));
  }

  public removePasskeys(req: MessageInitShape<typeof RemovePasskeyRequestSchema>): Promise<RemovePasskeyResponse> {
    return this.grpcService.userNew.removePasskey(create(RemovePasskeyRequestSchema, req));
  }

  public createPasskeyRegistrationLink(
    req: MessageInitShape<typeof CreatePasskeyRegistrationLinkRequestSchema>,
  ): Promise<CreatePasskeyRegistrationLinkResponse> {
    return this.grpcService.userNew.createPasskeyRegistrationLink(create(CreatePasskeyRegistrationLinkRequestSchema, req));
  }

  public removePhone(userId: string): Promise<RemovePhoneResponse> {
    return this.grpcService.userNew.removePhone(create(RemovePhoneRequestSchema, { userId }));
  }

  public setPhone(req: MessageInitShape<typeof SetPhoneRequestSchema>): Promise<SetPhoneResponse> {
    return this.grpcService.userNew.setPhone(create(SetPhoneRequestSchema, req));
  }

  public setEmail(req: MessageInitShape<typeof SetEmailRequestSchema>): Promise<SetEmailResponse> {
    return this.grpcService.userNew.setEmail(create(SetEmailRequestSchema, req));
  }

  public removeTOTP(userId: string): Promise<RemoveTOTPResponse> {
    return this.grpcService.userNew.removeTOTP(create(RemoveTOTPRequestSchema, { userId }));
  }

  public removeU2F(userId: string, u2fId: string): Promise<RemoveU2FResponse> {
    return this.grpcService.userNew.removeU2F(create(RemoveU2FRequestSchema, { userId, u2fId }));
  }

  public removeOTPSMS(userId: string): Promise<RemoveOTPSMSResponse> {
    return this.grpcService.userNew.removeOTPSMS(create(RemoveOTPSMSRequestSchema, { userId }));
  }

  public removeOTPEmail(userId: string): Promise<RemoveOTPEmailResponse> {
    return this.grpcService.userNew.removeOTPEmail(create(RemoveOTPEmailRequestSchema, { userId }));
  }

  public resendInviteCode(userId: string): Promise<ResendInviteCodeResponse> {
    return this.grpcService.userNew.resendInviteCode(create(ResendInviteCodeRequestSchema, { userId }));
  }

  public createInviteCode(req: MessageInitShape<typeof CreateInviteCodeRequestSchema>): Promise<CreateInviteCodeResponse> {
    return this.grpcService.userNew.createInviteCode(create(CreateInviteCodeRequestSchema, req));
  }

  public passwordReset(req: MessageInitShape<typeof PasswordResetRequestSchema>) {
    return this.grpcService.userNew.passwordReset(create(PasswordResetRequestSchema, req));
  }

  public setPassword(req: MessageInitShape<typeof SetPasswordRequestSchema>): Promise<SetPasswordResponse> {
    return this.grpcService.userNew.setPassword(create(SetPasswordRequestSchema, req));
  }
}

function userToV2(user: User): UserV2 {
  const details = user.getDetails();
  return create(UserSchema, {
    userId: user.getId(),
    details: details && detailsToV2(details),
    state: user.getState() as number as UserState,
    username: user.getUserName(),
    loginNames: user.getLoginNamesList(),
    preferredLoginName: user.getPreferredLoginName(),
    type: typeToV2(user),
  });
}

function detailsToV2(details: ObjectDetails): Details {
  const changeDate = details.getChangeDate();
  return create(DetailsSchema, {
    sequence: BigInt(details.getSequence()),
    changeDate: changeDate && timestampToV2(changeDate),
    resourceOwner: details.getResourceOwner(),
  });
}

function timestampToV2(timestamp: Timestamp): TimestampV2 {
  return create(TimestampSchema, {
    seconds: BigInt(timestamp.getSeconds()),
    nanos: timestamp.getNanos(),
  });
}

function typeToV2(user: User): UserV2['type'] {
  const human = user.getHuman();
  if (human) {
    return { case: 'human', value: humanToV2(user, human) };
  }

  const machine = user.getMachine();
  if (machine) {
    return { case: 'machine', value: machineToV2(machine) };
  }

  return { case: undefined };
}

function humanToV2(user: User, human: Human): HumanUser {
  const profile = human.getProfile();
  const email = human.getEmail()?.getEmail();
  const phone = human.getPhone();
  const passwordChanged = human.getPasswordChanged();

  return create(HumanUserSchema, {
    userId: user.getId(),
    state: user.getState() as number as UserState,
    username: user.getUserName(),
    loginNames: user.getLoginNamesList(),
    preferredLoginName: user.getPreferredLoginName(),
    profile: profile && humanProfileToV2(profile),
    email: { email },
    phone: phone && humanPhoneToV2(phone),
    passwordChangeRequired: false,
    passwordChanged: passwordChanged && timestampToV2(passwordChanged),
  });
}

function humanProfileToV2(profile: Profile): HumanProfile {
  return create(HumanProfileSchema, {
    givenName: profile.getFirstName(),
    familyName: profile.getLastName(),
    nickName: profile.getNickName(),
    displayName: profile.getDisplayName(),
    preferredLanguage: profile.getPreferredLanguage(),
    gender: profile.getGender() as number as Gender,
    avatarUrl: profile.getAvatarUrl(),
  });
}

function humanPhoneToV2(phone: Phone): HumanPhone {
  return create(HumanPhoneSchema, {
    phone: phone.getPhone(),
    isVerified: phone.getIsPhoneVerified(),
  });
}

function machineToV2(machine: Machine): MachineUser {
  return create(MachineUserSchema, {
    name: machine.getName(),
    description: machine.getDescription(),
    hasSecret: machine.getHasSecret(),
    accessTokenType: machine.getAccessTokenType() as number as AccessTokenType,
  });
}

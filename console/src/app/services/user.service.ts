import { Injectable, signal } from '@angular/core';
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
import { create } from '@bufbuild/protobuf';
import { OAuthService } from 'angular-oauth2-oidc';
import { EMPTY, of, switchMap } from 'rxjs';
import { filter, map, startWith } from 'rxjs/operators';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, queryOptions, skipToken } from '@tanstack/angular-query-experimental';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private userId = this.getUserId();

  public userQuery() {
    return injectQuery(() => this.userQueryOptions());
  }

  public userQueryOptions() {
    const userId = this.userId();
    return queryOptions({
      queryKey: ['user', userId],
      queryFn: userId ? () => this.getUserById(userId).then((resp) => resp.user) : skipToken,
    });
  }

  constructor(
    private readonly grpcService: GrpcService,
    private readonly oauthService: OAuthService,
  ) {}

  private getUserId() {
    const userId$ = this.oauthService.events.pipe(
      filter((event) => event.type === 'token_received'),
      // can actually return null
      // https://github.com/manfredsteyer/angular-oauth2-oidc/blob/c724ad73eadbb28338b084e3afa5ed49a0ea058c/projects/lib/src/oauth-service.ts#L2365
      map(() => this.oauthService.getIdToken() as string | null),
      startWith(this.oauthService.getIdToken() as string | null),
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

    return toSignal(userId$, { initialValue: undefined });
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

  public setPassword(req: MessageInitShape<typeof SetPasswordRequestSchema>): Promise<SetPasswordResponse> {
    return this.grpcService.userNew.setPassword(create(SetPasswordRequestSchema, req));
  }
}

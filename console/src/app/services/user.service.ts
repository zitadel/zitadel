import { computed, Injectable, Signal } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  AddHumanUserRequestSchema,
  AddHumanUserResponse,
  CreateInviteCodeRequestSchema,
  CreateInviteCodeResponse,
  CreatePasskeyRegistrationLinkRequestSchema,
  CreatePasskeyRegistrationLinkResponse,
  DeactivateUserResponse,
  DeleteUserResponse,
  GetUserByIDResponse,
  ListAuthenticationFactorsRequestSchema,
  ListAuthenticationFactorsResponse,
  ListPasskeysRequestSchema,
  ListPasskeysResponse,
  ListUsersRequestSchema,
  ListUsersResponse,
  ReactivateUserResponse,
  RemoveOTPEmailResponse,
  RemoveOTPSMSResponse,
  RemovePasskeyRequestSchema,
  RemovePasskeyResponse,
  RemovePhoneResponse,
  RemoveTOTPResponse,
  RemoveU2FResponse,
  ResendInviteCodeResponse,
  SetEmailRequestSchema,
  SetEmailResponse,
  SetPasswordRequestSchema,
  SetPasswordResponse,
  SetPhoneRequestSchema,
  SetPhoneResponse,
  UnlockUserResponse,
  UpdateHumanUserRequestSchema,
  UpdateHumanUserResponse,
} from '@zitadel/proto/zitadel/user/v2/user_service_pb';
import type { MessageInitShape } from '@bufbuild/protobuf';
import { OAuthService } from 'angular-oauth2-oidc';
import { filter, map } from 'rxjs/operators';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, queryOptions, skipToken } from '@tanstack/angular-query-experimental';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private readonly payload: Signal<unknown | undefined>;
  public readonly userId: Signal<string | undefined>;
  public readonly isExpired: Signal<boolean>;

  public userQuery() {
    return injectQuery(() => this.userQueryOptions());
  }

  public userQueryOptions() {
    const userId = this.userId();
    return queryOptions({
      queryKey: [userId, 'user'],
      queryFn: userId ? () => this.getUserById(userId).then((resp) => resp.user) : skipToken,
    });
  }

  constructor(
    private readonly grpcService: GrpcService,
    private readonly oauthService: OAuthService,
  ) {
    this.payload = this.getPayload();
    this.userId = this.getUserId(this.payload);
    this.isExpired = this.getIsExpired(this.payload);
  }

  private getPayload() {
    const idToken$ = this.oauthService.events.pipe(
      filter((event) => event.type === 'token_received'),
      // can actually return null
      // https://github.com/manfredsteyer/angular-oauth2-oidc/blob/c724ad73eadbb28338b084e3afa5ed49a0ea058c/projects/lib/src/oauth-service.ts#L2365
      map(() => this.oauthService.getIdToken() as string | null),
    );
    const idToken = toSignal(idToken$, { initialValue: this.oauthService.getIdToken() as string | null });

    return computed(() => {
      try {
        // split jwt and get base64 encoded payload
        const unparsedPayload = atob((idToken() ?? '').split('.')[1]);
        // parse payload
        return JSON.parse(unparsedPayload) as unknown;
      } catch {
        return undefined;
      }
    });
  }

  private getUserId(payloadSignal: Signal<unknown | undefined>) {
    return computed(() => {
      const payload = payloadSignal();
      if (payload && typeof payload === 'object' && 'sub' in payload && typeof payload.sub === 'string') {
        return payload.sub;
      }
      return undefined;
    });
  }

  private getIsExpired(payloadSignal: Signal<unknown | undefined>) {
    const expSignal = computed(() => {
      const payload = payloadSignal();
      if (payload && typeof payload === 'object' && 'exp' in payload && typeof payload.exp === 'number') {
        return new Date(payload.exp * 1000);
      }
      return undefined;
    });

    return computed(() => {
      const exp = expSignal();
      return exp ? exp <= new Date() : true;
    });
  }

  public addHumanUser(req: MessageInitShape<typeof AddHumanUserRequestSchema>): Promise<AddHumanUserResponse> {
    return this.grpcService.userNew.addHumanUser(req);
  }

  public listUsers(req: MessageInitShape<typeof ListUsersRequestSchema>): Promise<ListUsersResponse> {
    return this.grpcService.userNew.listUsers(req);
  }

  public getUserById(userId: string): Promise<GetUserByIDResponse> {
    return this.grpcService.userNew.getUserByID({ userId });
  }

  public deactivateUser(userId: string): Promise<DeactivateUserResponse> {
    return this.grpcService.userNew.deactivateUser({ userId });
  }

  public reactivateUser(userId: string): Promise<ReactivateUserResponse> {
    return this.grpcService.userNew.reactivateUser({ userId });
  }

  public deleteUser(userId: string): Promise<DeleteUserResponse> {
    return this.grpcService.userNew.deleteUser({ userId });
  }

  public updateUser(req: MessageInitShape<typeof UpdateHumanUserRequestSchema>): Promise<UpdateHumanUserResponse> {
    return this.grpcService.userNew.updateHumanUser(req);
  }

  public unlockUser(userId: string): Promise<UnlockUserResponse> {
    return this.grpcService.userNew.unlockUser({ userId });
  }

  public listAuthenticationFactors(
    req: MessageInitShape<typeof ListAuthenticationFactorsRequestSchema>,
  ): Promise<ListAuthenticationFactorsResponse> {
    return this.grpcService.userNew.listAuthenticationFactors(req);
  }

  public listPasskeys(req: MessageInitShape<typeof ListPasskeysRequestSchema>): Promise<ListPasskeysResponse> {
    return this.grpcService.userNew.listPasskeys(req);
  }

  public removePasskeys(req: MessageInitShape<typeof RemovePasskeyRequestSchema>): Promise<RemovePasskeyResponse> {
    return this.grpcService.userNew.removePasskey(req);
  }

  public createPasskeyRegistrationLink(
    req: MessageInitShape<typeof CreatePasskeyRegistrationLinkRequestSchema>,
  ): Promise<CreatePasskeyRegistrationLinkResponse> {
    return this.grpcService.userNew.createPasskeyRegistrationLink(req);
  }

  public removePhone(userId: string): Promise<RemovePhoneResponse> {
    return this.grpcService.userNew.removePhone({ userId });
  }

  public setPhone(req: MessageInitShape<typeof SetPhoneRequestSchema>): Promise<SetPhoneResponse> {
    return this.grpcService.userNew.setPhone(req);
  }

  public setEmail(req: MessageInitShape<typeof SetEmailRequestSchema>): Promise<SetEmailResponse> {
    return this.grpcService.userNew.setEmail(req);
  }

  public removeTOTP(userId: string): Promise<RemoveTOTPResponse> {
    return this.grpcService.userNew.removeTOTP({ userId });
  }

  public removeU2F(userId: string, u2fId: string): Promise<RemoveU2FResponse> {
    return this.grpcService.userNew.removeU2F({ userId, u2fId });
  }

  public removeOTPSMS(userId: string): Promise<RemoveOTPSMSResponse> {
    return this.grpcService.userNew.removeOTPSMS({ userId });
  }

  public removeOTPEmail(userId: string): Promise<RemoveOTPEmailResponse> {
    return this.grpcService.userNew.removeOTPEmail({ userId });
  }

  public resendInviteCode(userId: string): Promise<ResendInviteCodeResponse> {
    return this.grpcService.userNew.resendInviteCode({ userId });
  }

  public createInviteCode(req: MessageInitShape<typeof CreateInviteCodeRequestSchema>): Promise<CreateInviteCodeResponse> {
    return this.grpcService.userNew.createInviteCode(req);
  }

  public setPassword(req: MessageInitShape<typeof SetPasswordRequestSchema>): Promise<SetPasswordResponse> {
    return this.grpcService.userNew.setPassword(req);
  }
}

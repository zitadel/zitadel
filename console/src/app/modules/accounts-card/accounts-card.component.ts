import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';
import { Session, SessionState as V1SessionState, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { SessionService } from 'src/app/services/session.service';
import { ListMyUserSessionsRequest } from '@zitadel/proto/zitadel/auth_pb';
import {
  catchError,
  defer,
  firstValueFrom,
  map,
  mergeWith,
  NEVER,
  Observable,
  of,
  shareReplay,
  switchMap,
  timeout,
  TimeoutError,
} from 'rxjs';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { ToastService } from 'src/app/services/toast.service';
import { SessionState as V2SessionState } from '@zitadel/proto/zitadel/user_pb';

interface V1AndV2Session {
  displayName: string;
  avatarUrl: string;
  loginName: string;
  userName: string;
  authState: V1SessionState | V2SessionState;
}

@Component({
  selector: 'cnsl-accounts-card',
  templateUrl: './accounts-card.component.html',
  styleUrls: ['./accounts-card.component.scss'],
})
export class AccountsCardComponent {
  @Input() public user?: User.AsObject;
  @Input() public iamuser: boolean | null = false;

  @Output() public closedCard: EventEmitter<void> = new EventEmitter();

  public UserState: any = UserState;
  private labelpolicy = toSignal(this.userService.labelpolicy$, { initialValue: undefined });
  public readonly sessions$: Observable<V1AndV2Session[] | undefined>;

  constructor(
    public authService: AuthenticationService,
    private router: Router,
    private userService: GrpcAuthService,
    private sessionService: SessionService,
    private readonly featureService: NewFeatureService,
    private toast: ToastService,
  ) {
    this.sessions$ = this.getUseLoginV2()
      .pipe(shareReplay({ refCount: true, bufferSize: 1 }))
      .pipe(
        switchMap((loginV2) => {
          if (!loginV2?.required) {
            return defer(() =>
              this.userService.listMyUserSessions().then((sessions) => {
                return sessions.resultList
                  .filter((user) => user.loginName !== this.user?.preferredLoginName)
                  .map((s) => {
                    return {
                      displayName: s.displayName,
                      avatarUrl: s.avatarUrl,
                      loginName: s.loginName,
                      authState: s.authState,
                      userName: s.userName,
                    };
                  });
              }),
            );
          } else {
            return defer(() =>
              this.sessionService.listMyUserSessions({}).then((sessions) => {
                return sessions.result
                  .filter((s) => s.loginName !== this.user?.preferredLoginName)
                  .map((s) => {
                    return {
                      displayName: s.displayName,
                      avatarUrl: s.avatarUrl,
                      loginName: s.loginName,
                      authState: s.authState,
                      userName: s.userName,
                    };
                  });
              }),
            );
          }
        }),
        catchError((err) => {
          this.toast.showError(err);
          return of([]);
        }),
      );
  }

  private getUseLoginV2() {
    return defer(() => this.featureService.getInstanceFeatures()).pipe(
      map(({ loginV2 }) => loginV2),
      timeout(1000),
      catchError((err) => {
        if (!(err instanceof TimeoutError)) {
          this.toast.showError(err);
        }
        return of(undefined);
      }),
      mergeWith(NEVER),
    );
  }

  public editUserProfile(): void {
    this.router.navigate(['users/me']);
    this.closedCard.emit();
  }

  public closeCard(element: HTMLElement): void {
    if (!element.classList.contains('dontcloseonclick')) {
      this.closedCard.emit();
    }
  }

  public selectAccount(loginHint: string): void {
    const configWithPrompt: Partial<AuthConfig> = {
      customQueryParams: {
        login_hint: loginHint,
      },
    };
    this.authService.authenticate(configWithPrompt);
  }

  public selectNewAccount(): void {
    const configWithPrompt: Partial<AuthConfig> = {
      customQueryParams: {
        prompt: 'login',
      } as any,
    };
    this.authService.authenticate(configWithPrompt);
  }

  public logout(): void {
    const lP = JSON.stringify(this.labelpolicy());
    localStorage.setItem('labelPolicyOnSignout', lP);

    this.authService.signout();
    this.closedCard.emit();
  }
}

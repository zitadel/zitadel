import { Component, EventEmitter, Input, NgIterable, Output } from '@angular/core';
import { Router } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';
import { SessionState as V1SessionState, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { SessionService } from 'src/app/services/session.service';
import {
  catchError,
  defer,
  from,
  map,
  mergeMap,
  Observable,
  of,
  ReplaySubject,
  shareReplay,
  switchMap,
  timeout,
  TimeoutError,
  toArray,
} from 'rxjs';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { ToastService } from 'src/app/services/toast.service';
import { SessionState as V2SessionState } from '@zitadel/proto/zitadel/user_pb';
import { filter, withLatestFrom } from 'rxjs/operators';

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
  @Input({ required: true })
  public set user(user: User.AsObject) {
    this.user$.next(user);
  }

  @Input() public iamuser: boolean | null = false;

  @Output() public closedCard = new EventEmitter<void>();

  protected readonly user$ = new ReplaySubject<User.AsObject>(1);
  protected readonly UserState = UserState;
  private readonly labelpolicy = toSignal(this.userService.labelpolicy$, { initialValue: undefined });
  protected readonly sessions$: Observable<V1AndV2Session[]>;

  constructor(
    protected readonly authService: AuthenticationService,
    private readonly router: Router,
    private readonly userService: GrpcAuthService,
    private readonly sessionService: SessionService,
    private readonly featureService: NewFeatureService,
    private readonly toast: ToastService,
  ) {
    this.sessions$ = this.getSessions().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
  }

  private getUseLoginV2() {
    return defer(() => this.featureService.getInstanceFeatures()).pipe(
      map(({ loginV2 }) => loginV2?.required ?? false),
      timeout(1000),
      catchError((err) => {
        if (!(err instanceof TimeoutError)) {
          this.toast.showError(err);
        }
        return of(false);
      }),
    );
  }

  private getSessions(): Observable<V1AndV2Session[]> {
    const useLoginV2$ = this.getUseLoginV2();

    return useLoginV2$.pipe(
      switchMap((useLoginV2) => {
        if (useLoginV2) {
          return this.getV2Sessions();
        } else {
          return this.getV1Sessions();
        }
      }),
      catchError((err) => {
        this.toast.showError(err);
        return of([]);
      }),
    );
  }

  private getV1Sessions(): Observable<V1AndV2Session[]> {
    return defer(() => this.userService.listMyUserSessions()).pipe(
      mergeMap(({ resultList }) => from(resultList)),
      withLatestFrom(this.user$),
      filter(([{ loginName }, user]) => loginName !== user.preferredLoginName),
      map(([s]) => ({
        displayName: s.displayName,
        avatarUrl: s.avatarUrl,
        loginName: s.loginName,
        authState: s.authState,
        userName: s.userName,
      })),
      toArray(),
    );
  }

  private getV2Sessions(): Observable<V1AndV2Session[]> {
    return defer(() =>
      this.sessionService.listSessions({
        queries: [
          {
            query: {
              case: 'userAgentQuery',
              value: {},
            },
          },
        ],
      }),
    ).pipe(
      mergeMap(({ sessions }) => from(sessions)),
      withLatestFrom(this.user$),
      filter(([s, user]) => s.factors?.user?.loginName !== user.preferredLoginName),
      map(([s]) => ({
        displayName: s.factors?.user?.displayName ?? '',
        avatarUrl: '',
        loginName: s.factors?.user?.loginName ?? '',
        authState: V2SessionState.ACTIVE,
        userName: s.factors?.user?.loginName ?? '',
      })),
      map((s) => [s.loginName, s] as const),
      toArray(),
      map((sessions) => Array.from(new Map(sessions).values())), // Ensure unique loginNames
    );
  }

  public editUserProfile(): void {
    this.router.navigate(['users/me']).then();
    this.closedCard.emit();
  }

  public selectAccount(loginHint: string): void {
    const configWithPrompt: Partial<AuthConfig> = {
      customQueryParams: {
        login_hint: loginHint,
      },
    };
    this.authService.authenticate(configWithPrompt).then();
  }

  public selectNewAccount(): void {
    const configWithPrompt: Partial<AuthConfig> = {
      customQueryParams: {
        prompt: 'login',
      } as any,
    };
    this.authService.authenticate(configWithPrompt).then();
  }

  public logout(): void {
    const lP = JSON.stringify(this.labelpolicy());
    localStorage.setItem('labelPolicyOnSignout', lP);

    this.authService.signout();
    this.closedCard.emit();
  }
}

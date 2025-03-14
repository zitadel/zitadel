import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';
import { Session, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { SessionService } from 'src/app/services/session.service';
import { ListMyUserSessionsRequest } from '@zitadel/proto/zitadel/auth_pb';

@Component({
  selector: 'cnsl-accounts-card',
  templateUrl: './accounts-card.component.html',
  styleUrls: ['./accounts-card.component.scss'],
})
export class AccountsCardComponent implements OnInit {
  @Input() public user?: User.AsObject;
  @Input() public iamuser: boolean | null = false;

  @Output() public closedCard: EventEmitter<void> = new EventEmitter();
  public sessions: Session.AsObject[] = [];
  public loadingUsers: boolean = false;
  public UserState: any = UserState;
  private labelpolicy = toSignal(this.userService.labelpolicy$, { initialValue: undefined });

  constructor(
    public authService: AuthenticationService,
    private router: Router,
    private userService: GrpcAuthService,
    private sessionService: SessionService,
  ) {
    this.userService
      .listMyUserSessions()
      .then((sessions) => {
        this.sessions = sessions.resultList.filter((user) => user.loginName !== this.user?.preferredLoginName);
        this.loadingUsers = false;
      })
      .catch(() => {
        this.loadingUsers = false;
      });

    this.sessionService.listMyUserSessions({}).then((sessions) => {
      console.log('sessions', sessions);
    });
  }

  ngOnInit(): void {
    this.loadingUsers = true;
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

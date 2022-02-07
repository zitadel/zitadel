import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';
import { Session, User } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-accounts-card',
  templateUrl: './accounts-card.component.html',
  styleUrls: ['./accounts-card.component.scss'],
})
export class AccountsCardComponent implements OnInit {
  @Input() public user!: User.AsObject;
  @Input() public iamuser: boolean | null = false;

  @Output() public closedCard: EventEmitter<void> = new EventEmitter();
  public sessions: Session.AsObject[] = [];
  public loadingUsers: boolean = false;
  constructor(public authService: AuthenticationService, private router: Router, private userService: GrpcAuthService) {
    this.userService
      .listMyUserSessions()
      .then((sessions) => {
        this.sessions = sessions.resultList;
        const index = this.sessions.findIndex((user) => user.loginName === this.user.preferredLoginName);
        if (index > -1) {
          this.sessions.splice(index, 1);
        }

        this.loadingUsers = false;
      })
      .catch(() => {
        this.loadingUsers = false;
      });
  }

  public ngOnInit(): void {
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

  public close(): void {
    this.closedCard.emit();
  }

  public selectAccount(loginHint?: string): void {
    const configWithPrompt: Partial<AuthConfig> = {
      customQueryParams: {
        // prompt: 'select_account',
      } as any,
    };
    if (loginHint) {
      (configWithPrompt as any).customQueryParams['login_hint'] = loginHint;
    }
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
    this.authService.signout();
    this.closedCard.emit();
  }

  public get isOnSystem(): boolean {
    return (
      ['/system', '/views', '/failed-events'].includes(this.router.url) ||
      new RegExp('/system/policy/*').test(this.router.url)
    );
  }
}

import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { AuthConfig } from 'angular-oauth2-oidc';
import { UserProfile, UserSessionView } from 'src/app/proto/generated/auth_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { AuthService } from 'src/app/services/auth.service';

@Component({
    selector: 'app-accounts-card',
    templateUrl: './accounts-card.component.html',
    styleUrls: ['./accounts-card.component.scss'],
})
export class AccountsCardComponent implements OnInit {
    @Input() public profile!: UserProfile.AsObject;
    @Input() public iamuser: boolean = false;

    @Output() public close: EventEmitter<void> = new EventEmitter();
    public users: UserSessionView.AsObject[] = [];
    public loadingUsers: boolean = false;
    constructor(public authService: AuthService, private router: Router, private userService: AuthUserService) {
        this.userService.getMyUserSessions().then(sessions => {
            this.users = sessions.toObject().userSessionsList;

            const index = this.users.findIndex(user => user.userName === this.profile.userName);
            this.users.splice(index, 1);

            this.loadingUsers = false;
        }).catch(() => {
            this.loadingUsers = false;
        });
    }

    public ngOnInit(): void {
        this.loadingUsers = true;
    }

    public editUserProfile(): void {
        this.router.navigate(['user/me']);
        this.close.emit();
    }

    public closeCard(element: HTMLElement): void {
        if (!element.classList.contains('dontcloseonclick')) {
            this.close.emit();
        }
    }

    public selectAccount(loginName?: string): void {
        const configWithPrompt: Partial<AuthConfig> = {
            customQueryParams: {
                // prompt: 'select_account',
            } as any,
        };
        if (loginName) {
            (configWithPrompt as any).customQueryParams['login_hint'] = loginName;
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
        this.close.emit();
    }
}

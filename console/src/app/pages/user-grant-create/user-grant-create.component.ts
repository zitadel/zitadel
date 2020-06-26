import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription } from 'rxjs';
import { Org } from 'src/app/proto/generated/auth_pb';
import { Project, ProjectRole, UserGrant } from 'src/app/proto/generated/management_pb';
import { AuthService } from 'src/app/services/auth.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-user-grant-create',
    templateUrl: './user-grant-create.component.html',
    styleUrls: ['./user-grant-create.component.scss'],
})
export class UserGrantCreateComponent implements OnDestroy {
    public org!: Org.AsObject;
    public userId: string = '';
    public projectId: string = '';
    public grantId: string = '';
    public rolesList: string[] = [];

    public STEPS: number = 3; // org, project, roles
    public currentCreateStep: number = 1;

    private subscription: Subscription = new Subscription();
    constructor(
        private authService: AuthService,
        private userService: MgmtUserService,
        private toast: ToastService,
        private _location: Location,
        private route: ActivatedRoute,
    ) {
        this.subscription = this.route.params.subscribe(({ id }: Params) => {
            this.userId = id;
        });

        this.authService.GetActiveOrg().then(org => {
            this.org = org;
        });
    }

    public close(): void {
        this._location.back();
    }

    public addGrant(): void {
        this.userService.CreateUserGrant(
            this.projectId,
            this.userId,
            this.rolesList,
        ).then((data: UserGrant) => {
            this.close();
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }

    public selectProject(project: Project.AsObject): void {
        this.projectId = project.id;
    }

    public selectRoles(roles: ProjectRole.AsObject[]): void {
        this.rolesList = roles.map(role => role.key);
    }

    public next(): void {
        this.currentCreateStep++;
    }

    public previous(): void {
        this.currentCreateStep--;
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }
}

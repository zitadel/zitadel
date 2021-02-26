import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription } from 'rxjs';
import { UserTarget } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { ProjectGrantView, ProjectRole, ProjectView, UserGrant, UserView } from 'src/app/proto/generated/zitadel/management_pb';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-user-grant-create',
    templateUrl: './user-grant-create.component.html',
    styleUrls: ['./user-grant-create.component.scss'],
})
export class UserGrantCreateComponent implements OnDestroy {
    public context!: UserGrantContext;

    public org!: Org.AsObject;
    public userId: string = '';

    public projectId: string = '';
    public project!: ProjectGrantView.AsObject | ProjectView.AsObject;

    public grantId: string = '';
    public rolesList: string[] = [];

    public STEPS: number = 2; // project, roles
    public currentCreateStep: number = 1;

    public filterValue: string = '';

    private subscription: Subscription = new Subscription();

    public UserGrantContext: any = UserGrantContext;

    public grantRolesKeyList: string[] = [];

    public user!: UserView.AsObject;
    public UserTarget: any = UserTarget;

    public ProjectGrantView: any = ProjectGrantView;
    public ProjectView: any = ProjectView;
    constructor(
        private userService: ManagementService,
        private toast: ToastService,
        private _location: Location,
        private route: ActivatedRoute,
        private authService: GrpcAuthService,
        private mgmtService: ManagementService,
    ) {
        this.subscription = this.route.params.subscribe((params: Params) => {
            const { projectid, grantid, userid } = params;
            this.context = UserGrantContext.NONE;

            this.projectId = projectid;
            this.grantId = grantid;
            this.userId = userid;

            if (this.projectId && !this.grantId) {
                this.context = UserGrantContext.OWNED_PROJECT;
            } else if (this.projectId && this.grantId) {
                this.context = UserGrantContext.GRANTED_PROJECT;
                this.mgmtService.GetGrantedProjectByID(this.projectId, this.grantId).then(resp => {
                    this.grantRolesKeyList = resp.toObject().roleKeysList;
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
            } else if (this.userId) {
                this.context = UserGrantContext.USER;
                this.mgmtService.GetUserByID(this.userId).then(resp => {
                    this.user = resp.toObject();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
            }
        });

        this.authService.getActiveOrg().then(org => {
            this.org = org;
        });
    }

    public close(): void {
        this._location.back();
    }

    public addGrant(): void {
        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                this.userService.CreateUserGrant(
                    this.userId,
                    this.rolesList,
                    this.projectId,
                ).then((data: UserGrant) => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
            case UserGrantContext.GRANTED_PROJECT:
                this.userService.CreateUserGrant(
                    this.userId,
                    this.rolesList,
                    this.projectId,
                    this.grantId,
                ).then((data: UserGrant) => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
            case UserGrantContext.USER:
                let grantId;

                if ((this.project as ProjectGrantView.AsObject)?.id) {
                    grantId = (this.project as ProjectGrantView.AsObject).id;
                }

                this.userService.CreateUserGrant(
                    this.userId,
                    this.rolesList,
                    this.project.projectId,
                    grantId,
                ).then((data: UserGrant) => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
            case UserGrantContext.NONE:
                let tempGrantId;

                if ((this.project as ProjectGrantView.AsObject)?.id) {
                    tempGrantId = (this.project as ProjectGrantView.AsObject).id;
                }

                this.userService.CreateUserGrant(
                    this.userId,
                    this.rolesList,
                    this.project.projectId,
                    tempGrantId,
                ).then((data: UserGrant) => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
        }

    }

    public selectProject(project: ProjectView.AsObject | ProjectGrantView.AsObject | any): void {
        this.project = project;
        this.projectId = project.projectId;
        this.grantRolesKeyList = project.roleKeysList ?? [];
    }

    public selectUser(user: UserView.AsObject): void {
        this.userId = user.id;
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

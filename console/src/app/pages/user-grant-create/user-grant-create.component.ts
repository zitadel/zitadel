import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription } from 'rxjs';
import { UserTarget } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrantedProject, Project, Role } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageKey, StorageService } from 'src/app/services/storage.service';
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
    public project!: GrantedProject.AsObject | Project.AsObject;

    public grantId: string = '';
    public rolesList: string[] = [];

    public STEPS: number = 2; // project, roles
    public currentCreateStep: number = 1;

    public filterValue: string = '';

    private subscription: Subscription = new Subscription();

    public UserGrantContext: any = UserGrantContext;

    public grantedRoleKeysList: string[] = [];

    public user!: User.AsObject;
    public UserTarget: any = UserTarget;

    public ProjectGrantView: any = GrantedProject;
    public ProjectView: any = Project;
    constructor(
        private userService: ManagementService,
        private toast: ToastService,
        private _location: Location,
        private route: ActivatedRoute,
        private authService: GrpcAuthService,
        private mgmtService: ManagementService,
        private storage: StorageService,
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
                this.mgmtService.getGrantedProjectByID(this.projectId, this.grantId).then(resp => {
                    if (resp.grantedProject?.grantedRoleKeysList) {
                        this.grantedRoleKeysList = resp.grantedProject?.grantedRoleKeysList;
                    }
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
            } else if (this.userId) {
                this.context = UserGrantContext.USER;
                this.mgmtService.getUserByID(this.userId).then(resp => {
                    if (resp.user) {
                        this.user = resp.user;
                    }
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
            }
        });

        const temporg = this.storage.getItem<Org.AsObject>(StorageKey.organization);
        if (temporg) {
            this.org = temporg;
        }
    }

    public close(): void {
        this._location.back();
    }

    public addGrant(): void {
        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                this.userService.addUserGrant(
                    this.userId,
                    this.rolesList,
                    this.projectId,
                ).then(() => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
            case UserGrantContext.GRANTED_PROJECT:
                this.userService.addUserGrant(
                    this.userId,
                    this.rolesList,
                    this.projectId,
                    this.grantId,
                ).then(() => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
            case UserGrantContext.USER:
                let grantId;

                if ((this.project as GrantedProject.AsObject)?.grantId) {
                    grantId = (this.project as GrantedProject.AsObject).grantId;
                }

                this.userService.addUserGrant(
                    this.userId,
                    this.rolesList,
                    this.projectId,
                    grantId,
                ).then(() => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
            case UserGrantContext.NONE:
                console.log('none');
                let tempGrantId;

                if ((this.project as GrantedProject.AsObject)?.grantId) {
                    tempGrantId = (this.project as GrantedProject.AsObject).grantId;
                }

                this.userService.addUserGrant(
                    this.userId,
                    this.rolesList,
                    this.projectId,
                    tempGrantId,
                ).then(() => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
                    this.close();
                }).catch((error: any) => {
                    this.toast.showError(error);
                });
                break;
        }
    }

    public selectProject(project: Project.AsObject | GrantedProject.AsObject | any): void {
        this.project = project;
        this.projectId = project.id || project.projectId;

        this.grantedRoleKeysList = project.grantedRoleKeysList ?? [];
    }

    public selectUser(user: User.AsObject): void {
        this.userId = user.id;
    }

    public selectRoles(roles: Role.AsObject[]): void {
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

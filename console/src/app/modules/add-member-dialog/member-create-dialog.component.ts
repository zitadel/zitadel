import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { GrantedProject, Project, Role } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectAutocompleteType } from '../search-project-autocomplete/search-project-autocomplete.component';

export enum CreationType {
    PROJECT_OWNED = 0,
    PROJECT_GRANTED = 1,
    ORG = 2,
    IAM = 3,
}
@Component({
    selector: 'app-member-create-dialog',
    templateUrl: './member-create-dialog.component.html',
    styleUrls: ['./member-create-dialog.component.scss'],
})
export class MemberCreateDialogComponent {
    private projectId: string = '';
    private grantId: string = '';
    public preselectedUsers: Array<User.AsObject> = [];

    public creationType!: CreationType;

    /**
     *  Specifies options for creating members,
     *  without ending $, to enable write event permission even if user is allowed
     *  to create members for only one specific project.
     */
    public creationTypes: Array<{ type: CreationType, disabled$: Observable<boolean>; }> = [
        { type: CreationType.IAM, disabled$: this.authService.isAllowed(['iam.member.write$']) },
        { type: CreationType.ORG, disabled$: this.authService.isAllowed(['org.member.write$']) },
        { type: CreationType.PROJECT_OWNED, disabled$: this.authService.isAllowed(['project.member.write']) },
        { type: CreationType.PROJECT_GRANTED, disabled$: this.authService.isAllowed(['project.grant.member.write']) },
    ];
    public users: Array<User.AsObject> = [];
    public roles: Array<Role.AsObject> | string[] = [];
    public CreationType: any = CreationType;
    public ProjectAutocompleteType: any = ProjectAutocompleteType;
    public memberRoleOptions: string[] = [];

    public showCreationTypeSelector: boolean = false;
    constructor(
        private mgmtService: ManagementService,
        private adminService: AdminService,
        private authService: GrpcAuthService,
        public dialogRef: MatDialogRef<MemberCreateDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        private toastService: ToastService,
    ) {
        if (data?.projectId) {
            this.projectId = data.projectId;
        }
        if (data?.user) {
            this.preselectedUsers = [data.user];
            this.users = [data.user];
        }

        if (data?.creationType !== undefined) {
            this.creationType = data.creationType;
            this.loadRoles();
        } else {
            this.showCreationTypeSelector = true;
        }
    }

    public loadRoles(): void {
        switch (this.creationType) {
            case CreationType.PROJECT_GRANTED:
                this.mgmtService.listProjectGrantMemberRoles().then(resp => {
                    this.memberRoleOptions = resp.resultList;
                }).catch(error => {
                    this.toastService.showError(error);
                });
                break;
            case CreationType.PROJECT_OWNED:
                this.mgmtService.listProjectMemberRoles().then(resp => {
                    this.memberRoleOptions = resp.resultList;
                }).catch(error => {
                    this.toastService.showError(error);
                });
                break;
            case CreationType.IAM:
                this.adminService.listIAMMemberRoles().then(resp => {
                    this.memberRoleOptions = resp.rolesList;
                }).catch(error => {
                    this.toastService.showError(error);
                });
                break;
        }
    }

    public selectProject(project: Project.AsObject | GrantedProject.AsObject | any): void {
        if (project.projectId && project.grantId) {
            this.projectId = project.projectId;
            this.grantId = project.grantId;
        }
        else if (project.id) {
            this.projectId = project.id;
        }
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({
            users: this.users,
            roles: this.roles,
            creationType: this.creationType,
            projectId: this.projectId,
            grantId: this.grantId,
        });
    }

    public setOrgMemberRoles(roles: string[]): void {
        this.roles = roles;
    }
}

import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ProjectRole, User } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

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
    public projectId: string = '';
    public creationType!: CreationType;
    public users: Array<User.AsObject> = [];
    public roles: Array<ProjectRole.AsObject> | string[] = [];
    public CreationType: any = CreationType;
    public memberRoleOptions: string[] = [];
    constructor(
        private projectService: ProjectService,
        private adminService: AdminService,
        public dialogRef: MatDialogRef<MemberCreateDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        toastService: ToastService,
    ) {
        this.creationType = data.creationType;
        this.projectId = data.projectId;

        if (this.creationType === CreationType.PROJECT_GRANTED) {
            this.projectService.GetProjectGrantMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                toastService.showError(error.message);
            });
        } else if (this.creationType === CreationType.PROJECT_OWNED) {
            this.projectService.GetProjectMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                toastService.showError(error.message);
            });
        } else if (this.creationType === CreationType.IAM) {
            this.adminService.GetIamMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                toastService.showError(error.message);
            });
        }
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ users: this.users, roles: this.roles });
    }

    public setOrgMemberRoles(roles: string[]): void {
        this.roles = roles;
    }
}

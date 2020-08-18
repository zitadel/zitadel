import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ProjectGrantView, ProjectRole, ProjectView, User } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ProjectService } from 'src/app/services/project.service';
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
    public creationType!: CreationType;
    public creationTypes: CreationType[] = [
        CreationType.IAM,
        CreationType.ORG,
        CreationType.PROJECT_OWNED,
        CreationType.PROJECT_GRANTED,
    ];
    public users: Array<User.AsObject> = [];
    public roles: Array<ProjectRole.AsObject> | string[] = [];
    public CreationType: any = CreationType;
    public ProjectAutocompleteType: any = ProjectAutocompleteType;
    public memberRoleOptions: string[] = [];

    public showCreationTypeSelector: boolean = false;
    constructor(
        private projectService: ProjectService,
        private adminService: AdminService,
        public dialogRef: MatDialogRef<MemberCreateDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        private toastService: ToastService,
    ) {
        if (data?.projectId) {
            this.projectId = data.projectId;
        }
        if (data?.creationType) {
            this.creationType = data.creationType;
            this.loadRoles();
        } else {
            this.showCreationTypeSelector = true;
        }
    }

    public loadRoles(): void {
        console.log(this.creationType);
        switch (this.creationType) {
            case CreationType.PROJECT_GRANTED:
                this.projectService.GetProjectGrantMemberRoles().then(resp => {
                    this.memberRoleOptions = resp.toObject().rolesList;
                }).catch(error => {
                    this.toastService.showError(error);
                });
                break;
            case CreationType.PROJECT_GRANTED:
                this.projectService.GetProjectMemberRoles().then(resp => {
                    this.memberRoleOptions = resp.toObject().rolesList;
                }).catch(error => {
                    this.toastService.showError(error);
                });
                break;
            case CreationType.IAM:
                this.adminService.GetIamMemberRoles().then(resp => {
                    this.memberRoleOptions = resp.toObject().rolesList;
                }).catch(error => {
                    this.toastService.showError(error);
                });
                break;
        }
    }

    public selectProject(project: ProjectView.AsObject | ProjectGrantView.AsObject | any): void {
        this.projectId = project.projectId;
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ users: this.users, roles: this.roles, projectId: this.projectId });
    }

    public setOrgMemberRoles(roles: string[]): void {
        this.roles = roles;
    }
}

import { Component, Inject } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { User } from 'src/app/proto/generated/management_pb';

export interface ProjectGrantMembersCreateDialogExportType {
    userIds: string[];
    rolesList: string[];
}
@Component({
    selector: 'app-project-grant-members-create-dialog',
    templateUrl: './project-grant-members-create-dialog.component.html',
    styleUrls: ['./project-grant-members-create-dialog.component.scss'],
})
export class ProjectGrantMembersCreateDialogComponent {
    public form!: FormGroup;
    public userIds: string[] = [];
    public rolesList: string[] = [];

    constructor(
        public dialogRef: MatDialogRef<ProjectGrantMembersCreateDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) { }

    public selectUsers(users: User.AsObject[]): void {
        this.userIds = users.map(user => user.id);
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        const exportData: ProjectGrantMembersCreateDialogExportType = { userIds: this.userIds, rolesList: this.rolesList };
        this.dialogRef.close(exportData);
    }
}

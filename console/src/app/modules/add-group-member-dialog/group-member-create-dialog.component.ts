import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { User } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
  selector: 'cnsl-group-member-create-dialog',
  templateUrl: './group-member-create-dialog.component.html',
  styleUrls: ['./group-member-create-dialog.component.scss'],
})
export class GroupMemberCreateDialogComponent {
  private grantId: string = '';
  public preselectedUsers: Array<User.AsObject> = [];

  public users: Array<User.AsObject> = [];
  public roles: string[] = [];
  public memberRoleOptions: string[] = [];

  constructor(
    public dialogRef: MatDialogRef<GroupMemberCreateDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    if (data?.user) {
      this.preselectedUsers = [data.user];
      this.users = [data.user];
    }

  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({
      users: this.users,
      roles: this.roles,
      grantId: this.grantId,
    });
  }
}

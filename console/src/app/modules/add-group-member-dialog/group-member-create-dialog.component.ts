import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';


@Component({
  selector: 'cnsl-group-member-create-dialog',
  templateUrl: './group-member-create-dialog.component.html',
  styleUrls: ['./group-member-create-dialog.component.scss'],
})
export class GroupMemberCreateDialogComponent {
  private projectId: string = '';
  private grantId: string = '';
  public preselectedUsers: Array<User.AsObject> = [];

  public users: Array<User.AsObject> = [];
  public roles: string[] = [];
  public memberRoleOptions: string[] = [];

  constructor(
    private mgmtService: ManagementService,
    private adminService: AdminService,
    private authService: GrpcAuthService,
    public dialogRef: MatDialogRef<GroupMemberCreateDialogComponent>,
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

  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({
      users: this.users,
      roles: this.roles,
      projectId: this.projectId,
      grantId: this.grantId,
    });
  }
}

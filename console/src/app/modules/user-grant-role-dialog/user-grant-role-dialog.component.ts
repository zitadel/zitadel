import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { UserGrant } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
  selector: 'cnsl-user-grant-role-dialog',
  templateUrl: './user-grant-role-dialog.component.html',
  styleUrls: ['./user-grant-role-dialog.component.scss'],
})
export class UserGrantRoleDialogComponent {
  public grant!: UserGrant.AsObject;
  public selectedRoles: Role.AsObject[] = [];

  constructor(public dialogRef: MatDialogRef<UserGrantRoleDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.grant = data.grant;
  }

  public selectRoles(selected: any): void {
    this.selectedRoles = selected;
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({ roles: this.selectedRoles.map((r) => r.key) });
  }
}

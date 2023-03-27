import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { getMembershipColor } from 'src/app/utils/color';

@Component({
  selector: 'cnsl-add-member-roles-dialog',
  templateUrl: './add-member-roles-dialog.component.html',
  styleUrls: ['./add-member-roles-dialog.component.scss'],
})
export class AddMemberRolesDialogComponent {
  public allRoles: string[] = [];
  public selectedRoles: string[] = [];

  constructor(public dialogRef: MatDialogRef<AddMemberRolesDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.allRoles = Object.assign([], data.allRoles);
    this.selectedRoles = Object.assign([], data.selectedRoles);
  }

  public closeDialogWithRoles(): void {
    this.dialogRef.close(this.selectedRoles);
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public toggleRole(role: string): void {
    const index = this.selectedRoles.findIndex((r) => r === role);
    if (index > -1) {
      this.selectedRoles.splice(index, 1);
    } else {
      this.selectedRoles.push(role);
    }
  }

  public getColor(role: string) {
    return getMembershipColor(role)[500];
  }
}

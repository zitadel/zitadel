import { Component, inject, Inject, signal } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

export type UserGrantRoleDialogData = {
  projectId: string;
  grantId?: string;
  selectedRoleKeysList: string[];
  i18nTitle: string;
};

export type UserGrantRoleDialogResult = { roles: string[] };

@Component({
  selector: 'cnsl-user-grant-role-dialog',
  templateUrl: './user-grant-role-dialog.component.html',
  styleUrls: ['./user-grant-role-dialog.component.scss'],
  standalone: false,
})
export class UserGrantRoleDialogComponent {
  protected readonly data: UserGrantRoleDialogData = inject(MAT_DIALOG_DATA);
  private readonly dialogRef: MatDialogRef<UserGrantRoleDialogComponent, UserGrantRoleDialogResult> = inject(MatDialogRef);
  private readonly selectedRoleKeys = signal([...this.data.selectedRoleKeysList]);

  public selectRoles(selected: string[]): void {
    this.selectedRoleKeys.set(selected);
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({ roles: this.selectedRoleKeys() });
  }
}

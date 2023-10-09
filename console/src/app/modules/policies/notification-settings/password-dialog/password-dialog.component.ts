import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';

@Component({
  selector: 'cnsl-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.scss'],
})
export class PasswordDialogComponent {
  public password: string = '';
  constructor(public dialogRef: MatDialogRef<PasswordDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {}

  closeDialog(password: string = ''): void {
    this.dialogRef.close(password);
  }
}

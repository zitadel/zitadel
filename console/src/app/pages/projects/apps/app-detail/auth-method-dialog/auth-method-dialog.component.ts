import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';

@Component({
  selector: 'cnsl-auth-method-dialog',
  templateUrl: './auth-method-dialog.component.html',
  styleUrls: ['./auth-method-dialog.component.scss'],
})
export class AuthMethodDialogComponent {
  public authmethod: string = '';
  constructor(public dialogRef: MatDialogRef<AuthMethodDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.authmethod = data.initialAuthMethod;
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithMethod(): void {
    this.dialogRef.close(this.authmethod);
  }
}

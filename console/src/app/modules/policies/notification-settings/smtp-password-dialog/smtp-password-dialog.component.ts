import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-smtp-password-dialog',
  templateUrl: './smtp-password-dialog.component.html',
  styleUrls: ['./smtp-password-dialog.component.scss'],
})
export class SMTPPasswordDialogComponent {
  public password: string = '';
  constructor(public dialogRef: MatDialogRef<SMTPPasswordDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {}

  closeDialog(password: string = ''): void {
    this.dialogRef.close(password);
  }
}

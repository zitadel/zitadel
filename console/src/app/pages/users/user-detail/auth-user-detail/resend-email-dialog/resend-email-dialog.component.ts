import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

export type ResendEmailDialogData = {
  email: string | '';
};

export type ResendEmailDialogResult = { send: true; email: string } | { send: false };

@Component({
  selector: 'cnsl-resend-email-dialog',
  templateUrl: './resend-email-dialog.component.html',
  styleUrls: ['./resend-email-dialog.component.scss'],
})
export class ResendEmailDialogComponent {
  public email: string = '';
  constructor(
    public dialogRef: MatDialogRef<ResendEmailDialogComponent, ResendEmailDialogResult>,
    @Inject(MAT_DIALOG_DATA) public data: ResendEmailDialogData,
  ) {
    if (data.email) {
      this.email = data.email;
    }
  }

  closeDialog(): void {
    this.dialogRef.close({ send: false });
  }

  closeDialogWithSend(email: string = ''): void {
    this.dialogRef.close({ send: true, email });
  }
}

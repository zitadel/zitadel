import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-password-dialog-sms-provider',
  templateUrl: './password-dialog-sms-provider.component.html',
  styleUrls: ['./password-dialog-sms-provider.component.scss'],
})
export class PasswordDialogSMSProviderComponent {
  public password: string = '';
  constructor(
    public dialogRef: MatDialogRef<PasswordDialogSMSProviderComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  closeDialog(password: string = ''): void {
    this.dialogRef.close(password);
  }
}

import { Component, Inject } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { AddSMSProviderTwilioRequest, UpdateSMSProviderTwilioTokenRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { PasswordDialogComponent } from '../password-dialog/password-dialog.component';

enum SMSProviderType {
  Twilio = 1,
}

@Component({
  selector: 'cnsl-dialog-add-sms-provider',
  templateUrl: './dialog-add-sms-provider.component.html',
  styleUrls: ['./dialog-add-sms-provider.component.scss'],
})
export class DialogAddSMSProviderComponent {
  public SMSProviderType: any = SMSProviderType;
  public availableSMSProviders: SMSProviderType[] = [SMSProviderType.Twilio];
  public provider: SMSProviderType = SMSProviderType.Twilio;
  public req: AddSMSProviderTwilioRequest = new AddSMSProviderTwilioRequest();

  public twilioForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private service: AdminService,
    public dialogRef: MatDialogRef<DialogAddSMSProviderComponent>,
    private toast: ToastService,
    private dialog: MatDialog,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.twilioForm = this.fb.group({
      sid: ['', [Validators.required]],
      token: ['', [Validators.required]],
      senderNumber: ['', [Validators.required]],
    });
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithRequest(): void {
    this.req.setSid(this.sid?.value);
    this.req.setToken(this.token?.value);
    this.req.setSenderNumber(this.senderNumber?.value);

    this.dialogRef.close(this.req);
  }

  public changeToken(): void {
    const dialogRef = this.dialog.open(PasswordDialogComponent, {
      width: '400px',
      data: {
        i18nTitle: 'SETTING.SMS.TWILIO.SETTOKEN',
        i18nLabel: 'SETTING.SMS.TWILIO.TOKEN',
      },
    });

    dialogRef.afterClosed().subscribe((token: string) => {
      if (token) {
        const tokenReq = new UpdateSMSProviderTwilioTokenRequest();
        tokenReq.setToken(token);

        this.service
          .updateSMSProviderTwilioToken(tokenReq)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.TWILIO.TOKENSET', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public get senderNumber(): AbstractControl | null {
    return this.twilioForm.get('senderNumber');
  }

  public get token(): AbstractControl | null {
    return this.twilioForm.get('token');
  }

  public get sid(): AbstractControl | null {
    return this.twilioForm.get('sid');
  }
}

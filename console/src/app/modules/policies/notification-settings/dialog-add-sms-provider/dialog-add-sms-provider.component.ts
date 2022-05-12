import { Component, Inject } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { AddSMSProviderTwilioRequest } from 'src/app/proto/generated/zitadel/admin_pb';

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
    public dialogRef: MatDialogRef<DialogAddSMSProviderComponent>,
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

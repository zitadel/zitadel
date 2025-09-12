import { Component, Inject } from '@angular/core';
import { AbstractControl, FormControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import {
  AddSMSProviderTwilioRequest,
  UpdateSMSProviderTwilioRequest,
  UpdateSMSProviderTwilioTokenRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { SMSProvider, TwilioConfig } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { PasswordDialogSMSProviderComponent } from '../password-dialog-sms-provider/password-dialog-sms-provider.component';

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
  public req!: AddSMSProviderTwilioRequest | UpdateSMSProviderTwilioRequest;

  public twilioForm!: UntypedFormGroup;

  private smsProviders: SMSProvider.AsObject[] = [];

  constructor(
    private fb: UntypedFormBuilder,
    private service: AdminService,
    public dialogRef: MatDialogRef<DialogAddSMSProviderComponent>,
    private toast: ToastService,
    private dialog: MatDialog,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.twilioForm = this.fb.group({
      sid: ['', [requiredValidator]],
      senderNumber: [''],
      // NB: not required if not using verification service
      verifyServiceSid: [''],
    });

    this.smsProviders = data.smsProviders;
    if (!!this.twilio) {
      this.twilioForm.patchValue(this.twilio);
    } else {
      this.twilioForm.addControl('token', new FormControl('', requiredValidator));
    }
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithRequest(): void {
    if (!!this.twilio && this.twilioProvider && this.twilioProvider.id) {
      this.req = new UpdateSMSProviderTwilioRequest();
      this.req.setId(this.twilioProvider.id);
      this.req.setSid(this.sid?.value);
      this.req.setSenderNumber(this.senderNumber?.value);
      this.req.setVerifyServiceSid(this.verifyServiceSid?.value ?? '');
      this.dialogRef.close(this.req);
    } else {
      this.req = new AddSMSProviderTwilioRequest();
      this.req.setSid(this.sid?.value);
      this.req.setToken(this.token?.value);
      this.req.setSenderNumber(this.senderNumber?.value);
      this.req.setVerifyServiceSid(this.verifyServiceSid?.value ?? '');
      this.dialogRef.close(this.req);
    }
  }

  public changeToken(): void {
    const dialogRef = this.dialog.open(PasswordDialogSMSProviderComponent, {
      width: '400px',
      data: {
        i18nTitle: 'SETTING.SMS.TWILIO.SETTOKEN',
        i18nLabel: 'SETTING.SMS.TWILIO.TOKEN',
      },
    });

    dialogRef.afterClosed().subscribe((token: string) => {
      if (token && this.twilioProvider?.id) {
        const tokenReq = new UpdateSMSProviderTwilioTokenRequest();
        tokenReq.setToken(token);
        tokenReq.setId(this.twilioProvider.id);

        this.service
          .updateSMSProviderTwilioToken(tokenReq)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.TWILIO.TOKENSET', true);
            this.dialogRef.close();
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

  public get verifyServiceSid(): AbstractControl | null {
    return this.twilioForm.get('verifyServiceSid');
  }

  public get token(): AbstractControl | null {
    return this.twilioForm.get('token');
  }

  public get sid(): AbstractControl | null {
    return this.twilioForm.get('sid');
  }

  public get twilioProvider(): SMSProvider.AsObject | undefined {
    const twilioProvider: SMSProvider.AsObject | undefined = this.smsProviders.find((p) => p.twilio);
    if (twilioProvider) {
      return twilioProvider;
    } else {
      return undefined;
    }
  }

  public get twilio(): TwilioConfig.AsObject | undefined {
    const twilioProvider: SMSProvider.AsObject | undefined = this.smsProviders.find((p) => p.twilio);
    if (twilioProvider && !!twilioProvider.twilio) {
      return twilioProvider.twilio;
    } else {
      return undefined;
    }
  }
}

import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import {
    AddSMSProviderTwilioRequest,
    UpdateSMTPConfigPasswordResponse,
    UpdateSMTPConfigRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { SMSProvider, SMSProviderConfigState } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';

@Component({
  selector: 'cnsl-notification-settings',
  templateUrl: './notification-settings.component.html',
  styleUrls: ['./notification-settings.component.scss'],
})
export class NotificationSettingsComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public smsProviders: SMSProvider.AsObject[] = [];

  public loading: boolean = false;
  public form!: FormGroup;

  public SMSProviderConfigState: any = SMSProviderConfigState;
  public InfoSectionType: any = InfoSectionType;

  // show available providers

  constructor(
    private service: AdminService,
    private dialog: MatDialog,
    private toast: ToastService,
    private fb: FormBuilder,
  ) {
    this.form = this.fb.group({
      senderAddress: ['', [Validators.required]],
      senderName: ['', [Validators.required]],
      tls: [true, [Validators.required]],
      host: ['', [Validators.required]],
      user: ['', [Validators.required]],
      password: ['', [Validators.required]],
    });
  }

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.service
      .getSMTPConfig()
      .then((smtpConfig) => {
        if (smtpConfig.smtpConfig) {
          this.form.patchValue(smtpConfig.smtpConfig);
        }
      })
      .catch((error) => {
        if (error && error.code === 5) {
          console.log(error);
        }
      });

    this.service.listSMSProviders().then((smsProviders) => {
      if (smsProviders.resultList) {
        this.smsProviders = smsProviders.resultList;
        console.log(this.smsProviders);
      }
    });
  }

  private updateData(): Promise<UpdateSMTPConfigPasswordResponse.AsObject> | any {
    const req = new UpdateSMTPConfigRequest();
    req.setHost(this.host?.value ?? '');
    req.setSenderAddress(this.senderAddress?.value ?? '');
    req.setSenderName(this.senderName?.value ?? '');
    req.setTls(this.tls?.value ?? false);
    req.setUser(this.user?.value ?? '');

    console.log(req.toObject());

    // return this.service.updateSMTPConfig(req).then(() => {
    //   let passwordReq: UpdateSMTPConfigPasswordRequest;
    //   if (this.password) {
    //     passwordReq = new UpdateSMTPConfigPasswordRequest();
    //     passwordReq.setPassword(this.password.value);
    //     return this.service.updateSMTPConfigPassword(passwordReq);
    //   } else {
    //     return;
    //   }
    // });
  }

  public savePolicy(): void {
    const prom = this.updateData();
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
          this.loading = true;
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error: unknown) => {
          this.toast.showError(error);
        });
    }
  }

  public addSMSProvider(): void {
    const dialogRef = this.dialog.open(DialogAddSMSProviderComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'IDP.DELETE_TITLE',
        descriptionKey: 'IDP.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: AddSMSProviderTwilioRequest) => {
      if (req) {
        this.service
          .addSMSProviderTwilio(req)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.TWILIO.ADDED', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public get twilio(): SMSProvider.AsObject | undefined {
    return this.smsProviders.find((p) => p.twilio);
  }

  public get senderAddress(): AbstractControl | null {
    return this.form.get('senderAddress');
  }

  public get senderName(): AbstractControl | null {
    return this.form.get('senderName');
  }

  public get tls(): AbstractControl | null {
    return this.form.get('tls');
  }

  public get user(): AbstractControl | null {
    return this.form.get('user');
  }

  public get host(): AbstractControl | null {
    return this.form.get('host');
  }

  public get password(): AbstractControl | null {
    return this.form.get('password');
  }
}

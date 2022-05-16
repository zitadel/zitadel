import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import {
    AddSMSProviderTwilioRequest,
    UpdateSMTPConfigPasswordRequest,
    UpdateSMTPConfigPasswordResponse,
    UpdateSMTPConfigRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { DebugNotificationProvider, SMSProvider, SMSProviderConfigState } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';
import { SMTPPasswordDialogComponent } from './smtp-password-dialog/smtp-password-dialog.component';

@Component({
  selector: 'cnsl-notification-settings',
  templateUrl: './notification-settings.component.html',
  styleUrls: ['./notification-settings.component.scss'],
})
export class NotificationSettingsComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public smsProviders: SMSProvider.AsObject[] = [];
  public logNotificationProvider!: DebugNotificationProvider.AsObject;
  public fileNotificationProvider!: DebugNotificationProvider.AsObject;

  public smtpLoading: boolean = false;
  public smsProvidersLoading: boolean = false;
  public logProviderLoading: boolean = false;
  public fileProviderLoading: boolean = false;

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
    });
  }

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.smtpLoading = true;
    this.service
      .getSMTPConfig()
      .then((smtpConfig) => {
        this.smtpLoading = false;
        if (smtpConfig.smtpConfig) {
          this.form.patchValue(smtpConfig.smtpConfig);
        }
      })
      .catch((error) => {
        this.smtpLoading = false;
        if (error && error.code === 5) {
          console.log(error);
        }
      });

    this.smsProvidersLoading = true;
    this.service
      .listSMSProviders()
      .then((smsProviders) => {
        this.smsProvidersLoading = false;
        if (smsProviders.resultList) {
          this.smsProviders = smsProviders.resultList;
        }
      })
      .catch((error) => {
        this.smsProvidersLoading = false;
        this.toast.showError(error);
      });

    this.logProviderLoading = true;
    this.service
      .getLogNotificationProvider()
      .then((logNotificationProvider) => {
        this.logProviderLoading = false;
        if (logNotificationProvider.provider) {
          this.logNotificationProvider = logNotificationProvider.provider;
        }
      })
      .catch((error) => {
        this.logProviderLoading = false;
        this.toast.showError(error);
      });

    this.fileProviderLoading = true;
    this.service
      .getFileSystemNotificationProvider()
      .then((fileNotificationProvider) => {
        this.fileProviderLoading = false;
        if (fileNotificationProvider.provider) {
          this.fileNotificationProvider = fileNotificationProvider.provider;
        }
      })
      .catch((error) => {
        this.fileProviderLoading = false;
        this.toast.showError(error);
      });
  }

  private updateData(): Promise<UpdateSMTPConfigPasswordResponse.AsObject> | any {
    const req = new UpdateSMTPConfigRequest();
    req.setHost(this.host?.value ?? '');
    req.setSenderAddress(this.senderAddress?.value ?? '');
    req.setSenderName(this.senderName?.value ?? '');
    req.setTls(this.tls?.value ?? false);
    req.setUser(this.user?.value ?? '');

    return this.service.updateSMTPConfig(req).catch(this.toast.showError);
  }

  public savePolicy(): void {
    const prom = this.updateData();
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
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

  public setSMTPPassword(): void {
    const dialogRef = this.dialog.open(SMTPPasswordDialogComponent, {
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((password: string) => {
      if (password) {
        const passwordReq = new UpdateSMTPConfigPasswordRequest();
        passwordReq.setPassword(password);

        this.service
          .updateSMTPConfigPassword(passwordReq)
          .then(() => {
            this.toast.showInfo('SETTING.SMTP.PASSWORDSET', true);
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
}

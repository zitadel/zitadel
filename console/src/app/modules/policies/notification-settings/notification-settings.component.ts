import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { take } from 'rxjs';
import {
  AddSMSProviderTwilioRequest,
  AddSMTPConfigRequest,
  AddSMTPConfigResponse,
  UpdateSMSProviderTwilioRequest,
  UpdateSMTPConfigPasswordRequest,
  UpdateSMTPConfigRequest,
  UpdateSMTPConfigResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { DebugNotificationProvider, SMSProvider, SMSProviderConfigState } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';
import { PasswordDialogComponent } from './password-dialog/password-dialog.component';

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

  public form!: UntypedFormGroup;

  public SMSProviderConfigState: any = SMSProviderConfigState;
  public InfoSectionType: any = InfoSectionType;

  public hasSMTPConfig: boolean = false;

  // show available providers

  constructor(
    private service: AdminService,
    private dialog: MatDialog,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
    private authService: GrpcAuthService,
  ) {
    this.form = this.fb.group({
      senderAddress: [{ disabled: true, value: '' }, [requiredValidator]],
      senderName: [{ disabled: true, value: '' }, [requiredValidator]],
      tls: [{ disabled: true, value: true }, [requiredValidator]],
      hostAndPort: [{ disabled: true, value: '' }, [requiredValidator]],
      user: [{ disabled: true, value: '' }, [requiredValidator]],
    });
  }

  ngOnInit(): void {
    this.fetchData();
    this.authService
      .isAllowed(['iam.write'])
      .pipe(take(1))
      .subscribe((allowed) => {
        if (allowed) {
          this.form.enable();
        }
      });
  }

  private fetchData(): void {
    this.smtpLoading = true;
    this.service
      .getSMTPConfig()
      .then((smtpConfig) => {
        this.smtpLoading = false;
        if (smtpConfig.smtpConfig) {
          this.hasSMTPConfig = true;
          this.form.patchValue(smtpConfig.smtpConfig);
          this.form.patchValue({ ['hostAndPort']: smtpConfig.smtpConfig.host });
        }
      })
      .catch((error) => {
        this.smtpLoading = false;
        if (error && error.code === 5) {
          console.log(error);
          this.hasSMTPConfig = false;
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
      .catch(() => {
        this.logProviderLoading = false;
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
      .catch(() => {
        this.fileProviderLoading = false;
      });
  }

  private updateData(): Promise<UpdateSMTPConfigResponse.AsObject | AddSMTPConfigResponse> {
    if (this.hasSMTPConfig) {
      const req = new UpdateSMTPConfigRequest();
      req.setHost(this.hostAndPort?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');

      return this.service.updateSMTPConfig(req);
    } else {
      const req = new AddSMTPConfigRequest();
      req.setHost(this.hostAndPort?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');

      return this.service.addSMTPConfig(req);
    }
  }

  public savePolicy(): void {
    this.updateData()
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

  public editSMSProvider(): void {
    const dialogRef = this.dialog.open(DialogAddSMSProviderComponent, {
      width: '400px',
      data: {
        smsProviders: this.smsProviders,
      },
    });

    dialogRef.afterClosed().subscribe((req: AddSMSProviderTwilioRequest | UpdateSMSProviderTwilioRequest) => {
      if (req) {
        if (!!this.twilio) {
          this.service
            .updateSMSProviderTwilio(req as UpdateSMSProviderTwilioRequest)
            .then(() => {
              this.toast.showInfo('SETTING.SMS.TWILIO.ADDED', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          this.service
            .addSMSProviderTwilio(req as AddSMSProviderTwilioRequest)
            .then(() => {
              this.toast.showInfo('SETTING.SMS.TWILIO.ADDED', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public setSMTPPassword(): void {
    const dialogRef = this.dialog.open(PasswordDialogComponent, {
      width: '400px',
      data: {
        i18nTitle: 'SETTING.SMTP.SETPASSWORD',
        i18nLabel: 'SETTING.SMTP.PASSWORD',
      },
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

  public toggleSMSProviderState(id: string): void {
    const provider = this.smsProviders.find((p) => p.id === id);
    if (provider) {
      if (provider.state === SMSProviderConfigState.SMS_PROVIDER_CONFIG_ACTIVE) {
        this.service
          .deactivateSMSProvider(id)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.DEACTIVATED', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      } else if (provider.state === SMSProviderConfigState.SMS_PROVIDER_CONFIG_INACTIVE) {
        this.service
          .activateSMSProvider(id)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.ACTIVATED', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    }
  }

  public removeSMSProvider(id: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'SETTING.SMS.REMOVEPROVIDER',
        descriptionKey: 'SETTING.SMS.REMOVEPROVIDER_DESC',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.service
          .removeSMSProvider(id)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.TWILIO.REMOVED', true);
            this.fetchData();
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

  public get hostAndPort(): AbstractControl | null {
    return this.form.get('hostAndPort');
  }
}

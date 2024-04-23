import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { take } from 'rxjs';
import {
  AddSMTPConfigRequest,
  AddSMTPConfigResponse,
  UpdateSMTPConfigPasswordRequest,
  UpdateSMTPConfigRequest,
  UpdateSMTPConfigResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { InfoSectionType } from '../../info-section/info-section.component';
import { PasswordDialogComponent } from '../notification-sms-provider/password-dialog/password-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-notification-smtp-provider',
  templateUrl: './notification-smtp-provider.component.html',
  styleUrls: ['./notification-smtp-provider.component.scss'],
})
export class NotificationSMTPProviderComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;

  public smtpLoading: boolean = false;

  public form!: UntypedFormGroup;

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
      replyToAddress: [{ disabled: true, value: '' }],
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
  }

  private updateData(): Promise<UpdateSMTPConfigResponse.AsObject | AddSMTPConfigResponse> {
    if (this.hasSMTPConfig) {
      const req = new UpdateSMTPConfigRequest();
      req.setHost(this.hostAndPort?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setReplyToAddress(this.replyToAddress?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');

      return this.service.updateSMTPConfig(req);
    } else {
      const req = new AddSMTPConfigRequest();
      req.setHost(this.hostAndPort?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setReplyToAddress(this.replyToAddress?.value ?? '');
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

  public get senderAddress(): AbstractControl | null {
    return this.form.get('senderAddress');
  }

  public get senderName(): AbstractControl | null {
    return this.form.get('senderName');
  }

  public get replyToAddress(): AbstractControl | null {
    return this.form.get('replyToAddress');
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

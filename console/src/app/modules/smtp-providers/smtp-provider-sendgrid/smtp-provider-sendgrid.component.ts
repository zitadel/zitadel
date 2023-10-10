import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { Subject, take } from 'rxjs';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { MatLegacyCheckboxChange } from '@angular/material/legacy-checkbox';
import {
  AddSMTPConfigRequest,
  AddSMTPConfigResponse,
  UpdateSMTPConfigRequest,
  UpdateSMTPConfigResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { SMTPProviderType } from 'src/app/proto/generated/zitadel/settings_pb';
import { ToastService } from 'src/app/services/toast.service';
import { Router } from '@angular/router';

@Component({
  selector: 'cnsl-provider-sendgrid',
  templateUrl: './smtp-provider-sendgrid.component.html',
  styleUrls: ['./smtp-provider-sendgrid.component.scss'],
})
export class SMTPProviderSendgridComponent implements OnInit {
  public showOptional: boolean = false;
  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);
  public id: string | null = '';
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public smtpLoading: boolean = false;
  public hasSMTPConfig: boolean = false;

  public provider?: Provider.AsObject;
  public updateClientSecret: boolean = false;

  // stepper
  public currentCreateStep: number = 1;
  public requestRedirectValuesSubject$: Subject<void> = new Subject();
  public firstFormGroup!: UntypedFormGroup;
  public secondFormGroup!: UntypedFormGroup;

  private host: string = 'smtp.sendgrid.net';
  private unencryptedPort: number = 587;
  private encryptedPort: number = 465;

  constructor(
    private service: AdminService,
    private _location: Location,
    private fb: UntypedFormBuilder,
    private authService: GrpcAuthService,
    private toast: ToastService,
    private router: Router,
  ) {
    this.firstFormGroup = this.fb.group({
      tls: [false],
      hostAndPort: [`${this.host}:${this.unencryptedPort}`],
      user: ['apiKey'],
      password: [''],
    });

    this.secondFormGroup = this.fb.group({
      senderAddress: ['', [requiredValidator]],
      senderName: ['', [requiredValidator]],
      replyToAddress: [''],
    });
  }

  ngOnInit(): void {
    if (!this.router.url.endsWith('/create')) {
      this.fetchData();
      this.authService
        .isAllowed(['iam.write'])
        .pipe(take(1))
        .subscribe((allowed) => {
          if (allowed) {
            this.firstFormGroup.enable();
            this.secondFormGroup.enable();
          }
        });
    }
  }

  public changeStep(event: StepperSelectionEvent): void {
    this.currentCreateStep = event.selectedIndex + 1;

    if (event.selectedIndex >= 2) {
      this.requestRedirectValuesSubject$.next();
    }
  }

  public close(): void {
    this._location.back();
  }

  public toggleTLS(event: MatLegacyCheckboxChange) {
    this.hostAndPort?.setValue(`${this.host}:${event.checked ? this.encryptedPort : this.unencryptedPort}`);
  }

  private fetchData(): void {
    this.smtpLoading = true;
    this.service
      .getSMTPConfig()
      .then((smtpConfig) => {
        this.smtpLoading = false;
        if (smtpConfig.smtpConfig) {
          this.hasSMTPConfig = true;
          this.firstFormGroup.patchValue({
            ['tls']: smtpConfig.smtpConfig.tls,
            ['hostAndPort']: smtpConfig.smtpConfig.host,
            ['user']: smtpConfig.smtpConfig.user,
          });
          this.secondFormGroup.patchValue({
            ['senderAddress']: smtpConfig.smtpConfig.senderAddress,
            ['senderName']: smtpConfig.smtpConfig.senderName,
            ['replyToAddress']: smtpConfig.smtpConfig.replyToAddress,
          });
        }
      })
      .catch((error) => {
        this.smtpLoading = false;
        if (error && error.code === 5) {
          this.hasSMTPConfig = false;
        }
      });
  }

  private updateData(): Promise<UpdateSMTPConfigResponse.AsObject | AddSMTPConfigResponse> {
    if (this.hasSMTPConfig) {
      const req = new UpdateSMTPConfigRequest();
      req.setTls(this.tls?.value ?? false);

      if (this.hostAndPort && this.hostAndPort.value) {
        req.setHost(this.hostAndPort.value);
      }
      if (this.user && this.user.value) {
        req.setUser(this.user.value);
      }
      if (this.password && this.password.value) {
        req.setPassword(this.password.value);
      }
      if (this.senderAddress && this.senderAddress.value) {
        req.setSenderAddress(this.senderAddress.value);
      }
      if (this.senderName && this.senderName.value) {
        req.setSenderName(this.senderName.value);
      }
      if (this.replyToAddress && this.replyToAddress.value) {
        req.setReplyToAddress(this.replyToAddress.value);
      }
      req.setProviderType(SMTPProviderType.SMTP_PROVIDER_TYPE_SENDGRID);
      return this.service.updateSMTPConfig(req);
    } else {
      const req = new AddSMTPConfigRequest();
      req.setHost(this.hostAndPort?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setReplyToAddress(this.replyToAddress?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');
      req.setPassword(this.password?.value ?? '');
      req.setIsActive(true);
      req.setProviderType(SMTPProviderType.SMTP_PROVIDER_TYPE_SENDGRID);
      return this.service.addSMTPConfig(req);
    }
  }

  public savePolicy(): void {
    this.updateData()
      .then(() => {
        this.toast.showInfo('SETTING.SMTP.SAVED', true);
        setTimeout(() => {
          this.close();
        }, 2000);
      })
      .catch((error: unknown) => {
        this.toast.showError(error);
      });
  }

  public get tls(): AbstractControl | null {
    return this.firstFormGroup.get('tls');
  }

  public get user(): AbstractControl | null {
    return this.firstFormGroup.get('user');
  }

  public get hostAndPort(): AbstractControl | null {
    return this.firstFormGroup.get('hostAndPort');
  }

  public get password(): AbstractControl | null {
    return this.firstFormGroup.get('password');
  }

  public get senderAddress(): AbstractControl | null {
    return this.secondFormGroup.get('senderAddress');
  }

  public get senderName(): AbstractControl | null {
    return this.secondFormGroup.get('senderName');
  }

  public get replyToAddress(): AbstractControl | null {
    return this.secondFormGroup.get('replyToAddress');
  }
}

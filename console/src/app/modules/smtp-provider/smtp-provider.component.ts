import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { Subject } from 'rxjs';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Options } from 'src/app/proto/generated/zitadel/idp_pb';
import { requiredValidator } from '../form-field/validators/validators';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import {
  AddSMTPConfigRequest,
  AddSMTPConfigResponse,
  UpdateSMTPConfigRequest,
  UpdateSMTPConfigResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { ActivatedRoute, Router } from '@angular/router';
import { MatCheckboxChange } from '@angular/material/checkbox';
import {
  AmazonSESDefaultSettings,
  BrevoDefaultSettings,
  GenericDefaultSettings,
  GoogleDefaultSettings,
  MailchimpDefaultSettings,
  MailgunDefaultSettings,
  MailjetDefaultSettings,
  PostmarkDefaultSettings,
  ProviderDefaultSettings,
  OutlookDefaultSettings,
  SendgridDefaultSettings,
} from './known-smtp-providers-settings';

@Component({
  selector: 'cnsl-smtp-provider',
  templateUrl: './smtp-provider.component.html',
  styleUrls: ['./smtp-provider.scss'],
})
export class SMTPProviderComponent {
  public showOptional: boolean = false;
  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);
  public id: string = '';
  public providerDefaultSetting: ProviderDefaultSettings = GenericDefaultSettings;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public smtpLoading: boolean = false;
  public hasSMTPConfig: boolean = false;

  public updateClientSecret: boolean = false;

  // stepper
  public currentCreateStep: number = 1;
  public requestRedirectValuesSubject$: Subject<void> = new Subject();
  public firstFormGroup!: UntypedFormGroup;
  public secondFormGroup!: UntypedFormGroup;

  public senderEmailPlaceholder = 'sender@example.com';

  constructor(
    private service: AdminService,
    private _location: Location,
    private fb: UntypedFormBuilder,
    private toast: ToastService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.route.parent?.url.subscribe((urlPath) => {
      const providerName = urlPath[urlPath.length - 1].path;
      switch (providerName) {
        case 'aws-ses':
          this.providerDefaultSetting = AmazonSESDefaultSettings;
          break;
        case 'google':
          this.providerDefaultSetting = GoogleDefaultSettings;
          break;
        case 'mailgun':
          this.providerDefaultSetting = MailgunDefaultSettings;
          break;
        case 'mailjet':
          this.providerDefaultSetting = MailjetDefaultSettings;
          break;
        case 'postmark':
          this.providerDefaultSetting = PostmarkDefaultSettings;
          break;
        case 'sendgrid':
          this.providerDefaultSetting = SendgridDefaultSettings;
          break;
        case 'mailchimp':
          this.providerDefaultSetting = MailchimpDefaultSettings;
          break;
        case 'brevo':
          this.providerDefaultSetting = BrevoDefaultSettings;
          break;
        case 'outlook':
          this.providerDefaultSetting = OutlookDefaultSettings;
          break;
      }

      this.firstFormGroup = this.fb.group({
        description: [this.providerDefaultSetting.name],
        tls: [{ value: this.providerDefaultSetting.requiredTls, disabled: this.providerDefaultSetting.requiredTls }],
        region: [''],
        hostAndPort: [
          this.providerDefaultSetting?.host
            ? `${this.providerDefaultSetting?.host}:${this.providerDefaultSetting?.unencryptedPort}`
            : '',
        ],
        user: [this.providerDefaultSetting?.user.value || ''],
        password: [this.providerDefaultSetting?.password.value || ''],
      });

      this.senderEmailPlaceholder = this.providerDefaultSetting?.senderEmailPlaceholder || 'sender@example.com';

      this.secondFormGroup = this.fb.group({
        senderAddress: ['', [requiredValidator]],
        senderName: ['', [requiredValidator]],
        replyToAddress: [''],
      });

      this.region?.valueChanges.subscribe((region: string) => {
        this.hostAndPort?.setValue(
          `${region}:${
            this.tls ? this.providerDefaultSetting?.encryptedPort : this.providerDefaultSetting?.unencryptedPort
          }`,
        );
      });

      if (!this.router.url.endsWith('/create')) {
        this.id = this.route.snapshot.paramMap.get('id') || '';
        if (this.id) {
          this.fetchData(this.id);
        }
      }
    });
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

  public toggleTLS(event: MatCheckboxChange) {
    if (this.providerDefaultSetting.host) {
      this.hostAndPort?.setValue(
        `${this.providerDefaultSetting?.host}:${
          event.checked ? this.providerDefaultSetting?.encryptedPort : this.providerDefaultSetting?.unencryptedPort
        }`,
      );
    }
  }

  private fetchData(id: string): void {
    this.smtpLoading = true;
    this.service
      .getSMTPConfigById(id)
      .then((data) => {
        this.smtpLoading = false;
        if (data.smtpConfig) {
          this.hasSMTPConfig = true;
          this.firstFormGroup.patchValue({
            ['description']: data.smtpConfig.description,
            ['tls']: data.smtpConfig.tls,
            ['hostAndPort']: data.smtpConfig.host,
            ['user']: data.smtpConfig.user,
          });
          this.secondFormGroup.patchValue({
            ['senderAddress']: data.smtpConfig.senderAddress,
            ['senderName']: data.smtpConfig.senderName,
            ['replyToAddress']: data.smtpConfig.replyToAddress,
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
      req.setId(this.id);
      req.setDescription(this.description?.value || '');
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
      return this.service.updateSMTPConfig(req);
    } else {
      const req = new AddSMTPConfigRequest();
      req.setDescription(this.description?.value ?? '');
      req.setHost(this.hostAndPort?.value ?? '');
      req.setSenderAddress(this.senderAddress?.value ?? '');
      req.setSenderName(this.senderName?.value ?? '');
      req.setReplyToAddress(this.replyToAddress?.value ?? '');
      req.setTls(this.tls?.value ?? false);
      req.setUser(this.user?.value ?? '');
      req.setPassword(this.password?.value ?? '');
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
        if (`${error}`.includes('No changes')) {
          this.toast.showInfo('SETTING.SMTP.NOCHANGES', true);
          setTimeout(() => {
            this.close();
          }, 2000);
        } else {
          this.toast.showError(error);
        }
      });
  }

  public get description(): AbstractControl | null {
    return this.firstFormGroup.get('description');
  }

  public get tls(): AbstractControl | null {
    return this.firstFormGroup.get('tls');
  }

  public get region(): AbstractControl | null {
    return this.firstFormGroup.get('region');
  }

  public get hostAndPort(): AbstractControl | null {
    return this.firstFormGroup.get('hostAndPort');
  }

  public get user(): AbstractControl | null {
    return this.firstFormGroup.get('user');
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

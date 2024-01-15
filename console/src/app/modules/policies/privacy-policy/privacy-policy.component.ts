import { Component, Injector, Input, OnDestroy, OnInit, Type } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { Observable, Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import {
  GetPrivacyPolicyResponse as AdminGetPrivacyPolicyResponse,
  UpdatePrivacyPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddCustomPrivacyPolicyRequest,
  GetPrivacyPolicyResponse,
  UpdateCustomPrivacyPolicyRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { PrivacyPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-privacy-policy',
  templateUrl: './privacy-policy.component.html',
  styleUrls: ['./privacy-policy.component.scss'],
})
export class PrivacyPolicyComponent implements OnInit, OnDestroy {
  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  private sub: Subscription = new Subscription();

  public privacyPolicy: PrivacyPolicy.AsObject | undefined = undefined;
  public form!: UntypedFormGroup;
  public InfoSectionType: any = InfoSectionType;

  public canWrite$: Observable<boolean> = this.authService.isAllowed([
    this.serviceType === PolicyComponentServiceType.ADMIN
      ? 'iam.policy.write'
      : this.serviceType === PolicyComponentServiceType.MGMT
        ? 'policy.write'
        : '',
  ]);

  public LANGPLACEHOLDER: string = '{{.Lang}}';
  public copied: string = '';

  constructor(
    private authService: GrpcAuthService,
    private injector: Injector,
    private dialog: MatDialog,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
  ) {
    this.form = this.fb.group({
      tosLink: ['', []],
      privacyLink: ['', []],
      helpLink: ['', []],
      supportEmail: ['', []],
    });

    this.canWrite$.pipe(take(1)).subscribe((canWrite) => {
      if (canWrite) {
        this.form.enable();
      } else {
        this.form.disable();
      }
    });
  }

  ngOnInit(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        this.loadData();
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);
        this.loadData();
        break;
    }
  }

  public addChip(formControlName: string, value: string): void {
    const c = this.form.get(formControlName)?.value;
    this.form.get(formControlName)?.setValue(`${c}${value}`);
  }

  public async loadData(): Promise<any> {
    const getData = (): Promise<AdminGetPrivacyPolicyResponse.AsObject | GetPrivacyPolicyResponse.AsObject> => {
      return this.service.getPrivacyPolicy();
    };

    getData()
      .then((resp) => {
        if (resp.policy) {
          this.privacyPolicy = resp.policy;
          this.form.patchValue(this.privacyPolicy);
        } else {
          this.privacyPolicy = undefined;
          this.form.patchValue({
            tosLink: '',
            privacyLink: '',
            helpLink: '',
            supportEmail: '',
          });
        }
      })
      .catch((error) => {
        this.privacyPolicy = undefined;
        this.form.patchValue({
          tosLink: '',
          privacyLink: '',
          helpLink: '',
          supportEmail: '',
        });
      });
  }

  public saveCurrentMessage(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      if (!this.privacyPolicy || (this.privacyPolicy as PrivacyPolicy.AsObject).isDefault) {
        const req = new AddCustomPrivacyPolicyRequest();
        req.setPrivacyLink(this.form.get('privacyLink')?.value);
        req.setTosLink(this.form.get('tosLink')?.value);
        req.setHelpLink(this.form.get('helpLink')?.value);
        req.setSupportEmail(this.form.get('supportEmail')?.value);
        (this.service as ManagementService)
          .addCustomPrivacyPolicy(req)
          .then(() => {
            this.toast.showInfo('POLICY.PRIVACY_POLICY.SAVED', true);
            this.loadData();
          })
          .catch((error) => this.toast.showError(error));
      } else {
        const req = new UpdateCustomPrivacyPolicyRequest();
        req.setPrivacyLink(this.form.get('privacyLink')?.value);
        req.setTosLink(this.form.get('tosLink')?.value);
        req.setHelpLink(this.form.get('helpLink')?.value);
        req.setSupportEmail(this.form.get('supportEmail')?.value);

        (this.service as ManagementService)
          .updateCustomPrivacyPolicy(req)
          .then(() => {
            this.toast.showInfo('POLICY.PRIVACY_POLICY.SAVED', true);
            this.loadData();
          })
          .catch((error) => this.toast.showError(error));
      }
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      const req = new UpdatePrivacyPolicyRequest();
      req.setPrivacyLink(this.form.get('privacyLink')?.value);
      req.setTosLink(this.form.get('tosLink')?.value);
      req.setHelpLink(this.form.get('helpLink')?.value);
      req.setSupportEmail(this.form.get('supportEmail')?.value);

      (this.service as AdminService)
        .updatePrivacyPolicy(req)
        .then(() => {
          this.toast.showInfo('POLICY.PRIVACY_POLICY.SAVED', true);
          this.loadData();
        })
        .catch((error) => this.toast.showError(error));
    }
  }

  public resetDefault(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        icon: 'las la-history',
        confirmKey: 'ACTIONS.RESTORE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'POLICY.PRIVACY_POLICY.RESET_TITLE',
        descriptionKey: 'POLICY.PRIVACY_POLICY.RESET_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          (this.service as ManagementService)
            .resetPrivacyPolicyToDefault()
            .then(() => {
              setTimeout(() => {
                this.loadData();
              }, 1000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public get isDefault(): boolean {
    if (this.privacyPolicy && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.privacyPolicy as PrivacyPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

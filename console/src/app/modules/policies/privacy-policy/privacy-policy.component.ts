import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
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
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { CnslLinks } from '../../links/links.component';
import { GridPolicy, PRIVACY_POLICY } from '../../policy-grid/policies';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'app-privacy-policy',
  templateUrl: './privacy-policy.component.html',
  styleUrls: ['./privacy-policy.component.scss'],
})
export class PrivacyPolicyComponent implements OnDestroy {
  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public nextLinks: CnslLinks[] = [];
  private sub: Subscription = new Subscription();

  public privacyPolicy!: PrivacyPolicy.AsObject;
  public form!: FormGroup;
  public currentPolicy: GridPolicy = PRIVACY_POLICY;

  constructor(
    private route: ActivatedRoute,
    private injector: Injector,
    private dialog: MatDialog,
    private toast: ToastService,
    private fb: FormBuilder,
  ) {

    this.form = this.fb.group({
      tosLink: ['', []],
      privacyLink: ['', []],
    });

    this.route.data.pipe(switchMap(data => {
      this.serviceType = data.serviceType;
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

      return this.route.params;
    })).subscribe();
  }

  public async loadData(): Promise<any> {
    const getData = ():
      Promise<AdminGetPrivacyPolicyResponse.AsObject | GetPrivacyPolicyResponse.AsObject> => {
      return (this.service as AdminService).getPrivacyPolicy();
    };

    getData().then(resp => {
      if (resp.policy) {
        this.privacyPolicy = resp.policy;
        this.form.patchValue(this.privacyPolicy);
      }
    });
  }

  public saveCurrentMessage(): void {
    console.log(this.form.get('privacyLink')?.value, this.form.get('tosLink')?.value);
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      if ((this.privacyPolicy as PrivacyPolicy.AsObject).isDefault) {
        const req = new AddCustomPrivacyPolicyRequest();
        req.setPrivacyLink(this.form.get('privacyLink')?.value);
        req.setTosLink(this.form.get('tosLink')?.value);
        (this.service as ManagementService).addCustomPrivacyPolicy(req).then(() => {
          this.toast.showInfo('POLICY.PRIVACY_POLICY.SAVED', true);
        }).catch(error => this.toast.showError(error));
      } else {
        const req = new UpdateCustomPrivacyPolicyRequest();
        req.setPrivacyLink(this.form.get('privacyLink')?.value);
        req.setTosLink(this.form.get('tosLink')?.value);
        (this.service as ManagementService).updateCustomPrivacyPolicy(req).then(() => {
          this.toast.showInfo('POLICY.PRIVACY_POLICY.SAVED', true);
        }).catch(error => this.toast.showError(error));
      }

    } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      const req = new UpdatePrivacyPolicyRequest();
      req.setPrivacyLink(this.form.get('privacyLink')?.value);
      req.setTosLink(this.form.get('tosLink')?.value);

      (this.service as AdminService).updatePrivacyPolicy(req).then(() => {
        this.toast.showInfo('POLICY.PRIVACY_POLICY.SAVED', true);
      }).catch(error => this.toast.showError(error));
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

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          (this.service as ManagementService).resetPrivacyPolicyToDefault().then(() => {
            setTimeout(() => {
              this.loadData();
            }, 1000);
          }).catch(error => {
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

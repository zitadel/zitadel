import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
  GetOrgFeaturesResponse,
  SetDefaultFeaturesRequest,
  SetOrgFeaturesRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Features } from 'src/app/proto/generated/zitadel/features_pb';
import { GetFeaturesResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { StripeCustomer, SubscriptionService } from 'src/app/services/subscription.service';
import { ToastService } from 'src/app/services/toast.service';

import { COUNTRIES, Country } from './country';
import { PaymentInfoDialogComponent } from './payment-info-dialog/payment-info-dialog.component';

export enum FeatureServiceType {
  MGMT,
  ADMIN,
}

@Component({
  selector: 'app-features',
  templateUrl: './features.component.html',
  styleUrls: ['./features.component.scss'],
})
export class FeaturesComponent implements OnDestroy {
  private managementService!: ManagementService;
  public serviceType!: FeatureServiceType;

  public features!: Features.AsObject;

  private sub: Subscription = new Subscription();
  public org!: Org.AsObject;

  public FeatureServiceType: any = FeatureServiceType;

  public customerLoading: boolean = false;
  public stripeLoading: boolean = false;
  public stripeURL: string = '';
  public stripeCustomer!: StripeCustomer;

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private sessionStorage: StorageService,
    private injector: Injector,
    private adminService: AdminService,
    private subService: SubscriptionService,
    private dialog: MatDialog,
  ) {
    const temporg = this.sessionStorage.getItem('organization') as Org.AsObject;
    if (temporg) {
      this.org = temporg;
    }
    this.sub = this.route.data.pipe(switchMap(data => {
      this.serviceType = data.serviceType;
      if (this.serviceType === FeatureServiceType.MGMT) {
        this.managementService = this.injector.get(ManagementService as Type<ManagementService>);
      }
      return this.route.params;
    })).subscribe(_ => {
      this.fetchData();
    });

    if (this.serviceType === FeatureServiceType.MGMT) {
      this.customerLoading = true;
      this.subService.getCustomer(this.org.id)
        .then(payload => {
          this.customerLoading = false;
          this.stripeCustomer = payload;
          if (this.customerValid) {
            this.getLinkToStripe();
          }
        })
        .catch(error => {
          this.customerLoading = false;
          console.error(error);
        });
    }
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public setCustomer(): void {
    const dialogRefPhone = this.dialog.open(PaymentInfoDialogComponent, {
      data: {
        customer: this.stripeCustomer,
      },
      width: '400px',
    });

    dialogRefPhone.afterClosed().subscribe(customer => {
      if (customer) {
        console.log(customer);
        this.stripeCustomer = customer;
        this.subService.setCustomer(this.org.id, customer).then(() => {
          this.getLinkToStripe();
        }).catch(console.error);
      }
    });
  }

  public getLinkToStripe(): void {
    if (this.serviceType === FeatureServiceType.MGMT) {
      this.stripeLoading = true;
      this.subService.getLink(this.org.id, window.location.href)
        .then(payload => {
          this.stripeLoading = false;
          console.log(payload);
          this.stripeURL = payload.redirect_url;
        })
        .catch(error => {
          this.stripeLoading = false;
          console.error(error);
        });
    }
  }

  public fetchData(): void {
    this.getData().then(resp => {
      if (resp?.features) {
        this.features = resp.features;
      }
    });
  }

  private async getData(): Promise<GetFeaturesResponse.AsObject | GetOrgFeaturesResponse.AsObject | undefined> {
    switch (this.serviceType) {
      case FeatureServiceType.MGMT:
        return this.managementService.getFeatures();
      case FeatureServiceType.ADMIN:
        if (this.org?.id) {
          return this.adminService.getDefaultFeatures();
        }
        break;
    }
  }

  public savePolicy(): void {
    switch (this.serviceType) {
      case FeatureServiceType.MGMT:
        const req = new SetOrgFeaturesRequest();
        req.setOrgId(this.org.id);

        req.setLoginPolicyUsernameLogin(this.features.loginPolicyUsernameLogin);
        req.setLoginPolicyPasswordReset(this.features.loginPolicyPasswordReset);
        req.setLoginPolicyRegistration(this.features.loginPolicyRegistration);
        req.setLoginPolicyIdp(this.features.loginPolicyIdp);
        req.setLoginPolicyFactors(this.features.loginPolicyFactors);
        req.setLoginPolicyPasswordless(this.features.loginPolicyPasswordless);
        req.setPasswordComplexityPolicy(this.features.passwordComplexityPolicy);
        req.setLabelPolicy(this.features.labelPolicy);
        req.setCustomDomain(this.features.customDomain);
        req.setCustomText(this.features.customText);

        this.adminService.setOrgFeatures(req).then(() => {
          this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
          this.toast.showError(error);
        });
        break;
      case FeatureServiceType.ADMIN:
        // update Default org iam policy?
        const dreq = new SetDefaultFeaturesRequest();
        dreq.setLoginPolicyUsernameLogin(this.features.loginPolicyUsernameLogin);
        dreq.setLoginPolicyPasswordReset(this.features.loginPolicyPasswordReset);
        dreq.setLoginPolicyRegistration(this.features.loginPolicyRegistration);
        dreq.setLoginPolicyIdp(this.features.loginPolicyIdp);
        dreq.setLoginPolicyFactors(this.features.loginPolicyFactors);
        dreq.setLoginPolicyPasswordless(this.features.loginPolicyPasswordless);
        dreq.setPasswordComplexityPolicy(this.features.passwordComplexityPolicy);
        dreq.setLabelPolicy(this.features.labelPolicy);
        dreq.setCustomDomain(this.features.customDomain);
        dreq.setCustomText(this.features.customText);

        this.adminService.setDefaultFeatures(dreq).then(() => {
          this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
          this.toast.showError(error);
        });
        break;
    }
  }

  public resetFeatures(): void {
    if (this.serviceType === FeatureServiceType.MGMT) {
      this.adminService.resetOrgFeatures(this.org.id).then(() => {
        this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
        setTimeout(() => {
          this.fetchData();
        }, 1000);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
  }

  public get isDefault(): boolean {
    if (this.features && this.serviceType === FeatureServiceType.MGMT) {
      return this.features.isDefault;
    } else {
      return false;
    }
  }

  get customerValid(): boolean {
    return !!this.stripeCustomer?.contact &&
      !!this.stripeCustomer?.address &&
      !!this.stripeCustomer?.city &&
      !!this.stripeCustomer?.postal_code;
  }

  get customerCountry(): Country | undefined {
    return COUNTRIES.find(country => country.isoCode === this.stripeCustomer.country);
  }
}

import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { UntypedFormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs';
import { Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

@Component({
  selector: 'cnsl-providers',
  templateUrl: './providers.component.html',
})
export class ProvidersComponent {
  public showOptional: boolean = false;
  public options: Options = new Options();

  public id: string | null = '';
  public providertype: string | null = '';

  public updateClientSecret: boolean = false;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public form!: UntypedFormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
  ) {
    // this.form = new UntypedFormGroup({
    //   name: new UntypedFormControl('', [requiredValidator]),
    //   clientId: new UntypedFormControl('', [requiredValidator]),
    //   clientSecret: new UntypedFormControl('', [requiredValidator]),
    //   issuer: new UntypedFormControl('', [requiredValidator]),
    //   scopesList: new UntypedFormControl(['openid', 'profile', 'email'], []),
    // });

    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data.serviceType;

      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);

          const bread: Breadcrumb = {
            type: BreadcrumbType.ORG,
            routerLink: ['/org'],
          };

          breadcrumbService.setBreadcrumb([bread]);
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);

          const iamBread = new Breadcrumb({
            type: BreadcrumbType.ORG,
            name: 'Instance',
            routerLink: ['/instance'],
          });
          breadcrumbService.setBreadcrumb([iamBread]);
          break;
      }

      this.id = this.route.snapshot.paramMap.get('id');
      this.providertype = this.route.snapshot.paramMap.get('providertype');
    });
  }

  // public getData(id: string): void;
  //    {
  //     this.loading = true;
  //     const req =
  //       this.serviceType === PolicyComponentServiceType.ADMIN
  //         ? new AdminGetProviderByIDRequest()
  //         : new MgmtGetProviderByIDRequest();
  //     req.setId(id);
  //     this.service
  //       .getProviderByID(req)
  //       .then((resp) => {
  //         this.provider = resp.idp;
  //         this.loading = false;
  //         if (this.provider?.config?.oidc) {
  //           this.oidcFormGroup.patchValue(this.provider.config.oidc);
  //           this.name?.setValue(this.provider.name);
  //         }
  //       })
  //       .catch((error) => {
  //         this.toast.showError(error);
  //         this.loading = false;
  //       });
  //   }

  public close(): void {
    this._location.back();
  }
}

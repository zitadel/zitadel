import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnInit, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { take } from 'rxjs/operators';
import {
  AddGoogleProviderRequest as AdminAddGoogleProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGoogleProviderRequest as MgmtAddGoogleProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Component({
  selector: 'cnsl-provider-google',
  templateUrl: './provider-google.component.html',
  styleUrls: ['./provider-google.component.scss'],
})
export class ProviderGoogleComponent implements OnInit {
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public projectId: string = '';

  public form!: FormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;
  public updateClientSecret: boolean = false;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
  ) {
    const idpId = this.route.snapshot.paramMap.get('id');

    console.log(idpId);

    this.form = new FormGroup({
      name: new FormControl('', [Validators.required]),
      clientId: new FormControl('', [Validators.required]),
      clientSecret: new FormControl('', [Validators.required]),
      scopesList: new FormControl(['openid', 'profile', 'email'], []),
    });

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
    });
  }

  public ngOnInit(): void {
    this.route.params.pipe(take(1)).subscribe((params) => this.getData(params));
  }

  private getData({ id }: Params): void {
    const req =
      this.serviceType === PolicyComponentServiceType.ADMIN
        ? new AdminGetProviderByIDRequest()
        : new MgmtGetProviderByIDRequest();
    req.setId(id);
    this.service.getProviderByID(req).then((resp) => {
      this.provider = resp.idp;
      console.log(this.provider);
      if (this.provider?.config?.google) {
        this.form.patchValue(this.provider.config.google);
        this.name?.setValue(this.provider.name);
      }
    });
  }

  public addOIDCIdp(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddGoogleProviderRequest();

      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setScopesList(this.scopesList?.value);
      //   req.setProviderOptions()

      this.loading = true;
      (this.service as ManagementService)
        .addGoogleProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.router.navigate(
              [
                this.serviceType === PolicyComponentServiceType.MGMT
                  ? '/org-settings'
                  : this.serviceType === PolicyComponentServiceType.ADMIN
                  ? '/settings'
                  : '',
              ],
              { queryParams: { id: 'idp' } },
            );
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    } else if (PolicyComponentServiceType.ADMIN) {
      const req = new AdminAddGoogleProviderRequest();
      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setScopesList(this.scopesList?.value);

      this.loading = true;
      (this.service as AdminService)
        .addGoogleProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.router.navigate(
              [
                this.serviceType === PolicyComponentServiceType.MGMT
                  ? '/org-settings'
                  : this.serviceType === PolicyComponentServiceType.ADMIN
                  ? '/settings'
                  : '',
              ],
              { queryParams: { id: 'idp' } },
            );
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public close(): void {
    this._location.back();
  }

  public addScope(event: MatChipInputEvent): void {
    const input = event.chipInput?.inputElement;
    const value = event.value.trim();

    if (value !== '') {
      if (this.scopesList?.value) {
        this.scopesList.value.push(value);
        if (input) {
          input.value = '';
        }
      }
    }
  }

  public removeScope(uri: string): void {
    if (this.scopesList?.value) {
      const index = this.scopesList.value.indexOf(uri);

      if (index !== undefined && index >= 0) {
        this.scopesList.value.splice(index, 1);
      }
    }
  }

  public get name(): AbstractControl | null {
    return this.form.get('name');
  }

  public get clientId(): AbstractControl | null {
    return this.form.get('clientId');
  }

  public get clientSecret(): AbstractControl | null {
    return this.form.get('clientSecret');
  }

  public get scopesList(): AbstractControl | null {
    return this.form.get('scopesList');
  }
}

import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import {
  AddGoogleProviderRequest as AdminAddGoogleProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateGoogleProviderRequest as AdminUpdateGoogleProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGoogleProviderRequest as MgmtAddGoogleProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateGoogleProviderRequest as MgmtUpdateGoogleProviderRequest,
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
export class ProviderGoogleComponent {
  public showOptional: boolean = false;
  public options: Options = new Options();
  public id: string | null = '';
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

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
    private breadcrumbService: BreadcrumbService,
  ) {
    this.form = new FormGroup({
      name: new FormControl('', []),
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

          this.breadcrumbService.setBreadcrumb([bread]);
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);

          const iamBread = new Breadcrumb({
            type: BreadcrumbType.ORG,
            name: 'Instance',
            routerLink: ['/instance'],
          });
          this.breadcrumbService.setBreadcrumb([iamBread]);
          break;
      }

      this.id = this.route.snapshot.paramMap.get('id');
      if (this.id) {
        this.clientSecret?.setValidators([]);
        this.getData(this.id);
      }
    });
  }

  private getData(id: string): void {
    const req =
      this.serviceType === PolicyComponentServiceType.ADMIN
        ? new AdminGetProviderByIDRequest()
        : new MgmtGetProviderByIDRequest();
    req.setId(id);
    this.service
      .getProviderByID(req)
      .then((resp) => {
        this.provider = resp.idp;
        this.loading = false;
        if (this.provider?.config?.google) {
          this.form.patchValue(this.provider.config.google);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider ? this.updateGoogleProvider() : this.addGoogleProvider();
  }

  public addGoogleProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddGoogleProviderRequest();

      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setScopesList(this.scopesList?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as ManagementService)
        .addGoogleProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.router.navigate(['/org-settings'], { queryParams: { id: 'idp' } });
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
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as AdminService)
        .addGoogleProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.router.navigate(['/settings'], { queryParams: { id: 'idp' } });
          }, 2000);
        })
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    }
  }

  public updateGoogleProvider(): void {
    if (this.provider) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateGoogleProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setScopesList(this.scopesList?.value);
        req.setProviderOptions(this.options);

        if (this.updateClientSecret) {
          req.setClientSecret(this.clientSecret?.value);
        }

        this.loading = true;
        (this.service as ManagementService)
          .updateGoogleProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.router.navigate(['/org-settings'], { queryParams: { id: 'idp' } });
            }, 2000);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loading = false;
          });
      } else if (PolicyComponentServiceType.ADMIN) {
        const req = new AdminUpdateGoogleProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setScopesList(this.scopesList?.value);
        req.setProviderOptions(this.options);

        if (this.updateClientSecret) {
          req.setClientSecret(this.clientSecret?.value);
        }

        this.loading = true;
        (this.service as AdminService)
          .updateGoogleProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.router.navigate(['/settings'], { queryParams: { id: 'idp' } });
            }, 2000);
          })
          .catch((error) => {
            this.loading = false;
            this.toast.showError(error);
          });
      }
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

import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import {
  AddGenericOAuthProviderRequest as AdminAddGenericOAuthProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateGenericOAuthProviderRequest as AdminUpdateGenericOAuthProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGenericOAuthProviderRequest as MgmtAddGenericOAuthProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateGenericOAuthProviderRequest as MgmtUpdateGenericOAuthProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Component({
  selector: 'cnsl-provider-oauth',
  templateUrl: './provider-oauth.component.html',
  styleUrls: ['./provider-oauth.component.scss'],
})
export class ProviderOAuthComponent {
  public showOptional: boolean = false;
  public options: Options = new Options();

  public id: string | null = '';
  public updateClientSecret: boolean = false;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public form!: UntypedFormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
  ) {
    this.form = new UntypedFormGroup({
      name: new UntypedFormControl('', [Validators.required]),
      clientId: new UntypedFormControl('', [Validators.required]),
      clientSecret: new UntypedFormControl('', [Validators.required]),
      authorizationEndpoint: new UntypedFormControl('', [Validators.required]),
      tokenEndpoint: new UntypedFormControl('', [Validators.required]),
      userEndpoint: new UntypedFormControl('', [Validators.required]),
      idAttribute: new UntypedFormControl('', [Validators.required]),
      scopesList: new UntypedFormControl(['openid', 'profile', 'email'], []),
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

      this.id = this.route.snapshot.paramMap.get('id');
      if (this.id) {
        this.clientSecret?.setValidators([]);
        this.getData(this.id);
      }
    });
  }

  private getData(id: string): void {
    this.loading = true;
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
        if (this.provider?.config?.oauth) {
          this.form.patchValue(this.provider.config.oauth);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider ? this.updateGenericOAuthProvider() : this.addGenericOAuthProvider();
  }

  public addGenericOAuthProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddGenericOAuthProviderRequest();

      req.setName(this.name?.value);
      req.setAuthorizationEndpoint(this.authorizationEndpoint?.value);
      req.setIdAttribute(this.idAttribute?.value);
      req.setTokenEndpoint(this.tokenEndpoint?.value);
      req.setUserEndpoint(this.userEndpoint?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setScopesList(this.scopesList?.value);

      this.loading = true;
      (this.service as ManagementService)
        .addGenericOAuthProvider(req)
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
      const req = new AdminAddGenericOAuthProviderRequest();
      req.setName(this.name?.value);
      req.setAuthorizationEndpoint(this.authorizationEndpoint?.value);
      req.setIdAttribute(this.idAttribute?.value);
      req.setTokenEndpoint(this.tokenEndpoint?.value);
      req.setUserEndpoint(this.userEndpoint?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setScopesList(this.scopesList?.value);

      this.loading = true;
      (this.service as AdminService)
        .addGenericOAuthProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.router.navigate(['/settings'], { queryParams: { id: 'idp' } });
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    }
  }

  public updateGenericOAuthProvider(): void {
    if (this.provider) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateGenericOAuthProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setAuthorizationEndpoint(this.authorizationEndpoint?.value);
        req.setIdAttribute(this.idAttribute?.value);
        req.setTokenEndpoint(this.tokenEndpoint?.value);
        req.setUserEndpoint(this.userEndpoint?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setScopesList(this.scopesList?.value);

        this.loading = true;
        (this.service as ManagementService)
          .updateGenericOAuthProvider(req)
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
        const req = new AdminUpdateGenericOAuthProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setAuthorizationEndpoint(this.authorizationEndpoint?.value);
        req.setIdAttribute(this.idAttribute?.value);
        req.setTokenEndpoint(this.tokenEndpoint?.value);
        req.setUserEndpoint(this.userEndpoint?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setScopesList(this.scopesList?.value);

        this.loading = true;
        (this.service as AdminService)
          .updateGenericOAuthProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.router.navigate(['/settings'], { queryParams: { id: 'idp' } });
            }, 2000);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loading = false;
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

  public get authorizationEndpoint(): AbstractControl | null {
    return this.form.get('authorizationEndpoint');
  }

  public get tokenEndpoint(): AbstractControl | null {
    return this.form.get('tokenEndpoint');
  }

  public get userEndpoint(): AbstractControl | null {
    return this.form.get('userEndpoint');
  }

  public get idAttribute(): AbstractControl | null {
    return this.form.get('idAttribute');
  }

  public get clientId(): AbstractControl | null {
    return this.form.get('clientId');
  }

  public get clientSecret(): AbstractControl | null {
    return this.form.get('clientSecret');
  }

  public get issuer(): AbstractControl | null {
    return this.form.get('issuer');
  }

  public get scopesList(): AbstractControl | null {
    return this.form.get('scopesList');
  }
}

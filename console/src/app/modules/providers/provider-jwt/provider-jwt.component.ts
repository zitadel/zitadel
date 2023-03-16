import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs/operators';
import {
  AddJWTProviderRequest as AdminAddJWTProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateJWTProviderRequest as AdminUpdateJWTProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddJWTProviderRequest as MgmtAddJWTProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateJWTProviderRequest as MgmtUpdateJWTProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Component({
  selector: 'cnsl-provider-jwt',
  templateUrl: './provider-jwt.component.html',
})
export class ProviderJWTComponent {
  public showOptional: boolean = false;
  public options: Options = new Options();

  public id: string | null = '';
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
    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data.serviceType;
      console.log(data.serviceType);

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
        this.getData(this.id);
      }
    });

    this.form = new UntypedFormGroup({
      name: new UntypedFormControl('', [requiredValidator]),
      headerName: new UntypedFormControl('', [requiredValidator]),
      issuer: new UntypedFormControl('', [requiredValidator]),
      jwtEndpoint: new UntypedFormControl('', [requiredValidator]),
      keysEndpoint: new UntypedFormControl('', [requiredValidator]),
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
        if (this.provider?.config?.jwt) {
          this.form.patchValue(this.provider.config.jwt);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider ? this.updateJWTProvider() : this.addJWTProvider();
  }

  public addJWTProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddJWTProviderRequest();

      req.setName(this.name?.value);
      req.setHeaderName(this.headerName?.value);
      req.setIssuer(this.issuer?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.keysEndpoint?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as ManagementService)
        .addJWTProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.close();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    } else if (PolicyComponentServiceType.ADMIN) {
      const req = new AdminAddJWTProviderRequest();

      req.setName(this.name?.value);
      req.setHeaderName(this.headerName?.value);
      req.setIssuer(this.issuer?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.keysEndpoint?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as AdminService)
        .addJWTProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.close();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    }
  }

  public updateJWTProvider(): void {
    if (this.provider) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateJWTProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setHeaderName(this.headerName?.value);
        req.setIssuer(this.issuer?.value);
        req.setJwtEndpoint(this.jwtEndpoint?.value);
        req.setKeysEndpoint(this.keysEndpoint?.value);
        req.setProviderOptions(this.options);

        this.loading = true;
        (this.service as ManagementService)
          .updateJWTProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.close();
            }, 2000);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loading = false;
          });
      } else if (PolicyComponentServiceType.ADMIN) {
        const req = new AdminUpdateJWTProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setHeaderName(this.headerName?.value);
        req.setIssuer(this.issuer?.value);
        req.setJwtEndpoint(this.jwtEndpoint?.value);
        req.setKeysEndpoint(this.keysEndpoint?.value);
        req.setProviderOptions(this.options);

        this.loading = true;
        (this.service as AdminService)
          .updateJWTProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.close();
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

  public get name(): AbstractControl | null {
    return this.form.get('name');
  }

  public get headerName(): AbstractControl | null {
    return this.form.get('headerName');
  }

  public get issuer(): AbstractControl | null {
    return this.form.get('issuer');
  }

  public get jwtEndpoint(): AbstractControl | null {
    return this.form.get('jwtEndpoint');
  }

  public get keysEndpoint(): AbstractControl | null {
    return this.form.get('keysEndpoint');
  }
}

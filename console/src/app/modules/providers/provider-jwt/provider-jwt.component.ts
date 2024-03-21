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
import {IDPOwnerType, Options, Provider} from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddJWTProviderRequest as MgmtAddJWTProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateJWTProviderRequest as MgmtUpdateJWTProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { MatDialog } from '@angular/material/dialog';
import { ProviderNextService } from '../provider-next/provider-next.service';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { ProviderNextDialogComponent } from '../provider-next/provider-next-dialog.component';
import {CopyUrl} from "../provider-next/provider-next.component";

@Component({
  selector: 'cnsl-provider-jwt',
  templateUrl: './provider-jwt.component.html',
})
export class ProviderJWTComponent {
  public showOptional: boolean = false;
  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);

  public id: string | null = '';
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public form!: UntypedFormGroup;
  public loading: boolean = false;

  public provider?: Provider.AsObject;

  public autofillLink$ = new BehaviorSubject<string>('');
  public activateLink$ = new BehaviorSubject<string>('');
  public isActive$ = new BehaviorSubject<boolean>(false)
  public expandWhatNow$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public configureProvider$ = new BehaviorSubject<boolean>(false);
  public isInstance: boolean = false;

  constructor(
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
    private nextSvc: ProviderNextService,
  ) {

    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data['serviceType'];

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
          this.isInstance = true;
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
      }else {
        this.expandWhatNow$.next(true);
        this.configureProvider$.next(true);
      }
    });

    this.form = new UntypedFormGroup({
      name: new UntypedFormControl('', [requiredValidator]),
      headerName: new UntypedFormControl('', [requiredValidator]),
      issuer: new UntypedFormControl('', [requiredValidator]),
      jwtEndpoint: new UntypedFormControl('', [requiredValidator]),
      keysEndpoint: new UntypedFormControl('', [requiredValidator]),
    });

    this.authService
      .isAllowed(
        this.serviceType === PolicyComponentServiceType.ADMIN
          ? ['iam.idp.write']
          : this.serviceType === PolicyComponentServiceType.MGMT
            ? ['org.idp.write']
            : [],
      )
      .pipe(take(1))
      .subscribe((allowed) => {
        if (allowed) {
          this.form.enable();
        } else {
          this.form.disable();
        }
      });
  }

  private showAutofillLink(): void {
    this.autofillLink$.next('https://zitadel.com/docs/guides/integrate/identity-providers/additional-information');
  }

  private setActivateable(id: string) {
    this.activateLink$.next(!id ? '' : 'https://zitadel.com/docs/guides/integrate/identity-providers/okta-oidc#activate-idp');
    if (id) {
      this.expandWhatNow$.next(true);
      this.id = id;
    }
  }

  public activate() {
    this.service.addIDPToLoginPolicy(this.id!, this.serviceType === PolicyComponentServiceType.ADMIN ? IDPOwnerType.IDP_OWNER_TYPE_SYSTEM : IDPOwnerType.IDP_OWNER_TYPE_ORG).then(() => {
      this.toast.showInfo('POLICY.TOAST.ADDIDP', true);
      this.isActive$.next(true);
      this.setActivateable('');
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
          this.showAutofillLink();
          this.form.patchValue(this.provider.config.jwt);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
    this.service.getLoginPolicy()
      .then((policy) => {
        this.isActive$.next(!!policy.policy?.idpsList.find(idp => idp.idpId === this.id));
        this.setActivateable(this.isActive$.value ? '' : id);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public submitForm(): void {
    this.provider ? this.updateJWTProvider() : this.addJWTProvider();
  }

  public addJWTProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddJWTProviderRequest()
        : new AdminAddJWTProviderRequest();

    req.setName(this.name?.value);
    req.setHeaderName(this.headerName?.value);
    req.setIssuer(this.issuer?.value);
    req.setJwtEndpoint(this.jwtEndpoint?.value);
    req.setKeysEndpoint(this.keysEndpoint?.value);
    req.setProviderOptions(this.options);
    this.loading = true;
    this.service
      .addJWTProvider(req)
      .then((addedIDP) => {
        this.showAutofillLink();
        this.setActivateable(addedIDP.id);
        this.configureProvider$.next(false);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public updateJWTProvider(): void {
    if (this.provider) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateJWTProviderRequest()
          : new AdminUpdateJWTProviderRequest();
      req.setId(this.provider.id);
      req.setName(this.name?.value);
      req.setHeaderName(this.headerName?.value);
      req.setIssuer(this.issuer?.value);
      req.setJwtEndpoint(this.jwtEndpoint?.value);
      req.setKeysEndpoint(this.keysEndpoint?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      this.service
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

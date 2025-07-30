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
import { AutoLinkingOption, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
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
import { BehaviorSubject } from 'rxjs';
import { ProviderNextService } from '../provider-next/provider-next.service';

@Component({
  selector: 'cnsl-provider-jwt',
  templateUrl: './provider-jwt.component.html',
})
export class ProviderJWTComponent {
  public showOptional: boolean = false;
  public options: Options = new Options()
    .setIsCreationAllowed(true)
    .setIsLinkingAllowed(true)
    .setAutoLinking(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);

  // DEPRECATED: use id$ instead
  public id: string | null = '';
  // DEPRECATED: assert service$ instead
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  // DEPRECATED: use service$ instead
  private service!: ManagementService | AdminService;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public form!: UntypedFormGroup;
  public loading: boolean = false;

  public provider?: Provider.AsObject;

  public justCreated$: BehaviorSubject<string> = new BehaviorSubject<string>('');
  public justActivated$ = new BehaviorSubject<boolean>(false);

  private service$ = this.nextSvc.service(this.route.data);
  private id$ = this.nextSvc.id(this.route.paramMap, this.justCreated$);
  public exists$ = this.nextSvc.exists(this.id$);
  public autofillLink$ = this.nextSvc.autofillLink(
    this.id$,
    `https://zitadel.com/docs/guides/integrate/identity-providers/additional-information`,
  );
  public activateLink$ = this.nextSvc.activateLink(
    this.id$,
    this.justActivated$,
    'https://zitadel.com/docs/guides/integrate/identity-providers/okta-oidc#activate-idp',
    this.service$,
  );
  public expandWhatNow$ = this.nextSvc.expandWhatNow(this.id$, this.activateLink$, this.justCreated$);

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

  public activate() {
    this.nextSvc.activate(this.id$, this.justActivated$, this.service$);
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
    this.provider || this.justCreated$.value ? this.updateJWTProvider() : this.addJWTProvider();
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
        this.justCreated$.next(addedIDP.id);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public updateJWTProvider(): void {
    if (this.provider || this.justCreated$.value) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateJWTProviderRequest()
          : new AdminUpdateJWTProviderRequest();
      req.setId(this.provider?.id || this.justCreated$.value);
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

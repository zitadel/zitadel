import {Component, Injector, Type} from '@angular/core';
import {Location} from '@angular/common';
import {IDPOwnerType, Options, Provider, SAMLBinding} from '../../../proto/generated/zitadel/idp_pb';
import {AbstractControl, FormGroup, UntypedFormControl, UntypedFormGroup} from '@angular/forms';
import {PolicyComponentServiceType} from '../../policies/policy-component-types.enum';
import {ManagementService} from '../../../services/mgmt.service';
import {AdminService} from '../../../services/admin.service';
import {ToastService} from '../../../services/toast.service';
import {GrpcAuthService} from '../../../services/grpc-auth.service';
import {BehaviorSubject, combineLatestWith, from, Observable, of, Subject, switchMap, take} from 'rxjs';
import {ActivatedRoute, Router} from '@angular/router';
import {Breadcrumb, BreadcrumbService, BreadcrumbType} from '../../../services/breadcrumb.service';
import {atLeastOneIsFilled, requiredValidator} from '../../form-field/validators/validators';
import {
  AddSAMLProviderRequest as AdminAddSAMLProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateSAMLProviderRequest as AdminUpdateSAMLProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddSAMLProviderRequest as MgmtAddSAMLProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateSAMLProviderRequest as MgmtUpdateSAMLProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import {Environment, EnvironmentService} from '../../../services/environment.service';
import {CopyUrl} from "../provider-next/provider-next.component";
import {combineLatest, filter, map, tap,} from "rxjs/operators";

@Component({
  selector: 'cnsl-provider-saml-sp',
  templateUrl: './provider-saml-sp.component.html',
  styleUrls: ['./provider-saml-sp.component.scss'],
})
export class ProviderSamlSpComponent {
  public id: string | null = '';
  public loading: boolean = false;
  public provider?: Provider.AsObject;
  public form!: FormGroup;
  public showOptional: boolean = false;
  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;

  private created$: Subject<string> = new BehaviorSubject<string>('');
  private id$: Observable<string|null> = this.route.paramMap.pipe(
    // The ID observable should also emit when the IDP was just created
    combineLatestWith(this.created$),
    map(([params, created]) => created ? created : params.get('id')),
  )
  public exists$: Observable<boolean> = this.id$.pipe(
    map(id => !!id),
  )
  public autofillLink$ = this.id$.pipe(
    filter(id => !!id),
    map(()=> `https://zitadel.com/docs/guides/integrate/identity-providers/mocksaml#optional-add-zitadel-action-to-autofill-userdata`)
  );
  public activated$ = new BehaviorSubject<boolean>(false);
  public activateLink$: Observable<string> = this.id$.pipe(
    combineLatestWith(this.activated$),
    // Because we also want to emit when the IDP is not active, we return an empty string if the IDP does not exist
    switchMap(([id, activated]) =>     (!id || activated ? of(false) : from(this.service.getLoginPolicy()).pipe(
      map(policy => !policy.policy?.idpsList.find(idp => idp.idpId === id)),
    )).pipe(
      map((show) => !show ? '' : 'https://zitadel.com/docs/guides/integrate/identity-providers/mocksaml#activate-idp'),
      tap(console.log),
    ),
  ))
  // we expand initially if the IDP does not exist or if the idp was just created
  public expandWhatNow$ = this.id$.pipe(
    combineLatestWith(this.activateLink$, this.created$),
    map(([id, activateLink, created]) => !id || activateLink || created),
  );
  public copyUrls$: Observable<CopyUrl[]> = this.id$.pipe(
    filter(id => !!id),
    switchMap(id => this.envSvc.env.pipe(
      map((environment: Environment) => {
        const idpBase = `${environment.issuer}/idps/${id}/saml`;
        return [
          {
            label: 'ZITADEL Metadata',
            url: `${idpBase}/metadata`,
            downloadable: true,
          },
          {
            label: 'ZITADEL Single Logout',
            url: `${idpBase}/slo`,
          },
          {
            label: 'ZITADEL ACS Login Form',
            url: `${environment.issuer}/ui/login/login/externalidp/saml/acs`,
          },
          {
            label: 'ZITADEL ACS Intent API',
            url: `${idpBase}/acs`,
          },
        ];
      })
    ))
  )

  public isInstance: boolean = false;

  bindingValues: string[] = Object.keys(SAMLBinding);

  constructor(
    private _location: Location,
    private toast: ToastService,
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private router: Router,
    private injector: Injector,
    private breadcrumbService: BreadcrumbService,
    private envSvc: EnvironmentService,
  ) {
    this._buildBreadcrumbs();
    this._initializeForm();
    this._checkFormPermissions();
  }

  private _initializeForm(): void {
    this.form = new UntypedFormGroup(
      {
        name: new UntypedFormControl('', [requiredValidator]),
        metadataXml: new UntypedFormControl('', []),
        metadataUrl: new UntypedFormControl('', []),
        binding: new UntypedFormControl(this.bindingValues[0], [requiredValidator]),
        withSignedRequest: new UntypedFormControl(true, [requiredValidator]),
      },
      atLeastOneIsFilled('metadataXml', 'metadataUrl'),
    );
  }

  private _checkFormPermissions(): void {
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

  private _buildBreadcrumbs(): void {
    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data['serviceType'];
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
          this.isInstance = true
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
        this.getData(this.id);
      }
   });
  }

  public updateSAMLProvider(): void {
    if (this.provider) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateSAMLProviderRequest()
          : new AdminUpdateSAMLProviderRequest();

      req.setId(this.provider.id);
      req.setName(this.name?.value);
      if (this.metadataXml?.value) {
        req.setMetadataXml(this.metadataXml?.value);
      } else {
        req.setMetadataUrl(this.metadataUrl?.value);
      }
      req.setWithSignedRequest(this.withSignedRequest?.value);
      // @ts-ignore
      req.setBinding(SAMLBinding[this.binding?.value]);
      req.setProviderOptions(this.options);

      this.loading = true;
      this.service
        .updateSAMLProvider(req)
        .then(() => {
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

  public addSAMLProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddSAMLProviderRequest()
        : new AdminAddSAMLProviderRequest();
    req.setName(this.name?.value);
    if (this.metadataXml?.value) {
      req.setMetadataXml(this.metadataXml?.value);
    } else {
      req.setMetadataUrl(this.metadataUrl?.value);
    }
    req.setProviderOptions(this.options);
    // @ts-ignore
    req.setBinding(SAMLBinding[this.binding?.value]);
    req.setWithSignedRequest(this.withSignedRequest?.value);
    this.loading = true;
    this.service
      .addSAMLProvider(req)
      .then((addedIDP) => {
        this.created$.next(addedIDP.id);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider ? this.updateSAMLProvider() : this.addSAMLProvider();
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
        this.name?.setValue(this.provider?.name);
        if (this.provider?.config?.saml) {
          this.form.patchValue(this.provider.config.saml);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  close(): void {
    this._location.back();
  }

  compareBinding(value: string, index: number) {
    if (value) {
      return value === Object.keys(SAMLBinding)[index];
    }
    return false;
  }

  public activate() {
    this.id$.pipe(
      take(1),
      switchMap(id => from(this.service.addIDPToLoginPolicy(id!, this.serviceType === PolicyComponentServiceType.ADMIN ? IDPOwnerType.IDP_OWNER_TYPE_SYSTEM : IDPOwnerType.IDP_OWNER_TYPE_ORG))),
    ).subscribe({
      next: () => {
        this.toast.showInfo('POLICY.TOAST.ADDIDP', true);
        this.activated$.next(true)
      },
      error: error => this.toast.showError(error),
    })
  }

  private get name(): AbstractControl | null {
    return this.form.get('name');
  }

  private get metadataXml(): AbstractControl | null {
    return this.form.get('metadataXml');
  }

  private get metadataUrl(): AbstractControl | null {
    return this.form.get('metadataUrl');
  }

  private get binding(): AbstractControl | null {
    return this.form.get('binding');
  }

  private get withSignedRequest(): AbstractControl | null {
    return this.form.get('withSignedRequest');
  }
}

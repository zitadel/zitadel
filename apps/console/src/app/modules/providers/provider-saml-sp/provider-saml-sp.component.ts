import { Component, Injector, Type } from '@angular/core';
import { Location } from '@angular/common';
import {
  AutoLinkingOption,
  Options,
  Provider,
  SAMLBinding,
  SAMLNameIDFormat,
} from '../../../proto/generated/zitadel/idp_pb';
import { AbstractControl, FormGroup, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ManagementService } from '../../../services/mgmt.service';
import { AdminService } from '../../../services/admin.service';
import { ToastService } from '../../../services/toast.service';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { BehaviorSubject, shareReplay, switchMap, take } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from '../../../services/breadcrumb.service';
import { atLeastOneIsFilled, requiredValidator } from '../../form-field/validators/validators';
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
import { Environment, EnvironmentService } from '../../../services/environment.service';
import { filter, map } from 'rxjs/operators';
import { ProviderNextService } from '../provider-next/provider-next.service';

@Component({
  selector: 'cnsl-provider-saml-sp',
  templateUrl: './provider-saml-sp.component.html',
  styleUrls: ['./provider-saml-sp.component.scss'],
})
export class ProviderSamlSpComponent {
  // DEPRECATED: use id$ instead
  public id: string | null = '';
  public loading: boolean = false;
  public provider?: Provider.AsObject;
  public form!: FormGroup;
  public showOptional: boolean = false;
  public options: Options = new Options()
    .setIsCreationAllowed(true)
    .setIsLinkingAllowed(true)
    .setAutoLinking(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);
  // DEPRECATED: assert service$ instead
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  // DEPRECATED: use service$ instead
  private service!: ManagementService | AdminService;
  bindingValues: string[] = Object.keys(SAMLBinding);
  nameIDFormatValues: string[] = Object.keys(SAMLNameIDFormat);

  public justCreated$: BehaviorSubject<string> = new BehaviorSubject<string>('');
  public justActivated$ = new BehaviorSubject<boolean>(false);

  private service$ = this.nextSvc.service(this.route.data);
  private id$ = this.nextSvc.id(this.route.paramMap, this.justCreated$);
  public exists$ = this.nextSvc.exists(this.id$);
  public autofillLink$ = this.nextSvc.autofillLink(
    this.id$,
    `https://zitadel.com/docs/guides/integrate/identity-providers/mocksaml#optional-add-zitadel-action-to-autofill-userdata`,
  );
  public activateLink$ = this.nextSvc.activateLink(
    this.id$,
    this.justActivated$,
    'https://zitadel.com/docs/guides/integrate/identity-providers/mocksaml#activate-idp',
    this.service$,
  );
  public expandWhatNow$ = this.nextSvc.expandWhatNow(this.id$, this.activateLink$, this.justCreated$);
  public copyUrls$ = this.id$.pipe(
    filter((id) => !!id),
    switchMap((id) =>
      this.envSvc.env.pipe(
        map((environment: Environment) => {
          const idpBase = `${environment.issuer}/idps/${id}/saml`;
          return [
            {
              label: 'ZITADEL Metadata',
              url: `${idpBase}/metadata`,
              downloadable: true,
            },
            {
              label: 'ZITADEL ACS Login Form',
              url: `${environment.issuer}/ui/login/login/externalidp/saml/acs`,
            },
            {
              label: 'ZITADEL ACS Intent API',
              url: `${idpBase}/acs`,
            },
            {
              label: 'ZITADEL Single Logout',
              url: `${idpBase}/slo`,
            },
          ];
        }),
      ),
    ),
    shareReplay(1),
  );

  constructor(
    private _location: Location,
    private toast: ToastService,
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private injector: Injector,
    private breadcrumbService: BreadcrumbService,
    private envSvc: EnvironmentService,
    private nextSvc: ProviderNextService,
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
        nameIdFormat: new UntypedFormControl(SAMLNameIDFormat.SAML_NAME_ID_FORMAT_PERSISTENT, []),
        transientMappingAttributeName: new UntypedFormControl('', []),
        federatedLogoutEnabled: new UntypedFormControl(false, []),
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

  public activate() {
    this.nextSvc.activate(this.id$, this.justActivated$, this.service$);
  }

  public updateSAMLProvider(): void {
    if (this.provider || this.justCreated$.value) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateSAMLProviderRequest()
          : new AdminUpdateSAMLProviderRequest();

      req.setId(this.provider?.id || this.justCreated$.value);
      req.setName(this.name?.value);
      if (this.metadataXml?.value) {
        req.setMetadataUrl('');
        req.setMetadataXml(this.metadataXml?.value);
      } else {
        req.setMetadataXml('');
        req.setMetadataUrl(this.metadataUrl?.value);
      }
      req.setWithSignedRequest(this.withSignedRequest?.value);
      // @ts-ignore
      req.setBinding(SAMLBinding[this.binding?.value]);
      // @ts-ignore
      req.setNameIdFormat(SAMLNameIDFormat[this.nameIDFormat?.value]);
      req.setTransientMappingAttributeName(this.transientMapping?.value);
      req.setFederatedLogoutEnabled(this.federatedLogoutEnabled?.value);
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
      req.setMetadataUrl('');
      req.setMetadataXml(this.metadataXml?.value);
    } else {
      req.setMetadataXml('');
      req.setMetadataUrl(this.metadataUrl?.value);
    }
    req.setProviderOptions(this.options);
    // @ts-ignore
    req.setBinding(SAMLBinding[this.binding?.value]);
    req.setWithSignedRequest(this.withSignedRequest?.value);
    if (this.nameIDFormat) {
      // @ts-ignore
      req.setNameIdFormat(SAMLNameIDFormat[this.nameIDFormat.value]);
    }
    req.setTransientMappingAttributeName(this.transientMapping?.value);
    req.setFederatedLogoutEnabled(this.federatedLogoutEnabled?.value);
    this.loading = true;
    this.service
      .addSAMLProvider(req)
      .then((addedIDP) => {
        this.justCreated$.next(addedIDP.id);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider || this.justCreated$.value ? this.updateSAMLProvider() : this.addSAMLProvider();
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

  compareNameIDFormat(value: string, index: number) {
    console.log(value, index);
    if (value) {
      return value === Object.keys(SAMLNameIDFormat)[index];
    }
    return false;
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

  private get nameIDFormat(): AbstractControl | null {
    return this.form.get('nameIdFormat');
  }

  private get transientMapping(): AbstractControl | null {
    return this.form.get('transientMappingAttributeName');
  }

  private get federatedLogoutEnabled(): AbstractControl | null {
    return this.form.get('federatedLogoutEnabled');
  }
}

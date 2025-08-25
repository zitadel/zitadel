import { Component, Injector, Type } from '@angular/core';
import { Location } from '@angular/common';
import {
  AutoLinkingOption,
  Options,
  Provider,
  SAMLBinding,
  SAMLNameIDFormat,
} from '../../../proto/generated/zitadel/idp_pb';
import { AbstractControl, FormControl, FormGroup } from '@angular/forms';
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
import { getEnumKeys, getEnumKeyFromValue, convertEnumValuesToKeys } from '../../../utils/enum.utils';

interface SAMLProviderForm {
  name: FormControl<string>;
  metadataXml: FormControl<string>;
  metadataUrl: FormControl<string>;
  binding: FormControl<string>;
  withSignedRequest: FormControl<boolean>;
  nameIdFormat: FormControl<string>;
  transientMappingAttributeName: FormControl<string>;
  federatedLogoutEnabled: FormControl<boolean>;
}

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
  public form!: FormGroup<SAMLProviderForm>;
  public showOptional: boolean = false;
  public options: Options = new Options()
    .setIsCreationAllowed(true)
    .setIsLinkingAllowed(true)
    .setAutoLinking(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);
  // DEPRECATED: assert service$ instead
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  // DEPRECATED: use service$ instead
  private service!: ManagementService | AdminService;
  bindingValues: string[] = getEnumKeys(SAMLBinding);
  nameIDFormatValues: string[] = getEnumKeys(SAMLNameIDFormat);

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
    const defaultBinding = getEnumKeyFromValue(SAMLBinding, SAMLBinding.SAML_BINDING_POST) || this.bindingValues[0];
    const defaultNameIdFormat = getEnumKeyFromValue(SAMLNameIDFormat, SAMLNameIDFormat.SAML_NAME_ID_FORMAT_PERSISTENT) || this.nameIDFormatValues[0];
    
    this.form = new FormGroup<SAMLProviderForm>(
      {
        name: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
        metadataXml: new FormControl('', { nonNullable: true }),
        metadataUrl: new FormControl('', { nonNullable: true }),
        binding: new FormControl(defaultBinding, { nonNullable: true, validators: [requiredValidator] }),
        withSignedRequest: new FormControl(true, { nonNullable: true, validators: [requiredValidator] }),
        nameIdFormat: new FormControl(defaultNameIdFormat, { nonNullable: true }),
        transientMappingAttributeName: new FormControl('', { nonNullable: true }),
        federatedLogoutEnabled: new FormControl(false, { nonNullable: true }),
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
      req.setName(this.name.value);
      if (this.metadataXml.value) {
        req.setMetadataUrl('');
        req.setMetadataXml(this.metadataXml.value);
      } else {
        req.setMetadataXml('');
        req.setMetadataUrl(this.metadataUrl.value);
      }
      req.setWithSignedRequest(this.withSignedRequest.value);
      req.setBinding(SAMLBinding[this.binding.value as keyof typeof SAMLBinding]);
      req.setNameIdFormat(SAMLNameIDFormat[this.nameIDFormat.value as keyof typeof SAMLNameIDFormat]);
      req.setTransientMappingAttributeName(this.transientMapping.value);
      req.setFederatedLogoutEnabled(this.federatedLogoutEnabled.value);
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
    req.setName(this.name.value);
    if (this.metadataXml.value) {
      req.setMetadataUrl('');
      req.setMetadataXml(this.metadataXml.value);
    } else {
      req.setMetadataXml('');
      req.setMetadataUrl(this.metadataUrl.value);
    }
    req.setProviderOptions(this.options);
    req.setBinding(SAMLBinding[this.binding.value as keyof typeof SAMLBinding]);
    req.setWithSignedRequest(this.withSignedRequest.value);
    if (this.nameIDFormat) {
      req.setNameIdFormat(SAMLNameIDFormat[this.nameIDFormat.value as keyof typeof SAMLNameIDFormat]);
    }
    req.setTransientMappingAttributeName(this.transientMapping.value);
    req.setFederatedLogoutEnabled(this.federatedLogoutEnabled.value);
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
        this.name.setValue(this.provider?.name || '');
        if (this.provider?.config?.saml) {
          const samlConfig = this.provider.config.saml;
          const bindingKey = getEnumKeyFromValue(SAMLBinding, samlConfig.binding) || '';
          const nameIdFormatKey = getEnumKeyFromValue(SAMLNameIDFormat, samlConfig.nameIdFormat) || '';
          
          this.form.patchValue({
            metadataXml: typeof samlConfig.metadataXml === 'string' ? samlConfig.metadataXml : '',
            binding: bindingKey,
            withSignedRequest: samlConfig.withSignedRequest,
            nameIdFormat: nameIdFormatKey,
            transientMappingAttributeName: samlConfig.transientMappingAttributeName || '',
            federatedLogoutEnabled: samlConfig.federatedLogoutEnabled || false,
          });
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

  private get name(): FormControl<string> {
    return this.form.controls.name;
  }

  private get metadataXml(): FormControl<string> {
    return this.form.controls.metadataXml;
  }

  private get metadataUrl(): FormControl<string> {
    return this.form.controls.metadataUrl;
  }

  private get binding(): FormControl<string> {
    return this.form.controls.binding;
  }

  private get withSignedRequest(): FormControl<boolean> {
    return this.form.controls.withSignedRequest;
  }

  private get nameIDFormat(): FormControl<string> {
    return this.form.controls.nameIdFormat;
  }

  private get transientMapping(): FormControl<string> {
    return this.form.controls.transientMappingAttributeName;
  }

  private get federatedLogoutEnabled(): FormControl<boolean> {
    return this.form.controls.federatedLogoutEnabled;
  }
}

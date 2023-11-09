import { Component, Injector, Type } from '@angular/core';
import { Location } from '@angular/common';
import { Options, Provider } from '../../../proto/generated/zitadel/idp_pb';
import { AbstractControl, FormGroup, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ManagementService } from '../../../services/mgmt.service';
import { AdminService } from '../../../services/admin.service';
import { ToastService } from '../../../services/toast.service';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { take } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from '../../../services/breadcrumb.service';
import { requiredValidator } from '../../form-field/validators/validators';
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
import * as zitadel_idp_pb from '../../../proto/generated/zitadel/idp_pb';

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

  bindingValues: string[] = Object.keys(zitadel_idp_pb.SAMLBinding);

  constructor(
    private _location: Location,
    private toast: ToastService,
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private injector: Injector,
    private breadcrumbService: BreadcrumbService,
  ) {
    this._buildBreadcrumbs();
    this._initializeForm();
    this._checkFormPermissions();
  }

  private _initializeForm(): void {
    this.form = new UntypedFormGroup({
      name: new UntypedFormControl('', [requiredValidator]),
      metadataXml: new UntypedFormControl('', [requiredValidator]),
      metadataUrl: new UntypedFormControl('', [requiredValidator]),
      binding: new UntypedFormControl(this.bindingValues[0], [requiredValidator]),
      withSignedRequest: new UntypedFormControl(true, [requiredValidator]),
    });
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

  public updateSAMLProvider(): void {
    if (this.provider) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateSAMLProviderRequest()
          : new AdminUpdateSAMLProviderRequest();
      req.setId(this.provider.id);
      req.setName(this.name?.value);
      req.setMetadataUrl(this.metadataUrl?.value);
      req.setMetadataXml(this.metadataXml?.value);
      // @ts-ignore
      req.setBinding(zitadel_idp_pb.SAMLBinding[`${this.biding?.value}`]);
      req.setProviderOptions(this.options);

      this.loading = true;
      this.service
        .updateSAMLProvider(req)
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

  public addSAMLProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddSAMLProviderRequest()
        : new AdminAddSAMLProviderRequest();
    req.setName(this.name?.value);
    req.setMetadataUrl(this.metadataUrl?.value);
    req.setMetadataXml(this.metadataXml?.value);
    req.setProviderOptions(this.options);
    // @ts-ignore
    req.setBinding(zitadel_idp_pb.SAMLBinding[`${this.biding?.value}`]);
    req.setWithSignedRequest(this.withSignedRequest?.value);
    this.loading = true;
    this.service
      .addSAMLProvider(req)
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
        if (this.provider?.config?.saml) {
          this.form.patchValue(this.provider.config.saml);
          this.name?.setValue(this.provider.name);
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

  private get name(): AbstractControl | null {
    return this.form.get('name');
  }

  private get metadataXml(): AbstractControl | null {
    return this.form.get('metadataXml');
  }

  private get metadataUrl(): AbstractControl | null {
    return this.form.get('metadataUrl');
  }

  private get biding(): AbstractControl | null {
    return this.form.get('binding');
  }

  private get withSignedRequest(): AbstractControl | null {
    return this.form.get('withSignedRequest');
  }
}

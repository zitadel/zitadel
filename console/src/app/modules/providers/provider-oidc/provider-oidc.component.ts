import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute } from '@angular/router';
import { BehaviorSubject, take } from 'rxjs';
import {
  AddGenericOIDCProviderRequest as AdminAddGenericOIDCProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateGenericOIDCProviderRequest as AdminUpdateGenericOIDCProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AutoLinkingOption, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGenericOIDCProviderRequest as MgmtAddGenericOIDCProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateGenericOIDCProviderRequest as MgmtUpdateGenericOIDCProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ProviderNextService } from '../provider-next/provider-next.service';

@Component({
  selector: 'cnsl-provider-oidc',
  templateUrl: './provider-oidc.component.html',
  styleUrls: ['./provider-oidc.component.scss'],
})
export class ProviderOIDCComponent {
  public showOptional: boolean = false;
  public options: Options = new Options()
    .setIsCreationAllowed(true)
    .setIsLinkingAllowed(true)
    .setAutoLinking(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);

  // DEPRECATED: use id$ instead
  public id: string | null = '';
  public updateClientSecret: boolean = false;
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
  public copyUrls$ = this.nextSvc.callbackUrls();

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    breadcrumbService: BreadcrumbService,
    private nextSvc: ProviderNextService,
  ) {
    this.form = new UntypedFormGroup({
      name: new UntypedFormControl('', [requiredValidator]),
      clientId: new UntypedFormControl('', [requiredValidator]),
      clientSecret: new UntypedFormControl('', [requiredValidator]),
      issuer: new UntypedFormControl('', [requiredValidator]),
      scopesList: new UntypedFormControl(['openid', 'profile', 'email'], []),
      isIdTokenMapping: new UntypedFormControl(),
      usePkce: new UntypedFormControl(false),
    });

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
        this.clientSecret?.setValidators([]);
        this.getData(this.id);
      }
    });
  }

  public activate() {
    this.nextSvc.activate(this.id$, this.justActivated$, this.service$);
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
        if (this.provider?.config?.oidc) {
          this.form.patchValue(this.provider.config.oidc);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider || this.justCreated$.value ? this.updateGenericOIDCProvider() : this.addGenericOIDCProvider();
  }

  public addGenericOIDCProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddGenericOIDCProviderRequest()
        : new AdminAddGenericOIDCProviderRequest();

    req.setName(this.name?.value);
    req.setClientId(this.clientId?.value);
    req.setClientSecret(this.clientSecret?.value);
    req.setIssuer(this.issuer?.value);
    req.setScopesList(this.scopesList?.value);
    req.setProviderOptions(this.options);
    req.setIsIdTokenMapping(this.isIdTokenMapping?.value);
    req.setUsePkce(this.usePkce?.value);

    this.loading = true;
    this.service
      .addGenericOIDCProvider(req)
      .then((addedIDP) => {
        this.justCreated$.next(addedIDP.id);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public updateGenericOIDCProvider(): void {
    if (this.provider || this.justCreated$.value) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateGenericOIDCProviderRequest()
          : new AdminUpdateGenericOIDCProviderRequest();
      req.setId(this.provider?.id || this.justCreated$.value);
      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setProviderOptions(this.options);
      req.setIsIdTokenMapping(this.isIdTokenMapping?.value);
      req.setUsePkce(this.usePkce?.value);

      this.loading = true;
      this.service
        .updateGenericOIDCProvider(req)
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

  public get issuer(): AbstractControl | null {
    return this.form.get('issuer');
  }

  public get scopesList(): AbstractControl | null {
    return this.form.get('scopesList');
  }

  public get isIdTokenMapping(): AbstractControl | null {
    return this.form.get('isIdTokenMapping');
  }

  public get usePkce(): AbstractControl | null {
    return this.form.get('usePkce');
  }
}

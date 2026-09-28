import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute } from '@angular/router';
import { BehaviorSubject, take } from 'rxjs';
import {
  AddZitadelProviderRequest as AdminAddZitadelProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateZitadelProviderRequest as AdminUpdateZitadelProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AutoLinkingOption, InstanceRolesInfo, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddZitadelProviderRequest as MgmtAddZitadelProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateZitadelProviderRequest as MgmtUpdateZitadelProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ProviderNextService } from '../provider-next/provider-next.service';

@Component({
  selector: 'cnsl-provider-zitadel',
  templateUrl: './provider-zitadel.component.html',
  styleUrls: ['./provider-zitadel.component.scss'],
  standalone: false,
})
export class ProviderZitadelComponent {
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

  // Instance role mapping (instance-level providers only): organizations of the
  // remote ZITADEL instance whose project-role grants with an IAM_ prefix are
  // applied as instance memberships after login.
  public instanceRolesInfoList: InstanceRolesInfo.AsObject[] = [];
  public rolesInfoOrgId: UntypedFormControl = new UntypedFormControl('');
  public rolesInfoOrgDomain: UntypedFormControl = new UntypedFormControl('');

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
    'https://zitadel.com/docs/guides/integrate/identity-providers',
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

  // The role mapping is only available on instance-level providers: the API
  // rejects organization-scoped ZITADEL providers carrying instanceRolesInfo.
  public get isAdmin(): boolean {
    return this.serviceType === PolicyComponentServiceType.ADMIN;
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
        if (this.provider?.config?.zitadel) {
          this.form.patchValue(this.provider.config.zitadel);
          this.name?.setValue(this.provider.name);
          this.instanceRolesInfoList = this.provider.config.zitadel.instanceRolesInfoList ?? [];
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider || this.justCreated$.value ? this.updateZitadelProvider() : this.addZitadelProvider();
  }

  public addZitadelProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddZitadelProviderRequest();
      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as ManagementService)
        .addZitadelProvider(req)
        .then((addedIDP) => {
          this.justCreated$.next(addedIDP.id);
          this.loading = false;
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    } else {
      const req = new AdminAddZitadelProviderRequest();
      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setProviderOptions(this.options);
      req.setInstanceRolesInfoList(this.instanceRolesInfoToProto());

      this.loading = true;
      (this.service as AdminService)
        .addZitadelProvider(req)
        .then((addedIDP) => {
          this.justCreated$.next(addedIDP.id);
          this.loading = false;
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    }
  }

  public updateZitadelProvider(): void {
    if (this.provider || this.justCreated$.value) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateZitadelProviderRequest();
        req.setId(this.provider?.id || this.justCreated$.value);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecretForUpdate);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);
        req.setProviderOptions(this.options);

        this.loading = true;
        (this.service as ManagementService)
          .updateZitadelProvider(req)
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
      } else {
        const req = new AdminUpdateZitadelProviderRequest();
        req.setId(this.provider?.id || this.justCreated$.value);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecretForUpdate);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);
        req.setProviderOptions(this.options);
        req.setInstanceRolesInfoList(this.instanceRolesInfoToProto());

        this.loading = true;
        (this.service as AdminService)
          .updateZitadelProvider(req)
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
  }

  /**
   * The client secret to send on update. An existing provider only receives a
   * new secret when the user explicitly opted in: the form control keeps its
   * value when the "update client secret" checkbox is unticked again, and
   * sending it would silently rotate the secret. An empty value is treated as
   * "unchanged" by the API.
   */
  private get clientSecretForUpdate(): string {
    return !this.provider || this.updateClientSecret ? (this.clientSecret?.value ?? '') : '';
  }

  private instanceRolesInfoToProto(): InstanceRolesInfo[] {
    return this.instanceRolesInfoList.map((info) =>
      new InstanceRolesInfo().setOrganizationId(info.organizationId).setOrganizationDomain(info.organizationDomain),
    );
  }

  public addRolesInfo(): void {
    const organizationId = (this.rolesInfoOrgId.value ?? '').trim();
    const organizationDomain = (this.rolesInfoOrgDomain.value ?? '').trim();
    // Both fields are required by the API and matched together on login.
    if (!organizationId || !organizationDomain) {
      return;
    }
    this.instanceRolesInfoList = [...this.instanceRolesInfoList, { organizationId, organizationDomain }];
    this.rolesInfoOrgId.setValue('');
    this.rolesInfoOrgDomain.setValue('');
  }

  public removeRolesInfo(index: number): void {
    this.instanceRolesInfoList = this.instanceRolesInfoList.filter((_, i) => i !== index);
  }

  public close(): void {
    this._location.back();
  }

  public addScope(event: MatChipInputEvent): void {
    this.addScopeFromInput(event.chipInput?.inputElement, event.value);
  }

  /**
   * Adds the scope currently typed into the chip input. Takes the value from
   * the input element because the add button emits a MouseEvent, which carries
   * no value of its own.
   */
  public addScopeFromInput(input?: HTMLInputElement, value?: string): void {
    const scope = (value ?? input?.value ?? '').trim();

    if (scope !== '' && this.scopesList?.value) {
      this.scopesList.value.push(scope);
      if (input) {
        input.value = '';
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
}

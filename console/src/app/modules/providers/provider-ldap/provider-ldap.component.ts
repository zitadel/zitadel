import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { BehaviorSubject, take } from 'rxjs';
import {
  AddLDAPProviderRequest as AdminAddLDAPProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateLDAPProviderRequest as AdminUpdateLDAPProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AutoLinkingOption, LDAPAttributes, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddLDAPProviderRequest as MgmtAddLDAPProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateLDAPProviderRequest as MgmtUpdateLDAPProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { minArrayLengthValidator, requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ProviderNextService } from '../provider-next/provider-next.service';

@Component({
  selector: 'cnsl-provider-ldap',
  templateUrl: './provider-ldap.component.html',
})
export class ProviderLDAPComponent {
  public updateBindPassword: boolean = false;
  public showOptional: boolean = false;
  public options: Options = new Options()
    .setIsCreationAllowed(true)
    .setIsLinkingAllowed(true)
    .setAutoLinking(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);
  public attributes: LDAPAttributes = new LDAPAttributes();
  // DEPRECATED: use id$ instead
  public id: string | null = '';
  // DEPRECATED: assert service$ instead
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  // DEPRECATED: use service$ instead
  private service!: ManagementService | AdminService;

  public form!: FormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;

  public justCreated$: BehaviorSubject<string> = new BehaviorSubject<string>('');
  public justActivated$ = new BehaviorSubject<boolean>(false);

  private service$ = this.nextSvc.service(this.route.data);
  private id$ = this.nextSvc.id(this.route.paramMap, this.justCreated$);
  public exists$ = this.nextSvc.exists(this.id$);
  public activateLink$ = this.nextSvc.activateLink(
    this.id$,
    this.justActivated$,
    'https://zitadel.com/docs/guides/integrate/identity-providers/google#activate-idp',
    this.service$,
  );
  public expandWhatNow$ = this.nextSvc.expandWhatNow(this.id$, this.activateLink$, this.justCreated$);

  constructor(
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
    private nextSvc: ProviderNextService,
  ) {
    this.form = new FormGroup({
      name: new FormControl('', [requiredValidator]),
      serversList: new FormControl<string[]>([''], [minArrayLengthValidator(1)]),
      baseDn: new FormControl('', [requiredValidator]),
      bindDn: new FormControl('', [requiredValidator]),
      bindPassword: new FormControl('', [requiredValidator]),
      userBase: new FormControl('', [requiredValidator]),
      userFiltersList: new FormControl<string[]>([''], [minArrayLengthValidator(1)]),
      userObjectClassesList: new FormControl<string[]>([''], [minArrayLengthValidator(1)]),
      timeout: new FormControl<number>(0),
      startTls: new FormControl<boolean>(false),
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
        this.bindPassword?.setValidators([]);
        this.bindPassword?.updateValueAndValidity();
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
        if (resp.idp) {
          this.provider = resp.idp;
          this.loading = false;

          this.name?.setValue(this.provider.name);

          const config = this.provider?.config?.ldap;
          if (config) {
            this.serversList?.setValue(config.serversList);
            this.startTls?.setValue(config.startTls);
            this.baseDn?.setValue(config.baseDn);
            this.bindDn?.setValue(config.bindDn);
            this.userBase?.setValue(config.userBase);
            this.userObjectClassesList?.setValue(config.userObjectClassesList);
            this.userFiltersList?.setValue(config.userFiltersList);
            if (this.provider?.config?.ldap?.timeout?.seconds) {
              this.timeout?.setValue(this.provider?.config?.ldap?.timeout?.seconds);
            }
          }
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider || this.justCreated$.value ? this.updateLDAPProvider() : this.addLDAPProvider();
  }

  public addLDAPProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddLDAPProviderRequest()
        : new AdminAddLDAPProviderRequest();

    req.setName(this.name?.value);
    req.setProviderOptions(this.options);
    req.setAttributes(this.attributes);

    req.setBaseDn(this.baseDn?.value);
    req.setBindDn(this.bindDn?.value);
    req.setBindPassword(this.bindPassword?.value);
    req.setServersList(this.serversList?.value); // list
    req.setStartTls(this.startTls?.value);
    req.setTimeout(new Duration().setSeconds(this.timeout?.value ?? 0));
    req.setUserBase(this.userBase?.value);
    req.setUserFiltersList(this.userFiltersList?.value); // list
    req.setUserObjectClassesList(this.userObjectClassesList?.value); // list

    this.loading = true;
    (this.service as ManagementService)
      .addLDAPProvider(req)
      .then((addedIDP) => {
        this.justCreated$.next(addedIDP.id);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public updateLDAPProvider(): void {
    if (this.provider) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateLDAPProviderRequest()
          : new AdminUpdateLDAPProviderRequest();
      req.setId(this.provider.id);
      req.setName(this.name?.value);
      req.setProviderOptions(this.options);

      req.setAttributes(this.attributes);

      req.setBaseDn(this.baseDn?.value);
      req.setBindDn(this.bindDn?.value);
      if (this.updateBindPassword) {
        req.setBindPassword(this.bindPassword?.value);
      }
      req.setServersList(this.serversList?.value);
      req.setStartTls(this.startTls?.value);
      req.setTimeout(new Duration().setSeconds(this.timeout?.value ?? 0));
      req.setUserBase(this.userBase?.value);
      req.setUserFiltersList(this.userFiltersList?.value);
      req.setUserObjectClassesList(this.userObjectClassesList?.value);

      this.loading = true;
      (this.service as ManagementService)
        .updateLDAPProvider(req)
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

  public get baseDn(): AbstractControl | null {
    return this.form.get('baseDn');
  }

  public get bindDn(): AbstractControl | null {
    return this.form.get('bindDn');
  }

  public get bindPassword(): AbstractControl | null {
    return this.form.get('bindPassword');
  }

  public get serversList(): AbstractControl | null {
    return this.form.get('serversList');
  }

  public get startTls(): AbstractControl | null {
    return this.form.get('startTls');
  }

  public get timeout(): AbstractControl | null {
    return this.form.get('timeout');
  }

  public get userBase(): AbstractControl | null {
    return this.form.get('userBase');
  }

  public get userFiltersList(): AbstractControl | null {
    return this.form.get('userFiltersList');
  }

  public get userObjectClassesList(): AbstractControl | null {
    return this.form.get('userObjectClassesList');
  }
}

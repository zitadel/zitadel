import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute } from '@angular/router';
import { BehaviorSubject, take } from 'rxjs';
import {
  AddAppleProviderRequest as AdminAddAppleProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateAppleProviderRequest as AdminUpdateAppleProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AutoLinkingOption, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddAppleProviderRequest as MgmtAddAppleProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateAppleProviderRequest as MgmtUpdateAppleProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ProviderNextService } from '../provider-next/provider-next.service';

const MAX_ALLOWED_SIZE = 5 * 1024;

@Component({
  selector: 'cnsl-provider-apple',
  templateUrl: './provider-apple.component.html',
})
export class ProviderAppleComponent {
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

  public form!: FormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;
  public updatePrivateKey: boolean = false;

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
    'https://zitadel.com/docs/guides/integrate/identity-providers/apple#activate-idp',
    this.service$,
  );
  public expandWhatNow$ = this.nextSvc.expandWhatNow(this.id$, this.activateLink$, this.justCreated$);
  public copyUrls$ = this.nextSvc.callbackUrls();

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
      name: new FormControl('', []),
      clientId: new FormControl('', [requiredValidator]),
      teamId: new FormControl('', [requiredValidator]),
      keyId: new FormControl('', [requiredValidator]),
      privateKey: new FormControl('', [requiredValidator]),
      scopesList: new FormControl(['name', 'email'], []),
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
        this.privateKey?.setValidators([]);
        this.getData(this.id);
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
        if (this.provider?.config?.apple) {
          this.form.patchValue(this.provider.config.apple);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider || this.justCreated$.value ? this.updateAppleProvider() : this.addAppleProvider();
  }

  public addAppleProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddAppleProviderRequest()
        : new AdminAddAppleProviderRequest();

    req.setName(this.name?.value);
    req.setClientId(this.clientId?.value);
    req.setTeamId(this.teamId?.value);
    req.setKeyId(this.keyId?.value);
    req.setPrivateKey(this.privateKey?.value);
    req.setScopesList(this.scopesList?.value);
    req.setProviderOptions(this.options);

    this.loading = true;
    this.service
      .addAppleProvider(req)
      .then((addedIDP) => {
        this.justCreated$.next(addedIDP.id);
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public updateAppleProvider(): void {
    if (this.provider || this.justCreated$.value) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateAppleProviderRequest();
        req.setId(this.provider?.id || this.justCreated$.value);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setTeamId(this.teamId?.value);
        req.setKeyId(this.keyId?.value);
        req.setScopesList(this.scopesList?.value);
        req.setProviderOptions(this.options);

        if (this.updatePrivateKey) {
          req.setPrivateKey(this.privateKey?.value);
        }

        this.loading = true;
        (this.service as ManagementService)
          .updateAppleProvider(req)
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
        const req = new AdminUpdateAppleProviderRequest();
        req.setId(this.provider?.id || this.justCreated$.value);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setTeamId(this.teamId?.value);
        req.setKeyId(this.keyId?.value);
        req.setScopesList(this.scopesList?.value);
        req.setProviderOptions(this.options);

        if (this.updatePrivateKey) {
          req.setPrivateKey(this.privateKey?.value);
        }

        this.loading = true;
        (this.service as AdminService)
          .updateAppleProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.close();
            }, 2000);
          })
          .catch((error) => {
            this.loading = false;
            this.toast.showError(error);
          });
      }
    }
  }

  public close(): void {
    this._location.back();
  }

  public onDropKey(filelist: FileList): void {
    const file = filelist.item(0);
    if (file) {
      if (file.size > MAX_ALLOWED_SIZE) {
        this.toast.showInfo('IDP.APPLE.KEYMAXSIZEEXCEEDED', true);
      } else {
        this.privateKey?.setValue('');
        const reader = new FileReader();
        reader.onload = ((aXML) => {
          return (e) => {
            const keyBase64 = e.target?.result;
            if (keyBase64 && typeof keyBase64 === 'string') {
              const contentType = file.type || 'application/octet-stream';
              const cropped = keyBase64.replace(`data:${contentType};base64,`, '');
              this.privateKey?.setValue(cropped);
            }
          };
        })(file);
        reader.readAsDataURL(file);
      }
    }
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

  public get teamId(): AbstractControl | null {
    return this.form.get('teamId');
  }

  public get keyId(): AbstractControl | null {
    return this.form.get('keyId');
  }

  public get privateKey(): AbstractControl | null {
    return this.form.get('privateKey');
  }

  public get scopesList(): AbstractControl | null {
    return this.form.get('scopesList');
  }
}

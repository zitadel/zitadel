import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute } from '@angular/router';
import { BehaviorSubject, Observable, take } from 'rxjs';
import {
  AddGitLabProviderRequest as AdminAddGitLabProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateGitLabProviderRequest as AdminUpdateGitLabProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {IDPOwnerType, Options, Provider} from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGitLabProviderRequest as MgmtAddGitLabProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateGitLabProviderRequest as MgmtUpdateGitLabProviderRequest,
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
import { ProviderNextDialogComponent } from '../provider-next/provider-next-dialog.component';

@Component({
  selector: 'cnsl-provider-gitlab',
  templateUrl: './provider-gitlab.component.html',
})
export class ProviderGitlabComponent {
  public showOptional: boolean = false;
  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);
  public id: string | null = '';
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public form!: FormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;
  public updateClientSecret: boolean = false;

  private autofillLink$ = new BehaviorSubject<string>('');
  public activateLink$ = new BehaviorSubject<string>('');

  constructor(
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
    private dialog: MatDialog,
    nextSvc: ProviderNextService,
  ) {


    this.form = new FormGroup({
      name: new FormControl('', []),
      clientId: new FormControl('', [requiredValidator]),
      clientSecret: new FormControl('', [requiredValidator]),
      scopesList: new FormControl(['openid', 'profile', 'email'], []),
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
        this.clientSecret?.setValidators([]);
        this.getData(this.id);
      }
    });

    this.setActive(false);
    nextSvc.next(
      'GitLab',
      this.activateLink$,
      this.serviceType === PolicyComponentServiceType.ADMIN,
      'DESCRIPTIONS.SETTINGS.IDPS.CALLBACK.TITLE',
      'DESCRIPTIONS.SETTINGS.IDPS.CALLBACK.DESCRIPTION',
      'https://zitadel.com/docs/guides/integrate/identity-providers/gitlab#gitlab-configuration',
      this.autofillLink$,
      () => [],
    );
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
        if (this.provider?.config?.gitlab) {
          this.form.patchValue(this.provider.config.gitlab);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
    this.service.getLoginPolicy()
      .then((policy) => {
        this.setActive(!!policy.policy?.idpsList.find(idp => idp.idpId === this.id));
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public submitForm(): void {
    this.provider ? this.updateGitlabProvider() : this.addGitlabProvider();
  }

  public addGitlabProvider(): void {
    const req =
      this.serviceType === PolicyComponentServiceType.MGMT
        ? new MgmtAddGitLabProviderRequest()
        : new AdminAddGitLabProviderRequest();

    req.setName(this.name?.value);
    req.setClientId(this.clientId?.value);
    req.setClientSecret(this.clientSecret?.value);
    req.setScopesList(this.scopesList?.value);
    req.setProviderOptions(this.options);

    this.loading = true;
    this.service
      .addGitLabProvider(req)
      .then((addedIDP) => {
        this.showAutofillGuide();
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public activate(id: string) {
    this.service.addIDPToLoginPolicy(id, this.serviceType === PolicyComponentServiceType.ADMIN ? IDPOwnerType.IDP_OWNER_TYPE_SYSTEM : IDPOwnerType.IDP_OWNER_TYPE_ORG).then(() => {
      this.toast.showInfo('POLICY.TOAST.ADDIDP', true);
      this.setActive(true);
      this.id = id;
    });
  }

  public updateGitlabProvider(): void {
    if (this.provider) {
      const req =
        this.serviceType === PolicyComponentServiceType.MGMT
          ? new MgmtUpdateGitLabProviderRequest()
          : new AdminUpdateGitLabProviderRequest();
      req.setId(this.provider.id);
      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setScopesList(this.scopesList?.value);
      req.setProviderOptions(this.options);

      if (this.updateClientSecret) {
        req.setClientSecret(this.clientSecret?.value);
      }

      this.loading = true;
      this.service
        .updateGitLabProvider(req)
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

  private showAutofillGuide(): void {
    this.autofillLink$.next(
      'https://zitadel.com/docs/guides/integrate/identity-providers/gitlab#optional-add-zitadel-action-to-autofill-userdata',
    );
  }

  private setActive(active: boolean) {
    this.activateLink$.next(active ? '' : 'https://zitadel.com/docs/guides/integrate/identity-providers/gitlab#activate-idp' );
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

  public get scopesList(): AbstractControl | null {
    return this.form.get('scopesList');
  }
}

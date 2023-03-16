import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Component } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import {
  AddGenericOIDCProviderRequest as AdminAddGenericOIDCProviderRequest,
  UpdateGenericOIDCProviderRequest as AdminUpdateGenericOIDCProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGenericOIDCProviderRequest as MgmtAddGenericOIDCProviderRequest,
  UpdateGenericOIDCProviderRequest as MgmtUpdateGenericOIDCProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { ProviderService } from '../provider.service';

@Component({
  selector: 'cnsl-provider-oidc',
  templateUrl: './provider-oidc.component.html',
})
export class ProviderOIDCComponent {
  public updateClientSecret: boolean = false;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public form!: UntypedFormGroup;
  public provider?: Provider.AsObject;

  constructor(private providerService: ProviderService, private toast: ToastService) {
    this.form = new UntypedFormGroup({
      name: new UntypedFormControl('', [requiredValidator]),
      clientId: new UntypedFormControl('', [requiredValidator]),
      clientSecret: new UntypedFormControl('', [requiredValidator]),
      issuer: new UntypedFormControl('', [requiredValidator]),
      scopesList: new UntypedFormControl(['openid', 'profile', 'email'], []),
    });
  }

  public submitForm(): void {
    this.provider ? this.updateGenericOIDCProvider() : this.addGenericOIDCProvider();
  }

  public addGenericOIDCProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddGenericOIDCProviderRequest();

      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);

      this.loading = true;
      (this.service as ManagementService)
        .addGenericOIDCProvider(req)
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
      const req = new AdminAddGenericOIDCProviderRequest();
      req.setName(this.name?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);

      this.loading = true;
      (this.service as AdminService)
        .addGenericOIDCProvider(req)
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

  public updateGenericOIDCProvider(): void {
    this.providerService
      .updateProvider('asd', this.form)
      .then(() => {
        this.providerService.navigateBack();
      })
      .catch((error) => {
        this.toast.showError(error);
      });

    if (this.provider) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateGenericOIDCProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);

        this.loading = true;
        (this.service as ManagementService)
          .updateGenericOIDCProvider(req)
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
        const req = new AdminUpdateGenericOIDCProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);

        this.loading = true;
        (this.service as AdminService)
          .updateGenericOIDCProvider(req)
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
  }

  public close(): void {
    this.providerService.navigateBack();
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
}

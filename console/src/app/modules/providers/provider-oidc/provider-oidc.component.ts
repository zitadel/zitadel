import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Component } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { Provider, ProviderType } from 'src/app/proto/generated/zitadel/idp_pb';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

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
    this.providerService
      .addProvider(this.form, ProviderType.PROVIDER_TYPE_OIDC)
      .then((idp) => {
        setTimeout(() => {
          this.close();
        }, 2000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public updateGenericOIDCProvider(): void {
    this.providerService
      .updateProvider('asd', this.form, ProviderType.PROVIDER_TYPE_OIDC)
      .then(() => {
        this.providerService.navigateBack();
      })
      .catch((error) => {
        this.toast.showError(error);
      });

    if (this.provider) {
      this.providerService
        .updateProvider('a', this.form, ProviderType.PROVIDER_TYPE_OIDC)
        .then((idp) => {
          setTimeout(() => {
            this.close();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
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

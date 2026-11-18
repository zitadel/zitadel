import { NgTemplateOutlet } from '@angular/common';
import { Component, Inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { NewAdminService } from 'src/app/services/new-admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { FormFieldModule } from '../form-field/form-field.module';
import { InfoSectionType } from '../info-section/info-section.component';
import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';

export type HttpProviderType = 'email' | 'sms';

export interface HttpProviderDialogData {
  type: HttpProviderType;
  // set when editing an existing provider
  provider?: {
    id: string;
    endpoint: string;
    description: string;
    signingKey?: string;
  };
}

/**
 * Adds or edits an HTTP email or SMS provider. Instead of SMTP or Twilio, ZITADEL then sends every
 * notification as a signed JSON request to the endpoint. A new signing key (on create or rotation)
 * is shown once the provider is saved.
 * Closes with `true` when a provider was saved.
 */
@Component({
  selector: 'cnsl-http-provider-dialog',
  templateUrl: './http-provider-dialog.component.html',
  styleUrls: ['./http-provider-dialog.component.scss'],
  imports: [
    NgTemplateOutlet,
    ReactiveFormsModule,
    MatButtonModule,
    MatCheckboxModule,
    MatDialogModule,
    MatTooltipModule,
    TranslateModule,
    CopyToClipboardModule,
    FormFieldModule,
    InfoSectionModule,
    InputModule,
  ],
})
export class HttpProviderDialogComponent {
  protected readonly InfoSectionType = InfoSectionType;
  protected readonly form;
  protected saving = false;
  protected saved = false;
  protected signingKey = '';
  protected showSigningKey = false;
  protected copied = '';

  constructor(
    private readonly fb: FormBuilder,
    private readonly newAdminService: NewAdminService,
    private readonly toast: ToastService,
    private readonly dialogRef: MatDialogRef<HttpProviderDialogComponent, boolean>,
    @Inject(MAT_DIALOG_DATA) protected readonly data: HttpProviderDialogData,
  ) {
    this.form = this.fb.nonNullable.group({
      endpoint: [
        data.provider?.endpoint ?? '',
        [Validators.required, Validators.maxLength(2048), Validators.pattern(/^https?:\/\/\S+$/)],
      ],
      description: [data.provider?.description ?? '', [Validators.maxLength(200)]],
      rotateSigningKey: [false],
    });
  }

  protected get isEdit(): boolean {
    return !!this.data.provider;
  }

  protected async save(): Promise<void> {
    if (this.form.invalid || this.saving) {
      return;
    }
    const { endpoint, description, rotateSigningKey } = this.form.getRawValue();
    this.saving = true;
    try {
      const signingKey = this.data.provider
        ? await this.update(this.data.provider.id, endpoint, description, rotateSigningKey)
        : await this.add(endpoint, description);
      this.toast.showInfo(this.isEdit ? 'SETTING.HTTPPROVIDER.UPDATED' : 'SETTING.HTTPPROVIDER.ADDED', true);
      this.saved = true;
      if (signingKey) {
        this.signingKey = signingKey;
      } else {
        this.close();
      }
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.saving = false;
    }
  }

  private async add(endpoint: string, description: string): Promise<string> {
    const req = { endpoint, description };
    const resp =
      this.data.type === 'email'
        ? await this.newAdminService.addEmailProviderHTTP(req)
        : await this.newAdminService.addSMSProviderHTTP(req);
    return resp.signingKey;
  }

  private async update(id: string, endpoint: string, description: string, rotateSigningKey: boolean): Promise<string> {
    // the API only accepts an expiration of 0s, which replaces the signing key right away
    const req = { id, endpoint, description, ...(rotateSigningKey ? { expirationSigningKey: { seconds: BigInt(0) } } : {}) };
    const resp =
      this.data.type === 'email'
        ? await this.newAdminService.updateEmailProviderHTTP(req)
        : await this.newAdminService.updateSMSProviderHTTP(req);
    return resp.signingKey ?? '';
  }

  protected close(): void {
    this.dialogRef.close(this.saved);
  }
}

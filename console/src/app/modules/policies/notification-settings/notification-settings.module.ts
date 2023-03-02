import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { InputModule } from '../../input/input.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';
import { NotificationSettingsComponent } from './notification-settings.component';
import { PasswordDialogComponent } from './password-dialog/password-dialog.component';

@NgModule({
  declarations: [NotificationSettingsComponent, DialogAddSMSProviderComponent, PasswordDialogComponent],
  imports: [
    CommonModule,
    CardModule,
    InfoSectionModule,
    FormsModule,
    ReactiveFormsModule,
    HasRolePipeModule,
    MatButtonModule,
    MatCheckboxModule,
    InputModule,
    MatIconModule,
    FormFieldModule,
    WarnDialogModule,
    MatSelectModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    TranslateModule,
  ],
  exports: [NotificationSettingsComponent],
})
export class NotificationSettingsModule {}

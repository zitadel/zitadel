import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { InputModule } from '../../input/input.module';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';
import { NotificationSettingsComponent } from './notification-settings.component';
import { SMTPPasswordDialogComponent } from './smtp-password-dialog/smtp-password-dialog.component';

@NgModule({
  declarations: [NotificationSettingsComponent, DialogAddSMSProviderComponent, SMTPPasswordDialogComponent],
  imports: [
    CommonModule,
    CardModule,
    InfoSectionModule,
    FormsModule,
    ReactiveFormsModule,
    MatButtonModule,
    MatCheckboxModule,
    InputModule,
    MatIconModule,
    FormFieldModule,
    MatSelectModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    TranslateModule,
  ],
  exports: [NotificationSettingsComponent],
})
export class NotificationSettingsModule {}

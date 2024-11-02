import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';

import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { InputModule } from '../../input/input.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';
import { NotificationSMSProviderComponent } from './notification-sms-provider.component';
import { PasswordDialogSMSProviderComponent } from './password-dialog-sms-provider/password-dialog-sms-provider.component';
import { MatDialogModule } from '@angular/material/dialog';

@NgModule({
  declarations: [NotificationSMSProviderComponent, DialogAddSMSProviderComponent, PasswordDialogSMSProviderComponent],
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
    MatDialogModule,
    TranslateModule,
  ],
  exports: [NotificationSMSProviderComponent],
})
export class NotificationSMSProviderModule {}

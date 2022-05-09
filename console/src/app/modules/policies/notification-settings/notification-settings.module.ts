import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { NotificationSettingsComponent } from './notification-settings.component';

@NgModule({
  declarations: [NotificationSettingsComponent],
  imports: [
    CommonModule,
    CardModule,
    FormsModule,
    MatButtonModule,
    FormFieldModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    TranslateModule,
  ],
  exports: [NotificationSettingsComponent],
})
export class NotificationSettingsModule {}

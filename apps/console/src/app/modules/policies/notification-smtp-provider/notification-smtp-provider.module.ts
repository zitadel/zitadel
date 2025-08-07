import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { NotificationSMTPProviderComponent } from './notification-smtp-provider.component';
import { InputModule } from '../../input/input.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { SMTPTableModule } from '../../smtp-table/smtp-table.module';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@NgModule({
  declarations: [NotificationSMTPProviderComponent],
  imports: [
    InputModule,
    FormFieldModule,
    CommonModule,
    MatButtonModule,
    CardModule,
    MatIconModule,
    SMTPTableModule,
    RouterModule,
    HasRolePipeModule,
    MatProgressSpinnerModule,
    TranslateModule,
  ],
  exports: [NotificationSMTPProviderComponent],
})
export class NotificationSMTPProviderModule {}

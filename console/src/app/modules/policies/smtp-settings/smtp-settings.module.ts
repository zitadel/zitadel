import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { SMTPSettingsComponent } from './smtp-settings.component';
import { InputModule } from '../../input/input.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { SMTPTableModule } from '../../smtp-table/smtp-table.module';

@NgModule({
  declarations: [SMTPSettingsComponent],
  imports: [
    InputModule,
    FormFieldModule,
    CommonModule,
    MatLegacyButtonModule,
    CardModule,
    MatIconModule,
    SMTPTableModule,
    RouterModule,
    HasRolePipeModule,
    MatProgressSpinnerModule,
    TranslateModule,
  ],
  exports: [SMTPSettingsComponent],
})
export class SMTPSettingsModule {}

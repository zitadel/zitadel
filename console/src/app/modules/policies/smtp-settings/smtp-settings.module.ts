import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { IdpTableModule } from '../../idp-table/idp-table.module';
import { SMTPSettingsComponent } from './smtp-settings.component';

@NgModule({
  declarations: [SMTPSettingsComponent],
  imports: [
    CommonModule,
    MatLegacyButtonModule,
    CardModule,
    MatIconModule,
    IdpTableModule,
    RouterModule,
    HasRolePipeModule,
    MatProgressSpinnerModule,
    TranslateModule,
  ],
  exports: [SMTPSettingsComponent],
})
export class SMTPSettingsModule {}

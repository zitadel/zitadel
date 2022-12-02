import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';

import { CardModule } from '../../card/card.module';
import { IdpTableModule } from '../../idp-table/idp-table.module';
import { IdpSettingsComponent } from './idp-settings.component';

@NgModule({
  declarations: [IdpSettingsComponent],
  imports: [CommonModule, CardModule, IdpTableModule, MatProgressSpinnerModule, TranslateModule],
  exports: [IdpSettingsComponent],
})
export class IdpSettingsModule {}

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { GeneralSettingsComponent } from './general-settings.component';

@NgModule({
  declarations: [GeneralSettingsComponent],
  imports: [
    CommonModule,
    CardModule,
    FormsModule,
    MatButtonModule,
    FormFieldModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    HasRolePipeModule,
    TranslateModule,
  ],
  exports: [GeneralSettingsComponent],
})
export class GeneralSettingsModule {}

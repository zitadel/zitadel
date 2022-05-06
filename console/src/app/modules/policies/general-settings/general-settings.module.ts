import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { GeneralSettingsComponent } from './general-settings.component';

@NgModule({
  declarations: [GeneralSettingsComponent],
  imports: [
    CommonModule,
    CardModule,
    FormsModule,
    FormFieldModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    TranslateModule,
  ],
  exports: [GeneralSettingsComponent],
})
export class GeneralSettingsModule {}

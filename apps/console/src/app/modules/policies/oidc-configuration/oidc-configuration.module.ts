import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InputModule } from '../../input/input.module';
import { OIDCConfigurationComponent } from './oidc-configuration.component';

@NgModule({
  declarations: [OIDCConfigurationComponent],
  imports: [
    CommonModule,
    CardModule,
    FormsModule,
    MatButtonModule,
    FormFieldModule,
    ReactiveFormsModule,
    InputModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    TranslateModule,
  ],
  exports: [OIDCConfigurationComponent],
})
export class OIDCConfigurationModule {}

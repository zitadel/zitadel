import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
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

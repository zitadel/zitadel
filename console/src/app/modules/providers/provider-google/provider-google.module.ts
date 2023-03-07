import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { CardModule } from '../../card/card.module';
import { CreateLayoutModule } from '../../create-layout/create-layout.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { ProviderOptionsModule } from '../../provider-options/provider-options.module';
import { ProviderGoogleRoutingModule } from './provider-google-routing.module';
import { ProviderGoogleComponent } from './provider-google.component';

@NgModule({
  declarations: [ProviderGoogleComponent],
  imports: [
    ProviderGoogleRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    CreateLayoutModule,
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatSelectModule,
    MatIconModule,
    MatChipsModule,
    CardModule,
    MatCheckboxModule,
    MatTooltipModule,
    TranslateModule,
    ProviderOptionsModule,
    MatLegacyProgressSpinnerModule,
  ],
})
export default class ProviderGoogleModule {}

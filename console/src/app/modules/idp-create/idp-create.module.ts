import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { CardModule } from '../card/card.module';
import { CreateLayoutModule } from '../create-layout/create-layout.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { IdpCreateRoutingModule } from './idp-create-routing.module';
import { IdpCreateComponent } from './idp-create.component';
import { IdpTypeRadioComponent } from './idp-type-radio/idp-type-radio.component';

@NgModule({
  declarations: [IdpCreateComponent, IdpTypeRadioComponent],
  imports: [
    IdpCreateRoutingModule,
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
    MatProgressBarModule,
  ],
})
export default class IdpCreateModule {}

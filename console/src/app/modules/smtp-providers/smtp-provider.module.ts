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

import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { MatStepperModule } from '@angular/material/stepper';
import { CardModule } from '../card/card.module';
import { CreateLayoutModule } from '../create-layout/create-layout.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { ProviderOptionsModule } from '../provider-options/provider-options.module';
import { StringListModule } from '../string-list/string-list.module';
import { SMTPProviderGoogleComponent } from './smtp-provider-google/smtp-provider-google.component';
import { SMTPProvidersRoutingModule } from './smtp-provider-routing.module';
import { SMTPProviderSendgridComponent } from './smtp-provider-sendgrid/smtp-provider-sendgrid.component';

@NgModule({
  declarations: [SMTPProviderGoogleComponent, SMTPProviderSendgridComponent],
  imports: [
    SMTPProvidersRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    CreateLayoutModule,
    StringListModule,
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatProgressBarModule,
    MatSelectModule,
    MatIconModule,
    MatChipsModule,
    CardModule,
    MatCheckboxModule,
    MatTooltipModule,
    MatStepperModule,
    TranslateModule,
    ProviderOptionsModule,
    MatLegacyProgressSpinnerModule,
  ],
})
export default class SMTPProviderModule {}

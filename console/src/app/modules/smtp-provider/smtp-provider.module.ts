import { CommonModule, TitleCasePipe } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { MatStepperModule } from '@angular/material/stepper';
import { CardModule } from '../card/card.module';
import { CreateLayoutModule } from '../create-layout/create-layout.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { ProviderOptionsModule } from '../provider-options/provider-options.module';
import { StringListModule } from '../string-list/string-list.module';
import { SMTPProvidersRoutingModule } from './smtp-provider-routing.module';
import { SMTPProviderComponent } from './smtp-provider.component';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@NgModule({
  declarations: [SMTPProviderComponent],
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
    HasRolePipeModule,
    MatProgressSpinnerModule,
  ],
})
export default class SMTPProviderModule {}

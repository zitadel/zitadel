import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

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
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatSelectModule,
    MatIconModule,
    MatChipsModule,
    MatCheckboxModule,
    MatTooltipModule,
    TranslateModule,
    MatProgressBarModule,
  ],
})
export class IdpCreateModule { }

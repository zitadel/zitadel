import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import {
  TimestampToRetentionPipeModule,
} from 'src/app/pipes/timestamp-to-retention-pipe/timestamp-to-retention-pipe.module';

import { FormFieldModule } from '../form-field/form-field.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { FeaturesRoutingModule } from './features-routing.module';
import { FeaturesComponent } from './features.component';
import { PaymentInfoDialogComponent } from './payment-info-dialog/payment-info-dialog.component';

@NgModule({
  declarations: [
    FeaturesComponent,
    PaymentInfoDialogComponent,
  ],
  imports: [
    FeaturesRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    InputModule,
    MatButtonModule,
    FormFieldModule,
    InputModule,
    HasRoleModule,
    MatSlideToggleModule,
    MatSelectModule,
    MatIconModule,
    HasRoleModule,
    HasRolePipeModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
    InfoSectionModule,
    TranslateModule,
    DetailLayoutModule,
    TimestampToRetentionPipeModule,
  ],
  exports: [
    FeaturesComponent,
  ]
})
export class FeaturesModule { }

import { TextFieldModule } from '@angular/cdk/text-field';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { HasRoleModule } from '../../../directives/has-role/has-role.module';
import { DetailLayoutModule } from '../../../modules/detail-layout/detail-layout.module';
import { InputModule } from '../../../modules/input/input.module';
import { HasFeaturePipeModule } from '../../../pipes/has-feature-pipe/has-feature-pipe.module';
import { HasRolePipeModule } from '../../../pipes/has-role-pipe/has-role-pipe.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { PolicyGridModule } from '../../policy-grid/policy-grid.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { PrivacyPolicyRoutingModule } from './privacy-policy-routing.module';
import { PrivacyPolicyComponent } from './privacy-policy.component';

@NgModule({
  declarations: [PrivacyPolicyComponent],
  imports: [
    PrivacyPolicyRoutingModule,
    MatSelectModule,
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    InputModule,
    FormFieldModule,
    MatButtonModule,
    HasFeaturePipeModule,
    MatIconModule,
    HasRoleModule,
    HasRolePipeModule,
    MatTooltipModule,
    TranslateModule,
    MatTooltipModule,
    DetailLayoutModule,
    MatProgressSpinnerModule,
    TextFieldModule,
    MatDialogModule,
    WarnDialogModule,
    PolicyGridModule,
    InfoSectionModule,
  ],
})
export class PrivacyPolicyModule { }

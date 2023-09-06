import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { InputModule } from '../../input/input.module';
import { SecurityPolicyComponent } from './security-policy.component';

@NgModule({
  declarations: [SecurityPolicyComponent],
  imports: [
    CommonModule,
    CardModule,
    FormsModule,
    InfoSectionModule,
    MatCheckboxModule,
    FormsModule,
    ReactiveFormsModule,
    MatButtonModule,
    FormFieldModule,
    InputModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    HasRolePipeModule,
    MatTooltipModule,
    TranslateModule,
  ],
  exports: [SecurityPolicyComponent],
})
export class SecurityPolicyModule {}

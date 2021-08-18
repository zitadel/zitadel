import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';

import { InfoSectionModule } from '../../info-section/info-section.module';
import { PolicyGridModule } from '../../policy-grid/policy-grid.module';
import { PasswordComplexityPolicyRoutingModule } from './password-complexity-policy-routing.module';
import { PasswordComplexityPolicyComponent } from './password-complexity-policy.component';

@NgModule({
  declarations: [PasswordComplexityPolicyComponent],
  imports: [
    PasswordComplexityPolicyRoutingModule,
    CommonModule,
    FormsModule,
    InputModule,
    MatButtonModule,
    MatSlideToggleModule,
    MatIconModule,
    HasRoleModule,
    MatTooltipModule,
    TranslateModule,
    DetailLayoutModule,
    MatProgressSpinnerModule,
    PolicyGridModule,
    InfoSectionModule,
  ],
})
export class PasswordComplexityPolicyModule { }

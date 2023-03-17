import { A11yModule } from '@angular/cdk/a11y';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';

import { InputModule } from '../input/input.module';
import { OrgContextComponent } from './org-context.component';

@NgModule({
  declarations: [OrgContextComponent],
  imports: [
    CommonModule,
    FormsModule,
    A11yModule,
    ReactiveFormsModule,
    MatIconModule,
    RouterModule,
    MatProgressSpinnerModule,
    MatButtonModule,
    InputModule,
    MatTooltipModule,
    TranslateModule,
    MatButtonModule,
    HasRoleModule,
  ],
  exports: [OrgContextComponent],
})
export class OrgContextModule {}

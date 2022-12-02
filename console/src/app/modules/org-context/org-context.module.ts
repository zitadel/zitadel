import { A11yModule } from '@angular/cdk/a11y';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
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

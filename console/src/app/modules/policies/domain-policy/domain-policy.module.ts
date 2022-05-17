import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { DomainPolicyComponent } from './domain-policy.component';

@NgModule({
  declarations: [DomainPolicyComponent],
  imports: [
    CommonModule,
    FormsModule,
    CardModule,
    InputModule,
    MatButtonModule,
    HasRolePipeModule,
    MatIconModule,
    HasRoleModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    InfoSectionModule,
    MatCheckboxModule,
    TranslateModule,
    DetailLayoutModule,
  ],
  exports: [DomainPolicyComponent],
})
export class DomainPolicyModule {}

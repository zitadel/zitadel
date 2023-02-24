import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { InfoSectionModule } from '../info-section/info-section.module';
import { SettingsGridComponent } from './settings-grid.component';

@NgModule({
  declarations: [SettingsGridComponent],
  imports: [
    CommonModule,
    HasRolePipeModule,
    HasRoleModule,
    TranslateModule,
    RouterModule,
    MatButtonModule,
    MatIconModule,
    MatTooltipModule,
    InfoSectionModule,
  ],
  exports: [SettingsGridComponent],
})
export class SettingsGridModule {}

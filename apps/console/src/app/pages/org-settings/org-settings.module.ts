import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { SettingsListModule } from 'src/app/modules/settings-list/settings-list.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { OrgSettingsRoutingModule } from './org-settings-routing.module';
import { OrgSettingsComponent } from './org-settings.component';

@NgModule({
  declarations: [OrgSettingsComponent],
  imports: [CommonModule, OrgSettingsRoutingModule, SettingsListModule, HasRolePipeModule, TranslateModule],
})
export default class OrgSettingsModule {}

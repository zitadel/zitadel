import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { SettingsListModule } from 'src/app/modules/settings-list/settings-list.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { InstanceSettingsRoutingModule } from './instance-settings-routing.module';
import { InstanceSettingsComponent } from './instance-settings.component';

@NgModule({
  declarations: [InstanceSettingsComponent],
  imports: [CommonModule, InstanceSettingsRoutingModule, SettingsListModule, HasRolePipeModule, TranslateModule],
})
export class InstanceSettingsModule {}

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { InstanceSettingsRoutingModule } from './instance-settings-routing.module';
import { InstanceSettingsComponent } from './instance-settings.component';

@NgModule({
  declarations: [InstanceSettingsComponent],
  imports: [CommonModule, InstanceSettingsRoutingModule, HasRolePipeModule, TranslateModule],
})
export class InstanceSettingsModule {}

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { InstanceSettingsComponent } from './instance-settings.component';

const routes: Routes = [
  {
    path: '',
    component: InstanceSettingsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class InstanceSettingsRoutingModule {}

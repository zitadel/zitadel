import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { InstanceMembersComponent } from './instance-members.component';

const routes: Routes = [
  {
    path: '',
    component: InstanceMembersComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IamMembersRoutingModule {}

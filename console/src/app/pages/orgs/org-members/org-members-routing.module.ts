import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OrgMembersComponent } from './org-members.component';

const routes: Routes = [
  {
    path: '',
    component: OrgMembersComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OrgMembersRoutingModule {}

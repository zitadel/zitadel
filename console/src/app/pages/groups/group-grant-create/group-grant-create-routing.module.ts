import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { GroupGrantCreateComponent } from './group-grant-create.component';

const routes: Routes = [
  {
    path: '',
    component: GroupGrantCreateComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GroupGrantCreateRoutingModule {}

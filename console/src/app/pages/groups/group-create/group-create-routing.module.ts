import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { GroupCreateComponent } from './group-create.component';

const routes: Routes = [
  {
    path: '',
    component: GroupCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GroupCreateRoutingModule {}

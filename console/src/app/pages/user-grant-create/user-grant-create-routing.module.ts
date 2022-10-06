import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { UserGrantCreateComponent } from './user-grant-create.component';

const routes: Routes = [
  {
    path: '',
    component: UserGrantCreateComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class UserGrantCreateRoutingModule {}

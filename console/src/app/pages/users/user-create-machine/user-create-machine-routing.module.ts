import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { UserCreateMachineComponent } from './user-create-machine.component';

const routes: Routes = [
  {
    path: '',
    component: UserCreateMachineComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class UserCreateMachineRoutingModule {}

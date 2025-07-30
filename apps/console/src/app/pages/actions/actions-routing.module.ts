import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ActionsComponent } from './actions.component';

const routes: Routes = [
  {
    path: '',
    component: ActionsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ActionsRoutingModule {}

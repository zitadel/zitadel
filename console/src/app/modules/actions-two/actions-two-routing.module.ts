import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ActionsTwoActionsComponent } from './actions-two-actions/actions-two-actions.component';

const routes: Routes = [
  {
    path: '',
    component: ActionsTwoActionsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ActionsTwoRoutingModule {}
